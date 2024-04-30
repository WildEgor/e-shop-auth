package dtos

type VerifyIdentityRequestDto struct {
	Identity string `json:"identity" validate:"required"` // email or phone
	Code     string `json:"code" validate:"required"`
}
