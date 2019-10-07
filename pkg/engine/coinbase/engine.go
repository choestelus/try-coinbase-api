package coinbase

import (
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
