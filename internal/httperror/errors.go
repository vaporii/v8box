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

type BadClientRequestError struct {
	Message string
}

func (e *BadClientRequestError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}
