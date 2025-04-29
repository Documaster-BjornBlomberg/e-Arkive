package graph

import (
	"context"
	"database/sql"
	"fmt"
	"graphql-backend/graph/model"
	"log"
)

// Group returns GroupResolver implementation
func (r *Resolver) Group() GroupResolver {
	return &groupResolver{r}
}

type GroupResolver interface {
	Members(ctx context.Context, obj *model.Group) ([]*model.User, error)
}

type groupResolver struct{ *Resolver }

// Members resolver for Group type
func (r *groupResolver) Members(ctx context.Context, obj *model.Group) ([]*model.User, error) {
	logAction(fmt.Sprintf("Fetching members for group ID: %s", obj.ID))

	// Query members of the group
	rows, err := r.DB.Query(`
		SELECT u.id, u.username
		FROM users u
		JOIN group_members gm ON u.id = gm.user_id
		WHERE gm.group_id = ?
		ORDER BY u.username ASC
	`, obj.ID)
	if err != nil {
		log.Printf("Error fetching group members: %v", err)
		return nil, fmt.Errorf("failed to fetch group members: %v", err)
	}
	defer rows.Close()

	var members []*model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			log.Printf("Error scanning user row: %v", err)
			return nil, fmt.Errorf("failed to scan user row: %v", err)
		}
		user.Name = user.Username // Using username as name for simplicity
		members = append(members, &user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over member rows: %v", err)
		return nil, fmt.Errorf("failed to iterate over member rows: %v", err)
	}

	return members, nil
}

// OwnerUser resolver for Node type
func (r *nodeResolver) OwnerUser(ctx context.Context, obj *model.Node) (*model.User, error) {
	if obj.OwnerUserID == nil || *obj.OwnerUserID == "" {
		return nil, nil
	}

	logAction(fmt.Sprintf("Fetching owner user for node ID: %s", obj.ID))

	var user model.User
	err := r.DB.QueryRow(`
		SELECT id, username
		FROM users
		WHERE id = ?
	`, *obj.OwnerUserID).Scan(&user.ID, &user.Username)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Printf("Error fetching owner user: %v", err)
		return nil, fmt.Errorf("failed to fetch owner user: %v", err)
	}

	user.Name = user.Username // Using username as name for simplicity
	return &user, nil
}

// OwnerGroup resolver for Node type
func (r *nodeResolver) OwnerGroup(ctx context.Context, obj *model.Node) (*model.Group, error) {
	if obj.OwnerGroupID == nil || *obj.OwnerGroupID == "" {
		return nil, nil
	}

	logAction(fmt.Sprintf("Fetching owner group for node ID: %s", obj.ID))

	var group model.Group
	err := r.DB.QueryRow(`
		SELECT id, name
		FROM groups
		WHERE id = ?
	`, *obj.OwnerGroupID).Scan(&group.ID, &group.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Printf("Error fetching owner group: %v", err)
		return nil, fmt.Errorf("failed to fetch owner group: %v", err)
	}

	return &group, nil
}

// Groups resolver for User type
func (r *userResolver) Groups(ctx context.Context, obj *model.User) ([]*model.Group, error) {
	logAction(fmt.Sprintf("Fetching groups for user ID: %s", obj.ID))

	// Get user ID from token
	currentUserID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// If user is requesting their own groups, or is an admin, proceed
	if currentUserID == obj.ID {
		// Self-lookup is always allowed
	} else {
		// Check if current user is an admin
		var isAdmin bool
		err = r.DB.QueryRow(`
			SELECT EXISTS(
				SELECT 1 
				FROM group_members gm
				JOIN groups g ON gm.group_id = g.id
				WHERE gm.user_id = ? AND g.name = 'Administrators'
			)
		`, currentUserID).Scan(&isAdmin)

		if err != nil {
			log.Printf("Error checking admin status: %v", err)
			return nil, fmt.Errorf("failed to check administrator status: %v", err)
		}

		if !isAdmin {
			return nil, fmt.Errorf("permission denied: can only view your own groups")
		}
	}

	return getUserGroups(ctx, r.DB, obj.ID)
}
