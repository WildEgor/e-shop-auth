package dtos

type RegistrationRequestDto struct {
	Email     string `json:"email" validate:"required,email,lte=255"`
	Phone     string `json:"phone" validate:"required,lte=255"`
	Password  string `json:"password" validate:"required,lte=255"`
	FirstName string `json:"first_name,omitempty" validate:"lte=255"`
	LastName  string `json:"last_name,omitempty" validate:"lte=255"`
}
