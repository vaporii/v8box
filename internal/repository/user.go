package repository

import (
	"database/sql"

	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/logging"
	"github.com/vaporii/v8box/internal/models"

	_ "modernc.org/sqlite"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByOAuthKey(oauthKey string) (*models.User, error)
	GetUserById(userId string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (UserRepository, error) {
	if config.LoadConfig().Environment == "dev" {
		logging.Info("dev environment, deleting table")
		db.Exec(`
			DROP TABLE IF EXISTS users;
		`)
	}
	logging.Info("creating users table")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id				VARCHAR(255) PRIMARY KEY,
			username		TEXT NOT NULL,
			password_hash	TEXT NOT NULL,
			oauth_key		TEXT UNIQUE,
			avatar_url		TEXT,
			created_at		TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at		TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		DROP TRIGGER IF EXISTS update_users_updated_at;
		
		CREATE TRIGGER update_users_updated_at
		AFTER UPDATE ON users
		FOR EACH ROW
		BEGIN
			UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;
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
			oauth_key,
			avatar_url
		) VALUES (
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
		user.OAuthKey,
		user.AvatarURL,
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
			oauth_key,
			avatar_url,
			created_at,
			updated_at
		FROM users WHERE username=?
	`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.OAuthKey,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
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
			oauth_key,
			avatar_url,
			created_at,
			updated_at
		FROM users WHERE oauth_key=?
	`, oauthKey).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.OAuthKey,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	logging.Verbose("got user by oauth key")

	return user, nil
}

func (r *userRepository) GetUserById(userId string) (*models.User, error) {
	user := &models.User{}
	logging.Verbose("getting user by id")
	err := r.db.QueryRow(`
		SELECT
			id,
			username,
			password_hash,
			oauth_key,
			avatar_url,
			created_at,
			updated_at
		FROM users WHERE id=?
	`, userId).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.OAuthKey,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	logging.Verbose("got user by id")

	return user, nil
}
