//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"transaction_routine/internal/domain"
)

func TestTransactionRepository_Integration(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	ctx := context.Background()
	accountRepo := NewAccountRepository(db)
	txRepo := NewTransactionRepository(db)

	t.Run("Insert transaction", func(t *testing.T) {
		account, err := accountRepo.Insert(ctx, "33333333333")
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
