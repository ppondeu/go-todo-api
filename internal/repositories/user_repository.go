package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/application/ports"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.WithContext(ctx).Omit("Todos", "TodoStates").Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, domain.ErrEmailConflict
		}
		return nil, err
	}

	return r.FindByID(ctx, user.ID)
}

func (r *userRepository) FindByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, userID uuid.UUID, fields map[string]any) (*domain.User, error) {
	result := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).Updates(fields)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrUserNotFound
	}

	return r.FindByID(ctx, userID)
}

func (r *userRepository) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, token *string) error {
	result := r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", userID).
		Update("refresh_token", token)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&domain.User{}).Where("id = ?", userID).Update("refresh_token", nil)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domain.ErrUserNotFound
		}

		result = tx.Where("id = ?", userID).Delete(&domain.User{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domain.ErrUserNotFound
		}

		return nil
	})
}
