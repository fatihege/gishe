package auth

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, user RegisterInputHashed) (User, error)

	FindUserByEmail(ctx context.Context, email string) (User, error)

	FindUserByID(ctx context.Context, id uuid.UUID) (User, error)
}
