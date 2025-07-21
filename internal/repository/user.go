package repository

import (
	"database/sql"

	"github.com/vaporii/v8box/internal/logging"
	"github.com/vaporii/v8box/internal/models"

	_ "modernc.org/sqlite"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByOAuthKey(oauthKey string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (UserRepository, error) {
	logging.Info("creating users table")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id             	VARCHAR(255) PRIMARY KEY,
			username       	VARCHAR(50) NOT NULL UNIQUE,
			password_hash  	TEXT NOT NULL,
			oauth_key		TEXT UNIQUE
		);
	`)
	if err != nil {
		return nil, err
	}
	logging.Verbose("created users table")

	return &userRepository{
		db: db,
	}, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	logging.Verbose("creating user")
	_, err := r.db.Exec(`
		INSERT INTO users (
			id,
			username,
			password_hash,
			oauth_key
		) VALUES (
			?,
			?,
			?,
			?
		)
	`,
		user.ID,
		user.Username,
		user.Password,
		user.OAuthKey,
	)
	if err != nil {
		return err
	}
	logging.Verbose("created user")

	return nil
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	logging.Verbose("getting user by username")
	err := r.db.QueryRow(`
		SELECT
			id,
			username,
			password_hash,
			oauth_key
		FROM users WHERE username=?
	`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.OAuthKey,
	)
	if err != nil {
		return nil, err
	}
	logging.Verbose("got user by username")

	return user, nil
}

func (r *userRepository) GetUserByOAuthKey(oauthKey string) (*models.User, error) {
	user := &models.User{}
	logging.Verbose("getting user by oauth key")
	err := r.db.QueryRow(`
		SELECT
			id,
			username,
			password_hash,
			oauth_key
		FROM users WHERE oauth_key=?
	`, oauthKey).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.OAuthKey,
	)
	if err != nil {
		return nil, err
	}
	logging.Verbose("got user by oauth key")

	return user, nil
}
