package handlers

import (
	"github.com/WildEgor/e-shop-auth/internal/configs"
	dtos "github.com/WildEgor/e-shop-auth/internal/dtos/auth"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
)

type RefreshHandler struct {
	ur        *repositories.UserRepository
	tr        *repositories.TokensRepository
	jwt       *services.JWTAuthenticator
	jwtConfig *configs.JWTConfig
}

func NewRefreshHandler(
	ur *repositories.UserRepository,
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
	jwtConfig *configs.JWTConfig,
) *RefreshHandler {
	return &RefreshHandler{
		ur:        ur,
		tr:        tr,
		jwt:       jwt,
		jwtConfig: jwtConfig,
	}
}

// TODO: add swagger response
// Refreshes access token using refresh token godoc
//
//	@Param			Authorization header	string	true	"123"
//	@Param			body body	dtos.OTPLoginRequestDto	true	"Body"
//	@Tags			Auth Controller
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} map[string]string
// @Router /api/v1/auth/refresh-token [post]
func (h *RefreshHandler) Handle(ctx fiber.Ctx) error {
	resp := core_dtos.NewResponse(ctx)

	rt := ctx.Cookies("refresh_token")
	if rt == "" {
		resp.SetStatus(fiber.StatusBadRequest)
		return resp.JSON()
	}

	token, err := h.jwt.ParseToken(rt)
	if err != nil {
		resp.SetStatus(fiber.StatusBadRequest)
		return resp.JSON()
	}

	userId, err := h.tr.GetRT(token.TokenUuid)
	if err != nil {
		resp.SetStatus(fiber.StatusBadRequest)
		return resp.JSON()
	}

	user, err := h.ur.FindById(userId)
	if err != nil {
		resp.SetStatus(fiber.StatusBadRequest)
		return resp.JSON()
	}

	nat, atErr := h.jwt.GenerateToken(user.Id.Hex(), h.jwtConfig.ATDuration)
	if atErr != nil {
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	sErr := h.tr.SetAT(nat)
	if sErr != nil {
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	resp.SetStatus(fiber.StatusOK)
	resp.SetData([]*dtos.TokenPairResponseDto{{
		UserID:       user.Id.Hex(),
		AccessToken:  nat.Token,
		RefreshToken: rt,
	}})

	return resp.JSON()
}
