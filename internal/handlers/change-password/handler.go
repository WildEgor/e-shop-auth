package handlers

import (
	"github.com/WildEgor/e-shop-auth/internal/configs"
	domains "github.com/WildEgor/e-shop-auth/internal/domain"
	dtos2 "github.com/WildEgor/e-shop-auth/internal/dtos/auth"
	dtos "github.com/WildEgor/e-shop-auth/internal/dtos/user"
	auth_middleware "github.com/WildEgor/e-shop-auth/internal/middlewares/auth"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	"github.com/WildEgor/e-shop-auth/internal/validators"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
)

type ChangePasswordHandler struct {
	ur        *repositories.UserRepository
	tr        *repositories.TokensRepository
	jwt       *services.JWTAuthenticator
	jwtConfig *configs.JWTConfig
}

func NewChangePasswordHandler(
	ur *repositories.UserRepository,
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
	jwtConfig *configs.JWTConfig,
) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		ur:        ur,
		tr:        tr,
		jwt:       jwt,
		jwtConfig: jwtConfig,
	}
}

// Change password godoc
//
//	@Param			Authorization header	string	true	"123"
//	@Param			body body	dtos.ChangePasswordRequestDto	true	"Body"
//	@Summary		Allow change authenticated user own password
//	@Description	Allow change authenticated user own password
//	@Tags			Auth Controller
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} dtos.ChangePasswordRequestDto
// @Router /api/v1/auth/change-password [post]
func (h *ChangePasswordHandler) Handle(ctx fiber.Ctx) error {
	dto := &dtos.ChangePasswordRequestDto{}
	resp := core_dtos.NewResponse(ctx)
	resp.SetStatus(fiber.StatusBadRequest)

	if err := validators.ParseAndValidate(ctx, dto); err != nil {
		return err
	}

	authUser := auth_middleware.ExtractUser(ctx)

	if _, err := authUser.ComparePassword(dto.OldPassword); err != nil {
		domains.SetInvalidCredentialError(resp)
		return resp.JSON()
	}

	if err := authUser.SetPassword(dto.NewPassword); err != nil {
		domains.SetInvalidCredentialError(resp)
		return resp.JSON()
	}

	if err := h.ur.UpdatePassword(authUser); err != nil {
		domains.SetInvalidCredentialError(resp)
		return resp.JSON()
	}

	// 3. GenShortCode tokens
	at, atErr := h.jwt.GenerateToken(authUser.Id.Hex(), h.jwtConfig.ATDuration)
	rt, rtErr := h.jwt.GenerateToken(authUser.Id.Hex(), h.jwtConfig.RTDuration)
	if atErr != nil || rtErr != nil {
		domains.SetInvalidCredentialError(resp)
		return resp.JSON()
	}

	if h.tr.SetAT(at) != nil || h.tr.SetRT(rt) != nil {
		domains.SetInvalidCredentialError(resp)
		return resp.JSON()
	}

	resp.SetStatus(fiber.StatusOK)

	resp.SetData([]*dtos2.TokenPairResponseDto{
		{
			UserID:       authUser.Id.Hex(),
			AccessToken:  at.Token,
			RefreshToken: rt.Token,
		},
	})

	return resp.JSON()
}
