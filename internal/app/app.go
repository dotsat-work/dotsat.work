package app

import (
	"dotsat.work/internal/config"
)

type App struct {
	Cfg *config.Config
}

func New(cfg *config.Config) (*App, error) {

	return &App{
		Cfg: cfg,
	}, nil
}

func (a *App) Close() error {
	return nil
}
