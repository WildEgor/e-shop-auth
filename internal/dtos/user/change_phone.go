package dtos

type ChangePhoneRequestDto struct {
	Phone string `json:"phone" validate:"required"`
}
