package coinbase

import (
	"strings"
	"time"

	"github.com/choestelus/super-duper-succotash/pkg/order"
	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// Engine holds necessary configuration for coinbase pro API
type Engine struct {
	APIURL       string        `mapstructure:"api_url"`
	APILevel     int64         `mapstructure:"api_level"`
	Pair         string        `mapstructure:"pair"`
	PollInterval time.Duration `mapstructure:"poll_interval"`
}

// MustParseConfig parse config from supplied map[string]string
// crash when failed to parse.
func MustParseConfig(engineConfig map[string]string) Engine {
	e := Engine{}
	spew.Dump(engineConfig)
	mstrConfig := mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: true,
		Result:           &e,
	}
	decoder, err := mapstructure.NewDecoder(&mstrConfig)
	if err != nil {
		logrus.Panic(err)
	}

	err = decoder.Decode(engineConfig)
	if err != nil {
		logrus.Panic(err)
	}

	return e
}

// OpenStream streams orderbook with supplied configuration
func (e Engine) OpenStream(cfg map[string]string) <-chan order.Book {
	return FetchStream(e.PollInterval, e.APIURL, e.APILevel, e.Pair)
}

// Configure set self configuration with supplied args
func (e Engine) Configure(cfg map[string]string) order.BookStreamer {
	return MustParseConfig(cfg)
}

// AssetPair returns main asset and exchanging asset of pair
// e.g. BTC-USD main asset would be BTC and exchanging asset would be USD
func (e Engine) AssetPair() (string, string) {
	splittedPair := strings.Split("-", e.Pair)
	main, exchanging := splittedPair[0], splittedPair[1]
	return main, exchanging
}

// PlaceSideToRetrieve returns which side to place in order to get
// asset specified.
func (e Engine) PlaceSideToRetrieve(asset string) string {
	main, exchanging := e.AssetPair()
	switch {
	case asset == main:
		return "bid"
	case asset == exchanging:
		return "ask"
	default:
		return "invalid"
	}
}
