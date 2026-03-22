package mocks

import (
	"context"
	"transaction_routine/internal/domain"
)

type AccountRepositoryMock struct {
	InsertFunc     func(ctx context.Context, documentHash, documentEncrypted string) (*domain.Account, error)
	GetByIdFunc    func(ctx context.Context, accountId int) (*domain.Account, error)
	ExistsByIdFunc func(ctx context.Context, accountId int) (bool, error)
}

func (m *AccountRepositoryMock) Insert(ctx context.Context, documentHash, documentEncrypted string) (*domain.Account, error) {
	return m.InsertFunc(ctx, documentHash, documentEncrypted)
}

func (m *AccountRepositoryMock) GetById(ctx context.Context, accountId int) (*domain.Account, error) {
	return m.GetByIdFunc(ctx, accountId)
}

func (m *AccountRepositoryMock) ExistsById(ctx context.Context, accountId int) (bool, error) {
	return m.ExistsByIdFunc(ctx, accountId)
}
