package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (int, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*AccountsRequest, error)
	TransferMoney(int, int, int) error
	GetAccountByUsername(string) (*Account, error)
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
	return s.createAccountTable()
}

func (s *PostgresStorage) createAccountTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS accounts (
			account_id SERIAL PRIMARY KEY,
			first_name VARCHAR(50) NOT NULL,
			last_name VARCHAR(50) NOT NULL,
			balance INTEGER NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			username VARCHAR(50) NOT NULL UNIQUE,
			password VARCHAR(100) NOT NULL,
			roles VARCHAR(10)[]
	)
	`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStorage) GetAccountByUsername(username string) (*Account, error) {
	row := s.db.QueryRow("SELECT * FROM accounts WHERE username = $1", username)

	account := &Account{}

	if err := row.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Balance, &account.CreatedAt, &account.Username, &account.Password, pq.Array(&account.Roles)); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *PostgresStorage) CreateAccount(account *Account) (int, error) {
	query := `
	INSERT INTO accounts (
			first_name,
			last_name,
			balance,
			created_at,
			username,
			password,
			roles
	)
	VALUES
			($1, $2, $3, $4, $5, $6, $7)
	RETURNING account_id
	`

	var accountID int
	err := s.db.QueryRow(query, account.FirstName, account.LastName, account.Balance, account.CreatedAt, account.Username, account.Password, pq.Array(account.Roles)).Scan(&accountID)
	if err != nil {
		return 0, err
	}

	account.ID = accountID
	return accountID, nil
}

func (s *PostgresStorage) DeleteAccount(id int) error {
	row := s.db.QueryRow("SELECT * FROM accounts WHERE account_id = $1", id)

	account := &Account{}

	if err := row.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Balance, &account.CreatedAt); err != nil {
		return err
	}

	_, err := s.db.Exec("DELETE FROM accounts WHERE account_id = $1", id)

	return err
}

func (s *PostgresStorage) UpdateAccount(account *Account) error {
	query := `
	UPDATE 
		accounts 
	SET 
		first_name = $1, 
		last_name = $2, 
		balance = $3 
	WHERE 
		account_id = $4
	`

	_, err := s.db.Exec(query, account.FirstName, account.LastName, account.Balance, account.ID)

	return err
}

func (s *PostgresStorage) GetAccounts() ([]*AccountsRequest, error) {
	query := `
	SELECT
		first_name,
		last_name,
		username,
		created_at,
		account_id,
		balance
	FROM
		accounts
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*AccountsRequest, 0)

	for rows.Next() {
		account := &AccountsRequest{}
		if err := rows.Scan(&account.FirstName, &account.LastName, &account.Username, &account.CreatedAt, &account.Id, &account.Balance); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStorage) TransferMoney(fromAccountID, toAccountID, amount int) error {
	result, err := s.db.Exec(`
	UPDATE accounts
	SET balance =
		CASE
			WHEN account_id = $1 AND balance >= $3 THEN balance - $3
			WHEN account_id = $2 THEN balance + $3
		END
	WHERE
		account_id IN ($1, $2)
		AND EXISTS (SELECT 1 FROM accounts WHERE account_id = $2)
	`, fromAccountID, toAccountID, amount)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transfer failed")
	}

	return nil
}
