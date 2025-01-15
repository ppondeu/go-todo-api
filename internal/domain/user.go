package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email        string      `json:"email" gorm:"unique;not null;type:varchar(64)"`
	FirstName    string      `json:"first_name" gorm:"type:varchar(32)"`
	LastName     string      `json:"last_name" gorm:"type:varchar(32)"`
	Password     string      `json:"password" gorm:"not null"`
	ImageURL     string      `json:"image_url" gorm:"default:'https://www.gravatar.com/avatar/?d=mp'"`
	RefreshToken *string     `json:"refresh_token" gorm:"default:null;"`
	Todos        []Todo      `json:"todos" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TodoStates   []TodoState `json:"todo_states" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt    time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time  `json:"deleted_at" gorm:"type:timestamp"`
}

type UserSession struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Token     *string    `json:"refresh_token" gorm:"unique"`
	Expiry    *time.Time `json:"expiry" gorm:"type:timestamp"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
