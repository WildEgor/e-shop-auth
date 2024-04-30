package configs

import (
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/caarlos0/env/v7"
	"log/slog"
)

type NotifierConfig struct {
	DSN      string `env:"AMQP_DSN,required"`
	Exchange string `env:"NOTIFIER_EXCHANGE,required"`
}

func NewNotifierConfig(c *Configurator) *NotifierConfig {
	cfg := NotifierConfig{}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("notifier config error", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
	}

	return &cfg
}
