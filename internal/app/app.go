package app

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"dotsat.work/internal/config"
	"dotsat.work/internal/db"
	"dotsat.work/internal/repository"
	"dotsat.work/internal/service"
)

type App struct {
	Cfg            *config.Config
	DB             *sqlx.DB
	TenantService  *service.TenantService
	UserService    *service.UserService
	ProfileService *service.ProfileService
	AuthService    *service.AuthService
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
		if closeErr := database.Close(); closeErr != nil {
			return nil, fmt.Errorf("failed to run migrations: %w (also failed to close DB: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize repositories
	tenantRepository := repository.NewTenantRepository(database)
	userRepository := repository.NewUserRepository(database)
	profileRepository := repository.NewProfileRepository(database)
	tokenRepository := repository.NewTokenRepository(database)

	// Initialize services
	tenantService := service.NewTenantService(tenantRepository)
	userService := service.NewUserService(userRepository)
	profileService := service.NewProfileService(profileRepository)
	authService := service.NewAuthService(
		userRepository,
		tokenRepository,
		cfg.JWTSecret,
		cfg.IsProduction(),
		cfg.JWTExpiry,
	)

	return &App{
		Cfg:            cfg,
		DB:             database,
		TenantService:  tenantService,
		UserService:    userService,
		ProfileService: profileService,
		AuthService:    authService,
	}, nil
}

func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
