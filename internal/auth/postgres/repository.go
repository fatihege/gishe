package authpostgres

import (
	"context"
	"fmt"

	"github.com/fatihege/gishe/internal/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func (r *Repository) New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user auth.User) (auth.User, error) {
	query := `
		INSERT INTO users (name, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, password_hash, created_at
	`

	var newUser auth.User

	err := r.db.QueryRow(
		ctx, query, user.Name, user.Email, user.PasswordHash, user.CreatedAt,
	).Scan(
		&newUser.ID, &newUser.Name, &newUser.Email, &newUser.PasswordHash, &newUser.CreatedAt,
	)
	if err != nil {
		return auth.User{}, fmt.Errorf("create new user: %v", err)
	}

	return newUser, nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (auth.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM users WHERE email = $1
	`

	var user auth.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt,
	)
	if err != nil {
		return auth.User{}, fmt.Errorf("find user by email: %v", err)
	}

	return user, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id string) (auth.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM users WHERE id = $1
	`

	var user auth.User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt,
	)
	if err != nil {
		return auth.User{}, fmt.Errorf("find user by id: %v", err)
	}

	return user, nil
}
