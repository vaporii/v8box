package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/httperror"
	"github.com/vaporii/v8box/internal/logging"
)

type userAuthKeyType string

const UserAuthContextKey userAuthKeyType = "user"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context

		conf := config.LoadConfig()

		cookie, err := r.Cookie("JWT")
		// if checkErr(err, w, "No JWT cookie provided", 401) {
		if checkErr(err, r) {
			return
		}

		parsedTok, err := jwt.ParseWithClaims(cookie.Value, &dto.UserJwtPackage{}, func(t *jwt.Token) (any, error) {
			return []byte(conf.JwtSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		// if checkErr(err, w, "Bad token", 401) {
		if checkErr(err, r) {
			return
		}

		if claims, ok := parsedTok.Claims.(*dto.UserJwtPackage); ok {
			ctx = context.WithValue(r.Context(), UserAuthContextKey, *claims)
		} else {
			checkErr(errors.New("claims not ok"), r)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkErr(err error, r *http.Request) bool {
	if err != nil {
		ctx := context.WithValue(r.Context(), httperror.ErrorKey, err)
		*r = *r.WithContext(ctx)
		logging.Warning("HTTP Auth error: %v", err)

		return true
	}
	return false
}
