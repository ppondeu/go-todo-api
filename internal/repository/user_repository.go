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
	Delete(ID uuid.UUID) error
	FindSession(where interface{}) (*domain.UserSession, error)
	SaveSession(newSession *domain.UserSession) (*domain.UserSession, error)
	UpdateSession(where interface{}, session *domain.UserSession) (*domain.UserSession, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) Save(newUser *domain.User) (*domain.User, error) {
	err := r.db.Omit("Todos", "TodoCategory").Create(newUser).Error
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = r.db.Where("id = ?", newUser.ID).Preload("Todos").Preload("TodoCategory").First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) Find(where interface{}) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.Where(where).Preload("Todos.Category").Preload("TodoCategory").First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Preload("Todos.Category").Preload("TodoCategory").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepositoryImpl) Update(ID uuid.UUID, user *domain.User) (*domain.User, error) {
	err := r.db.Model(user).Where("id = ?", ID).Updates(user).Error
	if err != nil {
		return nil, err
	}

	var updatedUser domain.User
	err = r.db.Where("id = ?", ID).Preload("Todos").Preload("TodoCategory").First(&updatedUser).Error
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (r *userRepositoryImpl) Delete(ID uuid.UUID) error {
	err := r.db.Where("id = ?", ID).Delete(&domain.User{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) FindSession(where interface{}) (*domain.UserSession, error) {
	session := &domain.UserSession{}
	err := r.db.Where(where).First(session).Error
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *userRepositoryImpl) SaveSession(newSession *domain.UserSession) (*domain.UserSession, error) {
	err := r.db.Create(newSession).Error
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (r *userRepositoryImpl) UpdateSession(where interface{}, session *domain.UserSession) (*domain.UserSession, error) {
	updates := map[string]interface{}{
		"token":  &session.Token,
		"expiry": &session.Expiry,
	}
	err := r.db.Model(&domain.UserSession{}).Where(where).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	var sessionUpdated domain.UserSession
	err = r.db.Where(where).First(&sessionUpdated).Error
	if err != nil {
		return nil, err
	}

	return &sessionUpdated, nil
}
