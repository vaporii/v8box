package repository

import (
	"database/sql"

	"github.com/vaporii/v8box/internal/config"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/httperror"
	"github.com/vaporii/v8box/internal/logging"
	"github.com/vaporii/v8box/internal/models"

	_ "modernc.org/sqlite"
)

type NoteRepository interface {
	CreateNote(note *models.Note) (*models.Note, error)
	GetNoteByID(id string) (*models.Note, error)
	GetUserNotes(userId string) ([]models.Note, error)
	UpdateNote(id string, request dto.CreateNoteRequest) (*models.Note, error)
}

type noteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) (NoteRepository, error) {
	if config.LoadConfig().Environment == "dev" {
		logging.Info("dev environment, deleting notes table")
		db.Exec(`
			DROP TABLE IF EXISTS notes;
		`)
	}
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id				VARCHAR(255) PRIMARY KEY,
			user_id			VARCHAR(255) NOT NULL,
			title			VARCHAR(255) NOT NULL,
			content			TEXT,
			created_at		TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at		TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);

		DROP TRIGGER IF EXISTS update_notes_updated_at;
		
		CREATE TRIGGER update_notes_updated_at
		AFTER UPDATE ON notes
		FOR EACH ROW
		BEGIN
			UPDATE notes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;
	`)
	if err != nil {
		return nil, err
	}

	return &noteRepository{
		db: db,
	}, nil
}

func (r *noteRepository) CreateNote(note *models.Note) (*models.Note, error) {
	var retNote models.Note
	err := r.db.QueryRow(`
		INSERT INTO notes (
			id, user_id, title, content
		) VALUES (?, ?, ?, ?) RETURNING
			id, user_id, title, content, created_at, updated_at;
	`, note.ID, note.UserID, note.Title, note.Content).Scan(&retNote.ID, &retNote.UserID, &retNote.Title, &retNote.Content, &retNote.CreatedAt, &retNote.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &retNote, nil
}

func (r *noteRepository) GetNoteByID(id string) (*models.Note, error) {
	note := &models.Note{}
	err := r.db.QueryRow("SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id=?", id).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
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

	rows, err := r.db.Query("SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE user_id=?", userId)
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

func (r *noteRepository) UpdateNote(id string, request dto.CreateNoteRequest) (*models.Note, error) {
	var note models.Note

	err := r.db.QueryRow(`
		UPDATE notes
		SET title=?,
			content=?
		WHERE id=?
		RETURNING
			id, user_id, title, content, created_at, updated_at;
	`, request.Title, request.Content, id).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &note, nil
}
