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
	RecordCode   int    `json:"record_code"`
	UserStatus   string `json:"user_status"`
	RegisterDate string `json:"register_date"`
	Password     string `json:"password"`
}

type UserCreatedResponse struct {
	UserName   string `json:"user_name" validate:"required,min=1,max=20"`
	Email      string `json:"email" validate:"required,min=1,max=20,email"`
	Phone      string `json:"phone" validate:"required,e164"`
	RecordCode int    `json:"record_code"`
	Password   string `json:"password"`
	UserStatus string `json:"user_status"`
}

type UserLogin struct {
	RecordCode int    `json:"record_code"`
	Password   string `json:"password"`
}

type UserInfo struct {
	Passoword string `json:"password"`
	Status    string `json:"status"`
}
