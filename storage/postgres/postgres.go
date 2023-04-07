package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/config"
)

func InitPostgresDBConn(config *config.Config) (*sqlx.DB, error) {
	// 	Example DSN
	// user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10

	// Example URL
	// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10

	// Current URL
	// postgres://postgres:postgres@postgresdb:5432/postgresdb

	url := fmt.Sprintf("postgres://%s:%s@%s%s/%s", config.Database.DBUser, config.Database.DBPass, config.Database.DBHost, config.Database.DBPort, config.Database.DBName)
	// dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", config.Database.DBUser, config.Database.DBPass, config.Database.DBHost, config.Database.DBPort, config.Database.DBName)

	db, err := sqlx.Connect("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	// TODO: call function to insert predefined data into database

	return db, nil
}
