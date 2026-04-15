package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRouter(c *Container) http.Handler {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", c.UserHandler.Register)
		r.Post("/login", c.UserHandler.Login)
	})

	return r
}
