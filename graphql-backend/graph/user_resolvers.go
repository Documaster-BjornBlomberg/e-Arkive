// Package graph contains all GraphQL-related implementations
package graph

import (
	"context"
	"database/sql"
	"fmt"
	"graphql-backend/graph/model"
	"log"
)

// getUserById retrieves a user by ID
// This is a common function that can be used by both query and mutation resolvers
func (r *mutationResolver) getUserById(ctx context.Context, id string) (*model.User, error) {
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
