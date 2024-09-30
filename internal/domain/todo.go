package domain

import (
	"time"

	"github.com/google/uuid"
)

type Priority string

const (
	High   Priority = "high"
	Medium Priority = "medium"
	Low    Priority = "low"
)

type TodoState string

const (
	NotStarted TodoState = "not_started"
	InProgress TodoState = "in_progress"
	Done       TodoState = "done"
)

type Todo struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Title       string       `json:"title" gorm:"not null;type:varchar(50)"`
	Description string       `json:"description" gorm:"type:text"`
	State       TodoState    `json:"state" gorm:"default:not_started"`
	Priority    Priority     `json:"priority" gorm:"default:medium"`
	IsCompleted bool         `json:"is_completed" gorm:"default:false"`
	DueDate     time.Time    `json:"due_date" gorm:"type:timestamp"`
	CategoryID  *uuid.UUID   `json:"category_id" gorm:"type:uuid"`
	Category    TodoCategory `json:"category" gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	IsDeleted   bool         `json:"is_deleted" gorm:"default:false"`
	IsOverdue   bool         `json:"is_overdue" gorm:"default:false"`
	UserID      uuid.UUID    `json:"user_id" gorm:"type:uuid;not null"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

type TodoCategory struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name   string    `json:"name" gorm:"not null;type:varchar(50)"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
}
