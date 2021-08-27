package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"

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

	if err = initDBIfNotExists(db); err != nil {
		return db, fmt.Errorf("sqlite: failed to initialize db schema: %w", err)
	}
	return db, nil
}

func initDBIfNotExists(db *sql.DB) error {
	dbDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	fileSource, err := (&file.File{}).Open("file://schema")
	if err != nil {
		return err
	}

	migrateInstance, err := migrate.NewWithInstance(
		"file", fileSource,
		"sqlite3", dbDriver,
	)
	if err != nil {
		return err
	}

	err = migrateInstance.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
