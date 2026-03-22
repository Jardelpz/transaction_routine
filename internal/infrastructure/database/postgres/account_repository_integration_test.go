//go:build integration

package postgres

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"transaction_routine/internal/domain"
	"transaction_routine/internal/infrastructure/security"
)

var integrationDocSeq atomic.Int64

// uniqueTestDocument returns an 11-digit string valid for domain.ValidateDocument (repeated runs / shared DB).
func uniqueTestDocument(t *testing.T) string {
	t.Helper()
	n := integrationDocSeq.Add(1) + time.Now().UnixNano()
	return fmt.Sprintf("%011d", n%100_000_000_000)
}

func setupProtector(t *testing.T) *security.DocumentProtector {
	t.Helper()
	os.Setenv("DOCUMENT_ENCRYPTION_KEY", "12345678901234567890123456789012")
	os.Setenv("DOCUMENT_HASH_SALT", "test-salt")
	p, err := security.NewDocumentProtector()
	require.NoError(t, err)
	return p
}

func TestAccountRepository_Integration(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	ctx := context.Background()
	repo := NewAccountRepository(db)
	protector := setupProtector(t)

	t.Run("Insert and GetById", func(t *testing.T) {
		doc := uniqueTestDocument(t)
		hash := protector.Hash(doc)
		encrypted, err := protector.Encrypt(doc)
		require.NoError(t, err)

		account, err := repo.Insert(ctx, hash, encrypted)
		require.NoError(t, err)
		require.NotNil(t, account)
		assert.Greater(t, account.AccountId, int64(0))
		assert.Equal(t, hash, account.DocumentHash)

		fetched, err := repo.GetById(ctx, account.AccountId)
		require.NoError(t, err)
		assert.Equal(t, account.AccountId, fetched.AccountId)
		assert.Equal(t, hash, fetched.DocumentHash)

		decrypted, err := protector.Decrypt(fetched.DocumentEncrypted)
		require.NoError(t, err)
		assert.Equal(t, doc, decrypted)
	})

	t.Run("Insert duplicate document returns ErrAccountAlreadyExists", func(t *testing.T) {
		doc := uniqueTestDocument(t)
		hash := protector.Hash(doc)
		encrypted, _ := protector.Encrypt(doc)

		_, err := repo.Insert(ctx, hash, encrypted)
		require.NoError(t, err)

		_, err = repo.Insert(ctx, hash, encrypted)
		assert.ErrorIs(t, err, domain.ErrAccountAlreadyExists)
	})

	t.Run("GetById not found returns ErrAccountNotFound", func(t *testing.T) {
		_, err := repo.GetById(ctx, 999999)
		assert.ErrorIs(t, err, domain.ErrAccountNotFound)
	})

	t.Run("ExistsById", func(t *testing.T) {
		doc := uniqueTestDocument(t)
		hash := protector.Hash(doc)
		encrypted, _ := protector.Encrypt(doc)
		account, err := repo.Insert(ctx, hash, encrypted)
		require.NoError(t, err)

		exists, err := repo.ExistsById(ctx, account.AccountId)
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = repo.ExistsById(ctx, 999999)
		require.NoError(t, err)
		assert.False(t, exists)
	})
}
