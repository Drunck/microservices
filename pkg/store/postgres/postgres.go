package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type DB struct {
	Pool *pgxpool.Pool
}

type ConnConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DbName       string
	MaxOpenConns int
	MaxIdleTime  string
}

func OpenDB(cfg ConnConfig) (*DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(cfg.MaxOpenConns)
	config.MaxConnLifetime = duration

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.Ping(ctx); err != nil {
		return nil, err
	}

	return &DB{Pool: db}, nil
}
