package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vaporii/v8box/internal/auth"
	"github.com/vaporii/v8box/internal/config"

	"github.com/vaporii/v8box/internal/handler"
	"github.com/vaporii/v8box/internal/repository"
	"github.com/vaporii/v8box/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("err: %v\n", err)
		return
	}

	r, err := setupRouter(cfg)
	if err != nil {
		log.Fatalf("err: %v\n", err)
		return
	}

	http.ListenAndServe(cfg.ServerAddress, r)
}

func setupRouter(cfg *config.Config) (*chi.Mux, error) {
	r := chi.NewRouter()

	db, err := sql.Open("sqlite", cfg.SQLitePath)
	if err != nil {
		return nil, err
	}
	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		return nil, err
	}
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	noteRepo, err := repository.NewNoteRepository(db)
	if err != nil {
		return nil, err
	}
	noteService := service.NewNoteService(noteRepo)
	noteHandler := handler.NewNoteHandler(noteService)

	service, err := auth.RegisterHandlers(r, cfg)
	if err != nil {
		return nil, err
	}

	m := service.Middleware()

	r.Post("/register", userHandler.Register)
	r.Post("/register/oauth", userHandler.RegisterOAuth)
	r.With(m.Auth).Post("/note", noteHandler.Create)

	return r, nil
}
