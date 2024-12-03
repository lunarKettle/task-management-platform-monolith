package dto

// RegisterUserRequestDTO - данные для регистрации пользователя
type RegisterUserRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Username string `json:"name" validate:"required"`
	Role     string `json:"role"`
}

// LoginUserRequestDTO - данные для входа пользователя
type LoginUserRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
