package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/internal/usecase"
	"github.com/ppondeu/go-todo-api/pkg/dto"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/logs"
	"github.com/ppondeu/go-todo-api/pkg/response"
)

type TodoHandler struct {
	todoService usecase.TodoService
	validator   *validator.Validate
}

func NewTodoHandler(todoService *usecase.TodoService, validator *validator.Validate) *TodoHandler {
	return &TodoHandler{
		todoService: *todoService,
		validator:   validator,
	}
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid user id"))
	}

	createTodoRequest := new(dto.CreateTodoDto)
	err = c.Bind(createTodoRequest)
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("JSON required"))
	}

	err = h.validator.Struct(createTodoRequest)
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid field"))
	}

	logs.Info(createTodoRequest)

	todoRes, err := h.todoService.Create(userID, createTodoRequest)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewCreatedAPIResponse(c, "create todo successfully", todoRes)

}

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid todo id"))
	}

	updateTodoRequest := new(dto.UpdateTodoDto)
	err = c.Bind(updateTodoRequest)
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("JSON required"))
	}

	err = h.validator.Struct(updateTodoRequest)
	if err != nil {
		logs.Error(err)
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return response.NewErrorResponse(c, errs.NewBadRequestError("invalid field"))
		}

		return response.NewErrorResponse(c, errs.NewInternalError("something went wrong while validate body"))
	}

	todoRes, err := h.todoService.Update(todoID, updateTodoRequest)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "update todo successfully", todoRes)
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid todo id"))
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
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid user id"))
	}

	todosRes, err := h.todoService.FindByUserID(userID)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "get users successfully", todosRes)
}

func (h *TodoHandler) CreateCategory(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid user id"))
	}

	createCategoryRequest := new(dto.CreateCategoryDto)
	err = c.Bind(createCategoryRequest)
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid JSON body"))
	}

	err = h.validator.Struct(createCategoryRequest)
	if err != nil {
		logs.Error(err)
		if _, ok := err.(*validator.ValidationErrors); ok {
			return response.NewErrorResponse(c, errs.NewBadRequestError("invalid field"))
		}

		return response.NewErrorResponse(c, errs.NewBadRequestError("something went wrong while validate body"))
	}

	res, err := h.todoService.CreateCategory(userID, createCategoryRequest.Name)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "create category successfully", res)
}

func (h *TodoHandler) GetCategoriesByUser(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid user id"))
	}

	res, err := h.todoService.FindCategories(userID)
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "get categories successfully", res)
}

func (h *TodoHandler) UpdateCategory(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid category id"))
	}

	updateCategoryRequest := new(dto.UpdateCategoryDto)
	err = c.Bind(updateCategoryRequest)
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid JSON"))
	}

	err = h.validator.Struct(updateCategoryRequest)
	if err != nil {
		logs.Error(err)
		if _, ok := err.(*validator.ValidationErrors); ok {
			return response.NewErrorResponse(c, errs.NewBadRequestError("invalid field"))
		}
		return response.NewErrorResponse(c, errs.NewBadRequestError("something went wrong while validate body"))
	}

	res, err := h.todoService.UpdateCategory(categoryID, updateCategoryRequest.Name)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	return response.NewSuccessAPIResponse(c, "update category successfully", res)
}

func (h *TodoHandler) DeleteCategory(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid category id"))
	}

	err = h.todoService.Delete(categoryID)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	return response.NewSuccessAPIResponse(c, "delete category successfully", err)
}

func (h *TodoHandler) GetTodosByCategory(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid category id"))
	}

	res, err := h.todoService.FindTodosByCategory(categoryID)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	return response.NewSuccessAPIResponse(c, "get todos by category successful", res)
}
