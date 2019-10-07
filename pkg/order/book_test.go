package order

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchAsk(t *testing.T) {
	a := require.New(t)
	testcases := []struct {
		name                 string
		inputSide            string
		o                    Order
		input                decimal.Decimal
		expectedSatisfactory bool
		expectedLeft         decimal.Decimal
		expectedLeftover     Order
	}{
		{
			name:                 "[ask] matching with leftover",
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
			name:                 "[ask] match whole order",
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
			name:                 "[ask] input is too large to match with order",
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
		t.Logf("testcase: %v", tc.name)
		left, leftover, satisfied := MatchAsk(tc.o, tc.input)
		a.Equal(tc.expectedSatisfactory, satisfied)
		a.True(tc.expectedLeft.Equals(left))
		a.True(tc.expectedLeftover.Size.Equals(leftover.Size))
	}
}

func TestMatchBid(t *testing.T) {
	a := require.New(t)
	testcases := []struct {
		name                 string
		inputSide            string
		o                    Order
		input                decimal.Decimal
		expectedSatisfactory bool
		expectedLeft         decimal.Decimal
		expectedLeftover     Order
	}{
		{
			name:                 "[bid] matching with leftover",
			o:                    Order{Price: decimal.NewFromFloat(3000), Size: decimal.NewFromFloat(2)},
			input:                decimal.NewFromFloat(1),
			expectedSatisfactory: true,
			expectedLeft:         decimal.NewFromFloat(1),
			expectedLeftover: Order{
				Price: decimal.NewFromFloat(3000),
				Size:  decimal.NewFromFloat(1),
			},
		},
		{
			name:                 "[bid] match whole order",
			o:                    Order{Price: decimal.NewFromFloat(100), Size: decimal.NewFromFloat(4)},
			input:                decimal.NewFromFloat(4),
			expectedSatisfactory: true,
			expectedLeft:         decimal.NewFromFloat(0),
			expectedLeftover: Order{
				Price: decimal.NewFromFloat(100),
				Size:  decimal.NewFromFloat(0),
			},
		},
		{
			name:                 "[bid] input is too large to match with order",
			o:                    Order{Price: decimal.NewFromFloat(100), Size: decimal.NewFromFloat(1)},
			input:                decimal.NewFromFloat(10),
			expectedSatisfactory: false,
			expectedLeft:         decimal.NewFromFloat(0),
			expectedLeftover: Order{
				Price: decimal.NewFromFloat(0),
				Size:  decimal.NewFromFloat(0),
			},
		},
	}

	for _, tc := range testcases {
		t.Logf("testcase: %v", tc.name)
		left, leftover, satisfied := MatchBid(tc.o, tc.input)
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
		inputSide      string
		expectConsumed decimal.Decimal
		expectMatched  decimal.Decimal
	}{
		{
			name:           "[ask] insufficient volume",
			orders:         book.Asks,
			input:          decimal.NewFromFloat(15000),
			inputSide:      "bid",
			expectConsumed: decimal.NewFromFloat(9000),
			expectMatched:  decimal.NewFromFloat(2),
		},
		{
			name:           "[ask] match with whole volume",
			orders:         book.Asks,
			input:          decimal.NewFromFloat(9000),
			inputSide:      "bid",
			expectConsumed: decimal.NewFromFloat(9000),
			expectMatched:  decimal.NewFromFloat(2),
		},
		{
			name:           "[ask] match with leftover",
			orders:         book.Asks,
			input:          decimal.NewFromFloat(3000),
			inputSide:      "bid",
			expectConsumed: decimal.NewFromFloat(3000),
			expectMatched:  decimal.NewFromFloat(0.75),
		},
		{
			name:           "[ask] partial matching",
			orders:         book.Asks,
			input:          decimal.NewFromFloat(6000),
			inputSide:      "bid",
			expectConsumed: decimal.NewFromFloat(6000),
			expectMatched:  decimal.NewFromFloat(1.4),
		},
		{
			name:           "[bid] insufficient volume",
			orders:         book.Bids,
			input:          decimal.NewFromFloat(5),
			inputSide:      "ask",
			expectConsumed: decimal.NewFromFloat(2),
			expectMatched:  decimal.NewFromFloat(5000),
		},
		{
			name:           "[bid] match with whole volume",
			orders:         book.Bids,
			input:          decimal.NewFromFloat(2),
			inputSide:      "ask",
			expectConsumed: decimal.NewFromFloat(2),
			expectMatched:  decimal.NewFromFloat(5000),
		},
		{
			name:           "[bid] match with leftover",
			orders:         book.Bids,
			input:          decimal.NewFromFloat(0.5),
			inputSide:      "ask",
			expectConsumed: decimal.NewFromFloat(0.5),
			expectMatched:  decimal.NewFromFloat(1500),
		},
		{
			name:           "[bid] partial matching",
			orders:         book.Bids,
			input:          decimal.NewFromFloat(1.5),
			inputSide:      "ask",
			expectConsumed: decimal.NewFromFloat(1.5),
			expectMatched:  decimal.NewFromFloat(4000),
		},
	}

	for _, tc := range testcases {
		switch tc.inputSide {
		case "bid":
			t.Logf("testcase: %v", tc.name)
			consumed, matched := MatchUntilSatisfied(tc.inputSide, tc.orders, tc.input)
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
		case "ask":
			t.Logf("testcase: %v", tc.name)
			consumed, matched := MatchUntilSatisfied(tc.inputSide, tc.orders, tc.input)
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
}
