package repository

import (
	"database/sql"

	"github.com/BenjaminCallahan/my-bank-service/internal/domain"
	"github.com/BenjaminCallahan/my-bank-service/internal/repository/sqlite"
	"github.com/shopspring/decimal"
)

// Account wraps the methods for working with account repository
type Account interface {
	// CreateTransfer creates a new transfer in log of transfer
	CreateTransfer(float64) error
	// ProcessedTransferWithFn processed all transfer of users on the current time.
	// Sets the updated balance for the provided condition
	ProcessedTransferWithFn(func(decimal.Decimal, []domain.Transfer) decimal.Decimal) error
	// GetAccountCurrencyRate get account information with currency rate of his account
	GetAccountCurrencyRate(string) (domain.AccExchangeRate, error)
	// GetUserAccount get all information about user account
	GetUserAccount() (domain.BankAccount, error)
	// UpdateWithFn update user account for the provided condition
	UpdateWithFn(decimal.Decimal, func(amountInBankAccount decimal.Decimal) (bool, error)) error
}

// Repository to interract with repository Bank Account
type Repository struct {
	Account
}

// NewRepository create new Repository to interract with Bank Account
func NewRepository(db *sql.DB) *Repository {
	return &Repository{sqlite.NewAccount(db)}
}
