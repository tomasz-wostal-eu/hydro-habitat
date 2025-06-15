package store

import (
	"hydro-habitat/backend/config"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func NewPostgresStore(cfg *config.Config) (*DB, error) {
	db, err := sqlx.Connect("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
