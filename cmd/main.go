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
		orderbook := engine.OneShot(cfg.EngineConfig)
		side := engine.PlaceSideToRetrieve(cfg.InputAsset)
		orders, err := orderbook.GetOrdersBySide(side)
		if err != nil {
			logrus.Panic(err)
		}
		amount := decimal.RequireFromString(cfg.Amount)

		consumed, matched := order.MatchUntilSatisfied(side, orders, amount)
		inputAsset, outputAsset := engine.AssetPair()
		logrus.Infof("consumed  [%v] %v", consumed.StringFixed(8), inputAsset)
		logrus.Infof("got       [%v] %v", matched.StringFixed(8), outputAsset)
		logrus.Infof("avg price [%v] %v/%v", consumed.Div(matched).StringFixed(8), inputAsset, outputAsset)
	case "service":
		logrus.Panicf("not implemented: %v", cfg.Mode)
		// bookstream := engine.OpenStream(cfg.EngineConfig)

	default:
		logrus.Warnf("unrecognized mode: %v", cfg.Mode)
	}

}
