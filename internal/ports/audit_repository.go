package ports

import (
	"context"
	"transaction_routine/internal/domain"
)

type AuditRepository interface {
	Create(ctx context.Context, log domain.AuditLog) error
}
