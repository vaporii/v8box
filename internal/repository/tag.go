package repository

import (
	"database/sql"
	"errors"

	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/httperror"
	"github.com/vaporii/v8box/internal/logging"
	"github.com/vaporii/v8box/internal/models"
)

type TagRepository interface {
	CreateTag(tagId string, userId string, tag dto.Tag) (*models.Tag, error)
	GetTagByName(tagName string) (*models.Tag, error)
	GetTagById(tagId string) (*models.Tag, error)
	GetTagsByUser(userId string) ([]models.Tag, error)
	UpdateTag(tagId string, tag *dto.Tag) (*models.Tag, error)
	DeleteTag(tagId string) error
	GetTagsOnNote(noteId string) ([]models.Tag, error)
	GetNotesWithTag(tagId string) ([]models.Note, error)
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
			user_id		VARCHAR(255),
			name		TEXT,
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

func (r *tagRepository) CreateTag(tagId string, userId string, tag dto.Tag) (*models.Tag, error) {
	var retTag models.Tag

	err := r.db.QueryRow(`
		INSERT INTO tags (
			id, user_id, name, description
		) VALUES (?, ?, ?, ?) RETURNING
			id, user_id, name, description, created_at, updated_at;
	`, tagId, userId, tag.Name, tag.Description).Scan(&retTag.ID, &retTag.UserID, &retTag.Name, &retTag.Description, &retTag.CreatedAt, &retTag.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &retTag, nil
}

func (r *tagRepository) GetTagByName(tagName string) (*models.Tag, error) {
	var retTag models.Tag

	err := r.db.QueryRow(`
		SELECT
			id, user_id, name, description, created_at, updated_at
		FROM tags
		WHERE name=?;
	`, tagName).Scan(&retTag.ID, &retTag.UserID, &retTag.Name, &retTag.Description, &retTag.CreatedAt, &retTag.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &httperror.NotFoundError{Entity: "Tag"}
		}

		return nil, err
	}

	return &retTag, nil
}

func (r *tagRepository) GetTagById(tagId string) (*models.Tag, error) {
	var retTag models.Tag

	err := r.db.QueryRow(`
		SELECT
			id, user_id, name, description, created_at, updated_at
		FROM tags
		WHERE id=?;
	`, tagId).Scan(&retTag.ID, &retTag.UserID, &retTag.Name, &retTag.Description, &retTag.CreatedAt, &retTag.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &httperror.NotFoundError{Entity: "Tag"}
		}

		return nil, err
	}

	return &retTag, nil
}

func (r *tagRepository) GetTagsByUser(userId string) ([]models.Tag, error) {
	var userCount int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE id=?", userId).Scan(&userCount)
	if err != nil {
		return nil, err
	}

	if userCount == 0 {
		return nil, &httperror.NotFoundError{Entity: "User"}
	}

	rows, err := r.db.Query("SELECT id, user_id, name, description, created_at, updated_at FROM tags WHERE user_id=?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag = make([]models.Tag, 0)

	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.Description, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}
	if err = rows.Err(); err != nil {
		return tags, err
	}
	return tags, nil
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
			id, user_id, name, description, created_at, updated_at;
	`).Scan(&retTag.ID, &retTag.UserID, &retTag.Name, &retTag.Description, &retTag.CreatedAt, &retTag.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &httperror.NotFoundError{Entity: "Tag"}
		}

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

func (r *tagRepository) GetTagsOnNote(noteId string) ([]models.Tag, error) {
	var noteCount int
	err := r.db.QueryRow("SELECT COUNT(*) FROM notes WHERE id=?", noteId).Scan(&noteCount)
	if err != nil {
		return nil, err
	}

	if noteCount == 0 {
		return nil, &httperror.NotFoundError{Entity: "Note"}
	}

	rows, err := r.db.Query(`
		SELECT
			t.id, t.user_id, t.name, t.description, t.created_at, t.updated_at
		FROM tags t
		JOIN note_tags nt ON t.id = nt.tag_id
		WHERE nt.note_id=?;
	`, noteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag = make([]models.Tag, 0)

	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.Description, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}
	if err = rows.Err(); err != nil {
		return tags, err
	}
	return tags, nil
}

func (r *tagRepository) GetNotesWithTag(tagId string) ([]models.Note, error) {
	var tagCount int
	err := r.db.QueryRow("SELECT COUNT(*) FROM tags WHERE id=?", tagId).Scan(&tagCount)
	if err != nil {
		return nil, err
	}

	if tagCount == 0 {
		return nil, &httperror.NotFoundError{Entity: "Tag"}
	}

	rows, err := r.db.Query(`
		SELECT
			n.id, n.user_id, n.title, n.content, n.created_at, n.updated_at
		FROM notes n
		JOIN note_tags nt ON n.id = nt.note_id
		WHERE nt.tag_id=?;
	`, tagId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note = make([]models.Note, 0)

	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return notes, err
		}
		notes = append(notes, note)
	}
	if err = rows.Err(); err != nil {
		return notes, err
	}
	return notes, nil
}
