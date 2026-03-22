package usecase

import (
	"context"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/domain"
	"transaction_routine/internal/ports"
)

type CreateAccountUseCase struct {
	accountRepo       ports.AccountRepository
	documentProtector ports.DocumentProtector
}

func NewCreateAccountUseCase(accountRepo ports.AccountRepository, protector ports.DocumentProtector) *CreateAccountUseCase {
	return &CreateAccountUseCase{accountRepo: accountRepo, documentProtector: protector}
}

func (c *CreateAccountUseCase) Create(ctx context.Context, request dto.AccountRequest) (*dto.AccountResponse, error) {
	err := domain.ValidateDocument(request.DocumentNumber)
	if err != nil {
		return nil, err
	}

	documentHash := c.documentProtector.Hash(request.DocumentNumber)              // usado para comparações (deterministico) se o document ja existe por exemplo
	documentEncrypted, err := c.documentProtector.Encrypt(request.DocumentNumber) // reversivel
	if err != nil {
		return nil, err
	}

	response, err := c.accountRepo.Insert(ctx, documentHash, documentEncrypted)
	if err != nil {
		return nil, err
	}

	decryptedDocument, err := c.documentProtector.Decrypt(response.DocumentEncrypted)
	if err != nil {
		return nil, err
	}
	return &dto.AccountResponse{
		AccountId:      response.AccountId,
		DocumentNumber: decryptedDocument,
	}, nil
}
