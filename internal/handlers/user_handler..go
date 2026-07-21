package handlers

import (
	"net/http"

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

type UserHandler struct {
	userService usecases.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService usecases.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{userService: userService, validator: validator}
}

func (h *UserHandler) Register(c echo.Context) error {
	createUserRequest := new(dtos.UserCreateDTO)
	if err := c.Bind(createUserRequest); err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid user id"))
	}

	if err := h.validator.Struct(createUserRequest); err != nil {
		logs.Error(err)
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid field"))
		}

		return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid field"))
	}

	user, err := h.userService.Save(createUserRequest)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewCreatedAPIResponse(c, "register successfully", user)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.userService.FindAll()
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "get user successfully", users)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid user id"))
	}

	user, err := h.userService.FindByUserID(userID)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "get user successfully", user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	user, ok := c.Get("user").(*domain.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not found in context")
	}

	updateUserRequest := new(dtos.UserUpdateDTO)
	if err := c.Bind(updateUserRequest); err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid JSON body"))
	}

	if err := h.validator.Struct(updateUserRequest); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid field"))
		}
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("validation error"))
	}

	res, err := h.userService.Update(user.ID, updateUserRequest)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "update user successfully", res)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid user id"))
	}

	err = h.userService.Delete(userID)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "delete user successfully", nil)
}

func (h *UserHandler) GetMe(c echo.Context) error {
	user, ok := c.Get("user").(*domain.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not found in context")
	}

	return c.JSON(http.StatusOK, user)
}
