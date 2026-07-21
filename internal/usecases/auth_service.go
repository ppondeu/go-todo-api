package usecases

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/utils"
)

type AuthService interface {
	Login(context.Context, *dtos.UserLoginDTO) (*dtos.AuthResponse, error)
	Register(context.Context, *dtos.UserCreateDTO) (*dtos.AuthResponse, error)
	RefreshToken(context.Context, uuid.UUID) (*dtos.TokenResponse, error)
	Logout(context.Context, *domain.User) error
}

type authServiceImpl struct {
	userService UserService
	todoService TodoService
	jwtService  JWTService
}

func NewAuthService(userService UserService, todoService TodoService, jwtService JWTService) AuthService {
	return &authServiceImpl{userService: userService, todoService: todoService, jwtService: jwtService}
}

func (s *authServiceImpl) Login(ctx context.Context, dto *dtos.UserLoginDTO) (*dtos.AuthResponse, error) {
	user, err := s.userService.FindByEmail(ctx, dto.Email)
	if err != nil || utils.ComparePassword(user.Password, dto.Password) != nil {
		return nil, errs.NewBadRequestError("email or password doesn't match")
	}

	tokens, err := s.RefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &dtos.AuthResponse{User: mapUserResponse(user), Token: *tokens}, nil
}

func (s *authServiceImpl) Register(ctx context.Context, dto *dtos.UserCreateDTO) (*dtos.AuthResponse, error) {
	user, err := s.userService.Save(ctx, dto)
	if err != nil {
		return nil, err
	}
	if _, err := s.todoService.InitTodoState(ctx, user.ID); err != nil {
		return nil, err
	}

	tokens, err := s.RefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &dtos.AuthResponse{User: mapUserResponse(user), Token: *tokens}, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, userID uuid.UUID) (*dtos.TokenResponse, error) {
	tokens, err := s.getTokens(userID.String())
	if err != nil {
		return nil, err
	}
	if err := s.userService.UpdateRefreshToken(ctx, userID, &tokens.RefreshToken); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *authServiceImpl) getTokens(userID string) (*dtos.TokenResponse, error) {
	claims := &dtos.UserClaims{RegisteredClaims: jwt.RegisteredClaims{
		Subject: userID, ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), IssuedAt: jwt.NewNumericDate(time.Now()),
	}}
	accessToken, err := s.jwtService.SignToken(claims, s.jwtService.GetAccess())
	if err != nil {
		return nil, errs.NewInternalError("could not create access token")
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))
	refreshToken, err := s.jwtService.SignToken(claims, s.jwtService.GetRefresh())
	if err != nil {
		return nil, errs.NewInternalError("could not create refresh token")
	}

	return &dtos.TokenResponse{AccessToken: *accessToken, RefreshToken: *refreshToken}, nil
}

func (s *authServiceImpl) Logout(ctx context.Context, user *domain.User) error {
	if err := s.userService.UpdateRefreshToken(ctx, user.ID, nil); err != nil {
		return errs.NewBadRequestError("logout failed")
	}
	return nil
}

func mapUserResponse(user *domain.User) dtos.UserResponse {
	return dtos.UserResponse{
		ID: user.ID, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName, ImageURL: user.ImageURL,
	}
}
