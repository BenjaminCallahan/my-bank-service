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

func (s *AccountService) AddFunds(amount float64) error {
	return nil
}

func (s *AccountService) sumProfit() error {
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
