package service

import (
	"github.com/BenjaminCallahan/my-bank-service/internal/repository"
)

type AccountService struct {
	repo repository.Account
}

func NewAccountService(repo repository.Account) *AccountService {
	return &AccountService{repo: repo}
}

// AddFunds Allows you to deposit the amount sum
func (s *AccountService) AddFunds(amount float64) error {
	if err := s.repo.CreateTransfer(amount); err != nil {
		return err
	}

	if err := s.sumProfit(); err != nil {
		return err
	}
	return nil
}

// sumProfit Calculates the interest on the deposit and deposits the received money into the account
func (s *AccountService) sumProfit() error {
	if err := s.repo.ProcessedTransferWithFn(func(balance decimal.Decimal, transfers []domain.Transfer) decimal.Decimal {
		for _, transfer := range transfers {
			balance = balance.Add(transfer.Amount)
			accumulationAmount := balance.Mul(decimal.NewFromFloat(accumPercent)).Round(fractPartSBP)
			balance = balance.Add(accumulationAmount)
		}
		return balance
	}); err != nil {
		return err
	}
	return nil
}

// Withdraw Performs withdrawals from the account according to the specified rules.
// If the write-off goes beyond the rules, it gives an error
func (s *AccountService) Withdraw(cashOutAmount float64) error {
	decCashOutAmount := decimal.NewFromFloat(cashOutAmount)
	withdrawCondition := func(amountInBankAccount decimal.Decimal) (bool, error) {
		maxAllowedAmountCash := amountInBankAccount.Mul(decimal.NewFromFloat(maxWithdrawAmountPercent))

		// as the condition requires
		// not more than 70% of the amount on the account
		if !decCashOutAmount.LessThanOrEqual(maxAllowedAmountCash) {
			return false, domain.ErrNotEnoughMoney
		}
		return true, nil
	}
	err := s.repo.UpdateWithFn(decCashOutAmount.Neg(), withdrawCondition)
	if err != nil {
		return err
	}
	return nil
}

// GetCurrency Returns the account currency
func (s *AccountService) GetCurrency() (string, error) {
	userBankAcoount, err := s.repo.GetUserAccount()
	if err != nil {
		return "", err
	}

	return userBankAcoount.Currency, nil
}

// GetAccountCurrencyRate Returns the account currency rate to the transmitted currency cur
func (s *AccountService) GetAccountCurrencyRate(currency string) (float64, error) {
	return 0.0, nil
}

// GetBalance Returns the account balance in the specified currency
func (s *AccountService) GetBalance(currency string) (float64, error) {
	accountWithCurrencyRate, err := s.repo.GetAccountCurrencyRate(currency)
	if err != nil {
		return 0, err
	}

	balanceInf64, _ := accountWithCurrencyRate.Balance.
		Mul(accountWithCurrencyRate.ExchageRate).
		Round(fractPartSBP).
		Float64()
	return balanceInf64, nil
}
