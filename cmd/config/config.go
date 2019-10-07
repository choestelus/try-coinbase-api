package config

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type Config struct {
	Amount       string            `short:"a" long:"amount" required:"true" description:"input amount to calculate"`
	InputAsset   string            `short:"i" long:"input-asset" required:"true" description:"input asset type, output asset type will be automatically set via pair config according to exchange engine, if available"`
	OutputAsset  string            `short:"o" long:"output-asset" required:"false" description:"output asset type, can be set if engine support exchange routing with more than 1 pair"`
	Mode         string            `short:"m" long:"mode" required:"true" choice:"oneshot" choice:"service" description:"select wheter to run as oneshot or until manually stop"`
	Engine       string            `short:"E" long:"engine" required:"true" choice:"coinbase_pro" description:"select exchange engine to use"`
	EngineConfig map[string]string `short:"e" long:"engine-config" required:"true" description:"configuration for exchange engine, in key:value format, one pair per each flag"`
}

func MustParseConfig() Config {
	cfg := Config{}

	parser := flags.NewParser(&cfg, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			parser.WriteHelp(os.Stderr)
			os.Exit(1)
		}
	}
	return cfg
}
