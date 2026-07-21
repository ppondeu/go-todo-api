package domain

import "errors"

var (
	ErrTodoNotFound      = errors.New("todo not found")
	ErrTodoStateNotFound = errors.New("todo state not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailConflict     = errors.New("email already exists")
)
