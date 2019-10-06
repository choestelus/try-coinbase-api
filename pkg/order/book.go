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
	Tick() <-chan Book
}
