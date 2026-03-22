package usecase

import (
	"context"
	"errors"
	"testing"
	"transaction_routine/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/application/usecase/mocks"
)

func TestCreateAccountUseCase_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success - creates account", func(t *testing.T) {
		repo := &mocks.AccountRepositoryMock{
			InsertFunc: func(ctx context.Context, document string) (*domain.Account, error) {
				assert.Equal(t, "12345678901", document)
				return &domain.Account{AccountId: 1, DocumentNumber: "12345678901"}, nil
			},
		}
		uc := NewCreateAccountUseCase(repo)

		resp, err := uc.Create(ctx, dto.AccountRequest{DocumentNumber: "12345678901"})
		require.NoError(t, err)
		assert.Equal(t, 1, resp.AccountId)
		assert.Equal(t, "12345678901", resp.DocumentNumber)
	})

	t.Run("error - invalid document length", func(t *testing.T) {
		repo := &mocks.AccountRepositoryMock{
			InsertFunc: func(ctx context.Context, document string) (*domain.Account, error) {
				t.Fatal("Insert should not be called")
				return nil, nil
			},
		}
		uc := NewCreateAccountUseCase(repo)

		resp, err := uc.Create(ctx, dto.AccountRequest{DocumentNumber: "123"})
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, domain.ErrDocumentInvalid))
	})

	t.Run("error - document with non-numeric characters", func(t *testing.T) {
		repo := &mocks.AccountRepositoryMock{
			InsertFunc: func(ctx context.Context, document string) (*domain.Account, error) {
				t.Fatal("Insert should not be called")
				return nil, nil
			},
		}
		uc := NewCreateAccountUseCase(repo)

		resp, err := uc.Create(ctx, dto.AccountRequest{DocumentNumber: "1234567890a"})
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, domain.ErrDocumentNonNumber))
	})

	t.Run("error - account already exists", func(t *testing.T) {
		repo := &mocks.AccountRepositoryMock{
			InsertFunc: func(ctx context.Context, document string) (*domain.Account, error) {
				return nil, domain.ErrAccountAlreadyExists
			},
		}
		uc := NewCreateAccountUseCase(repo)

		resp, err := uc.Create(ctx, dto.AccountRequest{DocumentNumber: "12345678901"})
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, domain.ErrAccountAlreadyExists))
	})
}
