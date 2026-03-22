//go:build integration

package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/application/usecase"
	"transaction_routine/internal/infrastructure/database/postgres"
	"transaction_routine/internal/infrastructure/http/handler"
)

func init() {
	_ = godotenv.Load("../../../.env")
}

func TestAPI_Integration(t *testing.T) {
	db := postgres.GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	accountRepo := postgres.NewAccountRepository(db)
	txRepo := postgres.NewTransactionRepository(db)
	opTypeRepo := postgres.NewOperationTypeRepository(db)

	createAccountUC := usecase.NewCreateAccountUseCase(accountRepo)
	retrieveAccountUC := usecase.NewRetrieveAccountUseCase(accountRepo)
	createTransactionUC := usecase.NewCreateTransactionUseCase(txRepo, accountRepo, opTypeRepo)

	accountHandler := handler.NewAccountHandler(createAccountUC, retrieveAccountUC)
	transactionHandler := handler.NewTransactionHandler(createTransactionUC)

	router := NewRouter(accountHandler, transactionHandler)

	t.Run("full flow: create account, get account, create transaction", func(t *testing.T) {
		doc := "55555555555"

		// Create account
		createBody, _ := json.Marshal(dto.AccountRequest{DocumentNumber: doc})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/account", bytes.NewReader(createBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code, "response: %s", w.Body.String())
		var accountResp dto.AccountResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &accountResp))
		accountID := accountResp.AccountId
		assert.Greater(t, accountID, 0)

		// Get account
		w = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, "/v1/account/"+fmt.Sprintf("%d", accountID), nil)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var getResp dto.AccountResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &getResp))
		assert.Equal(t, accountID, getResp.AccountId)
		assert.Equal(t, doc, getResp.DocumentNumber)

		// Create transaction
		txBody, _ := json.Marshal(dto.TransactionRequest{
			AccountId:       accountID,
			OperationTypeId: 1,
			Amount:          100.50,
		})
		w = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodPost, "/v1/transaction", bytes.NewReader(txBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		var txResp dto.TransactionResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &txResp))
		assert.Equal(t, "created", txResp.Status)
	})
}
