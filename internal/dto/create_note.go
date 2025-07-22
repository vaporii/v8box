package dto

type CreateNoteRequest struct {
	Title   string `json:"username" validate:"required,min=1,max=255"`
	UserID  string `json:"-"`
	Content string `json:"content"`
}
