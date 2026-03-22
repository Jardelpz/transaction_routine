package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"transaction_routine/internal/application/usecase/mocks"
	"transaction_routine/internal/domain"
)

func TestRetrieveAccountUseCase_Retrieve(t *testing.T) {
	ctx := context.Background()

	t.Run("success - returns account", func(t *testing.T) {
		repo := &mocks.AccountRepositoryMock{
			GetByIdFunc: func(ctx context.Context, accountId int) (*domain.Account, error) {
				assert.Equal(t, 1, accountId)
				return &domain.Account{AccountId: 1, DocumentNumber: "12345678901"}, nil
			},
		}
		uc := NewRetrieveAccountUseCase(repo)

		resp, err := uc.Retrieve(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, 1, resp.AccountId)
		assert.Equal(t, "12345678901", resp.DocumentNumber)
	})

	t.Run("error - account not found", func(t *testing.T) {
		repo := &mocks.AccountRepositoryMock{
			GetByIdFunc: func(ctx context.Context, accountId int) (*domain.Account, error) {
				return nil, domain.ErrAccountNotFound
			},
		}
		uc := NewRetrieveAccountUseCase(repo)

		resp, err := uc.Retrieve(ctx, 999)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, domain.ErrAccountNotFound))
	})
}
