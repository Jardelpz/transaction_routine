package usecase

import (
	"context"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/ports"
)

type CreateTransactionUseCase struct {
	transactionRepo ports.TransactionRepository
}

func NewCreateTransactionUseCase(transactionRepo ports.TransactionRepository) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{transactionRepo: transactionRepo}
}

func (c *CreateTransactionUseCase) Create(ctx context.Context, accountId string) (*dto.TransactionResponse, error) {
	return nil, nil
}
