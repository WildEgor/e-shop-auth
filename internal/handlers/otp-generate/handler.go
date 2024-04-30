package handlers

import (
	dtos "github.com/WildEgor/e-shop-auth/internal/dtos/auth"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	"github.com/WildEgor/e-shop-auth/internal/validators"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
)

type OTPGenHandler struct {
	ur  *repositories.UserRepository
	otp *services.NotificationService
}

func NewOTPGenHandler(
	ur *repositories.UserRepository,
	otp *services.NotificationService,
) *OTPGenHandler {
	return &OTPGenHandler{
		ur,
		otp,
	}
}

// Generate OTP code godoc
//
//	@Param			body body	dtos.OTPGenerateRequestDto	true	"Body"
//	@Tags			OTP Controller
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} dtos.OTPGenerateResponseDto
// @Router /api/v1/auth/otp/gen [post]
func (h *OTPGenHandler) Handle(ctx fiber.Ctx) error {
	dto := &dtos.OTPGenerateRequestDto{}
	resp := core_dtos.NewResponse(ctx)

	if err := validators.ParseAndValidate(ctx, dto); err != nil {
		return err
	}

	us, err := h.ur.FindByPhone(dto.Phone)
	if err != nil {
		resp.SetStatus(fiber.StatusUnauthorized)
		return resp.JSON()
	}

	if err := us.IsOTPResendAvailable(); err != nil {
		resp.SetStatus(fiber.StatusTooManyRequests)
		return resp.JSON()
	}

	code, err := h.otp.GenerateAndSendOTPSms(us.Phone)
	if err != nil {
		resp.SetStatus(fiber.StatusTooManyRequests)
		return resp.JSON()
	}

	us.UpdateOTP(us.Phone, code)

	if err := h.ur.UpdateOTP(us.Id, &us.OTP); err != nil {
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	resp.SetStatus(fiber.StatusOK)

	resp.SetData([]*dtos.OTPGenerateResponseDto{{
		Identity: dto.Phone,
		Code:     us.OTP.Code, // for debug
	}})

	return resp.JSON()
}
