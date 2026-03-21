package postgres

import (
	"context"
	"database/sql"
	"transaction_routine/internal/domain"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (ar *TransactionRepository) Insert(ctx context.Context, transaction domain.Transaction) error {
	return nil
}
