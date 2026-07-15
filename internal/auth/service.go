package auth

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	repository Repository
	passwords  *PasswordHasher
}

func NewService(repository Repository, passwords *PasswordHasher) *Service {
	return &Service{
		repository: repository,
		passwords:  passwords,
	}
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (User, error) {
	name := strings.TrimSpace(input.Name)
	email := normalizeEmail(input.Email)
	if email == "" {
		return User{}, fmt.Errorf("email required")
	}

	if len(input.Password) < 4 {
		return User{}, fmt.Errorf("weak password")
	}

	passwordHash, err := s.passwords.Hash(input.Password)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %v", err)
	}

	user, err := s.repository.CreateUser(ctx, User{
		Name: name, Email: email, PasswordHash: passwordHash, CreatedAt: time.Now(),
	})
	if err != nil {
		return User{}, fmt.Errorf("create user: %v", err)
	}

	return user, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
