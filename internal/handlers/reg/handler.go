package handlers

import (
	"github.com/WildEgor/e-shop-auth/internal/configs"
	domains "github.com/WildEgor/e-shop-auth/internal/domain"
	dtos "github.com/WildEgor/e-shop-auth/internal/dtos/auth"
	"github.com/WildEgor/e-shop-auth/internal/mappers"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	"github.com/WildEgor/e-shop-auth/internal/validators"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
	"log/slog"
)

type RegHandler struct {
	ur        *repositories.UserRepository
	tr        *repositories.TokensRepository
	jwt       *services.JWTAuthenticator
	jwtConfig *configs.JWTConfig
}

func NewRegHandler(
	ur *repositories.UserRepository,
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
	jwtConfig *configs.JWTConfig,
) *RegHandler {
	return &RegHandler{
		ur:        ur,
		tr:        tr,
		jwt:       jwt,
		jwtConfig: jwtConfig,
	}
}

// TODO: add swagger response
// Reg new user godoc
//
//	@Param			body body	dtos.RegistrationRequestDto	true	"Body"
//	@Tags			User Controller
//	@Accept			json
//	@Produce		json
//
// @Success 200 {object} map[string]string
// @Router /api/v1/user/reg [post]
func (h *RegHandler) Handle(ctx fiber.Ctx) error {
	dto := &dtos.RegistrationRequestDto{}
	if err := validators.ParseAndValidate(ctx, dto); err != nil {
		return err
	}

	resp := core_dtos.NewResponse(ctx)

	// 1. Check if user exists
	existed := h.ur.CheckExistsEmail(dto.Email)
	if existed {
		resp.SetStatus(fiber.StatusConflict)
		return resp.JSON()
	}

	// 2. Create user of not exists
	userModel := mappers.CreateUser(dto)
	newUser, mongoErr := h.ur.Create(userModel)
	if mongoErr != nil {
		domains.SetEmailAlreadyExistError(resp)
		domains.SetPhoneAlreadyExistError(resp)
		return resp.JSON()
	}

	// 3. GenShortCode tokens
	at, atErr := h.jwt.GenerateToken(newUser.Id.Hex(), h.jwtConfig.ATDuration)
	rt, rtErr := h.jwt.GenerateToken(newUser.Id.Hex(), h.jwtConfig.RTDuration)
	if atErr != nil || rtErr != nil {
		slog.Error("error gen token pairs", atErr, rtErr)
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	if h.tr.SetAT(at) != nil || h.tr.SetRT(rt) != nil {
		slog.Error("error cache tokens")
		resp.SetStatus(fiber.StatusInternalServerError)
		return resp.JSON()
	}

	resp.SetStatus(fiber.StatusCreated)
	resp.SetData([]*dtos.TokenPairResponseDto{{
		UserID:       newUser.Id.Hex(),
		AccessToken:  at.Token,
		RefreshToken: rt.Token,
	}})

	return resp.JSON()
}
