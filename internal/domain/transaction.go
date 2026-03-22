package domain

import (
	"math"
	"time"
)

type Transaction struct {
	TransactionId   int
	AccountId       int
	OperationTypeId int
	Amount          float64
	EventDate       time.Time
}

func NormalizeAmount(operationType int, amount float64) float64 {
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
