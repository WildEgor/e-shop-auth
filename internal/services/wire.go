package services

import (
	"github.com/google/wire"
)

var ServicesSet = wire.NewSet(
	NewJWTAuthenticator,
	NewNotificationService,
)
