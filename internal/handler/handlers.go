package handler

import (
	"database/sql"
	"log"

	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/repository"
	"github.com/vaporii/v8box/internal/service"
)

type Handlers struct {
	UserHandler UserHandler
	NoteHandler NoteHandler
	AuthHandler AuthHandler
}

func NewHandlers(db *sql.DB, cfg config.Config) *Handlers {
	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		log.Fatalf("err setting up user repository: %v\n", err)
		return nil
	}

	noteRepo, err := repository.NewNoteRepository(db)
	if err != nil {
		log.Fatalf("err setting up note repository: %v\n", err)
		return nil
	}

	return &Handlers{
		UserHandler: NewUserHandler(service.NewUserService(userRepo, cfg)),
		NoteHandler: NewNoteHandler(service.NewNoteService(noteRepo)),
		AuthHandler: NewAuthHandler(service.NewAuthService(userRepo, cfg)),
	}
}
