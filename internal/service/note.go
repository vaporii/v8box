package service

import (
	"github.com/google/uuid"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/models"
	"github.com/vaporii/v8box/internal/repository"
)

type NoteService interface {
	Create(request dto.CreateNoteRequest) (*models.Note, error)
	GetUserNotes(userId string) ([]models.Note, error)
}

type noteService struct {
	noteRepo repository.NoteRepository
}

func NewNoteService(noteRepo repository.NoteRepository) NoteService {
	return &noteService{
		noteRepo: noteRepo,
	}
}

func (s *noteService) Create(request dto.CreateNoteRequest) (*models.Note, error) {
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
