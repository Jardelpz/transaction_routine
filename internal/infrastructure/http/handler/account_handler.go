package handler

import (
	"github.com/gin-gonic/gin"
	"transaction_routine/internal/application/usecase"
)

type AccountHandler struct {
	createAccountUC   *usecase.CreateAccountUseCase
	retrieveAccountUC *usecase.RetrieveAccountUseCase
}

func NewAccountHandler(create *usecase.CreateAccountUseCase, retrieve *usecase.RetrieveAccountUseCase) *AccountHandler {
	return &AccountHandler{
		createAccountUC:   create,
		retrieveAccountUC: retrieve,
	}
}

func (ah *AccountHandler) GetAccount(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": c.Param("account_id"),
	})

}

func (ah *AccountHandler) CreateAccount(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "CreateAccount",
	})

}
