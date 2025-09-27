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
	if cfg.User == "" {
		return errors.New("empty user")
	}
	if cfg.Pass == "" {
		return errors.New("empty pass")
	}

	url := toURL(cfg.Addr, cfg.Port, cfg.DB, cfg.User, cfg.Pass)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return fmt.Errorf("sql open: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(fsys)

	if err = goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err = goose.Up(db, dir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}

func toURL(host, port, dbName, user, password string) string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbName, user, password,
	)
}
