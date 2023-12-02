package db

import "fmt"

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
