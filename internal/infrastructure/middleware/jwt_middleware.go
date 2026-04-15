package middleware

import (
	"context"
	"net/http"
	"strings"

	"hotel-loyalty/internal/infrastructure/jwt"
)

type contextKey string

const UserContextKey = contextKey("user")

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// format: Bearer xxx
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// simpan ke context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
