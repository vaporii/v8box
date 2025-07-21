package service

import (
	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/repository"
)

type UserService interface {
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
