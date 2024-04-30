package auth_middleware

import (
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"log/slog"
	"strings"
)

type AuthMiddlewareConfig struct {
	Filter       func(ctx fiber.Ctx) bool
	UR           *repositories.UserRepository
	JWT          *services.JWTAuthenticator
	Unauthorized fiber.Handler
	Decode       func(ctx fiber.Ctx) (*jwt.MapClaims, error)
}

var AuthMiddlewareConfigDefault = AuthMiddlewareConfig{
	Filter:       nil,
	Decode:       nil,
	Unauthorized: nil,
}

var (
	LocalsUserKey      = "__user__"
	AccessTokenUUIDKey = "__access_token_uuid__"
	JWTClaimsKey       = "__jwtClaims__"
)

func configAuthDefault(config ...AuthMiddlewareConfig) AuthMiddlewareConfig {
	if len(config) < 1 {
		return AuthMiddlewareConfigDefault
	}

	cfg := config[0]

	if cfg.Filter == nil {
		cfg.Filter = AuthMiddlewareConfigDefault.Filter
	}

	if cfg.Decode == nil {
		// Set default Decode function if not passed
		cfg.Decode = func(c fiber.Ctx) (*jwt.MapClaims, error) {
			var token string
			authHeader := c.Get("Authorization")

			if len(authHeader) > 0 {
				components := strings.SplitN(authHeader, " ", 2)
				if len(components) == 2 && components[0] == "Bearer" {
					token = components[1]
				}
			}

			if len(token) == 0 {
				return nil, errors.New("empty token")
			}

			// TODO: need check token in Redis too
			claims, err := cfg.JWT.ParseToken(token)
			jwtPayload := jwt.MapClaims{}

			if err == nil && claims != nil && claims.IsValid {
				jwtPayload = jwt.MapClaims{
					"sub":           claims.UserID,
					"typ":           "Bearer",
					"exp":           claims.ExpiresIn,
					"access_token":  authHeader,
					"refresh_token": "",
				}

				c.Locals(AccessTokenUUIDKey, claims.TokenUuid)

			} else {
				return nil, errors.Wrap(err, "token validation")
			}

			return &jwtPayload, nil
		}
	}

	// Set default Unauthorized if not passed
	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c fiber.Ctx) error {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}

	return cfg
}

// NewAuthMiddleware validate accessToken in Keycloak and parse it, extract user from DB
func NewAuthMiddleware(config AuthMiddlewareConfig) fiber.Handler {
	// For setting default config
	cfg := configAuthDefault(config)

	return func(c fiber.Ctx) error {
		resp := core_dtos.NewResponse(c)

		// Don't execute middleware if Filter returns true
		if cfg.Filter != nil && cfg.Filter(c) {
			slog.Debug("auth middleware was skipped")
			return c.Next()
		}
		slog.Debug("auth middleware was run")

		claims, err := cfg.Decode(c)
		if err == nil {
			c.Locals(JWTClaimsKey, *claims)

			id := (*claims)["sub"].(string)

			user, err := cfg.UR.FindById(id)
			if err == nil {
				c.Locals(LocalsUserKey, user)
				return c.Next()
			}
		}

		resp.SetStatus(fiber.StatusUnauthorized)

		return resp.JSON()
	}
}
