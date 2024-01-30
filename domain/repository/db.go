package repository

import (
	"context"
	"database/sql"
)

type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type TransactionRepository interface {
	Transaction(txFunc func(tx SQLExecutor) error) error
}

type QueryOptions struct {
	Executor SQLExecutor
}
