package service

import (
	"github.com/google/uuid"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/httperror"
	"github.com/vaporii/v8box/internal/models"
	"github.com/vaporii/v8box/internal/repository"
)

type NoteService interface {
	Create(request dto.CreateNoteRequest) (*models.Note, error)
	GetUserNotes(userId string) ([]models.Note, error)
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

	note := &models.Note{
		ID:      uuid.NewString(),
		UserID:  request.UserID,
		Title:   request.Title,
		Content: request.Content,
	}

	err := s.noteRepo.CreateNote(note)
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
