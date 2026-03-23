package mocks

import (
	"context"
	"transaction_routine/internal/domain"
)

type AuditRepositoryMock struct {
	CreateFunc func(ctx context.Context, log domain.AuditLog) error
}

func (m *AuditRepositoryMock) Create(ctx context.Context, log domain.AuditLog) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, log)
	}
	return nil
}
