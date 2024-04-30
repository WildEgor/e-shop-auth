package db

import (
	mongo "github.com/WildEgor/e-shop-auth/internal/db/mongodb"
	"github.com/WildEgor/e-shop-auth/internal/db/redis"
	"github.com/google/wire"
)

var DbSet = wire.NewSet(
	mongo.NewMongoConnection,
	redis.NewRedisDBConnection,
)
