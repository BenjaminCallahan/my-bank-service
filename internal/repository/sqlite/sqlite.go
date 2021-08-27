package sqlite

import (
	"database/sql"

	// package for using sqlite
	_ "modernc.org/sqlite"
)

type Config struct {
	DBName string
}

func NewConnectDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DBName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
