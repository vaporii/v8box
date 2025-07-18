package service

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/models"
	"github.com/vaporii/v8box/internal/repository"
	"github.com/vaporii/v8box/internal/security"
)

type UserService interface {
	Register(request dto.RegisterRequest) (*models.User, error)
	RegisterOAuth(request dto.RegisterOAuthRequest) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (r *userService) Register(request dto.RegisterRequest) (*models.User, error) {
	_, err := r.userRepo.GetUserByUsername(request.Username)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return nil, err
	}

	hashedPassword, err := security.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       uuid.NewString(),
		Username: request.Username,
		Password: hashedPassword,
	}

	err = r.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userService) RegisterOAuth(request dto.RegisterOAuthRequest) (*models.User, error) {
	_, err := r.userRepo.GetUserByUsername(request.Username)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return nil, err
	}

	user := &models.User{
		ID:            uuid.NewString(),
		Username:      request.Username,
		OAuthProvider: request.OAuthProvider,
		OAuthID:       request.OAuthID,
		AccessToken:   request.AccessToken,
		RefreshToken:  request.RefreshToken,
		TokenExpiry:   request.TokenExpiry,
	}

	err = r.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
