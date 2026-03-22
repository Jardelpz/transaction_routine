//go:build integration

package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"transaction_routine/internal/domain"
)

func TestAccountRepository_Integration(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	ctx := context.Background()
	repo := NewAccountRepository(db)

	t.Run("Insert and GetById", func(t *testing.T) {
		doc := "98765432100"
		account, err := repo.Insert(ctx, doc)
		require.NoError(t, err)
		require.NotNil(t, account)
		assert.Greater(t, account.AccountId, 0)
		assert.Equal(t, doc, account.DocumentNumber)

		fetched, err := repo.GetById(ctx, account.AccountId)
		require.NoError(t, err)
		assert.Equal(t, account.AccountId, fetched.AccountId)
		assert.Equal(t, account.DocumentNumber, fetched.DocumentNumber)
	})

	t.Run("Insert duplicate document returns ErrAccountAlreadyExists", func(t *testing.T) {
		const doc = "11111111111"
		_, err := repo.Insert(ctx, doc)
		require.NoError(t, err)

		_, err = repo.Insert(ctx, doc)
		assert.ErrorIs(t, err, domain.ErrAccountAlreadyExists)
	})

	t.Run("GetById not found returns ErrAccountNotFound", func(t *testing.T) {
		_, err := repo.GetById(ctx, 999999)
		assert.ErrorIs(t, err, domain.ErrAccountNotFound)
	})

	t.Run("ExistsById", func(t *testing.T) {
		doc := "22222222222"
		account, err := repo.Insert(ctx, doc)
		require.NoError(t, err)

		exists, err := repo.ExistsById(ctx, account.AccountId)
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = repo.ExistsById(ctx, 999999)
		require.NoError(t, err)
		assert.False(t, exists)
	})
}
