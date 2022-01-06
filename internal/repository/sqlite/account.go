package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/BenjaminCallahan/my-bank-service/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// Account provides all entities for execute SQL query or transcations
type Account struct {
	db *sql.DB
	*Queries
}

// NewAccount creates a new Account
func NewAccount(db *sql.DB) *Account {
	return &Account{
		db:      db,
		Queries: newQueriesTx(db),
	}
}

// execTx executes function in DB transcation
func (a *Account) execTx(fn func(*Queries) error) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}

	queriesTx := newQueriesTx(tx)
	err = fn(queriesTx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("sqlite: failed to execute transcation: %w, rollback error: %s", err, rbErr.Error())
		}
		return err
	}
	return tx.Commit()
}

// ProcessedTransferWith processed all transfer of users on the current time.
//
// Sets the updated balance for the provided condition
func (a *Account) ProcessedTransferWithFn(fn func(decimal.Decimal, []domain.Transfer) decimal.Decimal) error {
	if err := a.execTx(func(q *Queries) error {
		bankAccount, err := q.GetUserAccount()
		if err != nil {
			return err
		}

		transfersList, err := q.getNotProcessedTransfers()
		if err != nil {
			return err
		}

		updatedBalance := fn(bankAccount.Balance, transfersList)

		if err = q.commitSuccesTransfer(); err != nil {
			return err
		}

		balance, err := q.updateAccountBalance(updatedBalance)
		if err != nil {
			return err
		}

		logrus.Infof("balance after updateAccountBalance %v", balance)
		return nil
	}); err != nil {
		return err
	}
	return nil
}


// UpdateWithFn update user account for the provided condition
func (a *Account) UpdateWithFn(cashOutAmount decimal.Decimal, fn func(amountInBankAccount decimal.Decimal) (bool, error)) error {
	if err := a.execTx(func(q *Queries) error {
		bankAccount, err := q.GetUserAccount()
		if err != nil {
			return err
		}

		hasEnoughMoney, err := fn(bankAccount.Balance)
		if err != nil {
			return err
		}
		if hasEnoughMoney {
			balance, err := q.updateAccountBalance(bankAccount.Balance.Add(cashOutAmount))
			if err != nil {
				return err
			}
			logrus.Infof("balance after UpdateWithCallback %v", balance)
			return nil
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
