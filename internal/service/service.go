package service

type BankAccount interface {
	// AddFunds Allows you to deposit the amount sum
	AddFunds(sum float64) error
	// sumProfit Calculates the interest on the deposit and deposits the received money into the account
	sumProfit() error
	// Withdraw Performs withdrawals from the account according to the specified rules.
	// If the write-off goes beyond the rules, it gives an error
	Withdraw(sum float64) error
	// GetCurrency Returns the account currency
	GetCurrency() (string, error)
	// GetAccountCurrencyRate Returns the account currency rate to the transmitted currency cur
	GetAccountCurrencyRate(cur string) (float64, error)
	// GetBalance Returns the account balance in the specified currency
	GetBalance(cur string) (float64, error)
}
