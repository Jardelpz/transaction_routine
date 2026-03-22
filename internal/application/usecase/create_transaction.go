package usecase

import (
	"context"
	"time"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/domain"
	"transaction_routine/internal/ports"
)

type CreateTransactionUseCase struct {
	transactionRepo   ports.TransactionRepository
	accountRepository ports.AccountRepository
	operationTypeRepo ports.OperationTypeRepository
}

func NewCreateTransactionUseCase(transactionRepo ports.TransactionRepository, accountRepository ports.AccountRepository, operationTypeRepo ports.OperationTypeRepository) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		transactionRepo:   transactionRepo,
		accountRepository: accountRepository,
		operationTypeRepo: operationTypeRepo,
	}
}

func (ct *CreateTransactionUseCase) Create(ctx context.Context, input dto.TransactionRequest) (*dto.TransactionResponse, error) {
	if input.Amount == 0 {
		return nil, domain.ErrInvalidAmount
	}

	accountExists, err := ct.accountRepository.ExistsById(ctx, input.AccountId)
	if err != nil {
		return nil, err
	}
	if !accountExists {
		return nil, domain.ErrAccountNotFound
	}

	operationExists, err := ct.operationTypeRepo.ExistsByID(ctx, input.OperationTypeId)
	if err != nil {
		return nil, err
	}
	if !operationExists {
		return nil, domain.ErrOperationTypeNotFound
	}

	amount := domain.NormalizeAmount(input.OperationTypeId, input.Amount)
	tx := domain.Transaction{
		AccountId:       input.AccountId,
		OperationTypeId: input.OperationTypeId,
		Amount:          amount,
		EventDate:       time.Now().UTC(),
	}

	err = ct.transactionRepo.Insert(ctx, tx)
	if err != nil {
		return nil, err
	}

	return &dto.TransactionResponse{Status: "created"}, nil
}
