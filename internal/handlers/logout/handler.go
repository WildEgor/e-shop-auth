package handlers

import (
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
)

type LogoutHandler struct {
	tr  *repositories.TokensRepository
	jwt *services.JWTAuthenticator
}

func NewLogoutHandler(
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
) *LogoutHandler {
	return &LogoutHandler{
		tr:  tr,
		jwt: jwt,
	}
}

// Logout godoc
//
//	@Param			Authorization header	string	true	"123"
//	@Tags			Auth Controller
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} map[string]string
// @Router /api/v1/auth/logout [post]
func (h *LogoutHandler) Handle(ctx fiber.Ctx) error {
	resp := core_dtos.NewResponse(ctx)

	// TODO: impl

	resp.SetStatus(fiber.StatusOK)
	return resp.JSON()
}
