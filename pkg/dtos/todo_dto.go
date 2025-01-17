package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
)

type CreateTodoDto struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	TodoStateID uuid.UUID `json:"todo_state_id" validate:"required,uuid"`
}

type UpdateTodoDto struct {
	Title       *string `json:"title" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Priority    *string `json:"priority" validate:"omitempty,oneof=high medium low"`
	TodoStateID *string `json:"todo_state_id" validate:"omitempty,nullable_todo_state_id"`
	DueDate     *string `json:"due_date" validate:"omitempty,nullable_due_date"`
}

type UpdateTodoState struct {
	Name string `json:"name" validate:"required"`
}

type TodoResponse struct {
	ID        uuid.UUID        `json:"id"`
	Title     string           `json:"title"`
	StateID   uuid.UUID        `json:"state_id"`
	Priority  *domain.Priority `json:"priority"`
	DueDate   *time.Time       `json:"due_date"`
	IsDeleted bool             `json:"is_deleted"`
	UserID    uuid.UUID        `json:"user_id"`
}

type TodoStateResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
