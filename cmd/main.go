package main

import (
	"github.com/choestelus/super-duper-succotash/cmd/config"
	"github.com/choestelus/super-duper-succotash/pkg/engine/coinbase"
	"github.com/choestelus/super-duper-succotash/pkg/order"
)

// AvailableEngines contains mapping from exchanges to engines
var AvailableEngines = map[string]order.BookStreamer{
	"coinbase_pro": &coinbase.Engine{},
}

func main() {
	cfg := config.MustParseConfig()
	engine := AvailableEngines[cfg.Engine]
	streamer := engine.MustParseConfig(cfg.EngineConfig)

}
