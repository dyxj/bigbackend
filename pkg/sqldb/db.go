package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	_driverName             = "postgres"
	_defaultMaxOpenConns    = 25
	_defaultMaxIdleConns    = 3
	_defaultConnMaxIdleTime = 1 * time.Minute
	_defaultConnMaxLifetime = 1 * time.Hour
)

// NewDBConn creates a new database connection pool
func NewDBConn(ctx context.Context, cfg Config, opts ...Option) (*sql.DB, error) {
	var connectionString = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host(), cfg.Port(), cfg.User(), cfg.Password(), cfg.DBName())

	db, err := sql.Open(_driverName, connectionString)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(_defaultMaxOpenConns)
	db.SetMaxIdleConns(_defaultMaxIdleConns)
	db.SetConnMaxIdleTime(_defaultConnMaxIdleTime)
	db.SetConnMaxLifetime(_defaultConnMaxLifetime)

	for _, opt := range opts {
		opt(db)
	}

	return db, nil
}
