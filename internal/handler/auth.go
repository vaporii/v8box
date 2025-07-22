package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vaporii/v8box/internal/config/provider"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/security"
	"github.com/vaporii/v8box/internal/service"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	GitHubOAuthLogin(w http.ResponseWriter, r *http.Request)
	GitHubOAuthCallback(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var login dto.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	h.authService.Register(login)
}

func (h *authHandler) GitHubOAuthLogin(w http.ResponseWriter, r *http.Request) {
	cfg := provider.LoadGithubOAuthConfig()

	stateToken := security.GenerateStateToken()
	http.Redirect(w, r, cfg.AuthCodeURL(stateToken), http.StatusFound)
}

func (h *authHandler) GitHubOAuthCallback(w http.ResponseWriter, r *http.Request) {
	claims, err := h.authService.GetGitHubOAuthJwt(r.Context(), r.FormValue("code"))
	if checkErr(err, r) {
		return
	}

	jwtToken, err := h.authService.CreateJWT(claims)
	if checkErr(err, r) {
		return
	}

	w.Header().Set("Set-Cookie", fmt.Sprintf("JWT=%s; Path=/", jwtToken))
	w.WriteHeader(http.StatusOK)
}
