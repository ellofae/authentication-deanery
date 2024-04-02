package domain

import "github.com/ellofae/authentication-deanery/internal/dto"

type IUserUsecase interface {
	CreateUser(*dto.UserRegistration) ([]byte, error)
}
