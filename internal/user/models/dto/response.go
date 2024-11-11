package dto

type RegisterUserResponseDTO struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type LoginUserResponseDTO struct {
	AccessToken string `json:"access_token" validate:"required"`
}
