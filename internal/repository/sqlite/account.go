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

// GetUserAccount get all information about user account
func (q *Queries) GetUserAccount() (domain.BankAccount, error) {
	var dest domain.BankAccount
	err := q.dbTx.QueryRow("SELECT * FROM accounts").Scan(&dest.ID, &dest.Balance, &dest.Currency)
	if err != nil {
		return dest, err
	}
	return dest, err
}

// updateAccountBalance update balance of user account
func (q *Queries) updateAccountBalance(amount decimal.Decimal) (decimal.Decimal, error) {
	var balance decimal.Decimal
	if err := q.dbTx.QueryRow("UPDATE accounts SET balance = ? RETURNING balance", amount).Scan(&balance); err != nil {
		return balance, err
	}
	return balance, nil
}

// GetAccountCurrencyRate get account information with currency rate of his account
func (q *Queries) GetAccountCurrencyRate(currency string) (domain.AccExchangeRate, error) {
	var dest domain.AccExchangeRate
	query := `
	SELECT
		accounts.balance,
		currencies.name,
		(
			case
				WHEN ?1 == ''
				OR currencies.name == ?1 THEN 1
				WHEN ?1 == 'RUB' THEN (
					SELECT
						exchange_rate
					FROM
						currency_rate
				)
			END
		) currency_rate
	FROM
		accounts
		INNER JOIN currencies ON accounts.currency_id == currencies.id;
	`
	if err := q.dbTx.QueryRow(query, currency).Scan(&dest.Balance, &dest.CurrencyName, &dest.ExchageRate); err != nil {
		return dest, err
	}

	return dest, nil
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
