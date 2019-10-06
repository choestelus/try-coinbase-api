package main

import (
	"github.com/choestelus/super-duper-succotash/cmd/config"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	cfg := config.MustParseConfig()
	spew.Dump(cfg)
}
