package usecase

import (
	"context"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/ports"
)

type RetrieveAccountUseCase struct {
	accountRepo       ports.AccountRepository
	documentProtector ports.DocumentProtector
}

func NewRetrieveAccountUseCase(accountRepo ports.AccountRepository, documentProtector ports.DocumentProtector) *RetrieveAccountUseCase {
	return &RetrieveAccountUseCase{accountRepo: accountRepo, documentProtector: documentProtector}
}

func (c *RetrieveAccountUseCase) Retrieve(ctx context.Context, accountId int) (*dto.AccountResponse, error) {
	account, err := c.accountRepo.GetById(ctx, accountId)
	if err != nil {
		return nil, err
	}

	document, err := c.documentProtector.Decrypt(account.DocumentEncrypted)
	if err != nil {
		return nil, err
	}
	return &dto.AccountResponse{
		AccountId:      account.AccountId,
		DocumentNumber: document,
	}, nil
}
