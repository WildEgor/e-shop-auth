package handlers

import (
	domains "github.com/WildEgor/e-shop-auth/internal/domain"
	dtos "github.com/WildEgor/e-shop-auth/internal/dtos/user"
	auth_middleware "github.com/WildEgor/e-shop-auth/internal/middlewares/auth"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	"github.com/WildEgor/e-shop-auth/internal/validators"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
)

type ChangeEmailHandler struct {
	ur *repositories.UserRepository
	ns *services.NotificationService
}

func NewChangeEmailHandler(
	ur *repositories.UserRepository,
	ns *services.NotificationService,
) *ChangeEmailHandler {
	return &ChangeEmailHandler{
		ur,
		ns,
	}
}

// Change e-mail godoc
//
//	@Param			Authorization header	string	true	"123"
//	@Param			body body	dtos.ChangeEmailRequestDto	true	"Body"
//	@Summary		Change authenticated user e-mail
//	@Description	Change authenticated user e-mail
//	@Tags			Auth Controller
//	@Accept			json
//	@Produce		json
//	@Router			/api/v1/auth/change-email [post]
func (h *ChangeEmailHandler) Handle(ctx fiber.Ctx) error {
	dto := &dtos.ChangeEmailRequestDto{}
	resp := core_dtos.NewResponse(ctx)

	if err := validators.ParseAndValidate(ctx, dto); err != nil {
		return err
	}

	us := auth_middleware.ExtractUser(ctx)

	if h.ur.CheckExistsEmail(dto.Email) {
		domains.SetEmailAlreadyExistError(resp)
		return resp.JSON()
	}

	if us.IsEmailEqual(dto.Email) {
		domains.SetEmailEqualityError(resp)
		return resp.JSON()
	}

	if err := us.IsEmailConfirmResendAvailable(); err != nil {
		domains.SetSendCodeTimeoutError(resp)
		return resp.JSON()
	}

	// TODO: also generate unique token for link confirm (test)

	code, err := h.ns.GenerateAndSendEmailConfirm(dto.Email)
	if err != nil {
		domains.SetSendCodeError(resp)
		return resp.JSON()
	}

	us.UpdateEmailVerification(dto.Email, code)

	if err := h.ur.UpdateVerification(us.Id, &us.Verification); err != nil {
		resp.SetStatus(fiber.StatusBadRequest)
		return resp.JSON()
	}

	resp.SetStatus(fiber.StatusOK)

	return resp.JSON()
}
