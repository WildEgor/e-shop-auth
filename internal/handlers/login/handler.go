package handlers

import (
	"github.com/WildEgor/e-shop-auth/internal/configs"
	dtos "github.com/WildEgor/e-shop-auth/internal/dtos/auth"
	"github.com/WildEgor/e-shop-auth/internal/models"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	"github.com/WildEgor/e-shop-auth/internal/validators"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
)

type LoginHandler struct {
	ur        *repositories.UserRepository
	tr        *repositories.TokensRepository
	jwt       *services.JWTAuthenticator
	jwtConfig *configs.JWTConfig
}

func NewLoginHandler(
	ur *repositories.UserRepository,
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
	jwtConfig *configs.JWTConfig,
) *LoginHandler {
	return &LoginHandler{
		ur:        ur,
		tr:        tr,
		jwt:       jwt,
		jwtConfig: jwtConfig,
	}
}

// Login via email/password or phone/password godoc
//
//	@Param			body body	dtos.LoginRequestDto	true	"Body"
//	@Tags			Auth Controller
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} dtos.LoginResponseDto
// @Router /api/v1/auth/login [post]
func (h *LoginHandler) Handle(ctx fiber.Ctx) error {
	dto := &dtos.LoginRequestDto{}
	resp := core_dtos.NewResponse(ctx)

	if err := validators.ParseAndValidate(ctx, dto); err != nil {
		return err
	}

	user, err := h.ur.FindByLogin(dto.Login)
	if err != nil {
		resp.SetStatus(fiber.StatusUnauthorized)
		return resp.JSON()
	}

	isEqual, err := user.ComparePassword(dto.Password)
	if err != nil {
		resp.SetStatus(fiber.StatusUnauthorized)
		return resp.JSON()
	}

	if !isEqual {
		resp.SetStatus(fiber.StatusUnauthorized)
		return resp.JSON()
	}

	at, atErr := h.jwt.GenerateToken(user.Id.Hex(), h.jwtConfig.ATDuration)
	rt, rtErr := h.jwt.GenerateToken(user.Id.Hex(), h.jwtConfig.RTDuration)
	if atErr != nil || rtErr != nil {
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	if h.tr.SetAT(at) != nil || h.tr.SetRT(rt) != nil {
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	resp.SetStatus(fiber.StatusOK)
	resp.SetData([]*dtos.LoginResponseDto{
		{
			UserId:       user.Id.Hex(),
			AccessToken:  at.Token,
			RefreshToken: rt.Token,
		},
	})

	// TODO: change it
	h.jwt.SetJWTCookies(resp, &models.TokenPairs{
		AccessToken:  at,
		RefreshToken: rt,
	})

	return resp.JSON()
}
