package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/application/usecase/mocks"
	"transaction_routine/internal/domain"
)

func TestCreateTransactionUseCase_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success - creates transaction", func(t *testing.T) {
		accountRepo := &mocks.AccountRepositoryMock{
			ExistsByIdFunc: func(ctx context.Context, accountId int64) (bool, error) {
				return true, nil
			},
		}
		opTypeRepo := &mocks.OperationTypeRepositoryMock{
			ExistsByIDFunc: func(ctx context.Context, operationTypeID int64) (bool, error) {
				return true, nil
			},
		}
		txRepo := &mocks.TransactionRepositoryMock{
			InsertFunc: func(ctx context.Context, tx domain.Transaction) error {
				assert.Equal(t, int64(1), tx.AccountId)
				assert.Equal(t, int64(1), tx.OperationTypeId)
				assert.Equal(t, -100.50, tx.Amount)
				return nil
			},
		}
		uc := NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo, noopAudit())

		resp, err := uc.Create(ctx, dto.TransactionRequest{
			AccountId:       1,
			OperationTypeId: 1,
			Amount:          100.50,
		})
		require.NoError(t, err)
		assert.Equal(t, "created", resp.Status)
	})

	t.Run("success - writes audit log", func(t *testing.T) {
		var captured domain.AuditLog
		audit := &mocks.AuditRepositoryMock{
			CreateFunc: func(ctx context.Context, log domain.AuditLog) error {
				captured = log
				return nil
			},
		}
		accountRepo := &mocks.AccountRepositoryMock{
			ExistsByIdFunc: func(ctx context.Context, accountId int64) (bool, error) { return true, nil },
		}
		opTypeRepo := &mocks.OperationTypeRepositoryMock{
			ExistsByIDFunc: func(ctx context.Context, operationTypeID int64) (bool, error) { return true, nil },
		}
		txRepo := &mocks.TransactionRepositoryMock{
			InsertFunc: func(ctx context.Context, tx domain.Transaction) error { return nil },
		}
		uc := NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo, audit)

		_, err := uc.Create(ctx, dto.TransactionRequest{
			AccountId:       5,
			OperationTypeId: 2,
			Amount:          10,
		})
		require.NoError(t, err)
		assert.Equal(t, "transaction_created", captured.EventType)
		assert.Equal(t, "transaction", captured.EntityType)
		assert.Equal(t, "5", captured.EntityID)
		require.NotEmpty(t, captured.Payload)
		var payload map[string]any
		require.NoError(t, json.Unmarshal(captured.Payload, &payload))
		assert.EqualValues(t, float64(5), payload["account_id"])
		assert.EqualValues(t, float64(2), payload["operation_type_id"])
	})

	t.Run("error - invalid amount zero", func(t *testing.T) {
		accountRepo := &mocks.AccountRepositoryMock{}
		opTypeRepo := &mocks.OperationTypeRepositoryMock{}
		txRepo := &mocks.TransactionRepositoryMock{
			InsertFunc: func(ctx context.Context, tx domain.Transaction) error {
				t.Fatal("Insert should not be called")
				return nil
			},
		}
		uc := NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo, noopAudit())

		resp, err := uc.Create(ctx, dto.TransactionRequest{
			AccountId:       1,
			OperationTypeId: 1,
			Amount:          0,
		})
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, domain.ErrInvalidAmount))
	})

	t.Run("error - account not found", func(t *testing.T) {
		accountRepo := &mocks.AccountRepositoryMock{
			ExistsByIdFunc: func(ctx context.Context, accountId int64) (bool, error) {
				return false, nil
			},
		}
		opTypeRepo := &mocks.OperationTypeRepositoryMock{}
		txRepo := &mocks.TransactionRepositoryMock{
			InsertFunc: func(ctx context.Context, tx domain.Transaction) error {
				t.Fatal("Insert should not be called")
				return nil
			},
		}
		uc := NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo, noopAudit())

		resp, err := uc.Create(ctx, dto.TransactionRequest{
			AccountId:       999,
			OperationTypeId: 1,
			Amount:          100,
		})
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, domain.ErrAccountNotFound))
	})

	t.Run("error - operation type not found", func(t *testing.T) {
		accountRepo := &mocks.AccountRepositoryMock{
			ExistsByIdFunc: func(ctx context.Context, accountId int64) (bool, error) {
				return true, nil
			},
		}
		opTypeRepo := &mocks.OperationTypeRepositoryMock{
			ExistsByIDFunc: func(ctx context.Context, operationTypeID int64) (bool, error) {
				return false, nil
			},
		}
		txRepo := &mocks.TransactionRepositoryMock{
			InsertFunc: func(ctx context.Context, tx domain.Transaction) error {
				t.Fatal("Insert should not be called")
				return nil
			},
		}
		uc := NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo, noopAudit())

		resp, err := uc.Create(ctx, dto.TransactionRequest{
			AccountId:       1,
			OperationTypeId: 99,
			Amount:          100,
		})
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, domain.ErrOperationTypeNotFound))
	})

	t.Run("success - credit voucher keeps positive amount", func(t *testing.T) {
		accountRepo := &mocks.AccountRepositoryMock{
			ExistsByIdFunc: func(ctx context.Context, accountId int64) (bool, error) {
				return true, nil
			},
		}
		opTypeRepo := &mocks.OperationTypeRepositoryMock{
			ExistsByIDFunc: func(ctx context.Context, operationTypeID int64) (bool, error) {
				return true, nil
			},
		}
		var capturedTx domain.Transaction
		txRepo := &mocks.TransactionRepositoryMock{
			InsertFunc: func(ctx context.Context, tx domain.Transaction) error {
				capturedTx = tx
				return nil
			},
		}
		uc := NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo, noopAudit())

		resp, err := uc.Create(ctx, dto.TransactionRequest{
			AccountId:       1,
			OperationTypeId: 4,
			Amount:          50.25,
		})
		require.NoError(t, err)
		assert.Equal(t, "created", resp.Status)
		assert.Equal(t, 50.25, capturedTx.Amount)
	})

	t.Run("success - audit failure does not fail request", func(t *testing.T) {
		audit := &mocks.AuditRepositoryMock{
			CreateFunc: func(ctx context.Context, log domain.AuditLog) error {
				return errors.New("audit error")
			},
		}
		accountRepo := &mocks.AccountRepositoryMock{
			ExistsByIdFunc: func(ctx context.Context, accountId int64) (bool, error) { return true, nil },
		}
		opTypeRepo := &mocks.OperationTypeRepositoryMock{
			ExistsByIDFunc: func(ctx context.Context, operationTypeID int64) (bool, error) { return true, nil },
		}
		txRepo := &mocks.TransactionRepositoryMock{
			InsertFunc: func(ctx context.Context, tx domain.Transaction) error { return nil },
		}
		uc := NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo, audit)

		resp, err := uc.Create(ctx, dto.TransactionRequest{
			AccountId:       1,
			OperationTypeId: 1,
			Amount:          10,
		})
		require.NoError(t, err)
		assert.Equal(t, "created", resp.Status)
	})
}
