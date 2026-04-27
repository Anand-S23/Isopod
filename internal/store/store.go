package store

import (
	"github.com/jmoiron/sqlx"
)

type Store struct {
	DB   *sqlx.DB
	User UserRepo
}

func NewStore(db *sqlx.DB, ur UserRepo) *Store {
	return &Store{
		DB:   db,
		User: ur,
	}
}
