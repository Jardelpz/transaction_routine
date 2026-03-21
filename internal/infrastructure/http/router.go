package http

import (
	"github.com/gin-gonic/gin"
	"transaction_routine/internal/infrastructure/http/handler"
)

func NewRouter(accountHandler *handler.AccountHandler, transactionHandler *handler.TransactionHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	// middleware

	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "alive and kicking",
		})
	})

	v1 := r.Group("/v1")
	{
		v1.POST("/account", accountHandler.CreateAccount)
		v1.GET("/account/:account_id", accountHandler.GetAccount)

		v1.POST("/transaction", transactionHandler.CreateTransaction)
	}

	return r
}
