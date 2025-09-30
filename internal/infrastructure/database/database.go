package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sqlx.DB
}

func Connect(ctx context.Context, cfg *Config) (*DB, error) {
	if cfg.Driver == "" {
		return nil, errors.New("empty driver")
	}
	if cfg.DSN == "" {
		return nil, errors.New("empty dsn")
	}

	db, err := sqlx.ConnectContext(ctx, cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("sqlx connect: %w, dsn: %v", err, cfg.DSN)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &DB{DB: db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}
