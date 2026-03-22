//go:build integration

package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"transaction_routine/internal/domain"
	"transaction_routine/internal/infrastructure/security"
)

func TestTransactionRepository_Integration(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	os.Setenv("DOCUMENT_ENCRYPTION_KEY", "12345678901234567890123456789012")
	os.Setenv("DOCUMENT_HASH_SALT", "test-salt")
	protector, err := security.NewDocumentProtector()
	require.NoError(t, err)

	ctx := context.Background()
	accountRepo := NewAccountRepository(db)
	txRepo := NewTransactionRepository(db)

	t.Run("Insert transaction", func(t *testing.T) {
		doc := "33333333333"
		hash := protector.Hash(doc)
		encrypted, err := protector.Encrypt(doc)
		require.NoError(t, err)

		account, err := accountRepo.Insert(ctx, hash, encrypted)
		require.NoError(t, err)
		require.NotNil(t, account)

		tx := domain.Transaction{
			AccountId:       account.AccountId,
			OperationTypeId: 1,
			Amount:          -100.50,
			EventDate:       time.Now().UTC(),
		}

		err = txRepo.Insert(ctx, tx)
		require.NoError(t, err)
	})
}
