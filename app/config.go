package app

import (
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/utils"
)

func (a *App) Config() *model.Config {
	if cfg := a.config.Load(); cfg != nil {
		return cfg.(*model.Config)
	}
	return &model.Config{}
}

func (a *App) LoadConfig(configFile string) *model.AppError {
	cfg, configPath, _, err := utils.LoadConfig(configFile)
	if err != nil {
		return err
	}

	if err = cfg.IsValid(); err != nil {
		return err
	}

	a.configFile = configPath

	a.config.Store(cfg)
	return nil
}
