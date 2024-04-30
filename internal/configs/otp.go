package configs

import (
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/caarlos0/env/v7"
	"log/slog"
)

type OTPConfig struct {
	Issuer string `env:"OTP_ISSUER" envDefault:"gAuth"`
	Length uint8  `env:"OTP_LENGTH" envDefault:"6"`
}

func NewOTPConfig(c *Configurator) *OTPConfig {
	cfg := OTPConfig{}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("otp config error", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
	}

	return &cfg
}
