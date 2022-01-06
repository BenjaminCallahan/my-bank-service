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

// commitSuccesTransfer mark transfer as a Success
func (q *Queries) commitSuccesTransfer() error {
	_, err := q.dbTx.Exec("UPDATE transfers SET is_processed = ? WHERE to_account_id == ?", transferCompleted, accountID)
	if err != nil {
		return err
	}
	return err
}

// getNotProcessedTransfers get all not processed transfer on the current time
func (q *Queries) getNotProcessedTransfers() ([]domain.Transfer, error) {
	rows, err := q.dbTx.Query("SELECT * FROM transfers WHERE is_processed == 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []domain.Transfer
	for rows.Next() {
		var transfer domain.Transfer
		if err = rows.Scan(
			&transfer.ID,
			&transfer.ToAccountID,
			&transfer.Amount,
			&transfer.IsProcessed,
		); err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transfers, err
}
