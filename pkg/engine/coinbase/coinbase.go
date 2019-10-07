package coinbase

import (
	"encoding/json"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/choestelus/super-duper-succotash/pkg/cast"
	"github.com/choestelus/super-duper-succotash/pkg/order"
	"github.com/davecgh/go-spew/spew"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// FetchOrderBook fetches coinbase pro orderbook according to supplied level
// see https://docs.pro.coinbase.com/#get-product-order-book for detailed information
// on API, level 3 API is not supported in this function because
// ratelimiting issue, to use level 3 API, use websocket-based implementation instead
//
// This function return raw JSON response as-is.
func FetchOrderBook(endpoint string, level int64, pair string) ([]byte, *time.Time, error) {
	// TODO: validate available pairs
	levelErr := validation.Validate(level, validation.Required, validation.Min(1), validation.Max(3))
	endpointErr := validation.Validate(endpoint, is.URL)
	err := errors.Combine(levelErr, endpointErr)
	if level == 3 {
		err = errors.Combine(err, fmt.Errorf("orderbook level 3 API is likely to be limited, use streaming instead"))
	}
	if err != nil {
		return nil, nil, errors.Wrap(err, "[coinbase] malformed params")
	}

	queryURL := fmt.Sprintf("%s/products/%s/book", endpoint, pair)
	client := resty.New()
	resp, err := client.R().Get(queryURL)

	if err != nil {
		return nil, nil, errors.Wrapf(err, "[coinbase] failed to call GET %v", queryURL)
	}

	updatedAt := time.Now()
	return resp.Body(), &updatedAt, nil
}

// ToOrder traverse through order list and transform data into order struct
func ToOrder(rawOrder interface{}) (*order.Order, error) {
	tuple3, ok := rawOrder.([]interface{})
	if !ok {
		return nil, errors.Wrap(fmt.Errorf("failed to convert order interface to slice of interface"), "ToOrder error")
	}
	priceField, err := cast.InterfaceStringToDecimal(tuple3[0])
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert price field to decimal")
	}
	sizeField, err := cast.InterfaceStringToDecimal(tuple3[1])
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert size field to decimal")
	}

	od := order.Order{Price: priceField, Size: sizeField}

	switch t3 := tuple3[2].(type) {
	case string:
		od.OrderID = t3
	case int:
		od.NumOrders = int64(t3)
	case float64:
		od.NumOrders = int64(t3)
	default:
		spew.Dump(t3)
		return nil, errors.Wrap(fmt.Errorf("failed to convert num_order/order_id to appropiate type"), "ToOrder error")
	}

	return &od, nil
}

// ToOrderBook transform raw response bytes into OrderBook struct
func ToOrderBook(raw []byte, timestamp time.Time) (*order.Book, error) {
	bookMap := map[string]interface{}{}

	if err := json.Unmarshal(raw, &bookMap); err != nil {
		return nil, errors.Wrap(err, "failed to transform coinbase API response to map[string]interface{}")
	}

	bidSlice, err := cast.InterfaceToSlice(bookMap["bids"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert raw bids interface to []interface")
	}
	askSlice, err := cast.InterfaceToSlice(bookMap["asks"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert raw asks interface to []interface")
	}

	bids := []order.Order{}
	asks := []order.Order{}

	for _, orderTuple := range bidSlice {
		od, err := ToOrder(orderTuple)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert order field into struct")
		}
		bids = append(bids, *od)
	}

	for _, orderTuple := range askSlice {
		od, err := ToOrder(orderTuple)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert order field into struct")
		}
		asks = append(asks, *od)
	}

	// while API docs example sequence is string
	// actual response is just number
	// This issue need to be addressed when use in critical
	// production system
	// TODO: try to cast sequence field as string
	seq, ok := bookMap["sequence"].(float64)
	if !ok {
		spew.Dump(bookMap["sequence"])
		return nil, errors.Wrap(fmt.Errorf("failed assert type on sequence as float64"), "interface->float64")
	}

	book := order.Book{
		Sequence:  fmt.Sprintf("%.0f", seq),
		Bids:      bids,
		Asks:      asks,
		UpdatedAt: timestamp,
	}

	return &book, nil

}

// StreamOrderBook streams orderbook from websocket
// and wrap into channel
func StreamOrderBook(endpoint string, pair string) <-chan order.Book {
	logrus.Warn("not implemented")
	return nil
}

// MustFetch fetches orderbook and transform into Book struct
// panic when failed
func MustFetch(endpoint string, level int64, pair string) order.Book {
	resp, updatedAt, err := FetchOrderBook(endpoint, level, pair)
	if err != nil {
		logrus.Panicf("failed to fetch order book: %v", err)
	}
	book, err := ToOrderBook(resp, *updatedAt)
	if err != nil {
		logrus.Panicf("failed to transform response to order book: %v", err)
	}
	return *book
}

// FetchStream wrap FetchOrderbook and return order book channel
func FetchStream(interval time.Duration, endpoint string, level int64, pair string) <-chan order.Book {
	bookStream := make(chan order.Book)
	go func(endpoint string, level int64, pair string) {
		for range time.Tick(interval) {
			bookStream <- MustFetch(endpoint, level, pair)
		}
	}(endpoint, level, pair)
	return bookStream
}
