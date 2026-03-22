package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"transaction_routine/api"
	"transaction_routine/internal/infrastructure/http/handler"
)

func NewRouter(accountHandler *handler.AccountHandler, transactionHandler *handler.TransactionHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	// middleware

	// OpenAPI spec (YAML)
	r.GET("/openapi.yaml", func(c *gin.Context) {
		c.Data(200, "application/x-yaml", api.OpenAPISpec)
	})

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/openapi.yaml")))

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
