package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BenjaminCallahan/my-bank-service/internal/domain"
	"github.com/BenjaminCallahan/my-bank-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_deposit(t *testing.T) {

	tests := []struct {
		name                 string
		inputBody            string
		shouldServiceFail    bool
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Ok",
			inputBody:            `{"amount":12}`,
			shouldServiceFail:    false,
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			name:                 "Invalid Amount",
			inputBody:            `{"amount":"1"}`,
			shouldServiceFail:    false,
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Service Failure",
			inputBody:            `{"amount":1}`,
			shouldServiceFail:    true,
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to add amount"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockedBankAccountService := &service.BankAccountMock{AddFundsFunc: func(sum float64) error {
				if tt.shouldServiceFail {
					return errors.New("failed")
				}
				return nil
			}}

			services := &service.Service{mockedBankAccountService}
			handler := NewHandler(services)

			g := gin.New()
			g.POST("/deposit", handler.deposit)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/deposit", bytes.NewBufferString(tt.inputBody))

			// make request
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}

}

func TestHandler_withdraw(t *testing.T) {

	tests := []struct {
		name                     string
		inputBody                string
		mockedBankAccountService *service.BankAccountMock
		expectedStatusCode       int
		expectedResponseBody     string
	}{
		{
			name:      "Ok",
			inputBody: `{"amount":12}`,
			mockedBankAccountService: &service.BankAccountMock{WithdrawFunc: func(sum float64) error {
				return nil
			}},
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			name:      "Invalid Amount",
			inputBody: `{"amount":"1"}`,
			mockedBankAccountService: &service.BankAccountMock{WithdrawFunc: func(sum float64) error {
				return nil
			}},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"amount":1}`,
			mockedBankAccountService: &service.BankAccountMock{WithdrawFunc: func(sum float64) error {
				return errors.New("failed")
			}},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to withdraw amount"}`,
		},
		{
			name:      "Not Enough Money",
			inputBody: `{"amount":1}`,
			mockedBankAccountService: &service.BankAccountMock{WithdrawFunc: func(sum float64) error {
				return domain.ErrNotEnoughMoney
			}},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"not enough money"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			services := &service.Service{tt.mockedBankAccountService}
			handler := NewHandler(services)

			g := gin.New()
			g.POST("/withdraw", handler.withdraw)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBufferString(tt.inputBody))

			// make request
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}

}

func TestHandler_balance(t *testing.T) {
	tests := []struct {
		name                 string
		currency             string
		shouldServiceFail    bool
		balance              float64
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Ok",
			currency:             "",
			shouldServiceFail:    false,
			balance:              1,
			expectedStatusCode:   200,
			expectedResponseBody: `1`,
		},
		{
			name:                 "Ok",
			currency:             "SBP",
			shouldServiceFail:    false,
			balance:              2,
			expectedStatusCode:   200,
			expectedResponseBody: `2`,
		},
		{
			name:                 "Ok",
			currency:             "RUB",
			shouldServiceFail:    false,
			balance:              1,
			expectedStatusCode:   200,
			expectedResponseBody: `1`,
		},
		{
			name:                 "Not correct currency",
			currency:             "c",
			shouldServiceFail:    false,
			balance:              0,
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Service Failure",
			currency:             "",
			shouldServiceFail:    true,
			balance:              0,
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get balance"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockedBankAccountService := &service.BankAccountMock{GetBalanceFunc: func(cur string) (float64, error) {
				if tt.shouldServiceFail {
					return tt.balance, errors.New("failed")
				}
				return tt.balance, nil
			}}

			services := &service.Service{mockedBankAccountService}
			handler := NewHandler(services)

			g := gin.New()
			g.GET("/balance", handler.balance)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/balance?currency=%s", tt.currency), http.NoBody)

			// make request
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_currency(t *testing.T) {

	tests := []struct {
		name                 string
		currency             string
		shouldServiceFail    bool
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Ok",
			currency:             "SBP",
			shouldServiceFail:    false,
			expectedStatusCode:   200,
			expectedResponseBody: `"SBP"`,
		},
		{
			name:                 "Service Failure",
			currency:             "SBP",
			shouldServiceFail:    true,
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get currency"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockedBankAccountService := &service.BankAccountMock{GetCurrencyFunc: func() (string, error) {
				if tt.shouldServiceFail {
					return tt.currency, errors.New("failed")
				}
				return tt.currency, nil
			}}

			services := &service.Service{mockedBankAccountService}
			handler := NewHandler(services)

			g := gin.New()
			g.GET("/currency", handler.currency)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/currency", http.NoBody)

			// make request
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}

}
