package usecase

import (
	"context"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/ports"
)

type CreateAccountUseCase struct {
	accountRepo ports.AccountRepository
}

func NewCreateAccountUseCase(accountRepo ports.AccountRepository) *CreateAccountUseCase {
	return &CreateAccountUseCase{accountRepo: accountRepo}
}

func (c *CreateAccountUseCase) Create(ctx context.Context) (*dto.AccountResponse, error) {
	return nil, nil
}
