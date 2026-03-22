package mocks

import (
	"context"
	"transaction_routine/internal/domain"
)

type TransactionRepositoryMock struct {
	InsertFunc func(ctx context.Context, tx domain.Transaction) error
}

func (m *TransactionRepositoryMock) Insert(ctx context.Context, tx domain.Transaction) error {
	return m.InsertFunc(ctx, tx)
}
