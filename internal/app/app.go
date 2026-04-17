package app

import (
	"net/http"

	"hotel-loyalty/internal/infrastructure/logger"
)

type App struct {
	Container *Container
}

func NewApp(container *Container) *App {
	return &App{Container: container}
}

func (a *App) Start() {
	// router
	router := SetupRouter(a.Container)
	logger.InfoLogger.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)
}
