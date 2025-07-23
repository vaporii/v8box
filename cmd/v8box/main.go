package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/middleware"

	"github.com/vaporii/v8box/internal/handler"
)

func main() {
	cfg := config.LoadConfig()

	r := chi.NewRouter()

	r.Use(middleware.ErrorHandler)

	routes, err := setupRouter(cfg)
	if err != nil {
		log.Fatalf("err: %v\n", err)
		return
	}

	r.Mount("/api/v1", routes)

	http.ListenAndServe(cfg.ServerAddress, r)
}

func setupRouter(cfg *config.Config) (*chi.Mux, error) {
	r := chi.NewRouter()

	db, err := sql.Open("sqlite", cfg.SQLitePath)
	if err != nil {
		return nil, err
	}

	handlers := handler.NewHandlers(db, *config.LoadConfig())

	r.Mount("/auth", setupAuthRoutes(handlers.AuthHandler))
	r.Mount("/me", setupMeRoutes(handlers))

	return r, nil
}

func setupMeRoutes(handlers *handler.Handlers) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Auth)

	r.Get("/", handlers.UserHandler.GetCurrentUser)
	r.Get("/note", handlers.NoteHandler.GetNotes)
	r.Get("/note/{id}", handlers.NoteHandler.GetNoteByID)
	r.Post("/note", handlers.NoteHandler.Create)
	r.Put("/note/{id}", handlers.NoteHandler.EditNoteByID)

	return r
}

func setupAuthRoutes(authHandler handler.AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/login/github", authHandler.GitHubOAuthLogin)
	r.Get("/callback/github", authHandler.GitHubOAuthCallback)
	r.Get("/login/google", authHandler.GoogleOAuthLogin)
	r.Get("/callback/google", authHandler.GoogleOAuthCallback)

	return r
}
