package dto

import "github.com/golang-jwt/jwt/v5"

type UserJwtPackage struct {
	Username  string `json:"username"`
	UserID    string `json:"user_id"`
	AvatarURL string `json:"avatar_url"`
	OAuthKey  string `json:"oauth_key,omitempty"`
	jwt.RegisteredClaims
}
