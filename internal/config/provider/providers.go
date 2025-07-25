package provider

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type OauthConfig struct {
	ProviderName string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

var githubOAuthConfig *oauth2.Config

// TODO: errors for null vals
func LoadGithubConfig() OauthConfig {
	return OauthConfig{
		ProviderName: "github",
		ClientID:     getEnv("V8BOX_GITHUB_CLIENT_ID", ""),
		ClientSecret: getEnv("V8BOX_GITHUB_CLIENT_SECRET", ""),
		RedirectURL:  getEnv("V8BOX_GITHUB_REDIRECT_URL", ""),
	}
}

func LoadGithubOAuthConfig() *oauth2.Config {
	if githubOAuthConfig != nil {
		return githubOAuthConfig
	}
	cfg := LoadGithubConfig()

	config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{},
		Endpoint:     github.Endpoint,
	}
	githubOAuthConfig = config

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
