package service

import (
	"github.com/google/uuid"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/httperror"
	"github.com/vaporii/v8box/internal/models"
	"github.com/vaporii/v8box/internal/repository"
)

type TagService interface {
	CreateTag(tagID string, tag dto.Tag) (*models.Tag, error)
	GetTagsOnNote(noteID string) ([]models.Tag, error)
	GetNotesWithTag(tagID string) ([]models.Note, error)
	EditTag(tagID string, tagDTO dto.Tag) (*models.Tag, error)
	GetTag(tagID string) (*models.Tag, error)
}

type tagService struct {
	tagRepo     repository.TagRepository
	userService UserService
}

func NewTagService(tagRepo repository.TagRepository, userService UserService) TagService {
	return &tagService{
		tagRepo:     tagRepo,
		userService: userService,
	}
}

func (s *tagService) CreateTag(userID string, tagDTO dto.Tag) (*models.Tag, error) {
	if !s.userService.CheckUserExists(userID) {
		return nil, &httperror.UnauthorizedError{Message: "User doesn't exist"}
	}

	tag, err := s.tagRepo.CreateTag(uuid.NewString(), userID, tagDTO)
	return tag, err
}

func (s *tagService) EditTag(tagID string, tagDTO dto.Tag) (*models.Tag, error) {
	return s.tagRepo.UpdateTag(tagID, &tagDTO)
}

func (s *tagService) GetTag(tagID string) (*models.Tag, error) {
	return s.tagRepo.GetTagById(tagID)
}

func (s *tagService) GetTagsOnNote(noteID string) ([]models.Tag, error) {
	tags, err := s.tagRepo.GetTagsOnNote(noteID)

	return tags, err
}

func (s *tagService) GetNotesWithTag(tagID string) ([]models.Note, error) {
	notes, err := s.tagRepo.GetNotesWithTag(tagID)

	return notes, err
}
