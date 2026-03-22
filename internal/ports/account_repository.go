package ports

import (
	"context"
	"transaction_routine/internal/domain"
)

type AccountRepository interface {
	Insert(ctx context.Context, documentHash string, documentEncrypted string) (*domain.Account, error)
	GetById(ctx context.Context, accountId int) (*domain.Account, error)
	ExistsById(ctx context.Context, accountId int) (bool, error)
}
