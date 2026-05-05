package app

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/cors"

	"hotel-loyalty/internal/infrastructure/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRouter(c *Container) http.Handler {
	r := chi.NewRouter()

	// 🔥 ambil dari env + normalize
	raw := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")

	var origins []string
	for _, o := range raw {
		origins = append(origins, strings.TrimSpace(o))
	}

	// 🔥 pasang CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// 🔥 static files
	fs := http.FileServer(http.Dir("./storage"))
	r.Handle("/storage/*", http.StripPrefix("/storage/", fs))

	// 🔥 AUTH
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", c.UserHandler.Register)
		r.Post("/login", c.UserHandler.Login)
		r.Post("/refresh", c.TokenHandler.RefreshToken)

		r.With(middleware.JWTMiddleware).Post("/logout", c.UserHandler.Logout)
	})

	// 🔥 MEMBER
	r.Route("/member", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)

		r.Get("/profile", c.UserHandler.Profile)
		r.Get("/qr", c.UserHandler.GetMemberQR)
	})

	// 🔥 STAFF
	r.Route("/staff", func(r chi.Router) {
		r.Use(
			middleware.JWTMiddleware,
			middleware.RequireRole("staff", "admin"),
		)

		r.Route("/hotels", func(r chi.Router) {

			// create hotel
			r.Post("/", c.HotelHandler.Create)
			// delete hotel
			r.Delete("/{hotel_id}", c.HotelHandler.Delete)

			// images
			r.Post("/{hotel_id}/images", c.HotelImageHandler.Create)
			r.Patch("/{hotel_id}/images/{image_id}/primary", c.HotelImageHandler.SetPrimary)
			// delete image
			r.Delete("/{hotel_id}/images/{image_id}", c.HotelImageHandler.Delete)

			// facilities
			r.Post("/{hotel_id}/facilities", c.HotelFacilityHandler.Assign)

			// create rooms
			r.Post("/{hotel_id}/rooms", c.RoomHandler.Create)
			r.Delete("/rooms/{room_id}", c.RoomHandler.Delete)

		})
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware, middleware.RequireRole("admin"))

		r.Post("/facilities", c.HotelFacilityHandler.Create)
		r.Delete("/facilities/{facility_id}", c.HotelFacilityHandler.Delete)

		r.Patch("/hotels/{hotel_id}/publish", c.HotelHandler.Publish)
	})

	// 🔥 PUBLIC (no auth)
	r.Route("/hotels", func(r chi.Router) {
		r.Get("/", c.HotelHandler.List)
		r.Get("/{hotel_id}", c.HotelHandler.GetByID)

		r.Get("/{hotel_id}/facilities", c.HotelFacilityHandler.GetByHotel)

		r.Get("/{hotel_id}/rooms", c.RoomHandler.List)
	})

	// 🔥 MASTER DATA
	r.Get("/facilities", c.HotelFacilityHandler.GetAll)

	return r
}
