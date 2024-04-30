package models

type TokenDetails struct {
	Token     string
	TokenUuid string
	UserID    string
	ExpiresIn int64
	IsValid   bool
}

type TokenPairs struct {
	AccessToken  *TokenDetails
	RefreshToken *TokenDetails
}
