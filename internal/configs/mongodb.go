package configs

import (
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/caarlos0/env/v7"
	"log/slog"
)

type MongoConfig struct {
	DBName string `env:"MONGODB_NAME,required"`
	URI    string `env:"MONGODB_URI,required"`
}

func NewMongoConfig(c *Configurator) *MongoConfig {
	cfg := MongoConfig{}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("mongo config error", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
	}

	return &cfg
}
