package configs

import (
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/caarlos0/env/v7"
	"log/slog"
)

type AppConfig struct {
	Name    string `env:"APP_NAME" envDefault:"app"`
	Port    string `env:"APP_PORT" envDefault:"8888"`
	Mode    string `env:"APP_MODE,required"`
	RPCPort string `env:"GRPC_PORT" envDefault:"8887"`
	GoEnv   string `env:"GO_ENV" envDefault:"local"`
	Version string `env:"VERSION" envDefault:"local"`
}

func NewAppConfig(c *Configurator) *AppConfig {
	cfg := AppConfig{}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("app config error", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
	}

	return &cfg
}

func (ac AppConfig) IsProduction() bool {
	return ac.Mode != "develop"
}
