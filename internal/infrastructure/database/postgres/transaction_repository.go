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
	query := `
		INSERT INTO transactions(account_id, operation_type_id, amount, event_date) 
		VALUES ($1, $2, $3, $4)
    `

	err := ar.db.QueryRowContext(ctx, query, transaction.AccountId, transaction.OperationTypeId, transaction.Amount, transaction.EventDate)
	if err != nil {
		return err.Err()
	}

	return nil
}
