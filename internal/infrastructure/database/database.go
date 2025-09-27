package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func Connect(ctx context.Context, cfg *Config) (*DB, error) {
	if cfg.User == "" {
		return nil, errors.New("empty user")
	}
	if cfg.Pass == "" {
		return nil, errors.New("empty pass")
	}

	url := toURL(cfg.Addr, cfg.Port, cfg.DB, cfg.User, cfg.Pass)

	db, err := sqlx.ConnectContext(ctx, "postgres", url)
	if err != nil {
		return nil, fmt.Errorf("sqlx connect: %w, url: %v", err, url)
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
