package domain

import (
	"errors"

	"github.com/shopspring/decimal"
)

// ErrNotEnoughMoney is returned if a user doesn't have enough money in a bank account
var ErrNotEnoughMoney = errors.New("not enough money")

type Transfer struct {
	ID          int             `json:"id"`
	ToAccountID int             `json:"to_account_id"`
	Amount      decimal.Decimal `json:"amount"`
	IsProcessed int             `json:"is_processed"`
}

type BankAccount struct {
	ID       int             `json:"id"`
	Balance  decimal.Decimal `json:"balance"`
	Currency string          `json:"Currency"`
}

type AccExchangeRate struct {
	Balance      decimal.Decimal `json:"balance"`
	CurrencyName string          `json:"currency_name"`
	ExchageRate  decimal.Decimal `json:"exhchange_rate"`
}
