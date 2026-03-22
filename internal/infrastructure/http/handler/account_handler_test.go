package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/application/usecase"
	"transaction_routine/internal/application/usecase/mocks"
	"transaction_routine/internal/domain"
)

func setupAccountHandlers(t *testing.T) (*AccountHandler, *mocks.AccountRepositoryMock) {
	repo := &mocks.AccountRepositoryMock{}
	createUC := usecase.NewCreateAccountUseCase(repo)
	retrieveUC := usecase.NewRetrieveAccountUseCase(repo)
	return NewAccountHandler(createUC, retrieveUC), repo
}

func TestAccountHandler_CreateAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success - creates account", func(t *testing.T) {
		ah, repo := setupAccountHandlers(t)
		repo.InsertFunc = func(ctx context.Context, document string) (*domain.Account, error) {
			return &domain.Account{AccountId: 1, DocumentNumber: document}, nil
		}

		body := dto.AccountRequest{DocumentNumber: "12345678901"}
		bodyBytes, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/account", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		ah.CreateAccount(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp dto.AccountResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, 1, resp.AccountId)
		assert.Equal(t, "12345678901", resp.DocumentNumber)
	})

	t.Run("error - invalid document", func(t *testing.T) {
		ah, repo := setupAccountHandlers(t)
		repo.InsertFunc = func(ctx context.Context, document string) (*domain.Account, error) {
			t.Fatal("Insert should not be called")
			return nil, nil
		}

		body := dto.AccountRequest{DocumentNumber: "123"}
		bodyBytes, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/account", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		ah.CreateAccount(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error - invalid JSON", func(t *testing.T) {
		ah, _ := setupAccountHandlers(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/account", bytes.NewReader([]byte("invalid")))
		c.Request.Header.Set("Content-Type", "application/json")

		ah.CreateAccount(c)

		assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusUnprocessableEntity, "expected 400 or 422, got %d", w.Code)
	})
}

func TestAccountHandler_GetAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success - returns account", func(t *testing.T) {
		ah, repo := setupAccountHandlers(t)
		repo.GetByIdFunc = func(ctx context.Context, accountId int) (*domain.Account, error) {
			return &domain.Account{AccountId: 1, DocumentNumber: "12345678901"}, nil
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/v1/account/1", nil)
		c.Params = gin.Params{{Key: "account_id", Value: "1"}}

		ah.GetAccount(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.AccountResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, 1, resp.AccountId)
		assert.Equal(t, "12345678901", resp.DocumentNumber)
	})

	t.Run("error - account not found", func(t *testing.T) {
		ah, repo := setupAccountHandlers(t)
		repo.GetByIdFunc = func(ctx context.Context, accountId int) (*domain.Account, error) {
			return nil, domain.ErrAccountNotFound
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/v1/account/999", nil)
		c.Params = gin.Params{{Key: "account_id", Value: "999"}}

		ah.GetAccount(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("error - invalid account_id", func(t *testing.T) {
		ah, _ := setupAccountHandlers(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/v1/account/abc", nil)
		c.Params = gin.Params{{Key: "account_id", Value: "abc"}}

		ah.GetAccount(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
