package dto

import "github.com/shopspring/decimal"

type TransactionRequest struct {
	AccountId     string          `json:"account_id"`
	OperationType string          `json:"operation_type_id"`
	Amount        decimal.Decimal `json:"amount"`
}

type TransactionResponse struct {
	Status string `json:"status"`
}
