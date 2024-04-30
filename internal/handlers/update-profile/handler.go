package handlers

import (
	dtos "github.com/WildEgor/e-shop-auth/internal/dtos/user"
	auth_middleware "github.com/WildEgor/e-shop-auth/internal/middlewares/auth"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/validators"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
)

type UpdateProfileHandler struct {
	ur *repositories.UserRepository
}

func NewUpdateProfileHandler(
	ur *repositories.UserRepository,
) *UpdateProfileHandler {
	return &UpdateProfileHandler{
		ur,
	}
}

// TODO: add swagger response
// Update profile godoc
//
//	@Param			Authorization header	string	true	"123"
//	@Param			body body	dtos.UpdateProfileRequestDto	true	"Body"
//	@Tags			User Controller
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} dtos.MeResponseDto
// @Router /api/v1/user/profile [put]
func (h *UpdateProfileHandler) Handle(ctx fiber.Ctx) error {
	dto := &dtos.UpdateProfileRequestDto{}
	resp := core_dtos.NewResponse(ctx)

	if err := validators.ParseAndValidate(ctx, dto); err != nil {
		return err
	}

	us := auth_middleware.ExtractUser(ctx)
	us.SetInfo(dto.FirstName, dto.LastName)

	if err := h.ur.UpdateInfo(us); err != nil {
		resp.SetStatus(fiber.StatusBadRequest)
		return resp.JSON()
	}

	resp.SetStatus(fiber.StatusOK)
	resp.SetData([]*dtos.MeResponseDto{
		{
			ID:     us.Id.Hex(),
			Mobile: us.Phone,
			Email:  us.Email,
		},
	})

	return resp.JSON()
}
