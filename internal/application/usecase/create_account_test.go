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

func TestCreateAccountUseCase_Create(t *testing.T) {
	ctx := context.Background()
	protector := &mocks.DocumentProtectorMock{
		HashFunc:    func(d string) string { return "hash-" + d },
		EncryptFunc: func(d string) (string, error) { return "enc-" + d, nil },
		DecryptFunc: func(c string) (string, error) { return "12345678901", nil },
	}

	t.Run("success - creates account with hash and encrypt", func(t *testing.T) {
		repo := &mocks.AccountRepositoryMock{
			InsertFunc: func(ctx context.Context, hash, enc string) (*domain.Account, error) {
				assert.Equal(t, "hash-12345678901", hash)
				assert.Equal(t, "enc-12345678901", enc)
				return &domain.Account{
					AccountId:         1,
					DocumentHash:      hash,
					DocumentEncrypted: enc,
				}, nil
			},
		}
		uc := NewCreateAccountUseCase(repo, protector)

		resp, err := uc.Create(ctx, dto.AccountRequest{DocumentNumber: "12345678901"})
		require.NoError(t, err)
		assert.Equal(t, 1, resp.AccountId)
		assert.Equal(t, "12345678901", resp.DocumentNumber)
	})

	t.Run("error - invalid document length", func(t *testing.T) {
		repo := &mocks.AccountRepositoryMock{
			InsertFunc: func(ctx context.Context, hash, enc string) (*domain.Account, error) {
				t.Fatal("Insert should not be called")
				return nil, nil
			},
		}
		uc := NewCreateAccountUseCase(repo, protector)

		resp, err := uc.Create(ctx, dto.AccountRequest{DocumentNumber: "123"})
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, domain.ErrDocumentInvalid))
	})

	t.Run("error - document with non-numeric characters", func(t *testing.T) {
		repo := &mocks.AccountRepositoryMock{
			InsertFunc: func(ctx context.Context, hash, enc string) (*domain.Account, error) {
				t.Fatal("Insert should not be called")
				return nil, nil
			},
		}
		uc := NewCreateAccountUseCase(repo, protector)

		resp, err := uc.Create(ctx, dto.AccountRequest{DocumentNumber: "1234567890a"})
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, domain.ErrDocumentNonNumber))
	})

	t.Run("error - encrypt fails", func(t *testing.T) {
		encProtector := &mocks.DocumentProtectorMock{
			HashFunc:    func(d string) string { return "hash" },
			EncryptFunc: func(d string) (string, error) { return "", errors.New("encrypt error") },
		}
		repo := &mocks.AccountRepositoryMock{
			InsertFunc: func(ctx context.Context, hash, enc string) (*domain.Account, error) {
				t.Fatal("Insert should not be called")
				return nil, nil
			},
		}
		uc := NewCreateAccountUseCase(repo, encProtector)

		resp, err := uc.Create(ctx, dto.AccountRequest{DocumentNumber: "12345678901"})
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "encrypt")
	})

	t.Run("error - account already exists", func(t *testing.T) {
		repo := &mocks.AccountRepositoryMock{
			InsertFunc: func(ctx context.Context, hash, enc string) (*domain.Account, error) {
				return nil, domain.ErrAccountAlreadyExists
			},
		}
		uc := NewCreateAccountUseCase(repo, protector)

		resp, err := uc.Create(ctx, dto.AccountRequest{DocumentNumber: "12345678901"})
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, domain.ErrAccountAlreadyExists))
	})

	t.Run("error - decrypt fails after insert", func(t *testing.T) {
		decProtector := &mocks.DocumentProtectorMock{
			HashFunc:    func(d string) string { return "hash" },
			EncryptFunc: func(d string) (string, error) { return "enc", nil },
			DecryptFunc: func(c string) (string, error) { return "", errors.New("decrypt error") },
		}
		repo := &mocks.AccountRepositoryMock{
			InsertFunc: func(ctx context.Context, hash, enc string) (*domain.Account, error) {
				return &domain.Account{AccountId: 1, DocumentEncrypted: enc}, nil
			},
		}
		uc := NewCreateAccountUseCase(repo, decProtector)

		resp, err := uc.Create(ctx, dto.AccountRequest{DocumentNumber: "12345678901"})
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "decrypt")
	})
}
