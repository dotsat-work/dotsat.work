package config

import (
	"log/slog"
	"os"
	"time"

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

	// Authentication
	JWTSecret string
	JWTExpiry time.Duration
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

		// Authentication
		JWTSecret: envRequired("JWT_SECRET"),
		JWTExpiry: envDuration("JWT_EXPIRY", 168*time.Hour), // 7-day default
	}

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

func envDuration(key string, def time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		slog.Warn("config invalid duration, using default", "key", key, "value", v, "default", def)
		return def
	}
	return d
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}
