package dto

type RegisterOAuthRequest struct {
	Username      string `json:"username" validate:"required"`
	OAuthProvider string `json:"oauth_provider,omitempty"`
	OAuthID       string `json:"oauth_id" validate:"required"`
	AccessToken   string `json:"access_token" validate:"required"`
	RefreshToken  string `json:"refresh_token" validate:"required"`
	TokenExpiry   int64  `json:"token_expiry" validate:"required"`
}
