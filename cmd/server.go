package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/ellofae/authentication-deanery/config"
	"github.com/ellofae/authentication-deanery/internal/controller/middleware"
	"github.com/ellofae/authentication-deanery/migrations/initialization"
	"github.com/ellofae/authentication-deanery/pkg/logger"
	"github.com/ellofae/authentication-deanery/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func establishHandlers(mux *http.ServeMux, connPool *pgxpool.Pool) {

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

	serveMux := http.NewServeMux()
	establishHandlers(serveMux, connPool)

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
