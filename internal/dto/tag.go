package dto

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
