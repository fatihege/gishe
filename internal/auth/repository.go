package auth

import (
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, user RegisterInputHashed) (User, error)

	FindUserByEmail(ctx context.Context, email string) (User, error)

	FindUserByID(ctx context.Context, id string) (User, error)
}
