package handlers

import (
	"github.com/WildEgor/e-shop-auth/internal/configs"
	domains "github.com/WildEgor/e-shop-auth/internal/domain"
	dtos "github.com/WildEgor/e-shop-auth/internal/dtos/auth"
	auth_middleware "github.com/WildEgor/e-shop-auth/internal/middlewares/auth"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	"github.com/WildEgor/e-shop-auth/internal/validators"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
)

type ConfirmPhoneHandler struct {
	ur        *repositories.UserRepository
	tr        *repositories.TokensRepository
	jwt       *services.JWTAuthenticator
	jwtConfig *configs.JWTConfig
}

func NewConfirmPhoneHandler(
	ur *repositories.UserRepository,
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
	jwtConfig *configs.JWTConfig,
) *ConfirmPhoneHandler {
	return &ConfirmPhoneHandler{
		ur,
		tr,
		jwt,
		jwtConfig,
	}
}

// Confirm phone godoc
//
//	@Param			body body	dtos.ConfirmPhoneRequestDto	true	"Body"
//	@Tags			Auth Controller
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} dtos.TokenPairResponseDto
// @Router /api/v1/auth/confirm-phone [post]
func (h *ConfirmPhoneHandler) Handle(ctx fiber.Ctx) error {
	dto := &dtos.ConfirmPhoneRequestDto{}
	resp := core_dtos.NewResponse(ctx)

	if err := validators.ParseAndValidate(ctx, dto); err != nil {
		return err
	}

	us := auth_middleware.ExtractUser(ctx)

	if err := us.VerifyIdentity(dto.Phone, dto.Code); err != nil {
		domains.SetMalformedCodeError(resp)
		return resp.JSON()
	}

	us.UpdatePhone(dto.Phone)
	us.ClearPhoneVerification()

	if err := h.ur.UpdateContacts(us); err != nil {
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	if err := h.ur.UpdateVerification(us.Id, &us.Verification); err != nil {
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	at, atErr := h.jwt.GenerateToken(us.Id.Hex(), h.jwtConfig.ATDuration)
	rt, rtErr := h.jwt.GenerateToken(us.Id.Hex(), h.jwtConfig.RTDuration)
	if atErr != nil || rtErr != nil {
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	if h.tr.SetAT(at) != nil || h.tr.SetRT(rt) != nil {
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	resp.SetStatus(fiber.StatusOK)

	resp.SetData([]*dtos.TokenPairResponseDto{
		{
			UserID:       us.Id.Hex(),
			AccessToken:  at.Token,
			RefreshToken: rt.Token,
		},
	})

	return resp.JSON()
}
