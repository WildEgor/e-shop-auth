package services

import (
	"fmt"
	"github.com/WildEgor/e-shop-auth/internal/configs"
	"github.com/WildEgor/e-shop-auth/internal/models"
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type JWTAuthenticator struct {
	jwtConfig *configs.JWTConfig
}

func NewJWTAuthenticator(jwtConfig *configs.JWTConfig) *JWTAuthenticator {

	return &JWTAuthenticator{
		jwtConfig: jwtConfig,
	}
}

func (j *JWTAuthenticator) GenerateToken(
	userId string,
	duration time.Duration,
) (*models.TokenDetails, error) {
	now := time.Now().UTC()
	td := &models.TokenDetails{
		ExpiresIn: 0,
		Token:     "",
	}
	td.ExpiresIn = now.Add(duration).Unix()
	td.TokenUuid = uuid.NewString()
	td.UserID = userId

	claims := make(jwt.MapClaims)
	claims["sub"] = userId
	claims["token_uuid"] = td.TokenUuid
	claims["exp"] = td.ExpiresIn
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(j.jwtConfig.Secret))
	if err != nil {
		return nil, errors.Wrap(err, "generate access token")
	}

	td.Token = ss

	return td, nil
}

func (j *JWTAuthenticator) ParseToken(token string) (*models.TokenDetails, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		s := []byte(j.jwtConfig.Secret)

		return s, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "parse token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || !parsed.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	err = claims.Valid()

	td := &models.TokenDetails{
		TokenUuid: fmt.Sprint(claims["token_uuid"]),
		UserID:    fmt.Sprint(claims["sub"]),
		IsValid:   err == nil,
	}

	return td, err
}

func (j *JWTAuthenticator) SetJWTCookies(resp *core_dtos.ResponseDto, tokens *models.TokenPairs) {
	//atCookie := fiber.Cookie{
	//	Name:     "jwt_access_token",
	//	Value:    tokens.AccessToken.Token,
	//	Expires:  time.UnixMilli(tokens.AccessToken.ExpiresIn),
	//	HTTPOnly: true,
	//}
	//
	//rtCookie := fiber.Cookie{
	//	Name:     "jwt_refresh_token",
	//	Value:    tokens.RefreshToken.Token,
	//	Expires:  time.UnixMilli(tokens.RefreshToken.ExpiresIn),
	//	HTTPOnly: true,
	//}
	//
	//ctx.Cookie(&atCookie)
	//ctx.Cookie(&rtCookie)
}
