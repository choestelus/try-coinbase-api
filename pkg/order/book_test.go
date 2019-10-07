package order

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	r := require.New(t)
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
		r.Equal(tc.expectedSatisfactory, satisfied)
		spew.Dump(left)
		r.True(tc.expectedLeft.Equals(left))
		r.True(tc.expectedLeftover.Size.Equals(leftover.Size))
	}
}
