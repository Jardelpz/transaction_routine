package ports

import "context"

type OperationTypeRepository interface {
	ExistsByID(ctx context.Context, operationTypeID int64) (bool, error)
}
