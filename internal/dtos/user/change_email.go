package dtos

type ChangeEmailRequestDto struct {
	Email string `json:"new_email" validate:"email,required"`
}
