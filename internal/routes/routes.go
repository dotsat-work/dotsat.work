package routes

import (
	"net/http"

	"dotsat.work/internal/app"
	"dotsat.work/internal/handler"
)

func SetupRoutes(app *app.App) http.Handler {
	// Handlers
	home := handler.NewHomeHandler()

	return home
}
