package repository

import (
	"context"
	"log"

	"github.com/ellofae/authentication-deanery/internal/database"
	"github.com/ellofae/authentication-deanery/internal/dto"
	"github.com/ellofae/authentication-deanery/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	logger *log.Logger
	pool   *pgxpool.Pool
}

func NewUserRepository(connPool *pgxpool.Pool) database.IUserRepository {
	return &UserRepository{
		logger: logger.GetLogger(),
		pool:   connPool,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *dto.UserRegistration) ([]byte, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Printf("Error acquiring connection. Error: %v.\n", err.Error())
		return nil, err
	}
	defer conn.Release()

	var data []byte
	err = conn.QueryRow(ctx, "SELECT create_user($1, $2, $3, $4)", user.UserName, user.Email, user.Phone, user.UserStatus).Scan(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *UserRepository) SetEncryptedPassword(ctx context.Context, credentials_id int, encrypted_password string) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Printf("Error acquiring connection. Error: %v.\n", err.Error())
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "SELECT set_user_password($1, $2)", credentials_id, encrypted_password)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetPasswordByRecordCode(ctx context.Context, record_code int) ([]byte, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Printf("Error acquiring connection. Error: %v.\n", err.Error())
		return nil, err
	}
	defer conn.Release()

	var data []byte
	err = conn.QueryRow(ctx, "SELECT get_user_password($1)", record_code).Scan(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
