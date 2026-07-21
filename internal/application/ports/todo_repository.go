package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
)

type TodoListFilter struct {
	Page     int
	Limit    int
	Search   string
	Priority domain.Priority
	StateID  *uuid.UUID
	DueFrom  *time.Time
	DueTo    *time.Time
	Sort     string
	Order    string
}

type TodoPage struct {
	Items      []domain.Todo
	Page       int
	Limit      int
	Total      int64
	TotalPages int
}

type TodoRepository interface {
	Save(ctx context.Context, userID uuid.UUID, todo *domain.Todo) (*domain.Todo, error)
	FindByID(ctx context.Context, userID, todoID uuid.UUID) (*domain.Todo, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Todo, error)
	List(ctx context.Context, userID uuid.UUID, filter TodoListFilter) (*TodoPage, error)
	Update(ctx context.Context, userID, todoID uuid.UUID, fields map[string]any) (*domain.Todo, error)
	Delete(ctx context.Context, userID, todoID uuid.UUID) error
	FindStateByID(ctx context.Context, userID, stateID uuid.UUID) (*domain.TodoState, error)
	UpdateTodoState(ctx context.Context, userID, stateID uuid.UUID, fields map[string]any) (*domain.TodoState, error)
	InitTodoState(ctx context.Context, userID uuid.UUID) ([]domain.TodoState, error)
}
