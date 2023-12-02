package db

import (
	"database/sql"
	"os"

	"github.com/mounis-bhat/go-bank/types"
)

type Storage interface {
	CreateAccount(*types.Account) (int, error)
	DeleteAccount(int) error
	UpdateAccount(*types.Account) error
	GetAccounts() ([]*types.AccountsRequest, error)
	TransferMoney(int, int, int) error
	GetAccountByUsername(string) (*types.Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connectionString := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func InitStorage(s *PostgresStorage) error {
	return s.autoMigrate()
}
