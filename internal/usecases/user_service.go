package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/application/ports"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/utils"
)

type UserService interface {
	Save(context.Context, *dtos.UserCreateDTO) (*domain.User, error)
	FindByUserID(context.Context, uuid.UUID) (*domain.User, error)
	FindByEmail(context.Context, string) (*domain.User, error)
	Update(context.Context, uuid.UUID, *dtos.UserUpdateDTO) (*domain.User, error)
	UpdateRefreshToken(context.Context, uuid.UUID, *string) error
	Delete(context.Context, uuid.UUID) error
}

type userServiceImpl struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) Save(ctx context.Context, dto *dtos.UserCreateDTO) (*domain.User, error) {
	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, errs.NewInternalError("could not create user")
	}

	user, err := s.userRepo.Create(ctx, &domain.User{Email: dto.Email, Password: *hashedPassword})
	if err != nil {
		return nil, mapUserError(err)
	}

	return user, nil
}

func (s *userServiceImpl) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, mapUserError(err)
	}
	return user, nil
}

func (s *userServiceImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, mapUserError(err)
	}
	return user, nil
}

func (s *userServiceImpl) Update(ctx context.Context, userID uuid.UUID, dto *dtos.UserUpdateDTO) (*domain.User, error) {
	fields := make(map[string]any)
	if dto.FirstName != nil {
		fields["first_name"] = *dto.FirstName
	}
	if dto.LastName != nil {
		fields["last_name"] = *dto.LastName
	}
	if dto.Password != nil {
		hashedPassword, err := utils.HashPassword(*dto.Password)
		if err != nil {
			return nil, errs.NewInternalError("could not update user")
		}
		fields["password"] = *hashedPassword
	}
	if len(fields) == 0 {
		return nil, errs.NewBadRequestError("at least one field is required")
	}

	user, err := s.userRepo.Update(ctx, userID, fields)
	if err != nil {
		return nil, mapUserError(err)
	}
	return user, nil
}

func (s *userServiceImpl) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, token *string) error {
	return mapUserError(s.userRepo.UpdateRefreshToken(ctx, userID, token))
}

func (s *userServiceImpl) Delete(ctx context.Context, userID uuid.UUID) error {
	return mapUserError(s.userRepo.DeleteAccount(ctx, userID))
}

func mapUserError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, domain.ErrUserNotFound) {
		return errs.NewNotFoundError("user not found")
	}
	if errors.Is(err, domain.ErrEmailConflict) {
		return errs.NewConflictError("email already exists")
	}
	return errs.NewInternalError("user operation failed")
}
