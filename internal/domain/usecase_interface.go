package domain

import (
	"github.com/ellofae/authentication-deanery/internal/dto"
	"github.com/ellofae/authentication-deanery/internal/models"
)

type IUserUsecase interface {
	CreateUser(*dto.UserRegistration) ([]byte, error)
	SetEncryptedPassword(int) (string, error)
	UserLogin(*dto.UserLogin) (*models.Tokens, error)
	RetreiveRoles() ([]byte, error)
}
