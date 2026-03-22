package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
	"transaction_routine/internal/application/usecase"
	"transaction_routine/internal/infrastructure/database/postgres"
	"transaction_routine/internal/infrastructure/http/handler"
	"transaction_routine/internal/infrastructure/logger"
	"transaction_routine/internal/infrastructure/security"

	router "transaction_routine/internal/infrastructure/http"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error is occurred  on .env file please check")
	}
}

func main() {
	dbPostgres := postgres.ConnectionDatabase()
	defer dbPostgres.Close()

	documentProtector, err := security.NewDocumentProtector()
	if err != nil {
		log.Fatal(err)
	}

	accountRepository := postgres.NewAccountRepository(dbPostgres)
	transactionRepository := postgres.NewTransactionRepository(dbPostgres)
	operationTypeRepository := postgres.NewOperationTypeRepository(dbPostgres)
	auditRepository := postgres.NewAuditRepository(dbPostgres)

	createAccountUC := usecase.NewCreateAccountUseCase(accountRepository, documentProtector, auditRepository)
	retrieveAccountUC := usecase.NewRetrieveAccountUseCase(accountRepository, documentProtector)
	createTransactionUC := usecase.NewCreateTransactionUseCase(transactionRepository, accountRepository, operationTypeRepository, auditRepository)

	accountHandler := handler.NewAccountHandler(createAccountUC, retrieveAccountUC)
	transactionHandler := handler.NewTransactionHandler(createTransactionUC)

	appLogger := logger.NewSlog()
	appRouter := router.NewRouter(accountHandler, transactionHandler, appLogger)
	srv := &http.Server{
		Addr:           ":8080",
		Handler:        appRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	appRouter.Run()
}
