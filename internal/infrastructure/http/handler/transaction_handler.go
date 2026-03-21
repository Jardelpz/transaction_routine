package handler

import (
	"github.com/gin-gonic/gin"
	"transaction_routine/internal/application/usecase"
)

type TransactionHandler struct {
	createTransactionUC *usecase.CreateTransactionUseCase
}

func NewTransactionHandler(create *usecase.CreateTransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		createTransactionUC: create,
	}
}

func (ah *TransactionHandler) CreateTransaction(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "CreateTransaction",
	})

}
