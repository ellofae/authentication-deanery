package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/ellofae/authentication-deanery/config"
	"github.com/ellofae/authentication-deanery/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

var migrations_path string = "file://migrations"

func parseConnectionString(cfg *config.Config) string {
	poolConnString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.UserDatabase.User,
		cfg.UserDatabase.Password,
		cfg.UserDatabase.Host,
		cfg.UserDatabase.Port,
		cfg.UserDatabase.DBName,
		cfg.UserDatabase.SSLmode)

	return poolConnString
}

func GetPoolParseConfig(cfg *config.Config) *pgxpool.Config {
	logger := logger.GetLogger()

	pgxConfig, err := pgxpool.ParseConfig(parseConnectionString(cfg))
	if err != nil {
		logger.Printf("Unable to parse connection string. Error: %v.", err.Error())
		os.Exit(1)
	}

	return pgxConfig
}

func RunMigrationsUp(ctx context.Context, cfg *config.Config) {
	logger := logger.GetLogger()

	connString := parseConnectionString(cfg)
	migration, err := migrate.New(migrations_path, connString)
	if err != nil {
		logger.Printf("Unable to get a migrate instance. Error: %v.\n", err.Error())
		os.Exit(1)
	}

	err = migration.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			logger.Println("No changes while migrating.")
			return
		}

		logger.Printf("Unable to migrate up. Error: %v.\n", err.Error())
		os.Exit(1)
	}
	logger.Println("Migrations are up successfully.")
}
