package database

import (
	"context"

	"github.com/ellofae/authentication-deanery/internal/dto"
)

type IUserRepository interface {
	CreateUser(context.Context, *dto.UserRegistration) ([]byte, error)
}