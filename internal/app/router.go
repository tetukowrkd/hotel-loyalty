package app

import (
	"net/http"

	"hotel-loyalty/internal/infrastructure/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRouter(c *Container) http.Handler {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", c.UserHandler.Register)
		r.Post("/login", c.UserHandler.Login)
		r.Post("/logout", c.UserHandler.Logout)
		r.Post("/refresh", c.UserHandler.RefreshToken)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)

		r.Get("/profile", c.UserHandler.Profile)
	})

	return r
}
