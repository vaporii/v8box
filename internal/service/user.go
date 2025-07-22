package service

import (
	"database/sql"
	"errors"

	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/httperror"
	"github.com/vaporii/v8box/internal/models"
	"github.com/vaporii/v8box/internal/repository"
)

type UserService interface {
	GetUser(userId string) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	conf     config.Config
}

func NewUserService(userRepo repository.UserRepository, conf config.Config) UserService {
	return &userService{
		userRepo: userRepo,
		conf:     conf,
	}
}

func (s *userService) GetUser(userId string) (*models.User, error) {
	user, err := s.userRepo.GetUserById(userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &httperror.NotFoundError{Entity: "User"}
	}
	return user, err
}
