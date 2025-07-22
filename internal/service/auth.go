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

type AuthService interface {
	Register(request dto.RegisterRequest) (*models.User, error)
	GetGitHubOAuthJwt(ctx context.Context, code string) (dto.UserJwtPackage, error)
	CreateJWT(claims dto.UserJwtPackage) (string, error)
	// RegisterOAuthUser(jwt dto.UserJwtPackage) error
}

type authService struct {
	userRepo repository.UserRepository
	conf     config.Config
}

func NewAuthService(userRepo repository.UserRepository, conf config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		conf:     conf,
	}
}

func (r *authService) Register(request dto.RegisterRequest) (*models.User, error) {
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

func (r *authService) GetGitHubOAuthJwt(ctx context.Context, code string) (dto.UserJwtPackage, error) {
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

	dbUser, err := r.userRepo.GetUserByOAuthKey(fmt.Sprintf("github_%d", user.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user := &models.User{
				ID:        uuid.NewString(),
				Username:  user.Login,
				OAuthKey:  fmt.Sprintf("github_%d", user.ID),
				AvatarURL: user.AvatarURL,
			}

			err = r.userRepo.CreateUser(user)
			if err != nil {
				return dto.UserJwtPackage{}, err
			}
			dbUser = user
		} else {
			return dto.UserJwtPackage{}, err
		}
	}

	claims := dto.UserJwtPackage{
		Username:  dbUser.Username,
		UserID:    dbUser.ID,
		AvatarURL: dbUser.AvatarURL,
		OAuthKey:  fmt.Sprintf("github_%d", user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    r.conf.Issuer,
		},
	}

	return claims, nil
}

func (r *authService) CreateJWT(claims dto.UserJwtPackage) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(r.conf.JwtSecret))
}
