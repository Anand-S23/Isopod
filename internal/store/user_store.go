package store

import (
	"context"
	"fmt"

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

func (r *PgUserRepo) Upsert(ctx context.Context, user *models.User) error {
	stmt, err := r.DB.PreparexContext(
		ctx,
		`INSERT INTO users (id, github_id, username, email, avatar_url, access_token, created_at, last_login_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET
			github_id = EXCLUDED.github_id,
			username = EXCLUDED.username,
			email = EXCLUDED.email,
			avatar_url = EXCLUDED.avatar_url,
			access_token = EXCLUDED.access_token,
			last_login_at = EXCLUDED.last_login_at
		WHERE users.id = EXCLUDED.id`,
	)
	if err != nil {
		return fmt.Errorf("prepare upsert statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		user.ID,
		user.GitHubID,
		user.Username,
		user.Email,
		user.AvatarURL,
		user.AccessToken,
		user.CreatedAt,
		user.LastLoginAt,
	)
	if err != nil {
		return fmt.Errorf("execute upsert statement: %w", err)
	}

	return nil
}

func (r *PgUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	stmt, err := r.DB.PrepareContext(
		ctx,
		`SELECT id, github_id, username, email, avatar_url, access_token, created_at, last_login_at 
		FROM users WHERE id = ?`)
	if err != nil {
		return nil, fmt.Errorf("prepare get by id statement: %w", err)
	}
	defer stmt.Close()

	var user models.User
	err := (
		stmt.QueryRowContext(
			ctx, id,
		).Scan(
			&user.ID, 
			&user.GitHubID, 
			&user.Username, 
			&user.Email, 
			&user.AvatarURL, 
			&user.AccessToken, 
			&user.CreatedAt, 
			&user.LastLoginAt,
		)
	)
	if err != nil {
		return nil, fmt.Errorf("execute get by id statement: %w", err)
	}

	return &user, nil
}

func (r *PgUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	stmt, err := r.DB.PrepareContext(
		ctx, 
		`SELECT id, github_id, username, email, avatar_url, access_token, created_at, last_login_at 
		FROM users WHERE email = ?`
	)
	if err != nil {
		return nil, fmt.Errorf("prepare get by email statement: %w", err)
	}
	defer stmt.Close()

	var user models.User
	err := (
		stmt.QueryRowContext(
			ctx, email
		).Scan(
			&user.ID, 
			&user.GitHubID, 
			&user.Username, 
			&user.Email, 
			&user.AvatarURL, 
			&user.AccessToken, 
			&user.CreatedAt, 
			&user.LastLoginAt,
		)
	)
	if err != nil {
		return nil, fmt.Errorf("execute get by email statement: %w", err)
	}

	return &user, nil
}
