package database

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Anand-S23/isopod/internal/controller"
)

const MAX_OPEN_CONNS = 25
const MAX_IDLE_CONNS = 5
const CONN_MAX_LIFETIME = 5 * time.Minute

func InitDB(ctx context.Context, databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", databaseURL)
	if err != nil {
		return nil, controller.ErrMsg("database: connect: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, controller.ErrMsg("database: ping: %w", err)
	}

	db.SetMaxOpenConns(MAX_OPEN_CONNS)
	db.SetMaxIdleConns(MAX_IDLE_CONNS)
	db.SetConnMaxLifetime(CONN_MAX_LIFETIME)

	return db, nil
}
