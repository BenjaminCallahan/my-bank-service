package sqlite

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)


func TestAccountRepository_CreateTransfer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testCases := []struct {
		name         string
		amount       float64
		mockBehavior func(amount float64)
		shouldFail   bool
	}{
		{
			name:   "Ok",
			amount: 0,
			mockBehavior: func(amount float64) {
				mock.ExpectExec("^INSERT INTO transfers (.+) VALUES (.+)$").
					WithArgs(accountID, amount).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			shouldFail: false,
		},
		{
			name:   "Insert Failed",
			amount: 0,
			mockBehavior: func(amount float64) {
				mock.ExpectExec("^INSERT INTO transfers (.+) VALUES (.+)$").
					WithArgs(accountID, amount).WillReturnError(errors.New("failed"))
			},
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.amount)

			queriTx := newQueriesTx(db)

			err := queriTx.CreateTransfer(tc.amount)
			if tc.shouldFail {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestAccountRepository_CommitSuccesTransfer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testCases := []struct {
		name         string
		amount       float64
		mockBehavior func()
		shouldFail   bool
	}{
		{
			name: "Ok",
			mockBehavior: func() {
				mock.ExpectExec("^UPDATE (.+) SET is_processed = \\? WHERE to_account_id == \\?$").
					WithArgs(transferCompleted, accountID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			shouldFail: false,
		},
		{
			name: "Update Failed",
			mockBehavior: func() {
				mock.ExpectExec("^UPDATE (.+) SET is_processed = \\? WHERE to_account_id == \\?$").
					WillReturnError(errors.New("failed"))
			},
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			queriTx := newQueriesTx(db)

			err := queriTx.commitSuccesTransfer()
			if tc.shouldFail {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestAccountRepository_GetNotProcessedTransfers(t *testing.T) {
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
				rowsTransfer := sqlmock.NewRows([]string{"id", "to_account_id", "amount", "is_processed"})
				mock.ExpectQuery("^SELECT (.+) FROM transfers WHERE is_processed == 0$").
					WillReturnRows(rowsTransfer)
			},
			shouldFail: false,
		},
		{
			name: "Select Failed",
			mockBehavior: func() {
				mock.ExpectQuery("^SELECT (.+) FROM transfers WHERE is_processed == 0$").
					WillReturnError(errors.New("failed"))
			},
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			queriTx := newQueriesTx(db)

			_, err := queriTx.getNotProcessedTransfers()
			if tc.shouldFail {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
