package errs

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
	return NewAppError(404, message)
}

func NewInternalError(message string) *AppError {
	return NewAppError(500, message)
}

func NewBadRequestError(message string) *AppError {
	return NewAppError(400, message)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(401, message)
}

func NewForbiddenError(message string) *AppError {
	return NewAppError(403, message)
}
