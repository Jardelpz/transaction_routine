package usecase

import (
	"context"
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
			ExistsByIdFunc: func(ctx context.Context, accountId int) (bool, error) {
				return true, nil
			},
		}
		opTypeRepo := &mocks.OperationTypeRepositoryMock{
			ExistsByIDFunc: func(ctx context.Context, operationTypeID int) (bool, error) {
				return true, nil
			},
		}
		txRepo := &mocks.TransactionRepositoryMock{
			InsertFunc: func(ctx context.Context, tx domain.Transaction) error {
				assert.Equal(t, 1, tx.AccountId)
				assert.Equal(t, 1, tx.OperationTypeId)
				assert.Equal(t, -100.50, tx.Amount)
				return nil
			},
		}
		uc := NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo)

		resp, err := uc.Create(ctx, dto.TransactionRequest{
			AccountId:       1,
			OperationTypeId: 1,
			Amount:          100.50,
		})
		require.NoError(t, err)
		assert.Equal(t, "created", resp.Status)
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
		uc := NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo)

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
			ExistsByIdFunc: func(ctx context.Context, accountId int) (bool, error) {
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
		uc := NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo)

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
			ExistsByIdFunc: func(ctx context.Context, accountId int) (bool, error) {
				return true, nil
			},
		}
		opTypeRepo := &mocks.OperationTypeRepositoryMock{
			ExistsByIDFunc: func(ctx context.Context, operationTypeID int) (bool, error) {
				return false, nil
			},
		}
		txRepo := &mocks.TransactionRepositoryMock{
			InsertFunc: func(ctx context.Context, tx domain.Transaction) error {
				t.Fatal("Insert should not be called")
				return nil
			},
		}
		uc := NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo)

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
			ExistsByIdFunc: func(ctx context.Context, accountId int) (bool, error) {
				return true, nil
			},
		}
		opTypeRepo := &mocks.OperationTypeRepositoryMock{
			ExistsByIDFunc: func(ctx context.Context, operationTypeID int) (bool, error) {
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
		uc := NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo)

		resp, err := uc.Create(ctx, dto.TransactionRequest{
			AccountId:       1,
			OperationTypeId: 4,
			Amount:          50.25,
		})
		require.NoError(t, err)
		assert.Equal(t, "created", resp.Status)
		assert.Equal(t, 50.25, capturedTx.Amount)
	})
}
