package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vaporii/v8box/internal/logging"
	"github.com/vaporii/v8box/internal/middleware"
	"github.com/vaporii/v8box/internal/service"
)

type UserHandler interface {
	GetCurrentUser(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
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
