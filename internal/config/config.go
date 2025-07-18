package config

import (
	"os"
	"time"
)

type Config struct {
	TokenDuration  time.Duration
	CookieDuration time.Duration
	Issuer         string
	URL            string
	AvatarPath     string
	DisableXSRF    bool
	TokenSecret    string
	ServerAddress  string
	SQLitePath     string
	Environment    string
}

func LoadConfig() (*Config, error) {
	return &Config{
		TokenDuration:  5 * time.Minute,
		CookieDuration: 24 * time.Hour,
		Issuer:         getEnv("V8BOX_ISSUER", "v8box"),
		URL:            getEnv("V8BOX_URL", ""),
		AvatarPath:     getEnv("V8BOX_AVATAR_PATH", "/tmp"),
		DisableXSRF:    getEnvAsBool("V8BOX_DISABLE_XSRF", true),
		TokenSecret:    getEnv("V8BOX_TOKEN_SECRET", "secret"),
		ServerAddress:  getEnv("V8BOX_ADDRESS", ":3000"),
		SQLitePath:     getEnv("V8BOX_SQLITE_PATH", "./dev.db"),
		Environment:    getEnv("V8BOX_ENVIRONMENT", "dev"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true"
	}
	return defaultValue
}
