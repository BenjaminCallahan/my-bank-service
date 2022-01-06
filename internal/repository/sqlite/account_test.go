package sqlite

import (
	"errors"
	"testing"

	"github.com/BenjaminCallahan/my-bank-service/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// using this Err for simulate any error while during mock object
var ErrFailedMock = errors.New("failed")

func TestAccountRepository_ExecTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testCases := []struct {
		name           string
		fnWithErrOrNot func(*Queries) error
		mockBehavior   func()
		shouldFail     bool
	}{
		{
			name:           "Ok",
			fnWithErrOrNot: func(*Queries) error { return nil },
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			shouldFail: false,
		},
		{
			name:           "Transaction Rollback",
			fnWithErrOrNot: func(*Queries) error { return ErrFailedMock },
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			shouldFail: true,
		},
		{
			name:           "Transaction Begin Err",
			fnWithErrOrNot: func(*Queries) error { return nil },
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectRollback().WillReturnError(ErrFailedMock)
			},
			shouldFail: true,
		},
		{
			name:           "Transaction Rollback failed while trying handle err Queries fn",
			fnWithErrOrNot: func(*Queries) error { return ErrFailedMock },
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectRollback().WillReturnError(ErrFailedMock)
			},
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			account := NewAccount(db)

			err := account.execTx(tc.fnWithErrOrNot)

			if tc.shouldFail {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

		})
	}
}

func TestAccountRepository_ProcessedTransferWith(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testCases := []struct {
		name          string
		updateBalance func(decimal.Decimal, []domain.Transfer) decimal.Decimal
		mockBehavior  func()
		shouldFail    bool
	}{
		{
			name: "Ok",
			updateBalance: func(decimal.Decimal, []domain.Transfer) decimal.Decimal {
				return decimal.NewFromFloat(0)
			},
			mockBehavior: func() {
				mock.ExpectBegin()

				//GetUserAccount
				userRows := mock.NewRows([]string{"id", "balance", "currency_id"}).AddRow(1, 12.23, 2)
				mock.ExpectQuery("^SELECT (.+) FROM accounts$").WillReturnRows(userRows)

				//getNotProcessedTransfers
				transferRows := sqlmock.NewRows([]string{"id", "to_account_id", "amount", "is_processed"})
				mock.ExpectQuery("^SELECT (.+) FROM transfers WHERE is_processed == 0$").
					WillReturnRows(transferRows)

				//commitSuccesTransfer
				mock.ExpectExec("^UPDATE (.+) SET is_processed = \\? WHERE to_account_id == \\?$").
					WithArgs(transferCompleted, accountID).WillReturnResult(sqlmock.NewResult(1, 1))

				//updateAccountBalance
				balanceRow := sqlmock.NewRows([]string{"balance"}).AddRow("0")
				mock.ExpectQuery("^UPDATE (.+) SET balance = \\? RETURNING balance$").WillReturnRows(balanceRow)

				mock.ExpectCommit()
			},
			shouldFail: false,
		},
		{
			name: "Transaction Rollback on Failed in GetUserAccount",
			updateBalance: func(decimal.Decimal, []domain.Transfer) decimal.Decimal {
				return decimal.NewFromFloat(0)
			},
			mockBehavior: func() {
				mock.ExpectBegin()

				//GetUserAccount
				mock.ExpectQuery("^SELECT (.+) FROM accounts$").WillReturnError(ErrFailedMock)

				mock.ExpectRollback()
			},
			shouldFail: true,
		},
		{
			name: "Transaction Rollback on Failed in getNotProcessedTransfers",
			updateBalance: func(decimal.Decimal, []domain.Transfer) decimal.Decimal {
				return decimal.NewFromFloat(0)
			},
			mockBehavior: func() {
				mock.ExpectBegin()

				//GetUserAccount
				userRows := mock.NewRows([]string{"id", "balance", "currency_id"}).AddRow(1, 12.23, 2)
				mock.ExpectQuery("^SELECT (.+) FROM accounts$").WillReturnRows(userRows)

				//getNotProcessedTransfers
				mock.ExpectQuery("^SELECT (.+) FROM transfers WHERE is_processed == 0$").
					WillReturnError(ErrFailedMock)

				mock.ExpectRollback()
			},
			shouldFail: true,
		},
		{
			name: "Transaction Rollback on Failed in commitSuccesTransfer",
			updateBalance: func(decimal.Decimal, []domain.Transfer) decimal.Decimal {
				return decimal.NewFromFloat(0)
			},
			mockBehavior: func() {
				mock.ExpectBegin()

				//GetUserAccount
				userRows := mock.NewRows([]string{"id", "balance", "currency_id"}).AddRow(1, 12.23, 2)
				mock.ExpectQuery("^SELECT (.+) FROM accounts$").WillReturnRows(userRows)

				//getNotProcessedTransfers
				transferRows := sqlmock.NewRows([]string{"id", "to_account_id", "amount", "is_processed"})
				mock.ExpectQuery("^SELECT (.+) FROM transfers WHERE is_processed == 0$").
					WillReturnRows(transferRows)

				//commitSuccesTransfer
				mock.ExpectExec("^UPDATE (.+) SET is_processed = \\? WHERE to_account_id == \\?$").
					WithArgs(transferCompleted, accountID).
					WillReturnError(ErrFailedMock)

				mock.ExpectRollback()
			},
			shouldFail: true,
		},
		{
			name: "Transaction Rollback on Failed in updateAccountBalance",
			updateBalance: func(decimal.Decimal, []domain.Transfer) decimal.Decimal {
				return decimal.NewFromFloat(0)
			},
			mockBehavior: func() {
				mock.ExpectBegin()

				//GetUserAccount
				userRows := mock.NewRows([]string{"id", "balance", "currency_id"}).AddRow(1, 12.23, 2)
				mock.ExpectQuery("^SELECT (.+) FROM accounts$").WillReturnRows(userRows)

				//getNotProcessedTransfers
				transferRows := sqlmock.NewRows([]string{"id", "to_account_id", "amount", "is_processed"})
				mock.ExpectQuery("^SELECT (.+) FROM transfers WHERE is_processed == 0$").
					WillReturnRows(transferRows)

				//commitSuccesTransfer
				mock.ExpectExec("^UPDATE (.+) SET is_processed = \\? WHERE to_account_id == \\?$").
					WithArgs(transferCompleted, accountID).WillReturnResult(sqlmock.NewResult(1, 1))

				//updateAccountBalance
				mock.ExpectQuery("^UPDATE (.+) SET balance = \\? RETURNING balance$").
					WillReturnError(ErrFailedMock)

				mock.ExpectRollback()
			},
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			account := NewAccount(db)

			err := account.ProcessedTransferWithFn(tc.updateBalance)

			if tc.shouldFail {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

		})
	}
}

func TestAccountRepository_GetUserAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testCases := []struct {
		name         string
		mockBehavior func()
		shouldFail   bool
	}{
		{
			name: "Ok",
			mockBehavior: func() {
				rows := mock.NewRows([]string{"id", "balance", "currency_id"}).AddRow(1, 12.23, 2)
				mock.ExpectQuery("^SELECT (.+) FROM accounts$").WillReturnRows(rows)
			},
			shouldFail: false,
		},
		{
			name: "Failed Query",
			mockBehavior: func() {
				mock.ExpectQuery("^SELECT (.+) FROM accounts$").WillReturnError(ErrFailedMock)
			},
			shouldFail: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			account := NewAccount(db)

			_, err := account.GetUserAccount()

			if tt.shouldFail {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestAccountRepository_UpdateWithCallback(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testCases := []struct {
		name             string
		input            decimal.Decimal
		hasEnoughMoneyFn func(amountInBankAccount decimal.Decimal) (bool, error)
		mockBehavior     func()
		shouldFail       bool
	}{
		{
			name:  "Ok",
			input: decimal.NewFromFloat(1),
			hasEnoughMoneyFn: func(amountInBankAccount decimal.Decimal) (bool, error) {
				return true, nil
			},
			mockBehavior: func() {
				mock.ExpectBegin()

				//GetUserAccount
				accountRows := mock.NewRows([]string{"id", "balance", "currency_id"}).AddRow(0, 12.23, 2)
				mock.ExpectQuery("^SELECT (.+) FROM accounts$").WillReturnRows(accountRows).RowsWillBeClosed()

				// updateAccountBalance
				balanceAccountRow := mock.NewRows([]string{"balance"}).AddRow(13.23)
				mock.ExpectQuery("^UPDATE accounts SET balance = \\? RETURNING balance$").WithArgs(decimal.NewFromFloat(13.23)).WillReturnRows(balanceAccountRow)

				mock.ExpectCommit()
			},
			shouldFail: false,
		},
		{
			name:  "Transaction Rollback on Failed in GetUserAccount",
			input: decimal.NewFromFloat(1),
			hasEnoughMoneyFn: func(amountInBankAccount decimal.Decimal) (bool, error) {
				return true, nil
			},
			mockBehavior: func() {
				mock.ExpectBegin()

				//GetUserAccount
				mock.ExpectQuery("^SELECT (.+) FROM accounts$").WillReturnError(ErrFailedMock)

				mock.ExpectRollback()
			},
			shouldFail: true,
		},
		{
			name:  "Transaction Rollback on Failed in updateAccountBalance",
			input: decimal.NewFromFloat(1),
			hasEnoughMoneyFn: func(amountInBankAccount decimal.Decimal) (bool, error) {
				return true, nil
			},
			mockBehavior: func() {
				mock.ExpectBegin()

				//GetUserAccount
				accountRows := mock.NewRows([]string{"id", "balance", "currency_id"}).AddRow(0, 12.23, 2)
				mock.ExpectQuery("^SELECT (.+) FROM accounts$").WillReturnRows(accountRows).RowsWillBeClosed()

				// updateAccountBalance
				mock.ExpectQuery("^UPDATE accounts SET balance = \\? RETURNING balance$").WithArgs(decimal.NewFromFloat(13.23)).
					WillReturnError(ErrFailedMock)

				mock.ExpectRollback()
			},
			shouldFail: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			account := NewAccount(db)

			err := account.UpdateWithFn(tt.input, tt.hasEnoughMoneyFn)

			if tt.shouldFail {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
