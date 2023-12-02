package db

import "database/sql"

func autoMigrate(database *sql.DB) error {
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

	_, err := database.Exec(query)
	return err
}
