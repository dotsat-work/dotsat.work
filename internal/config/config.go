package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Application
	AppName string
	AppEnv  string
	AppURL  string
	Port    string

	// Database
	DBDriver     string
	DBConnection string
}

func Load() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		slog.Info("no .env file found, using environment variables")
	}

	cfg := &Config{
		// Application
		AppName: envString("APP_NAME", "dotsat.work.test"),
		AppEnv:  envRequired("APP_ENV"),
		AppURL:  envRequired("APP_URL"),
		Port:    envString("PORT", "8090"),

		// Database
		DBDriver:     envString("DB_DRIVER", "postgres"),
		DBConnection: envRequired("DB_CONNECTION"),
	}

	// TODO: Validate if prod or dev
	return cfg
}

func envString(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		value = def
	}
	return value
}

func envRequired(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	slog.Error("config required env var missing", "key", key)
	os.Exit(1)
	return ""
}
