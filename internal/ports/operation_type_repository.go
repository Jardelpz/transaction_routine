package ports

import "context"

type OperationTypeRepository interface {
	ExistsByID(ctx context.Context, operationTypeID int) (bool, error)
}
