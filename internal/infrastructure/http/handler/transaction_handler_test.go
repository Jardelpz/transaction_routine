package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/application/usecase"
	"transaction_routine/internal/application/usecase/mocks"
	"transaction_routine/internal/domain"
)

func setupTransactionHandler(t *testing.T) (*TransactionHandler, *mocks.AccountRepositoryMock, *mocks.OperationTypeRepositoryMock, *mocks.TransactionRepositoryMock) {
	accountRepo := &mocks.AccountRepositoryMock{}
	opTypeRepo := &mocks.OperationTypeRepositoryMock{}
	txRepo := &mocks.TransactionRepositoryMock{}
	audit := &mocks.AuditRepositoryMock{}
	uc := usecase.NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo, audit)
	return NewTransactionHandler(uc), accountRepo, opTypeRepo, txRepo
}

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success - creates transaction", func(t *testing.T) {
		th, accountRepo, opTypeRepo, txRepo := setupTransactionHandler(t)
		accountRepo.ExistsByIdFunc = func(ctx context.Context, accountId int64) (bool, error) { return true, nil }
		opTypeRepo.ExistsByIDFunc = func(ctx context.Context, id int64) (bool, error) { return true, nil }
		txRepo.InsertFunc = func(ctx context.Context, tx domain.Transaction) error { return nil }

		body := dto.TransactionRequest{
			AccountId:       1,
			OperationTypeId: 1,
			Amount:          100.50,
		}
		bodyBytes, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/transaction", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		th.CreateTransaction(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp dto.TransactionResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, "created", resp.Status)
	})

	t.Run("error - invalid amount zero", func(t *testing.T) {
		th, accountRepo, opTypeRepo, _ := setupTransactionHandler(t)
		accountRepo.ExistsByIdFunc = func(ctx context.Context, accountId int64) (bool, error) { return true, nil }
		opTypeRepo.ExistsByIDFunc = func(ctx context.Context, id int64) (bool, error) { return true, nil }

		body := dto.TransactionRequest{
			AccountId:       1,
			OperationTypeId: 1,
			Amount:          0,
		}
		bodyBytes, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/transaction", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		th.CreateTransaction(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error - account not found", func(t *testing.T) {
		th, accountRepo, opTypeRepo, _ := setupTransactionHandler(t)
		accountRepo.ExistsByIdFunc = func(ctx context.Context, accountId int64) (bool, error) { return false, nil }
		opTypeRepo.ExistsByIDFunc = func(ctx context.Context, id int64) (bool, error) { return true, nil }

		body := dto.TransactionRequest{
			AccountId:       999,
			OperationTypeId: 1,
			Amount:          100,
		}
		bodyBytes, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/transaction", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		th.CreateTransaction(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error - invalid JSON", func(t *testing.T) {
		th, _, _, _ := setupTransactionHandler(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/transaction", bytes.NewReader([]byte("invalid")))
		c.Request.Header.Set("Content-Type", "application/json")

		th.CreateTransaction(c)

		assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusUnprocessableEntity, "expected 400 or 422, got %d", w.Code)
	})
}
