package db

import (
	"context"
	"database/sql"
	"log"
)

type SQLExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type QueryOptions struct {
	Executor SQLExecutor
}

// トランザクション
type TransactionRepository interface {
	Transaction(txFunc func(tx SQLExecutor) error) error
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (tr *transactionRepository) Transaction(txFunc func(tx SQLExecutor) error) error {
	tx, err := tr.db.Begin()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Rollback error: %v", rollbackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			log.Printf("Commit error: %v", commitErr)
		}
	}()

	err = txFunc(tx)
	return err
}
