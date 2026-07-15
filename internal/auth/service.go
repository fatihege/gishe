package auth

import (
	"context"
	"errors"
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

	foundUser, err := s.repository.FindUserByEmail(ctx, email)
	if foundUser.ID != "" {
		return User{}, ErrEmailAlreadyExists
	}

	if len(input.Password) < 4 {
		return User{}, ErrWeakPassword
	}

	passwordHash, err := s.passwords.Hash(input.Password)
	if err != nil {
		return User{}, err
	}

	user, err := s.repository.CreateUser(ctx, User{
		Name: name, Email: email, PasswordHash: passwordHash, CreatedAt: time.Now(),
	})
	if err != nil {
		return User{}, err
	}

	return user, nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (s *Service) Login(ctx context.Context, input LoginInput) (User, error) {
	email := normalizeEmail(input.Email)
	if email == "" {
		return User{}, fmt.Errorf("email required")
	}

	if input.Password == "" {
		return User{}, fmt.Errorf("password required")
	}

	user, err := s.repository.FindUserByEmail(ctx, email)
	if errors.Is(err, ErrUserNotFound) {
		return User{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, err
	}

	match, err := s.passwords.Compare(input.Password, user.PasswordHash)
	if err != nil {
		return User{}, err
	}

	if !match {
		return User{}, ErrInvalidCredentials
	}

	return user, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
