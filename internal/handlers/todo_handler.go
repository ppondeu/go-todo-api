package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/internal/usecases"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/logs"
	"github.com/ppondeu/go-todo-api/pkg/response"
)

type TodoHandler struct {
	todoService usecases.TodoService
	validator   *validator.Validate
}

func NewTodoHandler(todoService *usecases.TodoService, validator *validator.Validate) *TodoHandler {
	return &TodoHandler{
		todoService: *todoService,
		validator:   validator,
	}
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	user, ok := c.Get("user").(*domain.User)
	if !ok {
		logs.Error("invalid get req.user")
		return response.NewErrorResponse(c, errs.NewUnauthorizedError("Unauthorized"))
	}

	createTodoRequest := new(dtos.CreateTodoDto)
	err := c.Bind(createTodoRequest)
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("JSON required"))
	}

	err = h.validator.Struct(createTodoRequest)
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid field"))
	}

	logs.Info(createTodoRequest)

	todoRes, err := h.todoService.Create(user.ID, createTodoRequest)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewCreatedAPIResponse(c, "create todo successfully", todoRes)

}

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid todo id"))
	}

	updateTodoRequest := new(dtos.UpdateTodoDto)
	err = c.Bind(updateTodoRequest)
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("JSON required"))
	}

	err = h.validator.Struct(updateTodoRequest)
	if err != nil {
		logs.Error(err)
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid field"))
		}

		return response.NewErrorResponse(c, errs.NewInternalError("something went wrong while validate body"))
	}

	todoRes, err := h.todoService.Update(todoID, updateTodoRequest)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "update todo successfully", todoRes)
}

func (h *TodoHandler) UpdateTodoState(c echo.Context) error {
	todoStateID, err := uuid.Parse(c.Param("stateId"))
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid todo id"))
	}

	updateTodoStateDTO := new(dtos.UpdateTodoState)
	err = c.Bind(updateTodoStateDTO)
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("JSON required"))
	}

	err = h.validator.Struct(updateTodoStateDTO)
	if err != nil {
		logs.Error(err)
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid field"))
		}

		return response.NewErrorResponse(c, errs.NewInternalError("something went wrong while validate body"))
	}

	todoRes, err := h.todoService.UpdateTodoState(todoStateID, updateTodoStateDTO)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "update todo state successfully", todoRes)
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid todo id"))
	}

	err = h.todoService.Delete(todoID)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, fmt.Sprintf("delete todo with id: %v successfully", todoID), nil)
}

func (h *TodoHandler) GetTodosByUser(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid user id"))
	}

	todosRes, err := h.todoService.FindByUserID(userID)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "get users successfully", todosRes)
}
