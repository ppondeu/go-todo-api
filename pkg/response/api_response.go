package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/pkg/errs"
)

type APIResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func NewAPIResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	response := APIResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	return c.JSON(statusCode, response)
}

func NewErrorResponse(c echo.Context, err error) error {
	switch e := err.(type) {
	case *errs.AppError:
		return NewAPIResponse(c, e.Code, e.Message, nil)
	default:
		return NewAPIResponse(c, http.StatusInternalServerError, "An unexpected error occurred", nil)
	}
}

func NewSuccessAPIResponse(c echo.Context, message string, data interface{}) error {
	return NewAPIResponse(c, http.StatusOK, message, data)
}

func NewCreatedAPIResponse(c echo.Context, message string, data interface{}) error {
	return NewAPIResponse(c, http.StatusCreated, message, data)
}
