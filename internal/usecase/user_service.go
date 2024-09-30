package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/internal/repository"
	"github.com/ppondeu/go-todo-api/pkg/dto"
	"github.com/ppondeu/go-todo-api/pkg/utils"
)

type UserService interface {
	Save(newUser *dto.CreateUserDto) (*domain.User, error)
	Delete(ID uuid.UUID) error
	FindByUserID(ID uuid.UUID) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Update(ID uuid.UUID, user *dto.UpdateUserDto) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindSession(ID uuid.UUID) (*domain.UserSession, error)
}

type userServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: *userRepo,
	}
}

func (s *userServiceImpl) Save(userDto *dto.CreateUserDto) (*domain.User, error) {
	hashedPassword, err := utils.HashPassword(userDto.Password)
	if err != nil {
		fmt.Printf("Error hashing password: %v", err)
		return nil, err
	}

	newUser := &domain.User{
		Name:     userDto.Name,
		Email:    userDto.Email,
		Password: *hashedPassword,
	}

	user, err := s.userRepo.Save(newUser)
	if err != nil {
		fmt.Printf("Error saving user: %v", err)
		return nil, err
	}

	session := &domain.UserSession{
		UserID: user.ID,
	}

	_, err = s.userRepo.SaveSession(session)
	if err != nil {
		fmt.Printf("Error saving session: %v", err)
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) Delete(ID uuid.UUID) error {
	err := s.userRepo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) FindByUserID(ID uuid.UUID) (*domain.User, error) {
	where := map[string]interface{}{"id": ID}
	user, err := s.userRepo.Find(where)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) FindAll() ([]domain.User, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userServiceImpl) Update(ID uuid.UUID, userDto *dto.UpdateUserDto) (*domain.User, error) {

	if userDto.Password != "" {
		hashedPassword, err := utils.HashPassword(userDto.Password)
		if err != nil {
			return nil, err
		}
		userDto.Password = *hashedPassword
	}

	updateUser := &domain.User{
		Name:     userDto.Name,
		Email:    userDto.Email,
		Password: userDto.Password,
	}

	users, err := s.userRepo.Update(ID, updateUser)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userServiceImpl) FindByEmail(email string) (*domain.User, error) {
	where := map[string]interface{}{"email": email}
	user, err := s.userRepo.Find(where)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) FindSession(ID uuid.UUID) (*domain.UserSession, error) {
	where := map[string]interface{}{"user_id": ID}
	session, err := s.userRepo.FindSession(where)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *userServiceImpl) UpdateSession(ID uuid.UUID, session *domain.UserSession) (*domain.UserSession, error) {
	where := map[string]interface{}{"user_id": ID}
	sess, err := s.userRepo.UpdateSession(where, session)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
