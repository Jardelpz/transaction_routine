package dto

type TransactionRequest struct {
	AccountId       int64   `json:"account_id"`
	OperationTypeId int64   `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

type TransactionResponse struct {
	Status string `json:"status"`
}
