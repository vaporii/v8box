package repository

import (
	"database/sql"

	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/logging"
	"github.com/vaporii/v8box/internal/models"
)

type TagRepository interface {
	CreateTag(tagId string, tag *dto.Tag) (*models.Tag, error)
	GetTagByName(tagName string) (*models.Tag, error)
	GetTagById(tagId string) (*models.Tag, error)
	UpdateTag(tagId string, tag *dto.Tag) (*models.Tag, error)
	DeleteTag(tagId string) error
}

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) (TagRepository, error) {
	if config.LoadConfig().Environment == "dev" {
		logging.Info("dev environment, deleting tags table")
		db.Exec(`
			DROP TABLE IF EXISTS tags;
			DROP TABLE IF EXISTS note_tags;
		`)
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			id			VARCHAR(255) PRIMARY KEY,
			name		TEXT NOT NULL UNIQUE,
			description	TEXT,
			created_at	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at	TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		DROP TRIGGER IF EXISTS update_tags_updated_at;
		
		CREATE TRIGGER update_tags_updated_at
		AFTER UPDATE ON tags
		FOR EACH ROW
		BEGIN
			UPDATE tags SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;

		CREATE TABLE IF NOT EXISTS note_tags (
			note_id		VARCHAR(255),
			tag_id		VARCHAR(255),
			PRIMARY KEY (note_id, tag_id),
			FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return nil, err
	}

	return &tagRepository{
		db: db,
	}, nil
}

func (r *tagRepository) CreateTag(tagId string, tag *dto.Tag) (*models.Tag, error) {
	var retTag models.Tag

	err := r.db.QueryRow(`
		INSERT INTO tags (
			id, name, description
		) VALUES (?, ?, ?) RETURNING
			id, name, description, created_at, updated_at;
	`, tagId, tag.Name, tag.Description).Scan(&retTag.ID, &retTag.Name, &retTag.Description, &retTag.CreatedAt, &retTag.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &retTag, nil
}

func (r *tagRepository) GetTagByName(tagName string) (*models.Tag, error) {
	var retTag models.Tag

	err := r.db.QueryRow(`
		SELECT
			id, name, description, created_at, updated_at
		FROM tags
		WHERE name=?;
	`, tagName).Scan(&retTag.ID, &retTag.Name, &retTag.Description, &retTag.CreatedAt, &retTag.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &retTag, nil
}

func (r *tagRepository) GetTagById(tagId string) (*models.Tag, error) {
	var retTag models.Tag

	err := r.db.QueryRow(`
		SELECT
			id, name, description, created_at, updated_at
		FROM tags
		WHERE id=?;
	`, tagId).Scan(&retTag.ID, &retTag.Name, &retTag.Description, &retTag.CreatedAt, &retTag.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &retTag, nil
}

func (r *tagRepository) UpdateTag(tagId string, tag *dto.Tag) (*models.Tag, error) {
	var retTag models.Tag

	err := r.db.QueryRow(`
		UPDATE tags
		SET
			name=?,
			description=?
		WHERE id=?
		RETURNING
			id, name, description, created_at, updated_at;
	`).Scan(&retTag.ID, &retTag.Name, &retTag.Description, &retTag.CreatedAt, &retTag.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &retTag, nil
}

func (r *tagRepository) DeleteTag(tagId string) error {
	_, err := r.db.Exec(`
		DELETE FROM tags
		WHERE id=?
	`, tagId)
	return err
}
