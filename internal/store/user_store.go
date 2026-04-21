package store

import (
	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	UpsertUser(user *User) error
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
}

type PgUserRepo struct {
	DB *sqlx.DB
}

func NewPgUserRepo(db *sqlx.DB) *PgUserRepo {
	return &PgUserRepo{
		DB: db,
	}
}

func (r *PgUserRepo) UpsertUser(user *User) error {
	return nil
}

func (r *PgUserRepo) GetUserByID(id string) (*User, error) {
	return nil, nil
}

func (r *PgUserRepo) GetUserByEmail(email string) (*User, error) {
	return nil, nil
}
