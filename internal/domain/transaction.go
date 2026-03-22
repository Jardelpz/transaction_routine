package domain

import (
	"math"
	"time"
)

type Transaction struct {
	TransactionId   int64
	AccountId       int64
	OperationTypeId int64
	Amount          float64
	EventDate       time.Time
}

func NormalizeAmount(operationType int64, amount float64) float64 {
	var value = math.Abs(amount)

	switch operationType {
	case 1, 2, 3:
		return -value
	case 4:
		return value
	default:
		return value
	}
}
