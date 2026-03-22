package usecase

import (
	"context"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/domain"
	"transaction_routine/internal/ports"
)

type CreateAccountUseCase struct {
	accountRepo ports.AccountRepository
}

func NewCreateAccountUseCase(accountRepo ports.AccountRepository) *CreateAccountUseCase {
	return &CreateAccountUseCase{accountRepo: accountRepo}
}

func (c *CreateAccountUseCase) Create(ctx context.Context, request dto.AccountRequest) (*dto.AccountResponse, error) {
	err := domain.ValidateDocument(request.DocumentNumber)
	if err != nil {
		return nil, err
	}

	response, err := c.accountRepo.Insert(ctx, request.DocumentNumber)
	if err != nil {
		return nil, err
	}

	return &dto.AccountResponse{
		AccountId:      response.AccountId,
		DocumentNumber: response.DocumentNumber,
	}, nil
}
