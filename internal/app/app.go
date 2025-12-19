package app

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"dotsat.work/internal/config"
	"dotsat.work/internal/db"
)

type App struct {
	Cfg *config.Config
	DB  *sqlx.DB
	// TODO: Add repositories and services as you build them
	// TenantRepository *repository.TenantRepository
	// UserRepository   *repository.UserRepository
	// TenantService    *service.TenantService
	// UserService      *service.UserService
}

func New(cfg *config.Config) (*App, error) {
	// Initialize database
	database, err := db.Init(cfg.DBDriver, cfg.DBConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Run database migrations
	err = db.RunMigrations(database.DB)
	if err != nil {
		database.Close() // Close DB on migration failure
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// TODO: Initialize repositories
	// tenantRepository := repository.NewTenantRepository(database)
	// userRepository := repository.NewUserRepository(database)

	// TODO: Initialize services
	// tenantService := service.NewTenantService(tenantRepository)
	// userService := service.NewUserService(userRepository)

	return &App{
		Cfg: cfg,
		DB:  database,
		// TenantRepository: tenantRepository,
		// UserRepository:   userRepository,
		// TenantService:    tenantService,
		// UserService:      userService,
	}, nil
}

func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
