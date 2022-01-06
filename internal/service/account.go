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

func (s *AccountService) Withdraw(cashOutAmount float64) error {
	return nil
}

func (s *AccountService) GetCurrency() (string, error) {
	return "", nil
}

func (s *AccountService) GetAccountCurrencyRate(currency string) (float64, error) {
	return 0.0, nil
}

func (s *AccountService) GetBalance(currency string) (float64, error) {
	return 0.0, nil
}
