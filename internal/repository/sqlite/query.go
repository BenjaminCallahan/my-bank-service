package sqlite

import (
	"database/sql"
)

// DBTX implemented by sql.DB and sql.Tx
type DBTX interface {
	// Query executes a query that returns rows
	Query(query string, args ...interface{}) (*sql.Rows, error)
	// QueryRow executes a query that is expected to return at most one row
	QueryRow(query string, args ...interface{}) *sql.Row
	// Exec executes a query without returning any rows
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Queries to interact with DB transactions
type Queries struct {
	dbTx DBTX
}

// newQueriesTx creates new Queries to interract with DB transcactions
func newQueriesTx(ext DBTX) *Queries {
	return &Queries{
		dbTx: ext,
	}
}
