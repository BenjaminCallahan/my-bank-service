package service

import (
	"errors"
	"testing"

	"github.com/BenjaminCallahan/my-bank-service/internal/domain"
	"github.com/BenjaminCallahan/my-bank-service/internal/repository/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// using this Err for simulate any error while during mock object
var ErrFailedMock = errors.New("failed")

func TestService_AddFunds(t *testing.T) {

	tests := []struct {
		name             string
		arg              float64
		mockBehaviorRepo *mock.AccountMock
		wantErr          bool
	}{
		{
			name:             "Ok",
			mockBehaviorRepo: &mock.AccountMock{},
			wantErr:          false,
		},
		{
			name: "Err when trying to calculate sumProfit",
			mockBehaviorRepo: &mock.AccountMock{ProcessedTransferWithFnFunc: func(fn func(decimal.Decimal, []domain.Transfer) decimal.Decimal) error {
				return ErrFailedMock
			}},
			wantErr: true,
		},
		{
			name: "Err when trying to CreateTransfer",
			mockBehaviorRepo: &mock.AccountMock{CreateTransferFunc: func(float64) error {
				return ErrFailedMock
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountService := NewAccountService(tt.mockBehaviorRepo)

			err := accountService.AddFunds(tt.arg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_sumProfit(t *testing.T) {

	tests := []struct {
		name             string
		mockBehaviorRepo func(expectedBalance decimal.Decimal, balance decimal.Decimal, transfers []domain.Transfer) *mock.AccountMock
		wantErr          bool
		balance          decimal.Decimal
		expectedBalance  decimal.Decimal
		transfers        []domain.Transfer
	}{
		{
			name: "Ok",
			mockBehaviorRepo: func(expectedBalance, balance decimal.Decimal, transfers []domain.Transfer) *mock.AccountMock {
				return &mock.AccountMock{ProcessedTransferWithFnFunc: func(fn func(balance decimal.Decimal, transfers []domain.Transfer) decimal.Decimal) error {
					updatedBalance := fn(balance, transfers)
					assert.Truef(t, updatedBalance.Equal(expectedBalance), "updated balance equal to: %v want: %v", updatedBalance, expectedBalance)
					return nil
				},
				}
			},
			wantErr:         false,
			balance:         decimal.NewFromFloat(0),
			expectedBalance: decimal.NewFromFloat(14.26),
			transfers: []domain.Transfer{{
				Amount: decimal.NewFromFloat(13.45),
			}},
		},
		{
			name: "Ok",
			mockBehaviorRepo: func(expectedBalance, balance decimal.Decimal, transfers []domain.Transfer) *mock.AccountMock {
				return &mock.AccountMock{ProcessedTransferWithFnFunc: func(fn func(balance decimal.Decimal, transfers []domain.Transfer) decimal.Decimal) error {
					updatedBalance := fn(balance, transfers)
					assert.Truef(t, updatedBalance.Equal(expectedBalance), "updated balance equal to: %v want: %v", updatedBalance, expectedBalance)
					return nil
				},
				}
			},
			wantErr:         false,
			balance:         decimal.NewFromFloat(0),
			expectedBalance: decimal.NewFromFloat(14.26),
			transfers: []domain.Transfer{{
				Amount: decimal.NewFromFloat(13.45),
			}},
		},
		{
			name: "Err occured when tryin to repo.ProcessedTransferWith",
			mockBehaviorRepo: func(expectedBalance, balance decimal.Decimal, transfers []domain.Transfer) *mock.AccountMock {
				return &mock.AccountMock{ProcessedTransferWithFnFunc: func(fn func(balance decimal.Decimal, transfers []domain.Transfer) decimal.Decimal) error {
					updatedBalance := fn(balance, transfers)
					assert.Truef(t, updatedBalance.Equal(expectedBalance), "updated balance equal to: %v want: %v", updatedBalance, expectedBalance)
					return ErrFailedMock
				},
				}
			},
			wantErr:         true,
			balance:         decimal.NewFromFloat(0),
			expectedBalance: decimal.NewFromFloat(14.26),
			transfers: []domain.Transfer{{
				Amount: decimal.NewFromFloat(13.45),
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockedAccountRepo := tt.mockBehaviorRepo(tt.expectedBalance, tt.balance, tt.transfers)

			accountService := NewAccountService(mockedAccountRepo)

			err := accountService.sumProfit()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestService_Withdraw(t *testing.T) {

	tests := []struct {
		name                string
		amountInBankAccount decimal.Decimal
		cashOutAmount       float64
		wantErr             bool
	}{
		{
			name:                "Ok",
			amountInBankAccount: decimal.NewFromFloat(40),
			cashOutAmount:       16.23,
			wantErr:             false,
		},
		{
			name:                "Not enough money",
			amountInBankAccount: decimal.NewFromFloat(40),
			cashOutAmount:       28.1,
			wantErr:             true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockedAccountRepo := &mock.AccountMock{
				UpdateWithFnFunc: func(cachOutAmount decimal.Decimal, fn func(amountInBankAccount decimal.Decimal) (bool, error)) error {
					isAllowed, err := fn(tt.amountInBankAccount)
					if tt.wantErr {
						assert.False(t, isAllowed)
					} else {
						assert.True(t, isAllowed)
					}
					if err != nil {
						return err
					}

					return nil
				}}

			accountService := NewAccountService(mockedAccountRepo)

			err := accountService.Withdraw(tt.cashOutAmount)
			if tt.wantErr {
				assert.EqualError(t, err, domain.ErrNotEnoughMoney.Error())
			} else {
				assert.NoError(t, err)
			}
		})

	}

}

func TestService_GetCurrency(t *testing.T) {

	tests := []struct {
		name             string
		mockBehaviorRepo *mock.AccountMock
		wantErr          bool
	}{
		{
			name:             "Ok",
			mockBehaviorRepo: &mock.AccountMock{},
			wantErr:          false,
		},
		{
			name: "Err occured when tryin to repo.GetUserAccount",
			mockBehaviorRepo: &mock.AccountMock{GetUserAccountFunc: func() (domain.BankAccount, error) {
				return domain.BankAccount{}, ErrFailedMock
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			accountService := NewAccountService(tt.mockBehaviorRepo)

			_, err := accountService.GetCurrency()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetAccountCurrencyRate(t *testing.T) {

	tests := []struct {
		name             string
		mockBehaviorRepo *mock.AccountMock
		currency         string
		wantErr          bool
	}{
		{
			name:             "Ok",
			mockBehaviorRepo: &mock.AccountMock{},
			currency:         "RUB",
			wantErr:          false,
		},
		{
			name: "Err occured when tryin to repo.GetAccountCurrencyRate",
			mockBehaviorRepo: &mock.AccountMock{GetAccountCurrencyRateFunc: func(s string) (domain.AccExchangeRate, error) {
				return domain.AccExchangeRate{}, ErrFailedMock
			}},
			currency: "RUB",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountService := NewAccountService(tt.mockBehaviorRepo)

			_, err := accountService.GetAccountCurrencyRate(tt.currency)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetBalance(t *testing.T) {

	tests := []struct {
		name                      string
		destCurrency              string
		expectedBalance           float64
		mockedAccountRepoBehavior *mock.AccountMock
		wantErr                   bool
	}{
		{
			name:            "Balance in SBP",
			destCurrency:    "SBP",
			expectedBalance: 20.23,
			mockedAccountRepoBehavior: &mock.AccountMock{GetAccountCurrencyRateFunc: func(currency string) (domain.AccExchangeRate, error) {
				return domain.AccExchangeRate{
					Balance:     decimal.NewFromFloat(20.23),
					ExchageRate: decimal.NewFromFloat(1),
				}, nil
			}},
			wantErr: false,
		},
		{
			name:            "Balance in RUB",
			destCurrency:    "RUB",
			expectedBalance: 15.22,
			mockedAccountRepoBehavior: &mock.AccountMock{GetAccountCurrencyRateFunc: func(currency string) (domain.AccExchangeRate, error) {
				return domain.AccExchangeRate{
					Balance:     decimal.NewFromFloat(20.23),
					ExchageRate: decimal.NewFromFloat(0.7523),
				}, nil

			}},
			wantErr: false,
		},
		{
			name:            "Failed to get balance from AccountRepo",
			destCurrency:    "",
			expectedBalance: 0,
			mockedAccountRepoBehavior: &mock.AccountMock{GetAccountCurrencyRateFunc: func(currency string) (domain.AccExchangeRate, error) {
				return domain.AccExchangeRate{}, ErrFailedMock
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountService := NewAccountService(tt.mockedAccountRepoBehavior)

			balance, err := accountService.GetBalance(tt.destCurrency)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expectedBalance, balance)
				assert.NoError(t, err)
			}
		})
	}
}
