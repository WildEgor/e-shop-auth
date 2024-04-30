package handlers

import (
	"github.com/WildEgor/e-shop-auth/internal/configs"
	dtos "github.com/WildEgor/e-shop-auth/internal/dtos/auth"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	"github.com/WildEgor/e-shop-auth/internal/validators"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
)

type OTPLoginHandler struct {
	ur        *repositories.UserRepository
	tr        *repositories.TokensRepository
	jwt       *services.JWTAuthenticator
	jwtConfig *configs.JWTConfig
}

func NewOTPLoginHandler(
	ur *repositories.UserRepository,
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
	jwtConfig *configs.JWTConfig,
) *OTPLoginHandler {
	return &OTPLoginHandler{
		ur,
		tr,
		jwt,
		jwtConfig,
	}
}

// TODO: add swagger response
// Login via OTP code godoc
//
//	@Param			body body	dtos.OTPLoginRequestDto	true	"Body"
//	@Tags			OTP Controller
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} map[string]string
// @Router /api/v1/auth/otp/login [post]
func (h *OTPLoginHandler) Handle(ctx fiber.Ctx) error {
	dto := &dtos.OTPLoginRequestDto{}
	if err := validators.ParseAndValidate(ctx, dto); err != nil {
		return err
	}

	resp := core_dtos.NewResponse(ctx)

	us, err := h.ur.FindByPhone(dto.Phone)
	if err != nil {
		resp.SetStatus(fiber.StatusUnauthorized)
		return resp.JSON()
	}

	if err := us.VerifyOTP(us.Phone, dto.Code); err != nil {
		resp.SetStatus(fiber.StatusUnauthorized)
		return resp.JSON()
	}

	us.ClearOTP()

	if err := h.ur.UpdateOTP(us.Id, &us.OTP); err != nil {
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	at, atErr := h.jwt.GenerateToken(us.Id.Hex(), h.jwtConfig.ATDuration)
	rt, rtErr := h.jwt.GenerateToken(us.Id.Hex(), h.jwtConfig.RTDuration)
	if atErr != nil || rtErr != nil {
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	errAT := h.tr.SetAT(at)
	errRT := h.tr.SetRT(rt)
	if errAT != nil || errRT != nil {
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	resp.SetStatus(fiber.StatusOK)

	resp.SetData([]*dtos.TokenPairResponseDto{{
		UserID:       us.Id.Hex(),
		AccessToken:  at.Token,
		RefreshToken: rt.Token,
	}})

	return resp.JSON()
}
