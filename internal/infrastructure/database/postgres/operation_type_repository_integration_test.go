//go:build integration

package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationTypeRepository_Integration(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	ctx := context.Background()
	repo := NewOperationTypeRepository(db)

	t.Run("ExistsByID - existing types", func(t *testing.T) {
		for _, id := range []int64{1, 2, 3, 4} {
			exists, err := repo.ExistsByID(ctx, id)
			require.NoError(t, err)
			assert.True(t, exists, "operation type %d should exist", id)
		}
	})

	t.Run("ExistsByID - non-existing type", func(t *testing.T) {
		exists, err := repo.ExistsByID(ctx, 99)
		require.NoError(t, err)
		assert.False(t, exists)
	})
}
