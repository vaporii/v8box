package httperror

import "fmt"

type contextKey string

const ErrorKey contextKey = "appError"

type NotFoundError struct {
	Entity string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Entity)
}
