package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/internal/usecases"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/response"
)

type UserHandler struct {
	userService usecases.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService usecases.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{userService: userService, validator: validator}
}

func (h *UserHandler) GetMe(c echo.Context) error {
	authenticated, err := authenticatedUser(c)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	user, err := h.userService.FindByUserID(c.Request().Context(), authenticated.ID)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	return response.NewSuccessAPIResponse(c, "get user successfully", toUserResponse(user))
}

func (h *UserHandler) UpdateMe(c echo.Context) error {
	authenticated, err := authenticatedUser(c)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	request := new(dtos.UserUpdateDTO)
	if err := c.Bind(request); err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid JSON body"))
	}
	if err := h.validator.Struct(request); err != nil {
		return response.NewErrorResponse(c, errs.NewBadRequestError("invalid field"))
	}

	user, err := h.userService.Update(c.Request().Context(), authenticated.ID, request)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	return response.NewSuccessAPIResponse(c, "update user successfully", toUserResponse(user))
}

func (h *UserHandler) DeleteMe(c echo.Context) error {
	authenticated, err := authenticatedUser(c)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	if err := h.userService.Delete(c.Request().Context(), authenticated.ID); err != nil {
		return response.NewErrorResponse(c, err)
	}

	clearAuthCookies(c)
	return response.NewSuccessAPIResponse(c, "delete user successfully", nil)
}

func toUserResponse(user *domain.User) dtos.UserResponse {
	return dtos.UserResponse{
		ID: user.ID, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName, ImageURL: user.ImageURL,
	}
}
