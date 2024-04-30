package handlers

import (
	dtos "github.com/WildEgor/e-shop-auth/internal/dtos/user"
	auth_middleware "github.com/WildEgor/e-shop-auth/internal/middlewares/auth"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
)

type MeHandler struct {
	userRepository *repositories.UserRepository
}

func NewMeHandler(
	userRepository *repositories.UserRepository,
) *MeHandler {
	return &MeHandler{
		userRepository,
	}
}

// TODO: add swagger response
// Get profile godoc
//
//	@Param			Authorization header	string	true	"123"
//	@Tags			User Controller
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} dtos.MeResponseDto
// @Router /api/v1/user/me [get]
func (hch *MeHandler) Handle(ctx fiber.Ctx) error {

	resp := core_dtos.NewResponse(ctx)

	authUser := auth_middleware.ExtractUser(ctx)

	resp.SetStatus(fiber.StatusOK)
	resp.SetData([]*dtos.MeResponseDto{
		{
			ID:     authUser.Id.Hex(),
			Mobile: authUser.Phone,
			Email:  authUser.Email,
		},
	})

	return resp.JSON()
}
