package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"transaction_routine/internal/application/dto"
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
	var reqInput dto.TransactionRequest
	if err := c.BindJSON(&reqInput); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	// todo req with timeout
	response, err := ah.createTransactionUC.Create(c, reqInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)

}
