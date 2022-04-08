package graphql

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/jjvanvark/lor-deluxe/internal/models"
)

type userResolver struct {
	user *models.User
}

func (r *userResolver) Id(ctx context.Context) graphql.ID {
	var result graphql.ID = intToId(r.user.ID)
	return result
}

func (r *userResolver) Name(ctx context.Context) string {
	return r.user.Name
}

func (r *userResolver) Email(ctx context.Context) string {
	return r.user.Email
}

func (r *userResolver) Key(ctx context.Context) string {
	return r.user.Key
}
