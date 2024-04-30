package repositories

import (
	"github.com/WildEgor/e-shop-auth/internal/db"
	"github.com/google/wire"
)

var RepositoriesSet = wire.NewSet(
	db.DbSet,
	NewUserRepository,
	NewTokensRepository,
)
