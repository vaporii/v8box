package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/httperror"
	"github.com/vaporii/v8box/internal/models"
	"github.com/vaporii/v8box/internal/service"
)

type NoteHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetNotes(w http.ResponseWriter, r *http.Request)
	GetNoteByID(w http.ResponseWriter, r *http.Request)
}

type noteHandler struct {
	noteService service.NoteService
}

func NewNoteHandler(noteService service.NoteService) NoteHandler {
	return &noteHandler{
		noteService: noteService,
	}
}

func (h *noteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var noteRequest dto.CreateNoteRequest
	err := json.NewDecoder(r.Body).Decode(&noteRequest)
	if err != nil {
		err = &httperror.BadClientRequestError{Message: "Bad JSON request"}
	}
	if checkErr(err, r) {
		return
	}

	noteRequest.UserID = models.ExtractUser(r).UserID

	note, err := h.noteService.Create(noteRequest)
	if checkErr(err, r) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(*note)
}

func (h *noteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.noteService.GetUserNotes(models.ExtractUser(r).UserID)
	if checkErr(err, r) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(notes)
	if checkErr(err, r) {
		return
	}
}

func (h *noteHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	note, err := h.noteService.GetNoteByID(chi.URLParam(r, "id"))
	if checkErr(err, r) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(note)
	if checkErr(err, r) {
		return
	}
}
