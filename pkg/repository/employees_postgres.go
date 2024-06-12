package repository

import (
	"github.com/jmoiron/sqlx"
)

type EmployeesPostgres struct {
	db *sqlx.DB
}

func NewEmployeesPostgres(db *sqlx.DB) *EmployeesPostgres {
	return &EmployeesPostgres{db: db}
}
