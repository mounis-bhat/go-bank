package db

import (
	"github.com/lib/pq"
	"github.com/mounis-bhat/go-bank/types"
)

func (s *PostgresStorage) GetAccountByUsername(username string) (*types.Account, error) {
	query := `
	SELECT *
	FROM accounts
	WHERE username = $1
	`
	row := s.db.QueryRow(query, username)

	account := &types.Account{}

	if err := row.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Balance, &account.CreatedAt, &account.Username, &account.Password, pq.Array(&account.Roles)); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *PostgresStorage) CreateAccount(account *types.Account) (int, error) {
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
	selectQuery := `
	SELECT account_id
	FROM accounts
	WHERE account_id = $1
	`
	row := s.db.QueryRow(selectQuery, id)

	account := &types.DeleteAccountRequest{}

	if err := row.Scan(&account.ID); err != nil {
		return err
	}

	deleteQuery := `
	DELETE FROM accounts
	WHERE account_id = $1
	`
	_, err := s.db.Exec(deleteQuery, id)

	return err
}

func (s *PostgresStorage) UpdateAccount(account *types.Account) error {
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

func (s *PostgresStorage) GetAccounts() ([]*types.AccountsRequest, error) {
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

	accounts := make([]*types.AccountsRequest, 0)

	for rows.Next() {
		account := &types.AccountsRequest{}
		if err := rows.Scan(&account.FirstName, &account.LastName, &account.Username, &account.CreatedAt, &account.Id, &account.Balance); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
