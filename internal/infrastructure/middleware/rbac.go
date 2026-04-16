package middleware

import (
	"net/http"

	"hotel-loyalty/internal/infrastructure/jwt"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/model/response"
)

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, ok := r.Context().Value(UserContextKey).(*jwt.Claims)
			if !ok {
				logger.ErrorLogger.Println("[MIDDLEWARE][RBAC] unauthorized - no claims")

				response.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"message": "Unauthorized",
				})
				return
			}

			userRole := claims.Role
			userID := claims.UserID

			// 🔍 log request role check
			logger.InfoLogger.Println("[MIDDLEWARE][RBAC] checking role:",
				"userID=", userID,
				"role=", userRole,
				"allowed=", allowedRoles,
				"path=", r.URL.Path,
			)

			for _, role := range allowedRoles {
				if userRole == role {
					logger.InfoLogger.Println("[MIDDLEWARE][RBAC] access granted:",
						"userID=", userID,
						"role=", userRole,
						"path=", r.URL.Path,
					)

					next.ServeHTTP(w, r)
					return
				}
			}

			// ❌ forbidden
			logger.ErrorLogger.Println("[MIDDLEWARE][RBAC] forbidden:",
				"userID=", userID,
				"role=", userRole,
				"allowed=", allowedRoles,
				"path=", r.URL.Path,
			)

			response.WriteJSON(w, http.StatusForbidden, map[string]string{
				"message": "Forbidden",
			})
		})
	}
}
