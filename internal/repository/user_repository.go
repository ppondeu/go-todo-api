package repository

import (
	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(newUser *domain.User) (*domain.User, error)
	Find(where interface{}) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Update(ID uuid.UUID, user *domain.User) (*domain.User, error)
	Delete(user *domain.User) error
	FindSession(where interface{}) (*domain.UserSession, error)
	SaveSession(newSession *domain.UserSession) (*domain.UserSession, error)
	UpdateSession(where interface{}, session *domain.UserSession) (*domain.UserSession, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Save(newUser *domain.User) (*domain.User, error) {
	err := u.db.Omit("Todos", "TodoCategory").Create(newUser).Error
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (u *userRepository) Find(where interface{}) (*domain.User, error) {
	user := &domain.User{}
	err := u.db.Where(where).Preload("Todos").Preload("TodoCategory").First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := u.db.Preload("Todos").Preload("TodoCategory").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userRepository) Update(ID uuid.UUID, user *domain.User) (*domain.User, error) {
	err := u.db.Model(user).Where("id = ?", ID).Updates(user).Error
	if err != nil {
		return nil, err
	}

	var updatedUser domain.User
	err = u.db.Where("id = ?", ID).Preload("Todos").Preload("TodoCategory").First(&updatedUser).Error
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (u *userRepository) Delete(user *domain.User) error {
	err := u.db.Delete(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) FindSession(where interface{}) (*domain.UserSession, error) {
	session := &domain.UserSession{}
	err := u.db.Where(where).First(session).Error
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *userRepository) SaveSession(newSession *domain.UserSession) (*domain.UserSession, error) {
	err := u.db.Create(newSession).Error
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (u *userRepository) UpdateSession(where interface{}, session *domain.UserSession) (*domain.UserSession, error) {
	updates := map[string]interface{}{
		"token": &session.Token,
		"expiry":        &session.Expiry,
	}
	err := u.db.Model(&domain.UserSession{}).Where(where).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	var sessionUpdated domain.UserSession
	err = u.db.Where(where).First(&sessionUpdated).Error
	if err != nil {
		return nil, err
	}

	return &sessionUpdated, nil
}
