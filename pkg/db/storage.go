package db

import (
	"database/sql"
	"os"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connectionString := os.Getenv("DATABASE_URL")
	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := database.Ping(); err != nil {
		return nil, err
	}

	err = autoMigrate(database)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: database}, nil
}
