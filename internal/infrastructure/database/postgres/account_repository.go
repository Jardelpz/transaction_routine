package postgres

import (
	"context"
	"database/sql"

	"transaction_routine/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (ar *AccountRepository) Insert(ctx context.Context, account domain.Account) (*domain.Account, error) {
	return nil, nil
}

func (ar *AccountRepository) GetById(ctx context.Context, accountId int) (*domain.Account, error) {
	return &domain.Account{}, nil
}
