// Package graph innehåller alla GraphQL-relaterade implementationer
package graph

import (
	"context"
	"database/sql"
	"fmt"
	"graphql-backend/graph/model"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

// =============================================
// ========== TYPES OCH STRUKTURER ===========
// =============================================

// Resolver är den huvudsakliga resolver-typen för GraphQL
// Hantera alla grafrelaterade funktioner och databasanslutning
type Resolver struct {
	DB *sql.DB
}

// authTokenKey används för att lagra JWT token i context
type authTokenKey struct{}

// NodeResolver interface definition
type NodeResolver interface {
	Children(ctx context.Context, obj *model.Node) ([]*model.Node, error)
	Parent(ctx context.Context, obj *model.Node) (*model.Node, error)
	Files(ctx context.Context, obj *model.Node) ([]*model.File, error)
}

// FileResolver interface definition
type FileResolver interface {
	Node(ctx context.Context, obj *model.File) (*model.Node, error)
}

// Resolver implementations
type nodeResolver struct{ *Resolver }
type fileResolver struct{ *Resolver }
type userResolver struct{ *Resolver }

// =============================================
// ========== PERMISSIONS CONSTANTS ==========
// =============================================

// Permission bit masks
const (
	PERM_VIEW             = 1  // 000001 - Can see the node
	PERM_MODIFY           = 2  // 000010 - Can modify the node (name, etc.)
	PERM_DELETE           = 4  // 000100 - Can delete the node
	PERM_VIEW_PERMISSIONS = 8  // 001000 - Can see user permissions
	PERM_MODIFY_USER      = 16 // 010000 - Can modify user groups, names, passwords
	PERM_MANAGE_USER      = 32 // 100000 - Can delete/create users
)

// =============================================
// ========== GLOBALA VARIABLER ==============
// =============================================

// TokenBlacklist is a map to store invalidated tokens
var TokenBlacklist = make(map[string]bool)
var tokenBlacklistMutex sync.Mutex

// =============================================
// ========== RESOLVER CREATION ==============
// =============================================

// NewResolver skapar en ny resolver med en databasanslutning
func NewResolver(db *sql.DB) *Resolver {
	return &Resolver{DB: db}
}

// =============================================
// ========== LOGGNING =======================
// =============================================

// logAction loggar viktiga händelser i systemet
func logAction(action string) {
	log.Printf("[ACTION] %s", action)
}

// =============================================
// ========== PERMISSIONS FUNCTIONS ==========
// =============================================

// getUserGroups returns all groups that a user belongs to
func getUserGroups(ctx context.Context, db *sql.DB, userID string) ([]*model.Group, error) {
	logAction(fmt.Sprintf("Fetching groups for user ID: %s", userID))

	rows, err := db.Query(`
		SELECT g.id, g.name
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = ?
		ORDER BY g.name ASC
	`, userID)
	if err != nil {
		log.Printf("Error fetching user groups: %v", err)
		return nil, fmt.Errorf("failed to fetch user groups: %v", err)
	}
	defer rows.Close()

	var groups []*model.Group
	for rows.Next() {
		var group model.Group
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			log.Printf("Error scanning group row: %v", err)
			return nil, fmt.Errorf("failed to scan group row: %v", err)
		}
		groups = append(groups, &group)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over group rows: %v", err)
		return nil, fmt.Errorf("failed to iterate over group rows: %v", err)
	}

	return groups, nil
}

// checkPermission checks if a user has the specified permission for a node
func checkPermission(ctx context.Context, db *sql.DB, nodeID string, permissionBit int) (bool, error) {
	// Get user ID from JWT token
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return false, err
	}

	// First, check if the user is an admin
	var isAdmin bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM group_members gm
			JOIN groups g ON gm.group_id = g.id
			WHERE gm.user_id = ? AND g.name = 'Administrators'
		)
	`, userID).Scan(&isAdmin)

	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		return false, fmt.Errorf("failed to check administrator status: %v", err)
	}

	if isAdmin {
		// Admins have all permissions
		return true, nil
	}

	// Check if the user is the owner of the node
	var hasPermission bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM nodes 
			WHERE id = ? AND owner_user_id = ? AND (permissions & ?) > 0
		)
	`, nodeID, userID, permissionBit).Scan(&hasPermission)

	if err != nil {
		log.Printf("Error checking node ownership permission: %v", err)
		return false, fmt.Errorf("failed to check node ownership: %v", err)
	}

	if hasPermission {
		return true, nil
	}

	// Check if the user is in a group that has the required permission
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM nodes n
			JOIN group_members gm ON n.owner_group_id = gm.group_id
			WHERE n.id = ? AND gm.user_id = ? AND (n.permissions & ?) > 0
		)
	`, nodeID, userID, permissionBit).Scan(&hasPermission)

	if err != nil {
		log.Printf("Error checking group permission: %v", err)
		return false, fmt.Errorf("failed to check group permission: %v", err)
	}

	return hasPermission, nil
}

// getNodeWithPermissions fetches a node and checks if the user has permission to view it
func getNodeWithPermissions(ctx context.Context, db *sql.DB, id string) (*model.Node, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	// First check if the user has permission to view this node
	hasPermission, err := checkPermission(ctx, db, id, PERM_VIEW)
	if err != nil {
		return nil, err
	}

	if !hasPermission {
		return nil, fmt.Errorf("permission denied: cannot view this node")
	}

	// Query the node with ownership and permission information
	row := db.QueryRow(`
		SELECT id, name, parent_id, owner_user_id, owner_group_id, permissions, created_at, updated_at
		FROM nodes
		WHERE id = ?
	`, id)

	var node model.Node
	var parentID sql.NullString
	var ownerUserID sql.NullString
	var ownerGroupID sql.NullString

	err = row.Scan(&node.ID, &node.Name, &parentID, &ownerUserID, &ownerGroupID,
		&node.Permissions, &node.CreatedAt, &node.UpdatedAt)

	if err == sql.ErrNoRows {
		log.Printf("Node with ID %s not found", id)
		return nil, fmt.Errorf("node not found")
	} else if err != nil {
		log.Printf("Error fetching node with ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to fetch node: %v", err)
	}

	// Set nullable fields
	if parentID.Valid {
		parentIDStr := parentID.String
		node.ParentID = &parentIDStr
	}

	if ownerUserID.Valid {
		ownerUserIDStr := ownerUserID.String
		node.OwnerUserID = &ownerUserIDStr
	}

	if ownerGroupID.Valid {
		ownerGroupIDStr := ownerGroupID.String
		node.OwnerGroupID = &ownerGroupIDStr
	}

	return &node, nil
}

// getChildNodesWithPermissions fetches all child nodes of a parent that the user has permission to view
func getChildNodesWithPermissions(ctx context.Context, db *sql.DB, parentID string) ([]*model.Node, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	// First check if the parent node exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM nodes WHERE id = ?)", parentID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if parent node exists: %v", err)
		return nil, fmt.Errorf("failed to check if parent node exists: %v", err)
	}

	if !exists {
		log.Printf("Parent node with ID %s does not exist", parentID)
		return nil, fmt.Errorf("parent node not found")
	}

	// Get the user ID from context
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if the user is an admin
	var isAdmin bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM group_members gm
			JOIN groups g ON gm.group_id = g.id
			WHERE gm.user_id = ? AND g.name = 'Administrators'
		)
	`, userID).Scan(&isAdmin)

	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		return nil, fmt.Errorf("failed to check administrator status: %v", err)
	}

	var rows *sql.Rows

	if isAdmin {
		// Admins can see all nodes
		rows, err = db.Query(`
			SELECT id, name, parent_id, owner_user_id, owner_group_id, permissions, created_at, updated_at 
			FROM nodes 
			WHERE parent_id = ?
			ORDER BY name ASC
		`, parentID)
	} else {
		// Regular users can only see nodes they have permission to view
		rows, err = db.Query(`
			SELECT n.id, n.name, n.parent_id, n.owner_user_id, n.owner_group_id, n.permissions, n.created_at, n.updated_at 
			FROM nodes n
			WHERE n.parent_id = ? AND (
				(n.owner_user_id = ? AND (n.permissions & ?) > 0) OR
				EXISTS (
					SELECT 1 
					FROM group_members gm 
					WHERE gm.user_id = ? AND gm.group_id = n.owner_group_id AND (n.permissions & ?) > 0
				)
			)
			ORDER BY n.name ASC
		`, parentID, userID, PERM_VIEW, userID, PERM_VIEW)
	}

	if err != nil {
		log.Printf("Error fetching child nodes: %v", err)
		return nil, fmt.Errorf("failed to fetch child nodes: %v", err)
	}
	defer rows.Close()

	var nodes []*model.Node
	for rows.Next() {
		var node model.Node
		var pID sql.NullString
		var ownerUserID sql.NullString
		var ownerGroupID sql.NullString

		err := rows.Scan(&node.ID, &node.Name, &pID, &ownerUserID, &ownerGroupID,
			&node.Permissions, &node.CreatedAt, &node.UpdatedAt)

		if err != nil {
			log.Printf("Error scanning node: %v", err)
			return nil, fmt.Errorf("failed to process node data: %v", err)
		}

		if pID.Valid {
			parentIDStr := pID.String
			node.ParentID = &parentIDStr
		}

		if ownerUserID.Valid {
			ownerUserIDStr := ownerUserID.String
			node.OwnerUserID = &ownerUserIDStr
		}

		if ownerGroupID.Valid {
			ownerGroupIDStr := ownerGroupID.String
			node.OwnerGroupID = &ownerGroupIDStr
		}

		nodes = append(nodes, &node)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating through rows: %v", err)
		return nil, fmt.Errorf("error processing results: %v", err)
	}

	return nodes, nil
}

// =============================================
// ========== NODE FUNKTIONER ===============
// =============================================

// getNodeById hämtar en nod från databasen baserat på ID
func getNodeById(ctx context.Context, db *sql.DB, id string) (*model.Node, error) {
	// Use the new function that includes permissions
	return getNodeWithPermissions(ctx, db, id)
}

// scanNodeRows läser noddata från SQL-rader och konverterar till Node-objekt
func scanNodeRows(rows *sql.Rows) ([]*model.Node, error) {
	var nodes []*model.Node
	for rows.Next() {
		var node model.Node
		var parentID sql.NullString
		var createdAt, updatedAt string

		if err := rows.Scan(&node.ID, &node.Name, &parentID, &createdAt, &updatedAt); err != nil {
			log.Printf("Error scanning node row: %v", err)
			return nil, fmt.Errorf("failed to scan node row: %v", err)
		}

		node.CreatedAt = createdAt
		node.UpdatedAt = updatedAt
		if parentID.Valid {
			parentIDStr := parentID.String
			node.ParentID = &parentIDStr
		}

		nodes = append(nodes, &node)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over node rows: %v", err)
		return nil, fmt.Errorf("failed to iterate over node rows: %v", err)
	}

	return nodes, nil
}

// detectCycle checks if updating a node's parent would create a cycle in the node hierarchy
func (r *Resolver) detectCycle(nodeID string, newParentID string, isCycle *bool) error {
	*isCycle = false

	// If the new parent is the same as the node itself, it's a cycle
	if nodeID == newParentID {
		*isCycle = true
		return nil
	}

	// Check if any ancestors of the potential parent are actually the node we're moving
	currentID := newParentID
	for currentID != "" {
		// If currentID ever becomes the node we're moving, we have a cycle
		if currentID == nodeID {
			*isCycle = true
			return nil
		}

		// Get parent of current node
		var parentID sql.NullString
		err := r.DB.QueryRow("SELECT parent_id FROM nodes WHERE id = ?", currentID).Scan(&parentID)
		if err == sql.ErrNoRows {
			// Node not found, break the loop
			break
		} else if err != nil {
			return fmt.Errorf("error checking hierarchy: %v", err)
		}

		// If no parent, break the loop
		if !parentID.Valid {
			break
		}

		// Move up the hierarchy
		currentID = parentID.String
	}

	return nil
}

// =============================================
// ========== FILE FUNKTIONER ===============
// =============================================

// Implementering av GetFilesByNodeId för queryResolver
func (r *queryResolver) GetFilesByNodeId(ctx context.Context, nodeID string) ([]*model.File, error) {
	logAction(fmt.Sprintf("Fetching files for node ID: %s", nodeID))

	if r.DB == nil {
		log.Printf("Database connection is nil")
		return nil, fmt.Errorf("internal server error: database connection is not initialized")
	}

	// First check if user has permission to view this node
	hasPermission, err := checkPermission(ctx, r.DB, nodeID, PERM_VIEW)
	if err != nil {
		return nil, err
	}

	if !hasPermission {
		return nil, fmt.Errorf("permission denied: cannot view files in this node")
	}

	// Query files for this specific node
	rows, err := r.DB.Query(`
		SELECT id, name, size, content_type, created_at, file_data
		FROM files
		WHERE node_id = ?
		ORDER BY name ASC
	`, nodeID)
	if err != nil {
		log.Printf("Error fetching files for node ID %s: %v", nodeID, err)
		return nil, fmt.Errorf("failed to fetch files: %v", err)
	}
	defer rows.Close()

	var files []*model.File
	for rows.Next() {
		var file model.File
		var createdAt string
		var fileData []byte
		if err := rows.Scan(&file.ID, &file.Name, &file.Size, &file.ContentType, &createdAt, &fileData); err != nil {
			log.Printf("Error scanning file row: %v", err)
			return nil, fmt.Errorf("failed to scan file row: %v", err)
		}
		file.CreatedAt = createdAt

		// Don't send file data in listing, only metadata
		// We'll only include the file data when explicitly downloading the file

		// Set the nodeId on the file
		file.NodeID = &nodeID

		// Fetch metadata for each file
		metaRows, err := r.DB.Query("SELECT key, value FROM metadata WHERE file_id = ?", file.ID)
		if err != nil {
			log.Printf("Error fetching metadata for file ID %s: %v", file.ID, err)
			return nil, fmt.Errorf("failed to fetch metadata: %v", err)
		}

		var metadata []*model.Metadata
		for metaRows.Next() {
			var meta model.Metadata
			if err := metaRows.Scan(&meta.Key, &meta.Value); err != nil {
				log.Printf("Error scanning metadata row: %v", err)
				metaRows.Close()
				return nil, fmt.Errorf("failed to scan metadata row: %v", err)
			}
			metadata = append(metadata, &meta)
		}
		metaRows.Close()

		file.Metadata = metadata
		files = append(files, &file)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over file rows: %v", err)
		return nil, fmt.Errorf("failed to iterate over file rows: %v", err)
	}

	log.Printf("Successfully fetched %d files for node ID %s", len(files), nodeID)
	return files, nil
}

// =============================================
// ========== AUTENTISERING ==================
// =============================================

// WithAuthToken lägger till JWT token i context
func WithAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, authTokenKey{}, token)
}

// IsTokenBlacklisted kontrollerar om en token är i blacklist
func IsTokenBlacklisted(token string) bool {
	tokenBlacklistMutex.Lock()
	defer tokenBlacklistMutex.Unlock()
	return TokenBlacklist[token]
}

// BlacklistToken lägger till en token i blacklist
func BlacklistToken(token string) {
	tokenBlacklistMutex.Lock()
	defer tokenBlacklistMutex.Unlock()
	TokenBlacklist[token] = true
}

// generateJWT genererar en JWT-token för en användare
func generateJWT(userID, username string) (string, error) {
	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // Token expires in 7 days
	})

	// Sign the token with the secret key
	// In a production environment, this should be stored securely and not hardcoded
	jwtSecret := []byte("your-256-bit-secret") // Change this to a secure secret in production
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// validateJWT validerar en JWT-token och kontrollerar om den är blacklistad
func validateJWT(tokenString string) (jwt.MapClaims, error) {
	// Check if token is blacklisted
	if IsTokenBlacklisted(tokenString) {
		return nil, fmt.Errorf("token has been invalidated")
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		jwtSecret := []byte("your-256-bit-secret") // Change this to a secure secret in production
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetAuthToken hämtar JWT-token från context och validerar den
func GetAuthToken(ctx context.Context) (string, bool) {
	// Try both Authorization and Authenticate headers
	authHeader, ok := ctx.Value("Authorization").(string)
	if !ok || authHeader == "" {
		// Try Authenticate header if Authorization is not present
		authHeader, ok = ctx.Value("Authenticate").(string)
		if !ok || authHeader == "" {
			return "", false
		}
	}

	// Remove "Bearer " prefix if present
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", false
	}

	return token, true
}

// getUserIDFromContext extracts user ID from the JWT token in the context
func getUserIDFromContext(ctx context.Context) (string, error) {
	// Get the token from the context
	token, ok := GetAuthToken(ctx)
	if !ok {
		return "", fmt.Errorf("not authenticated")
	}

	// Parse and validate the token
	claims, err := validateJWT(token)
	if err != nil {
		log.Printf("Error validating JWT token: %v", err)
		return "", fmt.Errorf("not authenticated")
	}

	// Get user ID from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		log.Printf("Invalid token claims: user_id not found")
		return "", fmt.Errorf("invalid authentication token")
	}

	return userID, nil
}

// =============================================
// ========== USER SETTINGS =================
// =============================================

// Settings hämtar användarinställningar för en användare
func (r *userResolver) Settings(ctx context.Context, obj *model.User) ([]*model.UserSetting, error) {
	logAction(fmt.Sprintf("Fetching settings for user: %s", obj.ID))

	// Query the database for all settings for this user
	rows, err := r.DB.Query(`
		SELECT id, key, value, created_at, updated_at
		FROM user_settings
		WHERE user_id = ?
		ORDER BY key ASC
	`, obj.ID)

	if err != nil {
		log.Printf("Error fetching user settings: %v", err)
		return nil, fmt.Errorf("failed to fetch user settings: %v", err)
	}
	defer rows.Close()

	var settings []*model.UserSetting
	for rows.Next() {
		var setting model.UserSetting
		var createdAt, updatedAt string
		if err := rows.Scan(&setting.ID, &setting.Key, &setting.Value, &createdAt, &updatedAt); err != nil {
			log.Printf("Error scanning user setting row: %v", err)
			return nil, fmt.Errorf("failed to read user settings: %v", err)
		}
		setting.CreatedAt = createdAt
		setting.UpdatedAt = updatedAt
		settings = append(settings, &setting)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating through user settings rows: %v", err)
		return nil, fmt.Errorf("error reading user settings: %v", err)
	}

	return settings, nil
}

// =============================================
// ========== CORE RESOLVER GETTERS ===========
// =============================================

func (r *Resolver) GetNodeById(ctx context.Context, id string) (*model.Node, error) {
	return getNodeById(ctx, r.DB, id)
}

// =============================================
// ========== USER FUNKTIONER ===========
// =============================================

// getUserById retrieves a user by ID
func (r *queryResolver) getUserById(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRow("SELECT id, username, name FROM users WHERE id = ?", id).Scan(
		&user.ID, &user.Username, &user.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("Error querying user by ID: %v", err)
		return nil, fmt.Errorf("failed to query user: %v", err)
	}

	// Get user's groups
	groups, err := getUserGroups(ctx, r.DB, id)
	if err != nil {
		log.Printf("Warning: Failed to fetch user's groups: %v", err)
		// Continue without groups, it's not critical
	} else {
		user.Groups = groups
	}

	return &user, nil
}
