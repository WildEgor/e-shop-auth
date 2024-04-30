package dtos

// OTPGenerateRequestDto - OTP generation request
type OTPGenerateRequestDto struct {
	Phone string `json:"phone" validate:"required"`
}

type OTPGenerateResponseDto struct {
	Identity string `json:"identity_type"`
	Code     string `json:"code"`
}
