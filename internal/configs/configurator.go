package configs

import (
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/joho/godotenv"
	"log/slog"
)

type Configurator struct{}

func NewConfigurator() *Configurator {
	c := &Configurator{}
	c.Load()
	return c
}

func (c *Configurator) Load() {
	err := godotenv.Load(".env", ".env.local")
	if err != nil {
		slog.Error("error loading envs file", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
	}
}
