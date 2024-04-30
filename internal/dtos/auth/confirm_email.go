package dtos

type ConfirmEmailRequestDto struct {
	Email string `json:"email" validate:"required"`
	Code  string `json:"code" validate:"required"`
}
