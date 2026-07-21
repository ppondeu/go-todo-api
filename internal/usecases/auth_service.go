package usecases

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/logs"
	"github.com/ppondeu/go-todo-api/pkg/utils"
)

type AuthService interface {
	Login(loginDto *dtos.UserLoginDTO) (*dtos.AuthResponse, error)
	Register(userCreateDTO *dtos.UserCreateDTO) (*dtos.AuthResponse, error)
	RefreshToken(token uuid.UUID) (*dtos.TokenResponse, error)
	Logout(user *domain.User) error
}

type authServiceImpl struct {
	userService UserService
	todoService TodoService
	jwtService  JWTService
}

func NewAuthService(userService UserService, todoService TodoService, jwtService JWTService) AuthService {
	return &authServiceImpl{userService: userService, todoService: todoService, jwtService: jwtService}
}

func (s *authServiceImpl) Login(loginDto *dtos.UserLoginDTO) (*dtos.AuthResponse, error) {
	user, err := s.userService.FindByEmail(loginDto.Email)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewBadRequestError("email or password doesn't match")
	}

	err = utils.ComparePassword(user.Password, loginDto.Password)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewBadRequestError("email or password doesn't match")
	}

	tokenResponse, err := s.RefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	authResponse := &dtos.AuthResponse{
		User: dtos.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		Token: *tokenResponse,
	}

	return authResponse, nil
}

func (s *authServiceImpl) Register(userCreateDTO *dtos.UserCreateDTO) (*dtos.AuthResponse, error) {
	user, err := s.userService.Save(userCreateDTO)
	if err != nil {
		return nil, err
	}

	_, err = s.todoService.InitTodoState(user.ID)
	if err != nil {
		return nil, err
	}

	tokenResponse, err := s.RefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	authResponse := &dtos.AuthResponse{
		User: dtos.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		Token: *tokenResponse,
	}

	return authResponse, nil
}

func (s *authServiceImpl) RefreshToken(userID uuid.UUID) (*dtos.TokenResponse, error) {

	tokenResponse, err := s.getTokens(userID.String())
	if err != nil {
		return nil, err
	}

	userUpdateField := map[string]interface{}{
		"refresh_token": tokenResponse.RefreshToken,
	}
	if err := s.userService.UpdateRefreshToken(userID, userUpdateField); err != nil {
		return nil, err
	}

	return tokenResponse, nil
}

func (s *authServiceImpl) getTokens(userID string) (*dtos.TokenResponse, error) {
	claims := &dtos.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := s.jwtService.SignToken(claims, s.jwtService.GetAccess())
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))

	refreshToken, err := s.jwtService.SignToken(claims, s.jwtService.GetRefresh())
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	tokenResponse := &dtos.TokenResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}

	return tokenResponse, nil
}

func (s *authServiceImpl) Logout(user *domain.User) error {
	userUpdateField := map[string]interface{}{
		"refresh_token": nil,
	}
	err := s.userService.UpdateRefreshToken(user.ID, userUpdateField)
	if err != nil {
		return errs.NewBadRequestError("logout failed")
	}

	return nil
}
