package config

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type Config struct {
	Mode         string            `long:"mode" required:"true" choice:"oneshot" choice:"service" description:"select wheter to run as oneshot or until manually stop"`
	EngineConfig map[string]string `long:"engine-config" required:"true" description:"configuration for exchange engine, in key:value format"`
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
