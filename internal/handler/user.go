package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vaporii/v8box/internal/dto"
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

}
