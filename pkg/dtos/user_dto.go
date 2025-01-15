package dtos

import "github.com/google/uuid"

type UserCreateDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserUpdateDTO struct {
	FirstName string `json:"first_name" validate:"omitempty"`
	LastName  string `json:"last_name" validate:"omitempty"`
	Password  string `json:"password" validate:"omitempty,min=6"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}
