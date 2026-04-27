package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	_ = godotenv.Load()

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		return fmt.Errorf("DB_URI is required (set in environment or .env file)")
	}

	ctx := context.Background()
	command := flag.String("command", "up", "goose command, e.g. up, down, up-by-one, status, version")
	migrationsDir := flag.String("dir", "migrations", "migrations directory (relative to the current working directory; run from repo root)")
	flag.Parse()

	db, err := sql.Open("pgx", dbURI)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() { _ = db.Close() }()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	if err := goose.RunContext(ctx, *command, db, *migrationsDir, flag.Args()...); err != nil {
		return err
	}
	return nil
}
