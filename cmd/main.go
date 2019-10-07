package main

import (
	"github.com/choestelus/super-duper-succotash/cmd/config"
	"github.com/choestelus/super-duper-succotash/pkg/engine/coinbase"
	"github.com/choestelus/super-duper-succotash/pkg/order"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// AvailableEngines contains mapping from exchanges to engines
var AvailableEngines = map[string]order.BookStreamer{
	"coinbase_pro": coinbase.Engine{},
}

func main() {
	cfg := config.MustParseConfig()
	engine := AvailableEngines[cfg.Engine].Configure(cfg.EngineConfig)

	switch cfg.Mode {
	case "oneshot":
		ExchangeOneShot(cfg, engine)
	case "service":
		ExchangeStream(cfg, engine)
	default:
		logrus.Warnf("unrecognized mode: %v", cfg.Mode)
	}
}

// ExchangeOneShot groups exchanging operations for oneshot mode together
func ExchangeOneShot(cfg config.Config, engine order.BookStreamer) {
	inputAsset := cfg.InputAsset
	outputAsset := engine.PairOf(inputAsset)

	book := engine.OneShot(cfg.EngineConfig)
	side := engine.PlaceSideToRetrieve(cfg.InputAsset)
	orders, err := book.GetOrdersBySide(side)
	if err != nil {
		logrus.Panic(err)
	}

	amount := decimal.RequireFromString(cfg.Amount)
	consumed, matched := order.MatchUntilSatisfied(side, orders, amount)

	Report(book, amount, consumed, matched, inputAsset, outputAsset)
}

// ExchangeStream groups exchanging operations for service mode together
func ExchangeStream(cfg config.Config, engine order.BookStreamer) {
	stream := engine.OpenStream(cfg.EngineConfig)
	for book := range stream {
		inputAsset := cfg.InputAsset
		outputAsset := engine.PairOf(inputAsset)

		side := engine.PlaceSideToRetrieve(cfg.InputAsset)
		orders, err := book.GetOrdersBySide(side)
		if err != nil {
			logrus.Panic(err)
		}

		amount := decimal.RequireFromString(cfg.Amount)
		consumed, matched := order.MatchUntilSatisfied(side, orders, amount)

		Report(book, amount, consumed, matched, inputAsset, outputAsset)
	}
}

// Report pretty print summary of exchange conversion rate and transaction
func Report(book order.Book, inputAmount, consumed, matched decimal.Decimal, inputAsset, outputAsset string) {
	logrus.Infof("---------------------%v---------------------------------------------", book.UpdatedAt.UTC())
	logrus.Infof("attempt to trading with\t[%v] %v", inputAmount.StringFixed(8), inputAsset)
	logrus.Infof("consumed               \t[%v] %v", consumed.StringFixed(8), inputAsset)
	logrus.Infof("got                    \t[%v] %v", matched.StringFixed(8), outputAsset)
	logrus.Infof("avg price              \t[%v] %v/%v", consumed.Div(matched).StringFixed(8), inputAsset, outputAsset)
	logrus.Infof("---------------------------------------------------------------------------------------------------------")
}
