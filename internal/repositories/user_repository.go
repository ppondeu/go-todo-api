package repositories

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
	UpdateV2(ID uuid.UUID, userUpdateField map[string]interface{}) (*domain.User, error)
	Delete(ID uuid.UUID) error
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
	err := r.db.Omit("Todos").Create(newUser).Error
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = r.db.Where("id = ?", newUser.ID).Preload("Todos").First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) Find(where interface{}) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.Where(where).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
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
	err = r.db.Where("id = ?", ID).First(&updatedUser).Error
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (r *userRepositoryImpl) UpdateV2(ID uuid.UUID, userUpdateField map[string]interface{}) (*domain.User, error) {
	var user domain.User

	// Perform the update on the database with the WHERE condition based on the ID
	if err := r.db.Model(&user).Where("id = ?", ID).Updates(userUpdateField).Error; err != nil {
		return nil, err
	}

	// Fetch the updated user from the database
	if err := r.db.First(&user, "id = ?", ID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) Delete(ID uuid.UUID) error {
	err := r.db.Where("id = ?", ID).Delete(&domain.User{}).Error
	if err != nil {
		return err
	}

	return nil
}
