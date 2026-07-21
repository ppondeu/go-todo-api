package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, userID uuid.UUID, fields map[string]any) (*domain.User, error)
	UpdateRefreshToken(ctx context.Context, userID uuid.UUID, token *string) error
	DeleteAccount(ctx context.Context, userID uuid.UUID) error
}
