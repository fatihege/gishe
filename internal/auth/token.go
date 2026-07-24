package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secret    []byte
	issuer    string
	audience  string
	accessTTL time.Duration
}

func NewTokenManager(
	secret []byte,
	issuer string,
	audience string,
) *TokenManager {
	return &TokenManager{
		secret:    secret,
		issuer:    issuer,
		audience:  audience,
		accessTTL: 24 * time.Hour,
	}
}

func (m *TokenManager) NewAccessToken(user User) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(m.accessTTL)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		Issuer:    m.issuer,
		Audience:  jwt.ClaimStrings{m.audience},
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	})

	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %v", err)
	}

	return signedToken, expiresAt, nil
}

func (m *TokenManager) ParseAccessToken(tokenString string) (jwt.RegisteredClaims, error) {
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (any, error) {
			return m.secret, nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
		jwt.WithIssuer(m.issuer),
		jwt.WithAudience(m.audience),
	)
	if err != nil {
		return jwt.RegisteredClaims{}, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return jwt.RegisteredClaims{}, ErrInvalidCredentials
	}

	return claims, nil
}
