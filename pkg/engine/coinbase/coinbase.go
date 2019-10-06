package coinbase

import (
	"fmt"

	"emperror.dev/errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-resty/resty/v2"
)

// FetchOrderBook fetches coinbase pro orderbook according to supplied level
// see https://docs.pro.coinbase.com/#get-product-order-book for detailed information
// on API, level 3 API is not supported in this function because
// ratelimiting issue, to use level 3 API, use websocket-based implementation instead
func FetchOrderBook(endpoint string, level int64, pair string) error {
	levelErr := validation.Validate(level, validation.Length(1, 3))
	endpointErr := validation.Validate(endpoint, is.URL)
	err := errors.Combine(levelErr, endpointErr)
	if level == 3 {
		err = errors.Combine(err, fmt.Errorf("orderbook level 3 API is likely to be limited, use streaming instead"))
	}
	if err != nil {
		return errors.Wrap(err, "[coinbase] malformed params")
	}

	queryURL := fmt.Sprintf("%s/products/%s/book", endpoint, pair)
	client := resty.New()
	resp, err := client.R().Get(queryURL)
	// TODO: remove this after implemented
	_ = resp

	if err != nil {
		return errors.Wrapf(err, "[coinbase] failed to call GET %v", queryURL)
	}
	return nil
}
