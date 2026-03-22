package postgres

import (
	"context"
	"database/sql"
)

type OperationTypeRepository struct {
	db *sql.DB
}

func NewOperationTypeRepository(db *sql.DB) *OperationTypeRepository {
	return &OperationTypeRepository{db: db}
}

func (or *OperationTypeRepository) ExistsByID(ctx context.Context, operationTypeID int64) (bool, error) {
	var exists bool
	const query = `
        SELECT EXISTS (
            SELECT 1
            FROM operation_types
            WHERE operation_type_id = $1
        )
    `

	err := or.db.QueryRowContext(ctx, query, operationTypeID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
