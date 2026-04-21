package store

import (
	"github.com/jmoiron/sqlx"

	"github.com/Anand-S23/isopod/internal/models"
)

type UserRepo interface {
	Upsert(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
}

type PgUserRepo struct {
	DB *sqlx.DB
}

func NewPgUserRepo(db *sqlx.DB) *PgUserRepo {
	return &PgUserRepo{
		DB: db,
	}
}

func (r *PgUserRepo) Upsert(user *models.User) error {

	return nil
}

func (r *PgUserRepo) GetByID(id string) (*models.User, error) {
	return nil, nil
}

func (r *PgUserRepo) GetByEmail(email string) (*models.User, error) {
	return nil, nil
}
