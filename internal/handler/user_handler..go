package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/internal/usecase"
	"github.com/ppondeu/go-todo-api/pkg/dto"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/logs"
	"github.com/ppondeu/go-todo-api/pkg/response"
)

type UserHandler struct {
	userService usecase.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService usecase.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{userService: userService, validator: validator}
}

func (h *UserHandler) Register(c echo.Context) error {
	createUserRequest := new(dto.CreateUserDto)
	if err := c.Bind(createUserRequest); err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid user id"))
	}

	if err := h.validator.Struct(createUserRequest); err != nil {
		logs.Error(err)
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return response.NewErrorResponse(c, errs.NewBadRequestError("invalid field"))
		}

		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid field"))
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
		response.NewErrorResponse(c, errs.NewBadRequestError("invalid user id"))
	}

	user, err := h.userService.FindByUserID(userID)
	if err != nil {
		response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "get user successfully", user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid user id"))
	}

	updateUserRequest := new(dto.UpdateUserDto)
	if err := c.Bind(updateUserRequest); err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid JSON body"))
	}

	if err := h.validator.Struct(updateUserRequest); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return response.NewErrorResponse(c, errs.NewBadRequestError("invalid field"))
		}

		return response.NewErrorResponse(c, errs.NewInternalError("something wrong while validate body"))
	}

	user, err := h.userService.Update(userID, updateUserRequest)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "update user successfully", user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequestError("invalid user id"))
	}

	err = h.userService.Delete(userID)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "delete user successfully", nil)
}
