package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	OAuthKey string `json:"oauth_key,omitempty"`
}
