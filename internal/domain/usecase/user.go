package usecase

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"net/textproto"
	"time"

	"github.com/ellofae/authentication-deanery/internal/controller/middleware"
	"github.com/ellofae/authentication-deanery/internal/database"
	"github.com/ellofae/authentication-deanery/internal/dto"
	"github.com/ellofae/authentication-deanery/internal/models"
	"github.com/ellofae/authentication-deanery/internal/utils"
	"github.com/ellofae/authentication-deanery/pkg/gist"
	"github.com/ellofae/authentication-deanery/pkg/logger"
	"github.com/jordan-wright/email"
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

func (u *UserUsecase) UserLogin(user *dto.UserLogin) (*models.Tokens, error) {
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

	var stored_password []byte

	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		stored_password, err = u.repo.GetPasswordByRecordCode(ctx, user.RecordCode)
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

	compareResult, err := utils.ComparePasswords(user.Password, string(stored_password), u.cfgUsecase.AesEncryptionKey)
	if err != nil {
		return nil, err
	}

	if !compareResult {
		return nil, fmt.Errorf("wrong password for the passed record code")
	}

	accessToken, err := middleware.GenerateAccessToken(user.RecordCode)
	if err != nil {
		return nil, err
	}

	return &models.Tokens{
		AccessToken: accessToken,
	}, nil
}

func (u *UserUsecase) RetreiveRoles() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var err error

	var roles []byte

	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		roles, err = u.repo.RetreiveRoles(ctx)
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

	return roles, nil
}

func (u *UserUsecase) SendPassword(user *dto.EmailForm) error {
	e := &email.Email{
		To:      []string{user.Email},
		From:    fmt.Sprintf(gist.GistText.From_message, u.cfgUsecase.SmtpService.SmtpEmail),
		Subject: gist.GistText.Subject_message,
		Text:    []byte(fmt.Sprintf(gist.GistText.Email_message, user.UserName, user.Status, user.GeneratedPassword)),
		// HTML:    []byte("<h1>Fancy HTML is supported, too!</h1>"),
		Headers: textproto.MIMEHeader{},
	}

	if err := e.Send(u.cfgUsecase.SmtpService.SmtpAddress, smtp.PlainAuth("", u.cfgUsecase.SmtpService.SmtpEmail, u.cfgUsecase.SmtpService.SmtpPassword, u.cfgUsecase.SmtpService.SmtpService)); err != nil {
		u.logger.Printf("error shile sending: %v\n", err)
		return err
	}

	return nil
}
