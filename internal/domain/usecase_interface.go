package domain

import "github.com/ellofae/authentication-deanery/internal/dto"

type IUserUsecase interface {
	CreateUser(*dto.UserRegistration) ([]byte, error)
	SetEncryptedPassword(int) (string, error)
	UserLogin(*dto.UserLogin) (string, error)
}
