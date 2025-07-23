package service

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/httperror"
	"github.com/vaporii/v8box/internal/models"
	"github.com/vaporii/v8box/internal/repository"
)

type NoteService interface {
	Create(request dto.CreateNoteRequest) (*models.Note, error)
	GetUserNotes(userId string) ([]models.Note, error)
	GetNoteByID(id string) (*models.Note, error)
	EditNoteByID(id string, request dto.CreateNoteRequest) (*models.Note, error)
}

type noteService struct {
	noteRepo    repository.NoteRepository
	userService UserService
}

func NewNoteService(noteRepo repository.NoteRepository, userService UserService) NoteService {
	return &noteService{
		noteRepo:    noteRepo,
		userService: userService,
	}
}

func (s *noteService) Create(request dto.CreateNoteRequest) (*models.Note, error) {
	userExists := s.userService.CheckUserExists(request.UserID)
	if !userExists {
		return nil, &httperror.BadClientRequestError{Message: "User with ID doesn't exist"}
	}

	note, err := s.noteRepo.CreateNote(uuid.NewString(), &request)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *noteService) GetUserNotes(userId string) ([]models.Note, error) {
	notes, err := s.noteRepo.GetUserNotes(userId)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *noteService) GetNoteByID(id string) (*models.Note, error) {
	note, err := s.noteRepo.GetNoteByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &httperror.NotFoundError{Entity: "Note"}
		}
		return nil, err
	}

	return note, nil
}

func (s *noteService) EditNoteByID(id string, request dto.CreateNoteRequest) (*models.Note, error) {
	_, err := s.noteRepo.GetNoteByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &httperror.NotFoundError{Entity: "Note"}
		}
		return nil, err
	}

	note, err := s.noteRepo.UpdateNote(id, request)
	if err != nil {
		return nil, err
	}

	return note, nil
}
