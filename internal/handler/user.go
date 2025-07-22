package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vaporii/v8box/internal/httperror"
	"github.com/vaporii/v8box/internal/logging"
	"github.com/vaporii/v8box/internal/models"
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
	user, err := h.userService.GetUser(models.ExtractUser(r).UserID)
	if checkErr(err, r) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(*user)

	if checkErr(err, r) {
		return
	}
}

func checkErr(err error, r *http.Request) bool {
	if err != nil {
		errorVal := r.Context().Value(httperror.ErrorKey).(*error)
		*errorVal = err
		// ctx := context.WithValue(r.Context(), httperror.ErrorKey, err)
		// *r = *r.WithContext(ctx)
		logging.Warning("HTTP error: %v", err)

		return true
	}
	return false
}
