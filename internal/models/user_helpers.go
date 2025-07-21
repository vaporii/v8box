package models

import (
	"encoding/json"
	"net/http"

	"github.com/vaporii/v8box/internal/logging"
	"github.com/vaporii/v8box/internal/middleware"
)

func ExtractUser(r *http.Request) User {
	var user User

	temp, _ := json.Marshal(r.Context().Value(middleware.UserAuthContextKey))
	err := json.Unmarshal(temp, &user)

	if err != nil {
		logging.Error("couldn't get user from context: %v", err)
	}

	return user
}
