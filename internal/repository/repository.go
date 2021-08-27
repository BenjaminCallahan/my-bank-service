package repository

import (
	"database/sql"
)

type Account interface {
}

type Repository struct {
	Account
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{}
}
