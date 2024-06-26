package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ellofae/authentication-deanery/internal/controller"
	"github.com/ellofae/authentication-deanery/internal/domain"
	"github.com/ellofae/authentication-deanery/internal/dto"
	"github.com/ellofae/authentication-deanery/internal/utils"
	"github.com/ellofae/authentication-deanery/pkg/logger"
)

type UserHandler struct {
	logger  *log.Logger
	usecase domain.IUserUsecase
}

func NewUserHandler(userUsecase domain.IUserUsecase) controller.IHandler {
	return &UserHandler{
		logger:  logger.GetLogger(),
		usecase: userUsecase,
	}
}

func (h *UserHandler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		var err error

		parsed_url := strings.TrimPrefix(r.URL.Path, "/api/")
		url_parts := strings.Split(parsed_url, "/")

		switch r.Method {
		case http.MethodGet:
			if len(url_parts) == 1 {
				if url_parts[0] == "roles" {
					err = h.handleExistingRoles(w, r)
					if err != nil {
						http.Error(w, fmt.Sprintf("Unable to retreive existing roles. Error: %v.\n", err.Error()), http.StatusInternalServerError)
						return
					}

					return
				}
			}
		case http.MethodPost:
			if len(url_parts) == 1 {
				if url_parts[0] == "signup" {
					err = h.handleUserCreation(w, r)
					if err != nil {
						http.Error(w, fmt.Sprintf("Unable to create a user. Error: %v.\n", err.Error()), http.StatusInternalServerError)
						return
					}

					return
				} else if url_parts[0] == "login" {
					err = h.handleUserLogin(w, r)
					if err != nil {
						http.Error(w, fmt.Sprintf("Unable to login a user. Error: %v.\n", err.Error()), http.StatusInternalServerError)
						return
					}

					return
				} else if url_parts[0] == "get_username" {
					err = h.handleGetUsername(w, r)
					if err != nil {
						http.Error(w, fmt.Sprintf("Unable to retreive the username by record code. Error: %v.\n", err.Error()), http.StatusInternalServerError)
						return
					}

					return
				}
			}
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

}

func (h *UserHandler) handleUserCreation(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")

	var err error

	user := &dto.UserRegistration{}
	if err = utils.RequestDecode(r, user); err != nil {
		return err
	}

	json_data, err := h.usecase.CreateUser(user)
	if err != nil {
		return err
	}

	createdUser := &dto.UserCreated{}
	if err = json.Unmarshal(json_data, createdUser); err != nil {
		return err
	}
	generatedPassword, err := h.usecase.SetEncryptedPassword(createdUser.Credentials)
	if err != nil {
		return err
	}
	createdUser.Password = generatedPassword

	if err := h.usecase.SendPassword(&dto.EmailForm{
		UserName:          createdUser.UserName,
		Email:             createdUser.Email,
		Phone:             createdUser.Phone,
		Code:              createdUser.RecordCode,
		GeneratedPassword: generatedPassword,
		Status:            createdUser.UserStatus,
	}); err != nil {
		return err
	}

	response := &dto.UserCreatedResponse{
		UserName:   createdUser.UserName,
		Email:      createdUser.Email,
		Phone:      createdUser.Phone,
		RecordCode: createdUser.RecordCode,
		Password:   createdUser.Password,
		UserStatus: createdUser.UserStatus,
	}

	response_json, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response_json)
	return nil
}

func (h *UserHandler) handleUserLogin(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")

	var err error

	user := &dto.UserLogin{}
	if err = utils.RequestDecode(r, user); err != nil {
		return err
	}

	tokens, err := h.usecase.UserLogin(user)
	if err != nil {
		return err
	}

	json_data, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json_data)
	return nil
}

func (h *UserHandler) handleExistingRoles(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")

	var err error

	json_data, err := h.usecase.RetreiveRoles()
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json_data)
	return nil
}

func (h *UserHandler) handleGetUsername(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")

	var err error

	code := &dto.RecordCode{}
	if err = utils.RequestDecode(r, code); err != nil {
		return err
	}

	json_data, err := h.usecase.GetUsername(code)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json_data)
	return nil
}
