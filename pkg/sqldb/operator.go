package sqldb

import (
	"context"
	"database/sql"
)

// Queryable interface for sql QueryContext method
type Queryable interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// Executable interface for sql ExecContext method
type Executable interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type TransactionManager interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
