package store

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Anand-S23/isopod/internal/models"
)

type UserRepo interface {
	Upsert(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type PgUserRepo struct {
	DB *sqlx.DB
}

func NewPgUserRepo(db *sqlx.DB) *PgUserRepo {
	return &PgUserRepo{
		DB: db,
	}
}

const upsertUserQuery = `
    INSERT INTO users (id, github_id, username, email, avatar_url, access_token, created_at, last_login_at)
    VALUES (:id, :github_id, :username, :email, :avatar_url, :access_token, :created_at, :last_login_at)
    ON CONFLICT (id) DO UPDATE SET
        github_id     = EXCLUDED.github_id,
        username      = EXCLUDED.username,
        email         = EXCLUDED.email,
        avatar_url    = EXCLUDED.avatar_url,
        access_token  = EXCLUDED.access_token,
        last_login_at = EXCLUDED.last_login_at
`

func (r *PgUserRepo) Upsert(ctx context.Context, user *models.User) error {
	_, err := r.DB.NamedExecContext(ctx, upsertUserQuery, user)
	if err != nil {
		return fmt.Errorf("upserting user %s: %w", user.ID, err)
	}
	return nil
}

const userColumns = `
	id,
	github_id,
	username,
	email,
	avatar_url,
	access_token,
	created_at,
	last_login_at
`

func (r *PgUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT %s FROM users WHERE id = $1", userColumns)
	err := r.DB.GetContext(ctx, &user, query, id)
	return &user, err
}

func (r *PgUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT %s FROM users WHERE email = $1", userColumns)
	err := r.DB.GetContext(ctx, &user, query, email)
	return &user, err
}
