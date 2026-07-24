package auth

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrWeakPassword       = errors.New("password does not meet requirements")
	ErrInvalidToken       = errors.New("invalid token")
	ErrEmailRequired      = errors.New("email required")
	ErrPasswordRequired   = errors.New("password required")
)
