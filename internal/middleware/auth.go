package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/logging"
)

type userAuthKeyType string

const UserAuthContextKey userAuthKeyType = "user"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context

		conf := config.LoadConfig()

		cookie, err := r.Cookie("JWT")
		if checkErr(err, w, "No JWT cookie provided", 401) {
			return
		}

		parsedTok, err := jwt.ParseWithClaims(cookie.Value, &dto.UserJwtPackage{}, func(t *jwt.Token) (any, error) {
			return []byte(conf.JwtSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if checkErr(err, w, "Bad token", 401) {
			return
		}

		if claims, ok := parsedTok.Claims.(*dto.UserJwtPackage); ok {
			ctx = context.WithValue(r.Context(), UserAuthContextKey, *claims)
		} else {
			checkErr(errors.New("claims not ok"), w, "Internal Server Error", 500)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkErr(err error, w http.ResponseWriter, statusText string, statusCode int) bool {
	if err != nil {
		logging.Warning("HTTP error: %d %s err: %v", statusCode, statusText, err)
		http.Error(w, statusText, statusCode)
		return true
	}
	return false
}
