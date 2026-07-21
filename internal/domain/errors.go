package domain

import "errors"

var (
	ErrTodoNotFound      = errors.New("todo not found")
	ErrTodoStateNotFound = errors.New("todo state not found")
)
