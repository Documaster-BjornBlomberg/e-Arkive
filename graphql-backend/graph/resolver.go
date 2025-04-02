// Package graph innehåller alla GraphQL-relaterade implementationer
package graph

import "database/sql"

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

// Mutation returnerar implementationen av MutationResolver
// Detta används för alla mutationer (skrivoperationer) i GraphQL-schemat
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returnerar implementationen av QueryResolver
// Detta används för alla queries (läsoperationer) i GraphQL-schemat
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Todo returnerar implementationen av TodoResolver
// Detta är en del av exempel-implementationen och kan tas bort om det inte behövs
func (r *Resolver) Todo() TodoResolver {
	return &todoResolver{r}
}

// NewResolver skapar en ny instans av Resolver med den givna databasen.
func NewResolver(db *sql.DB) *Resolver {
	return &Resolver{DB: db}
}

// Ensure the Resolver struct is properly defined with a DB field.
type Resolver struct {
	DB *sql.DB
}
