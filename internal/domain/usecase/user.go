package usecase

import (
	"context"
	"log"
	"time"

	"github.com/ellofae/authentication-deanery/internal/database"
	"github.com/ellofae/authentication-deanery/internal/dto"
	"github.com/ellofae/authentication-deanery/internal/models"
	"github.com/ellofae/authentication-deanery/internal/utils"
	"github.com/ellofae/authentication-deanery/pkg/logger"
)

type UserUsecase struct {
	logger     *log.Logger
	repo       database.IUserRepository
	cfgUsecase *models.CfgUsecaseData
}

func NewUserUsecase(userRepository database.IUserRepository, cfgUsecase *models.CfgUsecaseData) *UserUsecase {
	return &UserUsecase{
		logger:     logger.GetLogger(),
		repo:       userRepository,
		cfgUsecase: cfgUsecase,
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

	return json_data, nil
}

func (u *UserUsecase) SetEncryptedPassword(credentials_id int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var err error

	generated_password := utils.GenerateRandomPassword(int(u.cfgUsecase.PasswordLength))
	encryptedPassword, err := utils.Encrypt([]byte(generated_password), []byte(u.cfgUsecase.AesEncryptionKey))
	if err != nil {
		u.logger.Printf("Unable to encrypt the generated password, error: %v\n", err)
		return "", err
	}

	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		err = u.repo.SetEncryptedPassword(ctx, credentials_id, string(encryptedPassword))
		errChan <- err
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errChan:
		if err != nil {
			u.logger.Printf("An error occured in repository while creating user. Error: %v.\n", err.Error())
			return "nil", err
		}
	}

	return generated_password, nil
}
