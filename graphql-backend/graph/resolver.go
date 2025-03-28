package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

// Mutation returns the MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returns the QueryResolver implementation.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Todo() TodoResolver {
	return &todoResolver{}
}
