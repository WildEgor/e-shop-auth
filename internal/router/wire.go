package router

import (
	"github.com/WildEgor/e-shop-auth/internal/handlers"
	"github.com/google/wire"
)

var RouterSet = wire.NewSet(
	handlers.HandlersSet,
	NewPublicRouter,
	NewPrivateRouter,
	NewSwaggerRouter,
)
