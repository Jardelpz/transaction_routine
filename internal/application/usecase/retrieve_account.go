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

func (c *RetrieveAccountUseCase) Retrieve(ctx context.Context) (*dto.AccountResponse, error) {
	return nil, nil
}
