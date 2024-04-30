package router

import (
	"github.com/WildEgor/e-shop-auth/internal/configs"
	hch "github.com/WildEgor/e-shop-auth/internal/handlers/health_check"
	login_handler "github.com/WildEgor/e-shop-auth/internal/handlers/login"
	logout_handler "github.com/WildEgor/e-shop-auth/internal/handlers/logout"
	otp_generate_handler "github.com/WildEgor/e-shop-auth/internal/handlers/otp-generate"
	otp_login_handler "github.com/WildEgor/e-shop-auth/internal/handlers/otp-login"
	rch "github.com/WildEgor/e-shop-auth/internal/handlers/ready_check"
	reg_handler "github.com/WildEgor/e-shop-auth/internal/handlers/reg"
	auth_middleware "github.com/WildEgor/e-shop-auth/internal/middlewares/auth"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"log/slog"
)

type PublicRouter struct {
	rc  *reg_handler.RegHandler
	lo  *login_handler.LoginHandler
	lt  *logout_handler.LogoutHandler
	og  *otp_generate_handler.OTPGenHandler
	lot *otp_login_handler.OTPLoginHandler
	hch *hch.HealthCheckHandler
	rch *rch.ReadyCheckHandler

	ur *repositories.UserRepository
	tr *repositories.TokensRepository

	jwt *services.JWTAuthenticator

	jwtc *configs.JWTConfig
}

func NewPublicRouter(
	rc *reg_handler.RegHandler,
	lo *login_handler.LoginHandler,
	lt *logout_handler.LogoutHandler,
	og *otp_generate_handler.OTPGenHandler,
	lot *otp_login_handler.OTPLoginHandler,
	hch *hch.HealthCheckHandler,
	rch *rch.ReadyCheckHandler,
	ur *repositories.UserRepository,
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
	jwtc *configs.JWTConfig,
) *PublicRouter {
	return &PublicRouter{
		rc,
		lo,
		lt,
		og,
		lot,
		hch,
		rch,
		ur,
		tr,
		jwt,
		jwtc,
	}
}

func (r *PublicRouter) Setup(app *fiber.App) {
	api := app.Group("/api", limiter.New(limiter.Config{
		Max:                    10,
		SkipSuccessfulRequests: true,
	}))
	v1 := api.Group("/v1")

	am := auth_middleware.NewAuthMiddleware(auth_middleware.AuthMiddlewareConfig{
		UR:  r.ur,
		JWT: r.jwt,
	})

	v1.Post("/auth/reg", r.rc.Handle)
	v1.Post("/auth/login", r.lo.Handle)
	v1.Post("/auth/logout", am, r.lt.Handle)
	v1.Post("/auth/otp/gen", r.og.Handle)
	v1.Post("/auth/otp/login", r.lot.Handle)

	v1.Get("/livez", healthcheck.NewHealthChecker(healthcheck.Config{
		Probe: func(ctx fiber.Ctx) bool {
			if err := r.hch.Handle(ctx); err != nil {
				slog.Error("error not healthy")
				return false
			}

			slog.Debug("is healthy")

			return true
		},
	}))
	v1.Get("/readyz", healthcheck.NewHealthChecker(healthcheck.Config{
		Probe: func(ctx fiber.Ctx) bool {
			if err := r.rch.Handle(ctx); err != nil {
				slog.Error("error not ready")
				return false
			}

			slog.Debug("is ready")

			return true
		},
	}))
}
