package order

import (
	"time"

	"github.com/shopspring/decimal"
)

// Order holds information of aggregated orderbook
// From exchange, NumOrders denotes number of aggregated orders
// and should not use for multiplying with Size
type Order struct {
	OrderID   string          `json:"order_id,omitempty"`
	Price     decimal.Decimal `json:"price"`
	Size      decimal.Decimal `json:"size"`
	NumOrders int64           `json:"num_orders,omitempty"`
}

// Volume return price*size of order
func (o Order) Volume() decimal.Decimal {
	return o.Price.Mul(o.Size)
}

// MatchAsk match input against ask order with supplied amount and returns taken amount, remaining order.
// returned boolean value denotes if matched order is satisfied for amount or not.
// unfortunately, iterator function or applying operation on functor is not idiomatic in go
// otherwise this could be rewrite to something simpler as fold operation.
func MatchAsk(o Order, input decimal.Decimal) (decimal.Decimal, Order, bool) {
	volume := o.Volume()
	if volume.GreaterThanOrEqual(input) {
		left := volume.Sub(input)

		leftSize := left.Div(o.Price)
		deductedOrder := o
		deductedOrder.Size = leftSize
		return left, deductedOrder, true
	}
	// if input amount is too large for order to match with
	left := decimal.Zero
	deductedOrder := o
	deductedOrder.Size = decimal.NewFromFloat(0)
	return left, deductedOrder, false
}

// MatchBid match input against bid order with supplied amount and returns taken amount, remaining order.
// returned boolean value denotes if matched order is satisfied for amount or not.
// unfortunately, iterator function or applying operation on functor is not idiomatic in go
// otherwise this could be rewrite to something simpler as fold operation.
func MatchBid(o Order, input decimal.Decimal) (decimal.Decimal, Order, bool) {
	input = input.Mul(o.Price)
	volume := o.Volume()
	if volume.GreaterThanOrEqual(input) {
		leftSize := volume.Sub(input)

		left := leftSize.Div(o.Price)
		deductedOrder := o
		deductedOrder.Size = left
		return left, deductedOrder, true
	}
	// if input amount is too large for order to match with
	left := decimal.Zero
	deductedOrder := o
	deductedOrder.Size = decimal.NewFromFloat(0)
	return left, deductedOrder, false
}

// MatchUntilSatisfied fold(left) over "sorted" orders and return consumed input and amount matched
func MatchUntilSatisfied(side string, ods []Order, amount decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	matched := decimal.Zero
	consumed := decimal.Zero
	for _, od := range ods {
		input := amount.Sub(consumed)

		left := decimal.Zero
		satisfied := false
		switch side {
		case "bid":
			left, _, satisfied = MatchAsk(od, input)
			consumed = consumed.Add(od.Volume().Sub(left))
			matched = matched.Add(od.Volume().Sub(left).Div(od.Price))
		case "ask":
			left, _, satisfied = MatchBid(od, input)
			consumed = consumed.Add(od.Size.Sub(left))
			matched = matched.Add(od.Volume().Sub(left.Mul(od.Price)))
		default:
			panic("unexpected side value")
		}

		if satisfied {
			break
		}
	}
	return consumed, matched
}

// Book holds general info of orderbook list
// both bid side and ask side along with retrieval timestamp
type Book struct {
	Sequence  string    `json:"sequence"`
	Bids      []Order   `json:"bids"`
	Asks      []Order   `json:"asks"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BookStreamer is main interface for using with
// streaming API part, underlying that's synchronous
// blocking API, should be wrapped with channel
type BookStreamer interface {
	OneShot(config map[string]string) Book
	OpenStream(config map[string]string) <-chan Book
	Configure(config map[string]string) BookStreamer
	PlaceSideToRetrieve() string
}
