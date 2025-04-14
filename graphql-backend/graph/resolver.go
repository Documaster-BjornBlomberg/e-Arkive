// Package graph innehåller alla GraphQL-relaterade implementationer
package graph

import (
	"context"
	"database/sql"
	"graphql-backend/graph/model"
)

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

// Node returnerar implementationen av NodeResolver
// Detta används för att lösa "parent" och "children" fält i Node-typen
func (r *Resolver) Node() NodeResolver {
	return &nodeResolver{r}
}

// File returnerar implementationen av FileResolver
// Detta används för att lösa "node" fält i File-typen
func (r *Resolver) File() FileResolver {
	return &fileResolver{r}
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
