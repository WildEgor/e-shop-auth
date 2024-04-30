package auth_middleware

import (
	"github.com/WildEgor/e-shop-auth/internal/models"
	"github.com/gofiber/fiber/v3"
	"time"
)

func ExtractUser(ctx fiber.Ctx) *models.UsersModel {
	data := ctx.Locals(LocalsUserKey)

	user, ok := data.(*models.UsersModel)
	if !ok {
		// TODO: ???
	}

	return user
}

func ExtractRefreshTokenFromCookies(ctx fiber.Ctx) string {
	rt := ctx.Cookies("refresh_token")
	return rt
}

func ExtractAssessTokenFromCookies(ctx fiber.Ctx) string {
	rt := ctx.Cookies(AccessTokenUUIDKey)
	return rt
}

func ResetCookies(c fiber.Ctx) {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: expired,
	})
	c.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Expires: expired,
	})
	c.Cookie(&fiber.Cookie{
		Name:    "logged_in",
		Value:   "",
		Expires: expired,
	})
}
