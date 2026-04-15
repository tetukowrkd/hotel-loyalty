package app

import (
	"net/http"

	"hotel-loyalty/internal/config"
	"hotel-loyalty/internal/infrastructure/database"
	"hotel-loyalty/internal/infrastructure/logger"
)

type App struct {
	Config *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{Config: cfg}
}

func (a *App) Start() {
	// DB
	db, err := database.NewPostgresDB(a.Config)
	if err != nil {
		logger.ErrorLogger.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	// container (dependency)
	container := NewContainer(db)

	// router
	router := SetupRouter(container)
	logger.InfoLogger.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)
}
