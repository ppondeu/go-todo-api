package dto

type CreateUserDto struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserDto struct {
	Name     string `json:"name" validate:"omitempty"`
	Password string `json:"password" validate:"omitempty,min=6"`
}
