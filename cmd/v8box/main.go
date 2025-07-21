package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/middleware"

	"github.com/vaporii/v8box/internal/handler"
	"github.com/vaporii/v8box/internal/repository"
	"github.com/vaporii/v8box/internal/service"
)

func main() {
	cfg := config.LoadConfig()

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
	userService := service.NewUserService(userRepo, *config.LoadConfig())
	userHandler := handler.NewUserHandler(userService)

	noteRepo, err := repository.NewNoteRepository(db)
	if err != nil {
		return nil, err
	}
	noteService := service.NewNoteService(noteRepo)
	_ = handler.NewNoteHandler(noteService)

	r.Post("/register", userHandler.Register)
	r.Get("/register/oauth", userHandler.GitHubOAuthLogin)
	r.Get("/callback", userHandler.GitHubOAuthCallback)
	r.With(middleware.Auth).Get("/user", userHandler.GetUser)

	return r, nil
}
