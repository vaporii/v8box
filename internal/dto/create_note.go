package dto

type CreateNoteRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Content string `json:"content"`
}
