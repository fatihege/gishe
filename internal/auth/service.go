package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Service struct {
	repository Repository
	passwords  *PasswordHasher
	tokens     *TokenManager
}

func NewService(repository Repository, passwords *PasswordHasher, tokens *TokenManager) *Service {
	return &Service{
		repository: repository,
		passwords:  passwords,
		tokens:     tokens,
	}
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (User, Token, error) {
	name := strings.TrimSpace(input.Name)
	email := normalizeEmail(input.Email)
	if email == "" {
		return User{}, Token{}, fmt.Errorf("email required")
	}

	foundUser, err := s.repository.FindUserByEmail(ctx, email)
	if foundUser.ID != "" {
		return User{}, Token{}, ErrEmailAlreadyExists
	}

	if len(input.Password) < 4 {
		return User{}, Token{}, ErrWeakPassword
	}

	passwordHash, err := s.passwords.Hash(input.Password)
	if err != nil {
		return User{}, Token{}, err
	}

	user, err := s.repository.CreateUser(ctx, User{
		Name: name, Email: email, PasswordHash: passwordHash, CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return User{}, Token{}, err
	}

	tokenString, expiresAt, err := s.tokens.NewAccessToken(user)
	if err != nil {
		return User{}, Token{}, err
	}

	token := Token{
		AccessToken: tokenString,
		ExpiresAt:   expiresAt,
	}

	return user, token, nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (s *Service) Login(ctx context.Context, input LoginInput) (User, Token, error) {
	email := normalizeEmail(input.Email)
	if email == "" {
		return User{}, Token{}, fmt.Errorf("email required")
	}

	if input.Password == "" {
		return User{}, Token{}, fmt.Errorf("password required")
	}

	user, err := s.repository.FindUserByEmail(ctx, email)
	if errors.Is(err, ErrUserNotFound) {
		return User{}, Token{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, Token{}, err
	}

	match, err := s.passwords.Compare(input.Password, user.PasswordHash)
	if err != nil {
		return User{}, Token{}, err
	}

	if !match {
		return User{}, Token{}, ErrInvalidCredentials
	}

	tokenString, expiresAt, err := s.tokens.NewAccessToken(user)
	if err != nil {
		return User{}, Token{}, err
	}

	token := Token{
		AccessToken: tokenString,
		ExpiresAt:   expiresAt,
	}

	claims, err := s.tokens.ParseAccessToken(token.AccessToken)
	if err != nil {
		log.Printf("parse acces token: %v", err)
		return User{}, Token{}, err
	}

	fmt.Println(claims.Subject)

	return user, token, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
