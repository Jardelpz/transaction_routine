package dto

type AccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type AccountResponse struct {
	AccountId      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}
