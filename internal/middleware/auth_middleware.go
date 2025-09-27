package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/deepraj02/go-postgres-starter/internal/store"
	"github.com/deepraj02/go-postgres-starter/internal/utils/auth"
	utils "github.com/deepraj02/go-postgres-starter/internal/utils/json"
)

type AuthMiddleware struct {
	AuthStore store.AuthStore
}

type contextKey string

const UserContextKey = contextKey("user")

func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{
				"error": "Authorization header required",
			})
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{
				"error": "Invalid authorization header format",
			})
			return
		}

		token := parts[1]
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{
				"error": "Invalid or expired token",
			})
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) *auth.Claims {
	if claims, ok := r.Context().Value(UserContextKey).(*auth.Claims); ok {
		return claims
	}
	return nil
}
