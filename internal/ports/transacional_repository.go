package ports

import (
	"context"
	"transaction_routine/internal/domain"
)

type TransactionRepository interface {
	Insert(ctx context.Context, account domain.Transaction) error
}
