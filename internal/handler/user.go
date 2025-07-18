package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vaporii/v8box/internal/config/provider"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/security"
	"github.com/vaporii/v8box/internal/service"
)

type UserHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	RegisterOAuth(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var login dto.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	h.userService.Register(login)
}

func (h *userHandler) RegisterOAuth(w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterOAuthRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	h.userService.RegisterOAuth(request)
}

func (h *userHandler) GitHubOAuthLogin(w http.ResponseWriter, r *http.Request) {
	cfg, err := provider.LoadGithubOAuthConfig()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	stateToken := security.GenerateStateToken()
	http.Redirect(w, r, cfg.AuthCodeURL(stateToken), 302)

	tok, err := cfg.Exchange(oa)
}
