package repository

import (
	"database/sql"

	"github.com/vaporii/v8box/internal/models"

	_ "modernc.org/sqlite"
)

type NoteRepository interface {
	CreateNote(user *models.Note) error
	GetNoteByID(string int) (*models.Note, error)
}

type noteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) (NoteRepository, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id				VARCHAR(255) PRIMARY KEY,
			user_id			VARCHAR(255),
			title			VARCHAR(255),
			content			TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
	`)
	if err != nil {
		return nil, err
	}

	return &noteRepository{
		db: db,
	}, nil
}

func (r *noteRepository) CreateNote(note *models.Note) error {
	_, err := r.db.Exec("INSERT INTO notes (id, user_id, title, content) VALUES (?, ?, ?, ?)", note.ID, note.UserID, note.Title, note.Content)
	if err != nil {
		return err
	}

	return nil
}

func (r *noteRepository) GetNoteByID(id int) (*models.Note, error) {
	note := &models.Note{}
	err := r.db.QueryRow("SELECT id, user_id, title, content FROM notes WHERE id=?", id).Scan(&note.ID, &note.UserID, &note.Title, &note.Content)
	if err != nil {
		return nil, err
	}

	return note, nil
}
