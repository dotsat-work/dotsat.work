package routes

import (
	"net/http"

	"dotsat.work/internal/app"
	"dotsat.work/internal/handler"
	"dotsat.work/internal/middleware"
)

func SetupRoutes(a *app.App) http.Handler {
	// Handlers
	home := handler.NewHomeHandler()

	mux := http.NewServeMux()

	// ============================================================================
	// PUBLIC ROUTES
	// ============================================================================

	// Home
	mux.Handle("GET /{$}", home)

	// ============================================================================
	// PROTECTED ROUTES (/app/*)
	// ============================================================================
	// TODO: Add protected routes here

	// ============================================================================
	// FALLBACK
	// ============================================================================

	// 404 - for now, just use home handler
	mux.Handle("/{path...}", home)

	// Global middleware - executed in order (top to bottom)
	handler := middleware.Chain(
		mux,
		middleware.AuthMiddleware(a.AuthService, a.UserService, a.ProfileService, a.TenantService),
	)

	return handler
}
