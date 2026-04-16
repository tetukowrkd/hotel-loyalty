package middleware

import (
	"context"
	"net/http"
	"strings"

	"hotel-loyalty/internal/infrastructure/jwt"
	"hotel-loyalty/internal/infrastructure/logger"
)

type contextKey string

const UserContextKey = contextKey("user")

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.ErrorLogger.Println("[MIDDLEWARE][JWT] missing Authorization header",
				"path=", r.URL.Path,
				"method=", r.Method,
			)

			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// 🔍 validasi format "Bearer <token>"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.ErrorLogger.Println("[MIDDLEWARE][JWT] invalid Authorization format",
				"path=", r.URL.Path,
				"method=", r.Method,
			)

			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			logger.ErrorLogger.Println("[MIDDLEWARE][JWT] invalid token",
				"error=", err,
				"path=", r.URL.Path,
				"method=", r.Method,
			)

			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 🔥 log sukses (tanpa token)
		logger.InfoLogger.Println("[MIDDLEWARE][JWT] token validated",
			"userID=", claims.UserID,
			"role=", claims.Role,
			"path=", r.URL.Path,
			"method=", r.Method,
		)

		// simpan ke context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
