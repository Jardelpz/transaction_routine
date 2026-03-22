package mocks

import "context"

type OperationTypeRepositoryMock struct {
	ExistsByIDFunc func(ctx context.Context, operationTypeID int64) (bool, error)
}

func (m *OperationTypeRepositoryMock) ExistsByID(ctx context.Context, operationTypeID int64) (bool, error) {
	return m.ExistsByIDFunc(ctx, operationTypeID)
}
