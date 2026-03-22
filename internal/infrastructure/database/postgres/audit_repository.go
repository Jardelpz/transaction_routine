package postgres

import (
	"context"
	"database/sql"

	"transaction_routine/internal/domain"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (ar *AuditRepository) Create(ctx context.Context, log domain.AuditLog) error {
	const query = `
		INSERT INTO audit_logs (event_type, entity_type, entity_id, payload)
		VALUES ($1, $2, $3, $4)
	`

	_, err := ar.db.ExecContext(
		ctx,
		query,
		log.EventType,
		log.EntityType,
		log.EntityID,
		log.Payload,
	)

	return err
}
