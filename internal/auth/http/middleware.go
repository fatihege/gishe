package authhttp

import (
	"context"
	"net/http"
	"strings"

	"github.com/fatihege/gishe/internal/auth"
	"github.com/fatihege/gishe/internal/httpx"
	"github.com/google/uuid"
)

type contextKey string

const (
	userIDKey contextKey = "authenticated-user-id"
)

type Middleware struct {
	tokens *auth.TokenManager
}

func NewMiddleware(tokens *auth.TokenManager) *Middleware {
	return &Middleware{
		tokens: tokens,
	}
}

func (m *Middleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		const prefix = "Bearer "

		if !strings.HasPrefix(header, prefix) {
			httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(header, prefix))

		claims, err := m.tokens.ParseAccessToken(tokenString)
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid access token")
			return
		}

		rawUserID := claims.Subject
		userID, err := uuid.Parse(rawUserID)
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid user id")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
