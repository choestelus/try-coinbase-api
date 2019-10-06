package coinbase

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/stretchr/testify/require"
)

func TestFetchOrderBook(t *testing.T) {
	coinbaseAPIURL := "https://api.pro.coinbase.com"
	tests := []struct {
		endpoint  string
		level     int64
		pair      string
		wantError bool
	}{
		{
			endpoint:  coinbaseAPIURL,
			level:     1,
			pair:      "ETH-USD",
			wantError: false,
		},
		{
			endpoint:  coinbaseAPIURL,
			level:     2,
			pair:      "BTC-USD",
			wantError: false,
		},
		{
			endpoint:  coinbaseAPIURL,
			level:     3,
			pair:      "BTC-USD",
			wantError: true,
		},
	}

	r := require.New(t)

	for _, tt := range tests {
		resp, err := FetchOrderBook(tt.endpoint, tt.level, tt.pair)
		if tt.wantError {
			r.Error(err)
		} else {
			r.NoError(err)
			r.NotNil(resp)
			JSONErr := validation.Validate(resp, is.JSON)
			r.NoError(JSONErr)
		}
	}
}
