package repository

import (
	"database/sql"

	"github.com/vaporii/v8box/internal/models"

	_ "modernc.org/sqlite"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (UserRepository, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id             	VARCHAR(255) PRIMARY KEY,
			username       	VARCHAR(50) NOT NULL UNIQUE,
			password_hash  	TEXT NOT NULL,
			oauth_provider 	VARCHAR(255),
			oauth_id		TEXT,
			access_token	TEXT,
			refresh_token	TEXT,
			token_expiry	INTEGER
		);
	`)
	if err != nil {
		return nil, err
	}

	return &userRepository{
		db: db,
	}, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (
			id,
			username,
			password_hash,
			oauth_provider,
			oauth_id,
			access_token,
			refresh_token,
			token_expiry
		) VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)
	`,
		user.ID,
		user.Username,
		user.Password,
		user.OAuthProvider,
		user.OAuthID,
		user.AccessToken,
		user.RefreshToken,
		user.TokenExpiry,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT
			id,
			username,
			password,
			oauth_provider,
			oauth_id,
			access_token,
			refresh_token,
			token_expiry
		FROM users WHERE username=?
	`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.OAuthProvider,
		&user.OAuthID,
		&user.AccessToken,
		&user.RefreshToken,
		&user.TokenExpiry,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
