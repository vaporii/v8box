package provider

import "os"

type OauthConfig struct {
	ProviderName string
	ClientID     string
	ClientSecret string
}

// TODO: errors for null vals
func LoadGithubConfig() (OauthConfig, error) {
	return OauthConfig{
		ProviderName: "github",
		ClientID:     getEnv("V8BOX_GITHUB_CLIENT_ID", ""),
		ClientSecret: getEnv("V8BOX_GITHUB_CLIENT_SECRET", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
