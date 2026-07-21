package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/internal/application/ports"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/internal/usecases"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/response"
)

type TodoHandler struct {
	todoService usecases.TodoService
	validator   *validator.Validate
}

func NewTodoHandler(todoService usecases.TodoService, validator *validator.Validate) *TodoHandler {
	return &TodoHandler{todoService: todoService, validator: validator}
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	user, err := authenticatedUser(c)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	request := new(dtos.CreateTodoDto)
	if err := c.Bind(request); err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid JSON body"))
	}
	if err := h.validator.Struct(request); err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid field"))
	}

	todo, err := h.todoService.Create(c.Request().Context(), user.ID, request)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	return response.NewCreatedAPIResponse(c, "create todo successfully", todo)
}

func (h *TodoHandler) ListTodos(c echo.Context) error {
	user, err := authenticatedUser(c)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	filter, err := parseTodoListFilter(c)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	page, err := h.todoService.List(c.Request().Context(), user.ID, filter)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	data := todoListResponse{
		Items: page.Items,
		Pagination: paginationResponse{
			Page:       page.Page,
			Limit:      page.Limit,
			Total:      page.Total,
			TotalPages: page.TotalPages,
		},
	}
	return response.NewSuccessAPIResponse(c, "get todos successfully", data)
}

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	user, err := authenticatedUser(c)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid todo id"))
	}

	request := new(dtos.UpdateTodoDto)
	if err := c.Bind(request); err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid JSON body"))
	}
	if err := h.validator.Struct(request); err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid field"))
	}

	todo, err := h.todoService.Update(c.Request().Context(), user.ID, todoID, request)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	return response.NewSuccessAPIResponse(c, "update todo successfully", todo)
}

func (h *TodoHandler) UpdateTodoState(c echo.Context) error {
	user, err := authenticatedUser(c)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	stateID, err := uuid.Parse(c.Param("stateId"))
	if err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid todo state id"))
	}

	request := new(dtos.UpdateTodoState)
	if err := c.Bind(request); err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid JSON body"))
	}
	if err := h.validator.Struct(request); err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid field"))
	}

	state, err := h.todoService.UpdateTodoState(c.Request().Context(), user.ID, stateID, request)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	return response.NewSuccessAPIResponse(c, "update todo state successfully", state)
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	user, err := authenticatedUser(c)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid todo id"))
	}

	if err := h.todoService.Delete(c.Request().Context(), user.ID, todoID); err != nil {
		return response.NewErrorResponse(c, err)
	}
	return response.NewSuccessAPIResponse(c, fmt.Sprintf("delete todo with id: %v successfully", todoID), nil)
}

func (h *TodoHandler) GetTodosByUser(c echo.Context) error {
	user, err := authenticatedUser(c)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	requestedUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid user id"))
	}
	if requestedUserID != user.ID {
		return response.NewErrorResponse(c, errs.NewForbiddenError("forbidden"))
	}

	todos, err := h.todoService.FindByUserID(c.Request().Context(), user.ID)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	return response.NewSuccessAPIResponse(c, "get todos successfully", todos)
}

type todoListResponse struct {
	Items      []domain.Todo      `json:"items"`
	Pagination paginationResponse `json:"pagination"`
}

type paginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

func authenticatedUser(c echo.Context) (*domain.User, error) {
	user, ok := c.Get("user").(*domain.User)
	if !ok {
		return nil, errs.NewUnauthorizedError("unauthorized")
	}
	return user, nil
}

func parseTodoListFilter(c echo.Context) (ports.TodoListFilter, error) {
	filter := ports.TodoListFilter{Page: 1, Limit: 20, Sort: "created_at", Order: "desc"}
	var err error
	if value := c.QueryParam("page"); value != "" {
		filter.Page, err = strconv.Atoi(value)
		if err != nil || filter.Page < 1 {
			return filter, errs.NewBadRequestError("page must be greater than zero")
		}
	}
	if value := c.QueryParam("limit"); value != "" {
		filter.Limit, err = strconv.Atoi(value)
		if err != nil || filter.Limit < 1 || filter.Limit > 100 {
			return filter, errs.NewBadRequestError("limit must be between 1 and 100")
		}
	}

	filter.Search = strings.TrimSpace(c.QueryParam("search"))
	if len(filter.Search) > 100 {
		return filter, errs.NewBadRequestError("search must not exceed 100 characters")
	}
	if value := c.QueryParam("priority"); value != "" {
		filter.Priority = domain.Priority(value)
		if filter.Priority != domain.High && filter.Priority != domain.Medium && filter.Priority != domain.Low {
			return filter, errs.NewBadRequestError("invalid priority")
		}
	}
	if value := c.QueryParam("state_id"); value != "" {
		stateID, parseErr := uuid.Parse(value)
		if parseErr != nil {
			return filter, errs.NewBadRequestError("invalid state id")
		}
		filter.StateID = &stateID
	}
	if filter.DueFrom, err = parseOptionalTime(c.QueryParam("due_from")); err != nil {
		return filter, errs.NewBadRequestError("invalid due_from")
	}
	if filter.DueTo, err = parseOptionalTime(c.QueryParam("due_to")); err != nil {
		return filter, errs.NewBadRequestError("invalid due_to")
	}
	if filter.DueFrom != nil && filter.DueTo != nil && filter.DueFrom.After(*filter.DueTo) {
		return filter, errs.NewBadRequestError("due_from must not be after due_to")
	}

	if value := c.QueryParam("sort"); value != "" {
		filter.Sort = value
	}
	if _, ok := allowedTodoSorts[filter.Sort]; !ok {
		return filter, errs.NewBadRequestError("invalid sort field")
	}
	if value := c.QueryParam("order"); value != "" {
		filter.Order = strings.ToLower(value)
	}
	if filter.Order != "asc" && filter.Order != "desc" {
		return filter, errs.NewBadRequestError("invalid sort order")
	}

	return filter, nil
}

var allowedTodoSorts = map[string]struct{}{
	"created_at": {}, "updated_at": {}, "due_date": {}, "priority": {}, "title": {},
}

func parseOptionalTime(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
