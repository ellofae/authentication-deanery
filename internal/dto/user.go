package dto

type UserRegistration struct {
	DisplayName string `json:"display_name" validate:"required,min=1,max=20"`
	Email       string `json:"email" validate:"required,min=1,max=20,email"`
	Phone       string `json:"phone" validate:"required,e164"`
	UserStatus  string `json:"user_status" validate:"required"`
}
