package usecase

import (
	"context"
	"log"
	"time"

	"github.com/ellofae/authentication-deanery/internal/database"
	"github.com/ellofae/authentication-deanery/internal/dto"
	"github.com/ellofae/authentication-deanery/internal/utils"
	"github.com/ellofae/authentication-deanery/pkg/logger"
)

const (
	PASSWORD_LENGTH = 14
)

type UserUsecase struct {
	logger *log.Logger
	repo   database.IUserRepository
}

func NewUserUsecase(userRepository database.IUserRepository) *UserUsecase {
	return &UserUsecase{
		logger: logger.GetLogger(),
		repo:   userRepository,
	}
}

func (u *UserUsecase) CreateUser(user *dto.UserRegistration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var err error

	validate := utils.NewValidator()
	if err = validate.Struct(user); err != nil {
		validation_errors := utils.ValidatorErrors(err)
		for _, error := range validation_errors {
			u.logger.Printf("User registration model validation error. Error: %v.\n", error)
		}

		return nil, err
	}

	var json_data []byte

	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		json_data, err = u.repo.CreateUser(ctx, user)
		errChan <- err
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errChan:
		if err != nil {
			u.logger.Printf("An error occured in repository while creating user. Error: %v.\n", err.Error())
			return nil, err
		}
	}

	// generated_password := utils.GenerateRandomPassword(PASSWORD_LENGTH)
	// encryptedPassword, err := utils.Encrypt([]byte(generated_password), []byte("passphrasewhichneedstobe32bytes2"))
	// if err != nil {
	// 	u.logger.Printf("Unable to enctype the generated password, error: %v\n", err)
	// 	return []byte{}, err
	// }

	return json_data, nil
}
