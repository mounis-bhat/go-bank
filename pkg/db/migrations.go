package db

func (s *PostgresStorage) autoMigrate() error {
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
