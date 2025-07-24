package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vaporii/v8box/internal/httperror"
	"github.com/vaporii/v8box/internal/logging"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.WithValue(r.Context(), httperror.ErrorKey, &err)
		next.ServeHTTP(w, r.WithContext(ctx))

		if err == nil {
			return
		}

		logging.Warning("err during HTTP request: %v", err)

		switch t := err.(type) {
		default:
			httpError(w, "Internal Server Error", 500)
		case *httperror.NotFoundError:
			httpError(w, t.Error(), 404)
		case *httperror.BadClientRequestError:
			httpError(w, t.Error(), 400)
		case *httperror.UnauthorizedError:
			httpError(w, t.Error(), 401)
		}
	})
}

func httpError(w http.ResponseWriter, errorMsg string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, "{\"error\":\"%s\"}", errorMsg)
}
