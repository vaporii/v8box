package repository

import (
	"database/sql"

	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/httperror"
	"github.com/vaporii/v8box/internal/models"

	_ "modernc.org/sqlite"
)

type NoteRepository interface {
	CreateNote(user *models.Note) error
	GetNoteByID(id string) (*models.Note, error)
	GetUserNotes(userId string) ([]models.Note, error)
	UpdateNote(id string, request dto.CreateNoteRequest) error
}

type noteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) (NoteRepository, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id				VARCHAR(255) PRIMARY KEY,
			user_id			VARCHAR(255) NOT NULL,
			title			VARCHAR(255) NOT NULL,
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

func (r *noteRepository) GetNoteByID(id string) (*models.Note, error) {
	note := &models.Note{}
	err := r.db.QueryRow("SELECT id, user_id, title, content FROM notes WHERE id=?", id).Scan(&note.ID, &note.UserID, &note.Title, &note.Content)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (r *noteRepository) GetUserNotes(userId string) ([]models.Note, error) {
	var userCount int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE id=?", userId).Scan(&userCount)
	if err != nil {
		return nil, err
	}

	if userCount == 0 {
		return nil, &httperror.NotFoundError{Entity: "User"}
	}

	rows, err := r.db.Query("SELECT id, user_id, title, content FROM notes WHERE user_id=?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note = make([]models.Note, 0)

	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content); err != nil {
			return notes, err
		}
		notes = append(notes, note)
	}
	if err = rows.Err(); err != nil {
		return notes, err
	}
	return notes, nil
}

func (r *noteRepository) UpdateNote(id string, request dto.CreateNoteRequest) error {
	_, err := r.db.Exec("UPDATE notes SET title=?, content=? WHERE id=?", request.Title, request.Content, id)
	if err != nil {
		return err
	}
	return nil
}
