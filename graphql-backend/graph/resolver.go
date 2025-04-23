// Package graph innehåller alla GraphQL-relaterade implementationer
package graph

import (
	"context"
	"database/sql"
	"fmt"
	"graphql-backend/graph/model"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver är den huvudsakliga resolver-typen för GraphQL
// Hantera alla grafrelaterade funktioner och databasanslutning
type Resolver struct {
	DB *sql.DB
}

// NewResolver skapar en ny resolver med en databasanslutning
func NewResolver(db *sql.DB) *Resolver {
	return &Resolver{DB: db}
}

// AuthTokenKey används för att lagra JWT token i context
type authTokenKey struct{}

// GetAuthToken extraherar JWT token från context
func GetAuthToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(authTokenKey{}).(string)
	return token, ok
}

// WithAuthToken lägger till JWT token i context
func WithAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, authTokenKey{}, token)
}

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

func logAction(action string) {
	log.Printf("[ACTION] %s", action)
}

func (r *queryResolver) GetFilesByNodeId(ctx context.Context, nodeID string) ([]*model.File, error) {
	logAction(fmt.Sprintf("Fetching files for node ID: %s", nodeID))

	if r.DB == nil {
		log.Printf("Database connection is nil")
		return nil, fmt.Errorf("internal server error: database connection is not initialized")
	}

	// Verify the node exists
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM nodes WHERE id = ?)", nodeID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if node exists: %v", err)
		return nil, fmt.Errorf("failed to verify node: %v", err)
	}
	if !exists {
		log.Printf("Node with ID %s does not exist", nodeID)
		return nil, fmt.Errorf("node not found")
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
func (r *queryResolver) GetNodeById(ctx context.Context, id string) (*model.Node, error) {
	return getNodeById(ctx, r.DB, id)
}
func getNodeById(ctx context.Context, db *sql.DB, id string) (*model.Node, error) {
	logAction(fmt.Sprintf("Fetching node with ID: %s", id))

	if db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	row := db.QueryRow(`
		SELECT id, name, parent_id, created_at, updated_at
		FROM nodes
		WHERE id = ?
	`, id)

	var node model.Node
	var parentID sql.NullString
	err := row.Scan(&node.ID, &node.Name, &parentID, &node.CreatedAt, &node.UpdatedAt)

	if err == sql.ErrNoRows {
		log.Printf("Node with ID %s not found", id)
		return nil, fmt.Errorf("node not found")
	} else if err != nil {
		log.Printf("Error fetching node with ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to fetch node: %v", err)
	}

	if parentID.Valid {
		parentIDStr := parentID.String
		node.ParentID = &parentIDStr
	}

	// Hämta alla barn till denna nod
	rows, err := db.Query(`
		SELECT id, name, parent_id, created_at, updated_at
		FROM nodes
		WHERE parent_id = ?
		ORDER BY name ASC
	`, node.ID)
	if err != nil {
		log.Printf("Error fetching children for node ID %s: %v", node.ID, err)
		return &node, nil // Returnera noden även om vi inte kunde hämta barnen
	}
	defer rows.Close()

	children, err := scanNodeRows(rows)
	if err != nil {
		log.Printf("Error scanning child rows: %v", err)
		return &node, nil // Returnera noden även om vi inte kunde läsa barnen
	}

	node.Children = children
	return &node, nil
}

type nodeResolver struct{ *Resolver }
type fileResolver struct{ *Resolver }

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
func validateJWT(tokenString string) (jwt.MapClaims, error) {
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
func (r *Resolver) detectCycle(nodeID string, newParentID string, isCycle *bool) error {
	// Om den nya föräldern är samma som noden själv, är det en cykel
	if nodeID == newParentID {
		*isCycle = true
		return nil
	}

	// Kontrollera om den nya föräldern har den aktuella noden som förfader
	var currentParentID sql.NullString
	err := r.DB.QueryRow("SELECT parent_id FROM nodes WHERE id = ?", newParentID).Scan(&currentParentID)
	if err == sql.ErrNoRows {
		// Om föräldern inte finns, är det inget problem
		*isCycle = false
		return nil
	} else if err != nil {
		return err
	}

	// Om förälderns förälder är NULL, är det slutet på kedjan - ingen cykel
	if !currentParentID.Valid {
		*isCycle = false
		return nil
	}

	// Om förälderns förälder är den aktuella noden, är det en cykel
	if currentParentID.String == nodeID {
		*isCycle = true
		return nil
	}

	// Fortsätt rekursivt upp i hierarkin
	return r.detectCycle(nodeID, currentParentID.String, isCycle)
}
