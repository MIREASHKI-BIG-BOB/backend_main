package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

var (
	//go:embed migrations/*.sql
	migrations embed.FS
)

func Migrate(_ context.Context, cfg *Config) error {
	if err := migrate(cfg, "migrations", migrations); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

func migrate(cfg *Config, dir string, fsys fs.FS) error {
	if cfg.Driver == "" {
		return errors.New("empty driver")
	}
	if cfg.DSN == "" {
		return errors.New("empty dsn")
	}

	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return fmt.Errorf("sql open: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(fsys)

	if err = goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err = goose.Up(db, dir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
