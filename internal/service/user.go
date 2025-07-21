package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/config/provider"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/models"
	githubprovider "github.com/vaporii/v8box/internal/models/github_provider"
	"github.com/vaporii/v8box/internal/repository"
	"github.com/vaporii/v8box/internal/security"
)

type UserService interface {
	Register(request dto.RegisterRequest) (*models.User, error)
	GetGitHubOAuthJwt(ctx context.Context, code string) (dto.UserJwtPackage, error)
	CreateJWT(claims dto.UserJwtPackage) (string, error)
	RegisterOAuthUser(jwt dto.UserJwtPackage) error
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

func (r *userService) RegisterOAuthUser(jwt dto.UserJwtPackage) error {
	_, err := r.userRepo.GetUserByOAuthKey(jwt.OAuthKey)
	if errors.Is(err, sql.ErrNoRows) {
		user := &models.User{
			ID:       uuid.NewString(),
			Username: jwt.Username,
			OAuthKey: jwt.OAuthKey,
		}

		r.userRepo.CreateUser(user)

		return nil
	}

	return err
}

func (s *userService) GetGitHubOAuthJwt(ctx context.Context, code string) (dto.UserJwtPackage, error) {
	cfg := provider.LoadGithubOAuthConfig()

	tok, err := cfg.Exchange(ctx, code)
	if err != nil {
		return dto.UserJwtPackage{}, err
	}

	client := cfg.Client(ctx, tok)
	res, err := client.Get("https://api.github.com/user")
	if err != nil {
		return dto.UserJwtPackage{}, err
	}
	defer res.Body.Close()

	var user githubprovider.GithubUser
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&user); err != nil {
		return dto.UserJwtPackage{}, err
	}

	claims := dto.UserJwtPackage{
		Username:  user.Login,
		UserID:    user.ID,
		AvatarURL: user.AvatarURL,
		OAuthKey:  fmt.Sprintf("github_%d", user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    s.conf.Issuer,
		},
	}

	return claims, nil
}

func (s *userService) CreateJWT(claims dto.UserJwtPackage) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.conf.JwtSecret))
}
