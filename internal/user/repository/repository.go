package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db     *sqlx.DB
	dbRead *sqlx.DB
}

func NewRepository(master, replica *sqlx.DB) *Repo {
	return &Repo{
		db:     master,
		dbRead: replica,
	}
}

func (r *Repo) WithTransaction() (*sqlx.Tx, error) {
	return r.db.Beginx()
}
