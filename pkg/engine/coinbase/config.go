package coinbase

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// Config holds necessary configuration for coinbase pro API
type Config struct {
	EngineName string `mapstructure:"name"`
	APIURL     string `mapstructure:"api_url"`
	APILevel   int64  `mapstructure:"api_level"`
}

// MustParseConfig parse config from supplied map[string]string
// crash when failed to parse.
func MustParseConfig(engineConfig map[string]string) Config {
	c := Config{}
	mstrConfig := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &c}
	decoder, err := mapstructure.NewDecoder(&mstrConfig)
	if err != nil {
		logrus.Panic(err)
	}

	err = decoder.Decode(engineConfig)
	if err != nil {
		logrus.Panic(err)
	}

	return c
}
