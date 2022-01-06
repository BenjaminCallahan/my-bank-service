package sqlite

import (
	"github.com/BenjaminCallahan/my-bank-service/internal/domain"
)

const (
	// represent default one user account id.
	accountID = 1

	// completed status of transfer
	transferCompleted = 1
)

// CreateTransfer creates a new transfer in log of transfer
func (q *Queries) CreateTransfer(amount float64) error {
	_, err := q.dbTx.Exec("INSERT INTO transfers (to_account_id, amount) VALUES (?, ?)", accountID, amount)
	if err != nil {
		return err
	}
	return nil
}

