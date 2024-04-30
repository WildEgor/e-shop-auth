package dtos

type LoginRequestDto struct {
	Login    string `json:"login" validate:"required,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
}

type LoginResponseDto struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
