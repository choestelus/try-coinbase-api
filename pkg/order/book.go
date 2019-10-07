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

// Match match order with supplied amount and returns taken amount and remaining order
func Match(o Order, amount decimal.Decimal) (decimal.Decimal, Order) {
	volume := o.Price.Mul(o.Size)
	if volume.GreaterThanOrEqual(amount) {
		taken := volume.Sub(amount)

		takenSize := taken.Div(o.Price)
		deductedOrder := o
		deductedOrder.Size = o.Size.Sub(takenSize)
		return taken, deductedOrder
	}
	// if amount is too large for order
	taken := amount
	deductedOrder := o
	deductedOrder.Size = decimal.NewFromFloat(0)
	return taken, deductedOrder
}

// Book holds general info of orderbook list
// both bid side and ask side along with retrieval timestamp
type Book struct {
	Sequence  string    `json:"sequence"`
	Bids      []Order   `json:"bids"`
	Asks      []Order   `json:"asks"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PlaceOrder simulates order placing into exchange, and return
// converted amount according to orderbook
// bid means buying order, ask means selling order, in user perspective
// but side argument determines which side order should be placed
// e.g. if we want to buy BTC in BTC-USD pair, we should place order at bid side
func (b Book) PlaceOrder(amount decimal.Decimal, side string) (decimal.Decimal, error) {
	switch side {
	case "bid":
	case "ask":
	default:
	}
	return decimal.Zero, nil
}

// BookStreamer is main interface for using with
// streaming API part, underlying that's synchronous
// blocking API, should be wrapped with channel
type BookStreamer interface {
	OpenStream(config map[string]string) <-chan Book
	Configure(config map[string]string) BookStreamer
	PlaceSideToRetrieve() string
}
