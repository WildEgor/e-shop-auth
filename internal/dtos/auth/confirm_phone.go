package dtos

type ConfirmPhoneRequestDto struct {
	Phone string `json:"phone" validate:"required"`
	Code  string `json:"code" validate:"required"`
}
