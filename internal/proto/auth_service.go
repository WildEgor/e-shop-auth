package proto

import (
	"context"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	"github.com/pkg/errors"
)

type AuthService struct {
	ur  *repositories.UserRepository
	jwt *services.JWTAuthenticator
}

func NewAuthService(
	ur *repositories.UserRepository,
	jwt *services.JWTAuthenticator,
) *AuthService {
	return &AuthService{
		ur,
		jwt,
	}
}

func (s *AuthService) ValidateToken(ctx context.Context, request *ValidateTokenRequest) (*UserData, error) {
	if len(request.Token) == 0 {
		// TODO: return rpc error?
		return nil, errors.New("empty token")
	}

	// TODO: need check token in Redis too
	claims, err := s.jwt.ParseToken(string(request.Token[:]))

	if claims == nil || !claims.IsValid {
		return nil, errors.Wrap(err, "token validation")
	}

	ur, err := s.ur.FindById(claims.UserID)
	if err != nil {
		// TODO: return rpc error?
		return nil, err
	}

	return &UserData{
		Id:        ur.Id.Hex(),
		FirstName: ur.FirstName,
		LastName:  ur.LastName,
		Email:     ur.Email,
		Phone:     ur.Phone,
		IsActive:  ur.IsActive(),
		// TODO: add ExpAt: ur.ExpiresIn
	}, nil
}

func (s *AuthService) FindByIds(ctx context.Context, request *FindByIdsRequest) (*FindByIdsResponse, error) {
	var response FindByIdsResponse

	if len(request.Ids) <= 0 {
		return &response, nil
	}

	users, err := s.ur.FindByIds(request.Ids)
	if err != nil {
		// TODO: return rpc error?
		return &response, err
	}

	response.Total, _ = s.ur.CountAllActive()

	// TODO: add mapper
	for _, model := range *users {
		response.Users = append(response.Users, &UserData{
			Id:        model.Id.Hex(),
			FirstName: model.FirstName,
			LastName:  model.LastName,
			Email:     model.Email,
			Phone:     model.Phone,
			IsActive:  model.IsActive(),
		})
	}

	return &response, nil
}

func (s *AuthService) mustEmbedUnimplementedAuthServiceServer() {
	panic("implement me")
}
