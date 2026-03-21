package ports

import (
	"context"
	"transaction_routine/internal/domain"
)

type AccountRepository interface {
	Insert(ctx context.Context, account domain.Account) (*domain.Account, error)
	GetById(ctx context.Context, accountId int) (*domain.Account, error)
}
