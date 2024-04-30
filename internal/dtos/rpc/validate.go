package dtos

type ValidateTokenPayloadDto struct {
	Token string `json:"token" validate:"required"`
}

type CentrifugeValidTokenResponseDto struct {
	User string `json:"user"`
}

type ValidateTokenResponseDto struct {
	Result CentrifugeValidTokenResponseDto `json:"result"`
}

type CentrifugeInvalidTokenResponseDto struct {
	Code   uint32 `json:"code"`
	Reason string `json:"reason"`
}

// ValidateTokenUnauthorizedResponseDto
//
//		{
//	 "disconnect": {
//	   "code": 4501,
//	   "reason": "unauthorized"
//	 }
//	}
//
// /*
type ValidateTokenUnauthorizedResponseDto struct {
	Disconnect CentrifugeInvalidTokenResponseDto `json:"disconnect"`
}

type CentrifugeRefreshResponseDto struct {
	Expired  bool `json:"expired"`
	ExpireAt uint `json:"expire_at"`
}

type RefreshResponseDto struct {
	Result CentrifugeRefreshResponseDto `json:"result"`
}
