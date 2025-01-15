package dtos

type UserLoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type AuthResponse struct {
	Token TokenResponse `json:"token"`
	User  UserResponse  `json:"user"`
}
