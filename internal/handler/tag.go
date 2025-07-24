package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vaporii/v8box/internal/dto"
	"github.com/vaporii/v8box/internal/models"
	"github.com/vaporii/v8box/internal/service"
)

type TagHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	EditTag(w http.ResponseWriter, r *http.Request)
	GetTag(w http.ResponseWriter, r *http.Request)
	GetTagsOnNote(w http.ResponseWriter, r *http.Request)
	GetNotesWithTag(w http.ResponseWriter, r *http.Request)
}

type tagHandler struct {
	tagService service.TagService
}

func NewTagHandler(tagService service.TagService) TagHandler {
	return &tagHandler{
		tagService: tagService,
	}
}

func (h *tagHandler) Create(w http.ResponseWriter, r *http.Request) {
	tagRequest := dto.Tag{}
	json.NewDecoder(r.Body).Decode(&tagRequest)

	tag, err := h.tagService.CreateTag(models.ExtractUser(r).UserID, tagRequest)
	if checkErr(err, r) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(tag)
	if checkErr(err, r) {
		return
	}
}

func (h *tagHandler) EditTag(w http.ResponseWriter, r *http.Request) {
	tagRequest := dto.Tag{}
	json.NewDecoder(r.Body).Decode(&tagRequest)

	tag, err := h.tagService.EditTag(chi.URLParam(r, "id"), tagRequest)
	if checkErr(err, r) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(tag)
	if checkErr(err, r) {
		return
	}
}

func (h *tagHandler) GetTag(w http.ResponseWriter, r *http.Request) {
	tag, err := h.tagService.GetTag(chi.URLParam(r, "id"))
	if checkErr(err, r) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(tag)
	if checkErr(err, r) {
		return
	}
}

func (h *tagHandler) GetTagsOnNote(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "note_id")
	tags, err := h.tagService.GetTagsOnNote(noteID)
	if checkErr(err, r) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(tags)
	if checkErr(err, r) {
		return
	}
}

func (h *tagHandler) GetNotesWithTag(w http.ResponseWriter, r *http.Request) {
	tagID := chi.URLParam(r, "tag_id")
	notes, err := h.tagService.GetNotesWithTag(tagID)
	if checkErr(err, r) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(notes)
	if checkErr(err, r) {
		return
	}
}
