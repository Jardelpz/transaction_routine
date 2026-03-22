package usecase

import (
	"context"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/ports"
)

type RetrieveAccountUseCase struct {
	accountRepo ports.AccountRepository
}

func NewRetrieveAccountUseCase(accountRepo ports.AccountRepository) *RetrieveAccountUseCase {
	return &RetrieveAccountUseCase{accountRepo: accountRepo}
}

func (c *RetrieveAccountUseCase) Retrieve(ctx context.Context, accountId int) (*dto.AccountResponse, error) {
	account, err := c.accountRepo.GetById(ctx, accountId)
	if err != nil {
		return nil, err
	}
	return &dto.AccountResponse{
		AccountId:      account.AccountId,
		DocumentNumber: account.DocumentNumber,
	}, nil
}
