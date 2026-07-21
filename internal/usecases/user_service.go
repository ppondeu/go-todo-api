package usecases

import (
	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/internal/repositories"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/logs"
	"github.com/ppondeu/go-todo-api/pkg/utils"
)

type UserService interface {
	Save(newUser *dtos.UserCreateDTO) (*domain.User, error)
	FindAll() ([]domain.User, error)
	FindByUserID(ID uuid.UUID) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindByRefreshToken(refreshToken string) (*domain.User, error)
	Update(ID uuid.UUID, user *dtos.UserUpdateDTO) (*domain.User, error)
	UpdateRefreshToken(userID uuid.UUID, usesUpdateField map[string]interface{}) error
	Delete(ID uuid.UUID) error
}

type userServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

func (s *userServiceImpl) Save(userCreateDTO *dtos.UserCreateDTO) (*domain.User, error) {
	hashedPassword, err := utils.HashPassword(userCreateDTO.Password)
	if err != nil {
		return nil, errs.NewBadRequestError("Error hashing password")
	}

	newUser := &domain.User{
		Email:    userCreateDTO.Email,
		Password: *hashedPassword,
	}

	user, err := s.userRepo.Save(newUser)
	if err != nil {
		logs.Error("duplicate email")
		return nil, errs.NewBadRequestError("Email already exists")
	}

	return user, nil
}

func (s *userServiceImpl) FindAll() ([]domain.User, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	return users, nil
}

func (s *userServiceImpl) FindByUserID(ID uuid.UUID) (*domain.User, error) {
	where := map[string]interface{}{"id": ID}
	user, err := s.userRepo.Find(where)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) FindByEmail(email string) (*domain.User, error) {
	where := map[string]interface{}{"email": email}
	user, err := s.userRepo.Find(where)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) FindByRefreshToken(email string) (*domain.User, error) {
	where := map[string]interface{}{"refresh_token": email}
	user, err := s.userRepo.Find(where)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) Update(ID uuid.UUID, userDto *dtos.UserUpdateDTO) (*domain.User, error) {

	if userDto.Password != "" {
		hashedPassword, err := utils.HashPassword(userDto.Password)
		if err != nil {
			logs.Error("Error hashing password")
			return nil, err
		}
		userDto.Password = *hashedPassword
	}

	updateUser := &domain.User{
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
		Password:  userDto.Password,
	}

	users, err := s.userRepo.Update(ID, updateUser)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	return users, nil
}

func (s *userServiceImpl) UpdateRefreshToken(userID uuid.UUID, userUpdateField map[string]interface{}) error {
	user, err := s.userRepo.UpdateV2(userID, userUpdateField)
	if err != nil {
		logs.Error(err)
		return err
	}
	logs.Info(user)

	return nil
}

func (s *userServiceImpl) Delete(ID uuid.UUID) error {
	err := s.userRepo.Delete(ID)
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	return nil
}
