package usecase

import (
	"context"
	"encoding/json"
	"strconv"
	"transaction_routine/internal/application/dto"
	"transaction_routine/internal/domain"
	"transaction_routine/internal/ports"
)

type CreateAccountUseCase struct {
	accountRepo       ports.AccountRepository
	documentProtector ports.DocumentProtector
	auditRepo         ports.AuditRepository
}

func NewCreateAccountUseCase(accountRepo ports.AccountRepository, protector ports.DocumentProtector, auditRepo ports.AuditRepository) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		accountRepo:       accountRepo,
		documentProtector: protector,
		auditRepo:         auditRepo,
	}
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

	payload, err := json.Marshal(map[string]any{
		"account_id":         response.AccountId,
		"document_encrypted": documentEncrypted,
	})
	if err == nil {
		_ = c.auditRepo.Create(ctx, domain.AuditLog{
			EventType:  "account_created",
			EntityType: "account",
			EntityID:   strconv.FormatInt(response.AccountId, 10),
			Payload:    payload,
		})
	}
	return &dto.AccountResponse{
		AccountId:      response.AccountId,
		DocumentNumber: decryptedDocument,
	}, nil
}
