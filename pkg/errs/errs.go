package errs

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(http.StatusNotFound, message)
}

func NewInternalError(message string) *AppError {
	return NewAppError(http.StatusInternalServerError, message)
}

func NewBadRequestError(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message)
}

func NewForbiddenError(message string) *AppError {
	return NewAppError(http.StatusForbidden, message)
}

func NewConflictError(message string) *AppError {
	return NewAppError(http.StatusConflict, message)
}
