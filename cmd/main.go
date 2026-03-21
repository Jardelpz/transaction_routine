package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
	"transaction_routine/internal/application/usecase"
	"transaction_routine/internal/infrastructure/database/postgres"
	"transaction_routine/internal/infrastructure/http/handler"

	infra "transaction_routine/internal/infrastructure/http"
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

	accountRepository := postgres.NewAccountRepository(dbPostgres)
	transactionRepository := postgres.NewTransactionRepository(dbPostgres)

	createAccountUC := usecase.NewCreateAccountUseCase(accountRepository)
	retrieveAccountUC := usecase.NewRetrieveAccountUseCase(accountRepository)
	createTransactionUC := usecase.NewCreateTransactionUseCase(transactionRepository)

	accountHandler := handler.NewAccountHandler(createAccountUC, retrieveAccountUC)
	transactionHandler := handler.NewTransactionHandler(createTransactionUC)

	router := infra.NewRouter(accountHandler, transactionHandler)
	srv := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	router.Run()
}
