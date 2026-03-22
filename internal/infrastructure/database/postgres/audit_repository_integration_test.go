//go:build integration

package postgres

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"transaction_routine/internal/domain"
)

func TestAuditRepository_Integration(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	ctx := context.Background()
	repo := NewAuditRepository(db)

	t.Run("Create inserts audit row", func(t *testing.T) {
		payload, err := json.Marshal(map[string]any{"test": true})
		require.NoError(t, err)

		log := domain.AuditLog{
			EventType:  "test_event",
			EntityType: "test_entity",
			EntityID:   "42",
			Payload:    payload,
		}

		err = repo.Create(ctx, log)
		require.NoError(t, err)

		var count int
		err = db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM audit_logs WHERE event_type = $1 AND entity_type = $2 AND entity_id = $3`,
			"test_event", "test_entity", "42",
		).Scan(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 1)
	})
}
