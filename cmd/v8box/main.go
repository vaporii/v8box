package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vaporii/v8box/internal/auth"
	"github.com/vaporii/v8box/internal/config"
)

func main() {
	r := chi.NewRouter()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("err: %v\n", err)
		return
	}

	service, err := auth.RegisterHandlers(r, cfg)
	if err != nil {
		log.Fatalf("err: %v\n", err)
		return
	}

	m := service.Middleware()

	r.With(m.Auth).Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	http.ListenAndServe(cfg.ServerAddress, r)
}
