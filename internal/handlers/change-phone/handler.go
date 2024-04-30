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

type ChangePhoneHandler struct {
	ur *repositories.UserRepository
	ns *services.NotificationService
}

func NewChangePhoneHandler(
	ur *repositories.UserRepository,
	ns *services.NotificationService,
) *ChangePhoneHandler {
	return &ChangePhoneHandler{
		ur,
		ns,
	}
}

// Change phone godoc
//
//	@Param			Authorization header	string	true	"123"
//	@Param			body body	dtos.ChangePhoneRequestDto	true	"Body"
//	@Summary		Change authenticated user phone
//	@Description	Change authenticated user phone
//	@Tags			Auth Controller
//	@Accept			json
//	@Produce		json
//	@Router			/api/v1/auth/change-phone [post]
func (h *ChangePhoneHandler) Handle(ctx fiber.Ctx) error {
	dto := &dtos.ChangePhoneRequestDto{}
	resp := core_dtos.NewResponse(ctx)

	if err := validators.ParseAndValidate(ctx, dto); err != nil {
		return err
	}

	us := auth_middleware.ExtractUser(ctx)

	if h.ur.CheckExistsPhone(dto.Phone) {
		domains.SetPhoneAlreadyExistError(resp)
		return resp.JSON()
	}

	if us.IsPhoneEqual(dto.Phone) {
		domains.SetPhoneEqualityError(resp)
		return resp.JSON()
	}

	if err := us.IsPhoneConfirmResendAvailable(); err != nil {
		domains.SetSendCodeError(resp)
		return resp.JSON()
	}

	code, err := h.ns.GenerateAndSendPhoneConfirm(dto.Phone)
	if err != nil {
		domains.SetSendCodeTimeoutError(resp)
		return resp.JSON()
	}

	us.UpdatePhoneVerification(dto.Phone, code)

	if err := h.ur.UpdateVerification(us.Id, &us.Verification); err != nil {
		domains.SetSendCodeError(resp)
		return resp.JSON()
	}

	resp.SetStatus(fiber.StatusOK)

	return resp.JSON()
}
