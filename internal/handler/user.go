package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vaporii/v8box/internal/config/provider"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/logging"
	"github.com/vaporii/v8box/internal/middleware"
	"github.com/vaporii/v8box/internal/security"
	"github.com/vaporii/v8box/internal/service"
)

type UserHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	GitHubOAuthLogin(w http.ResponseWriter, r *http.Request)
	GitHubOAuthCallback(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
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

func (h *userHandler) GitHubOAuthLogin(w http.ResponseWriter, r *http.Request) {
	cfg := provider.LoadGithubOAuthConfig()

	stateToken := security.GenerateStateToken()
	http.Redirect(w, r, cfg.AuthCodeURL(stateToken), http.StatusFound)
}

func (h *userHandler) GitHubOAuthCallback(w http.ResponseWriter, r *http.Request) {
	claims, err := h.userService.GetGitHubOAuthJwt(r.Context(), r.FormValue("code"))
	if err != nil {
		checkErr(err, w, http.StatusText(500), 500)
		return
	}

	err = h.userService.RegisterOAuthUser(claims)
	if checkErr(err, w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) {
		return
	}

	jwtToken, err := h.userService.CreateJWT(claims)
	if checkErr(err, w, http.StatusText(500), 500) {
		return
	}

	w.Header().Set("Set-Cookie", "JWT="+jwtToken)
	w.WriteHeader(http.StatusOK)
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(r.Context().Value(middleware.UserAuthContextKey))
	if checkErr(err, w, http.StatusText(500), 500) {
		return
	}
}

func checkErr(err error, w http.ResponseWriter, statusText string, statusCode int) bool {
	if err != nil {
		logging.Warning("HTTP error: %d %s err: %v", statusCode, statusText, err)
		http.Error(w, statusText, statusCode)
		return true
	}
	return false
}
