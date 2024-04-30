package mappers

import (
	dtos "github.com/WildEgor/e-shop-auth/internal/dtos/auth"
	"github.com/WildEgor/e-shop-auth/internal/models"
)

func CreateUser(dto *dtos.RegistrationRequestDto) *models.UsersModel {
	us := &models.UsersModel{
		Email:     dto.Email,
		Phone:     dto.Phone,
		Password:  dto.Password,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
	}

	err := us.SetPassword(dto.Password)
	if err != nil {
		return nil
	}

	return us
}
