package configs

import (
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/caarlos0/env/v7"
	"log/slog"
	"time"
)

type JWTConfig struct {
	Secret string `env:"JWT_SECRET,required"`
	at     string `env:"JWT_AT_TTL" envDefault:"10m"`
	rt     string `env:"JWT_RT_TTL" envDefault:"24h"`

	ATDuration time.Duration `env:"-"`
	RTDuration time.Duration `env:"-"`
}

func NewJWTConfig(c *Configurator) *JWTConfig {
	cfg := JWTConfig{}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("jwt config error", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
	}

	// FIXME: cannot parse
	//cfg.ATDuration, _ = time.ParseDuration(cfg.at)
	//cfg.RTDuration, _ = time.ParseDuration(cfg.rt)

	cfg.ATDuration, _ = time.ParseDuration("24h")
	cfg.RTDuration, _ = time.ParseDuration("10d")

	return &cfg
}
