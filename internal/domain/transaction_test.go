package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizeAmount(t *testing.T) {
	tests := []struct {
		name        string
		opType      int
		amount      float64
		want        float64
		description string
	}{
		{
			name:        "operation 1 normal purchase - negative",
			opType:      1,
			amount:      100.50,
			want:        -100.50,
			description: "Normal Purchase should negate amount",
		},
		{
			name:        "operation 2 installments - negative",
			opType:      2,
			amount:      200,
			want:        -200,
			description: "Purchase with installments should negate amount",
		},
		{
			name:        "operation 3 withdrawal - negative",
			opType:      3,
			amount:      50.25,
			want:        -50.25,
			description: "Withdrawal should negate amount",
		},
		{
			name:        "operation 4 credit voucher - positive",
			opType:      4,
			amount:      75.99,
			want:        75.99,
			description: "Credit Voucher should keep amount positive",
		},
		{
			name:        "negative input normalized to abs then negated",
			opType:      1,
			amount:      -100,
			want:        -100,
			description: "Negative input for op 1 should become -100 (abs then negate)",
		},
		{
			name:        "unknown operation type - positive",
			opType:      99,
			amount:      100,
			want:        100,
			description: "Unknown operation type should keep positive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeAmount(tt.opType, tt.amount)
			assert.Equal(t, tt.want, got)
		})
	}
}
