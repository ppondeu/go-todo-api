package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name         string         `json:"name" gorm:"not null;type:varchar(50)"`
	Email        string         `json:"email" gorm:"unique;not null;type:varchar(50)"`
	Password     string         `json:"password" gorm:"not null"`
	ImageURL     string         `json:"image_url" gorm:"default:'https://www.gravatar.com/avatar/?d=mp'"`
	Todos        []Todo         `json:"todos" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TodoCategory []TodoCategory `json:"todo_categories" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type UserSession struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Token     *string    `json:"refresh_token" gorm:"unique"`
	Expiry    *time.Time `json:"expiry" gorm:"type:timestamp"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
