package sqldb

import (
	"database/sql"
	"time"
)

type Option func(*sql.DB)

func WithMaxOpenConns(n int) Option {
	return func(db *sql.DB) {
		db.SetMaxOpenConns(n)
	}
}

func WithMaxIdleConns(n int) Option {
	return func(db *sql.DB) {
		db.SetMaxIdleConns(n)
	}
}

func WithConnMaxIdleTime(d time.Duration) Option {
	return func(db *sql.DB) {
		db.SetConnMaxIdleTime(d)
	}
}

func WithConnMaxLifetime(d time.Duration) Option {
	return func(db *sql.DB) {
		db.SetConnMaxLifetime(d)
	}
}
