package postgres

import (
	"context"
	"os"
	"time"

	"github.com/ellofae/authentication-deanery/config"
	"github.com/ellofae/authentication-deanery/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectionAttemps(conn_func func() error, attemps int, delay time.Duration) (err error) {
	logger := logger.GetLogger()
	for i := 0; i < attemps; i++ {
		err = conn_func()
		if err != nil {
			logger.Printf("Attempting to connect, current attemp: %v, appemps left: %v\n", i+1, attemps-i-1)
			time.Sleep(delay)
			continue
		}
	}
	return err
}

func OpenPoolConnection(ctx context.Context, cfg *config.Config, pgxConfig *pgxpool.Config) (poolConn *pgxpool.Pool) {
	logger := logger.GetLogger()

	err := ConnectionAttemps(func() error {
		var err error

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		poolConn, err = pgxpool.NewWithConfig(ctx, pgxConfig)

		return err
	}, 3, time.Duration(2)*time.Second)

	if err != nil {
		logger.Printf("Didn't manage to make connection with database. Error: %v.\n", err.Error())
		os.Exit(1)
	}

	logger.Println("Database connection is established successfully.")

	return
}
