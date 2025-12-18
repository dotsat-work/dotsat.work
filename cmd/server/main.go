package main

import (
	"log/slog"
	"net/http"

	"dotsat.work/internal/app"
	"dotsat.work/internal/config"
	"dotsat.work/internal/routes"
)

func main() {
	cfg := config.Load()

	app, err := app.New(cfg)
	if err != nil {
		slog.Error("failed to initialize app", "error", err)
		panic(err)
	}
	defer func() {
		closeErr := app.Close()
		if closeErr != nil {
			slog.Error("failed to close app", "error", closeErr)
		}
	}()

	handler := routes.SetupRoutes(app)
	slog.Info("Server starting", "port", cfg.Port, "env", cfg.AppEnv, "url", "http://localhost:"+cfg.Port)

	err = http.ListenAndServe(":"+cfg.Port, handler)
	if err != nil {
		slog.Error("failed to start server", "error", err)
		panic(err)
	}
}
