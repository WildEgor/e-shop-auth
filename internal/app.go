package pkg

import (
	"context"
	"fmt"
	"github.com/WildEgor/e-shop-auth/internal/configs"
	mongo2 "github.com/WildEgor/e-shop-auth/internal/db/mongodb"
	"github.com/WildEgor/e-shop-auth/internal/db/redis"
	eh "github.com/WildEgor/e-shop-auth/internal/handlers/errors"
	nfm "github.com/WildEgor/e-shop-auth/internal/middlewares/not_found"
	"github.com/WildEgor/e-shop-auth/internal/router"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/google/wire"
	"log/slog"
	"os"
)

var AppSet = wire.NewSet(
	NewApp,
	configs.ConfigsSet,
	router.RouterSet,
)

// Server represents the main server configuration.
type Server struct {
	App       *fiber.App
	AppConfig *configs.AppConfig

	Mongo       *mongo2.MongoConnection
	MongoConfig *configs.MongoConfig

	Redis       *redis.RedisConnection
	RedisConfig *configs.RedisConfig
}

func (srv *Server) Run(ctx context.Context) {
	slog.Info("server is listening")

	srv.Mongo.Connect(ctx)
	srv.Redis.Connect(ctx)

	if err := srv.App.Listen(fmt.Sprintf(":%s", srv.AppConfig.Port)); err != nil {
		slog.Error("unable to start server")
	}
}

func (srv *Server) Shutdown(ctx context.Context) {
	slog.Info("shutdown service")

	srv.Mongo.Disconnect(ctx)
	srv.Redis.Disconnect(ctx)

	if err := srv.App.Shutdown(); err != nil {
		slog.Error("unable to shutdown server")
	}
}

func NewApp(
	ac *configs.AppConfig,
	eh *eh.ErrorsHandler,
	prr *router.PrivateRouter,
	pbr *router.PublicRouter,
	sr *router.SwaggerRouter,

	mc *configs.MongoConfig,
	mongo *mongo2.MongoConnection,

	rc *configs.RedisConfig,
	redis *redis.RedisConnection,
) *Server {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	if ac.IsProduction() {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	slog.SetDefault(logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: eh.Handle,
		Views:        html.New("./views", ".html"),
	})

	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	app.Use(recover.New())

	prr.Setup(app)
	pbr.Setup(app)
	sr.Setup(app)

	// 404 handler
	app.Use(nfm.NewNotFound())

	return &Server{
		App:         app,
		AppConfig:   ac,
		Mongo:       mongo,
		MongoConfig: mc,
		Redis:       redis,
		RedisConfig: rc,
	}
}
