package order

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	a := require.New(t)
	testcases := []struct {
		name                 string
		o                    Order
		input                decimal.Decimal
		expectedSatisfactory bool
		expectedLeft         decimal.Decimal
		expectedLeftover     Order
	}{
		{
			name:                 "matching with leftover",
			o:                    Order{Price: decimal.NewFromFloat(3000), Size: decimal.NewFromFloat(2)},
			input:                decimal.NewFromFloat(500),
			expectedSatisfactory: true,
			expectedLeft:         decimal.NewFromFloat(5500),
			expectedLeftover: Order{
				Price: decimal.NewFromFloat(3000),
				Size:  decimal.NewFromFloat(5500).Div(decimal.NewFromFloat(3000)),
			},
		},
		{
			name:                 "match whole order",
			o:                    Order{Price: decimal.NewFromFloat(100), Size: decimal.NewFromFloat(4)},
			input:                decimal.NewFromFloat(400),
			expectedSatisfactory: true,
			expectedLeft:         decimal.NewFromFloat(0),
			expectedLeftover: Order{
				Price: decimal.NewFromFloat(3000),
				Size:  decimal.NewFromFloat(0),
			},
		},
		{
			name:                 "input is too large to match with order",
			o:                    Order{Price: decimal.NewFromFloat(100), Size: decimal.NewFromFloat(1)},
			input:                decimal.NewFromFloat(5000),
			expectedSatisfactory: false,
			expectedLeft:         decimal.NewFromFloat(0),
			expectedLeftover: Order{
				Price: decimal.NewFromFloat(0),
				Size:  decimal.NewFromFloat(0),
			},
		},
	}

	for _, tc := range testcases {
		left, leftover, satisfied := Match(tc.o, tc.input)
		t.Logf("testcase: %v", tc.name)
		a.Equal(tc.expectedSatisfactory, satisfied)
		a.True(tc.expectedLeft.Equals(left))
		a.True(tc.expectedLeftover.Size.Equals(leftover.Size))
	}
}

func TestMatchUntilSatisfied(t *testing.T) {
	a := assert.New(t)

	bids := []Order{
		Order{
			Price: decimal.NewFromFloat(3000),
			Size:  decimal.NewFromFloat(1),
		},
		Order{
			Price: decimal.NewFromFloat(2000),
			Size:  decimal.NewFromFloat(1),
		},
	}
	asks := []Order{
		Order{
			Price: decimal.NewFromFloat(4000),
			Size:  decimal.NewFromFloat(1),
		},
		Order{
			Price: decimal.NewFromFloat(5000),
			Size:  decimal.NewFromFloat(1),
		},
	}
	book := Book{
		Sequence:  "0",
		Bids:      bids,
		Asks:      asks,
		UpdatedAt: time.Now(),
	}

	testcases := []struct {
		name           string
		orders         []Order
		input          decimal.Decimal
		expectConsumed decimal.Decimal
		expectMatched  decimal.Decimal
	}{
		{
			name:           "insufficient volume",
			orders:         book.Asks,
			input:          decimal.NewFromFloat(15000),
			expectConsumed: decimal.NewFromFloat(9000),
			expectMatched:  decimal.NewFromFloat(2),
		},
		{
			name:           "match with whole volume",
			orders:         book.Asks,
			input:          decimal.NewFromFloat(9000),
			expectConsumed: decimal.NewFromFloat(9000),
			expectMatched:  decimal.NewFromFloat(2),
		},
		{
			name:           "match with leftover",
			orders:         book.Asks,
			input:          decimal.NewFromFloat(3000),
			expectConsumed: decimal.NewFromFloat(3000),
			expectMatched:  decimal.NewFromFloat(0.75),
		},
		{
			name:           "partial matching",
			orders:         book.Asks,
			input:          decimal.NewFromFloat(6000),
			expectConsumed: decimal.NewFromFloat(6000),
			expectMatched:  decimal.NewFromFloat(1.4),
		},
	}

	for _, tc := range testcases {
		t.Logf("testcase: %v", tc.name)
		consumed, matched := MatchUntilSatisfied(tc.orders, tc.input)
		a.Truef(tc.expectConsumed.Equals(consumed),
			"expect consumed amount %v, got %v",
			tc.expectConsumed.String(),
			consumed.String(),
		)
		a.Truef(tc.expectMatched.Equals(matched),
			"expect matched amount %v, got %v",
			tc.expectMatched.String(),
			matched.String(),
		)
	}
}
