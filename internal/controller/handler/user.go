package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	// endpoint /users
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		var err error

		switch r.Method {
		case http.MethodPost:
			err = h.handleUserCreation(w, r)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to create a user. Error: %v.\n", err.Error()), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
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

	created_user_data, err := json.Marshal(createdUser)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(created_user_data)
	return nil
}
