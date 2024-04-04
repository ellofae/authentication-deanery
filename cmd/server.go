package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/ellofae/authentication-deanery/config"
	"github.com/ellofae/authentication-deanery/internal/controller/handler"
	"github.com/ellofae/authentication-deanery/internal/controller/middleware"
	"github.com/ellofae/authentication-deanery/internal/database/repository"
	"github.com/ellofae/authentication-deanery/internal/domain/usecase"
	"github.com/ellofae/authentication-deanery/internal/models"
	"github.com/ellofae/authentication-deanery/migrations/initialization"
	"github.com/ellofae/authentication-deanery/pkg/logger"
	"github.com/ellofae/authentication-deanery/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func castPasswordLengthCgfTypes(cfg_str string) uint8 {
	logger := logger.GetLogger()

	num, err := strconv.ParseUint(cfg_str, 10, 8)
	if err != nil {
		logger.Printf("Unable to parse the config password length value to a numeric representation, error: %v\n", err)
		os.Exit(1)
	}

	if num > 255 {
		logger.Printf("Config password length value must be of type uint8, error: %v\n", err)
		os.Exit(1)
	}

	return uint8(num)
}

func establishHandlers(mux *http.ServeMux, connPool *pgxpool.Pool, cfg *config.Config) {

	user_repository := repository.NewUserRepository(connPool)
	user_usecase := usecase.NewUserUsecase(user_repository, &models.CfgUsecaseData{
		PasswordLength:   castPasswordLengthCgfTypes(cfg.Encryption.PasswordLength),
		AesEncryptionKey: cfg.Encryption.AesEncryptionKey,
	})
	user_handler := handler.NewUserHandler(user_usecase)

	user_handler.RegisterHandlers(mux)
}

func main() {
	logger := logger.GetLogger()
	cfg := config.ParseConfig(config.ConfigureViper())
	ctx := context.Background()

	connPool := postgres.OpenPoolConnection(ctx, cfg, postgres.GetPoolParseConfig(cfg))
	if err := connPool.Ping(ctx); err != nil {
		logger.Printf("Unable to ping the database connection. Error: %v.\n", err.Error())
		os.Exit(1)
	}
	postgres.RunMigrationsUp(ctx, cfg)

	err := initialization.InitializeDatabse(connPool)
	if err != nil {
		os.Exit(1)
	}

	idleTimeout, _ := strconv.Atoi(cfg.UserService.IdleTimeout)
	readTimeout, _ := strconv.Atoi(cfg.UserService.ReadTimeout)
	writeTimeout, _ := strconv.Atoi(cfg.UserService.WriteTimeout)

	middleware.InitJWTSecretKey(cfg.Authentication.JWTSecretKey)

	serveMux := http.NewServeMux()
	establishHandlers(serveMux, connPool, cfg)

	srv := &http.Server{
		Addr:         cfg.UserService.BindAddr,
		IdleTimeout:  time.Minute * time.Duration(idleTimeout),
		ReadTimeout:  time.Minute * time.Duration(readTimeout),
		WriteTimeout: time.Minute * time.Duration(writeTimeout),

		Handler: http.TimeoutHandler(middleware.RequestMiddleware(serveMux), 2*time.Minute, ""),
	}

	done := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		for {
			if <-c == os.Interrupt {
				if err := srv.Shutdown(context.Background()); err != nil {
					logger.Printf("Error while shutting down the server occured. Error: %v.\n", err.Error())
				}
				close(done)
				return
			}
		}
	}()

	err = srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logger.Printf("Error while starting the server occured. Error: %v.\n", err.Error())
	}
	<-done

	logger.Println("Server gracefully shutdown.")
}
