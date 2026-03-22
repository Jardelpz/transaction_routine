package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/application/usecase"
	"transaction_routine/internal/domain"
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
	accountId, err := strconv.ParseInt(c.Param("account_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidAccountType.Error()})
		return
	}

	account, err := ah.retrieveAccountUC.Retrieve(c, accountId)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrAccountNotFound.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (ah *AccountHandler) CreateAccount(c *gin.Context) {
	var reqInput dto.AccountRequest
	if err := c.BindJSON(&reqInput); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response, err := ah.createAccountUC.Create(c, reqInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}
