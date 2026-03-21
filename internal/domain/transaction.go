package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	TransactionId int
	AccountId     int
	OperationType int
	Amount        decimal.Decimal
	CreatedAt     time.Time
}
