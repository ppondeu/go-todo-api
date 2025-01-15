package handlers

import (
	"net/http"
	"time"

	v "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/internal/usecases"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/logs"
	"github.com/ppondeu/go-todo-api/pkg/response"
)

type AuthHandler struct {
	authService usecases.AuthService
	validator   *v.Validate
}

func NewAuthHandler(authService *usecases.AuthService, validator *v.Validate) *AuthHandler {
	return &AuthHandler{authService: *authService, validator: validator}
}

func (h *AuthHandler) Login(c echo.Context) error {
	loginRequest := new(dtos.UserLoginDTO)
	err := c.Bind(loginRequest)
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid JSON"))
	}

	logs.Info(loginRequest)

	err = h.validator.Struct(loginRequest)
	if err != nil {
		logs.Error(err)
		if _, ok := err.(*v.ValidationErrors); ok {
			return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid field"))
		}
		return response.NewErrorResponse(c, errs.NewBadRequestError("Field required"))
	}

	res, err := h.authService.Login(loginRequest)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = "refresh_token"
	refreshCookie.Value = res.Token.RefreshToken
	refreshCookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	refreshCookie.HttpOnly = true
	refreshCookie.Secure = false
	refreshCookie.SameSite = http.SameSiteLaxMode
	refreshCookie.Path = "/"

	accessCookie := new(http.Cookie)
	accessCookie.Name = "access_token"
	accessCookie.Value = res.Token.AccessToken
	accessCookie.Expires = time.Now().Add(15 * time.Minute)
	accessCookie.HttpOnly = true
	accessCookie.Secure = false
	accessCookie.SameSite = http.SameSiteLaxMode
	accessCookie.Path = "/"

	c.SetCookie(accessCookie)
	c.SetCookie(refreshCookie)

	return response.NewSuccessAPIResponse(c, "Login successfully", res)
}

func (h *AuthHandler) Register(c echo.Context) error {
	userUpdateDTO := new(dtos.UserCreateDTO)
	err := c.Bind(userUpdateDTO)
	if err != nil {
		logs.Error(err)
		return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid JSON"))
	}
	logs.Info(userUpdateDTO)

	err = h.validator.Struct(userUpdateDTO)
	if err != nil {
		logs.Error(err)
		if _, ok := err.(*v.ValidationErrors); ok {
			return response.NewErrorResponse(c, errs.NewBadRequestError("Invalid field"))
		}
		return response.NewErrorResponse(c, errs.NewBadRequestError("Field required"))
	}

	res, err := h.authService.Register(userUpdateDTO)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = "refresh_token"
	refreshCookie.Value = res.Token.RefreshToken
	refreshCookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	refreshCookie.HttpOnly = true
	refreshCookie.Secure = false
	refreshCookie.SameSite = http.SameSiteLaxMode
	refreshCookie.Path = "/"

	accessCookie := new(http.Cookie)
	accessCookie.Name = "access_token"
	accessCookie.Value = res.Token.AccessToken
	accessCookie.Expires = time.Now().Add(15 * time.Minute)
	accessCookie.HttpOnly = true
	accessCookie.Secure = false
	accessCookie.SameSite = http.SameSiteLaxMode
	accessCookie.Path = "/"

	c.SetCookie(accessCookie)
	c.SetCookie(refreshCookie)

	return response.NewSuccessAPIResponse(c, "Register successfully", res)
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	user, ok := c.Get("user").(*domain.User)
	if !ok {
		return response.NewErrorResponse(c, errs.NewUnauthorizedError("Unauthorized"))
	}

	tokenResponse, err := h.authService.RefreshToken(user.ID)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}

	return response.NewSuccessAPIResponse(c, "refresh token successfully", tokenResponse)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	logs.Info("\ntest\n")
	user, ok := c.Get("user").(*domain.User)
	if !ok {
		return response.NewErrorResponse(c, errs.NewUnauthorizedError("Unauthorized"))
	}

	logs.Info(user)

	err := h.authService.Logout(user)
	if err != nil {
		return response.NewErrorResponse(c, err)
	}
	return response.NewSuccessAPIResponse(c, "user logged out successfully", nil)
}
