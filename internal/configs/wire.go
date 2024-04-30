package configs

import (
	"github.com/google/wire"
)

var ConfigsSet = wire.NewSet(
	NewConfigurator,
	NewAppConfig,
	NewMongoConfig,
	NewRedisConfig,
	NewJWTConfig,
	NewOTPConfig,
	NewNotifierConfig,
)
