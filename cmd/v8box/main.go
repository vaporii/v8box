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
	r.Get("/notes", handlers.NoteHandler.GetNotes)
	r.Get("/notes/{id}", handlers.NoteHandler.GetNoteByID)
	r.Post("/notes", handlers.NoteHandler.Create)
	r.Put("/notes/{id}", handlers.NoteHandler.EditNoteByID)

	r.Get("/notes/{note_id}/tags", handlers.TagHandler.GetTagsOnNote)
	r.Get("/tags/{tag_id}/notes", handlers.TagHandler.GetNotesWithTag)
	r.Post("/tags", handlers.TagHandler.Create)
	r.Put("/tags/{id}", handlers.TagHandler.EditTag)
	r.Get("/tags/{id}", handlers.TagHandler.GetTag)

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
