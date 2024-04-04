package dto

type UserRegistration struct {
	UserName   string `json:"user_name" validate:"required,min=1,max=20"`
	Email      string `json:"email" validate:"required,min=1,max=20,email"`
	Phone      string `json:"phone" validate:"required,e164"`
	UserStatus string `json:"user_status" validate:"required"`
}

type UserCreated struct {
	UserName     string `json:"user_name" validate:"required,min=1,max=20"`
	Email        string `json:"email" validate:"required,min=1,max=20,email"`
	Phone        string `json:"phone" validate:"required,e164"`
	Credentials  int    `json:"credentials"`
	UserStatus   string `json:"user_status"`
	RegisterDate string `json:"register_date"`
	Password     string `json:"password"`
}
