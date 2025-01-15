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

const (
	Backlock   string = "BACKLOG"
	NotStarted string = "TODO"
	InProgress string = "IN_PROGRESS"
	Done       string = "DONE"
)

type Todo struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	StateID     uuid.UUID  `json:"state_id" gorm:"type:uuid;not null;index"`
	State       string     `json:"state" gorm:"default:not_started;foreignKey:StateID"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date" gorm:"type:timestamp;default:null"`
	IsDeleted   *bool      `json:"is_deleted" gorm:"default:false"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	User        User       `json:"user"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

type TodoState struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name   string    `json:"name" gorm:"not null;type:varchar(32)"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
}
