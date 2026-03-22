package dto

type TransactionRequest struct {
	AccountId       int     `json:"account_id"`
	OperationTypeId int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

type TransactionResponse struct {
	Status string `json:"status"`
}
