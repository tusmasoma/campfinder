package mysql

import (
	"database/sql"
	"log"

	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (tr *transactionRepository) Transaction(txFunc func(tx repository.SQLExecutor) error) error {
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
