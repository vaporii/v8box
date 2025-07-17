package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-pkgz/auth/v2"
	"github.com/go-pkgz/auth/v2/avatar"
	"github.com/go-pkgz/auth/v2/token"
	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/config/provider"
)

func RegisterHandlers(router chi.Router, cfg *config.Config) (*auth.Service, error) {
	options := auth.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) {
			return cfg.TokenSecret, nil
		}),
		TokenDuration:  cfg.TokenDuration,
		CookieDuration: cfg.CookieDuration,
		Issuer:         cfg.Issuer,
		URL:            cfg.URL,
		AvatarStore:    avatar.NewLocalFS(cfg.AvatarPath),
		DisableXSRF:    cfg.DisableXSRF,
	}
	service := auth.NewService(options)

	githubCfg, err := provider.LoadGithubConfig()
	if err != nil {
		return nil, err
	}
	service.AddProvider(githubCfg.ProviderName, githubCfg.ClientID, githubCfg.ClientSecret)

	authRoutes, avaRoutes := service.Handlers()
	router.Mount("/auth", authRoutes)
	router.Mount("/avatar", avaRoutes)

	return service, nil
}
