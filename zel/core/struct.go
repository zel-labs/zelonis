package core

import (
	"time"
	"zelonis/external"
)

type Config struct {
	Genesis *external.Block

	NetworkId uint64

	ZelGossipURLs []string

	TransactionHistory uint64

	SkipBcVersionCheck bool

	DocRoot string `toml:"-"`

	RPCTimeout time.Duration
}

var Defaults = Config{
	NetworkId:          0,
	RPCTimeout:         5 * time.Second,
	TransactionHistory: 10_000,
	SkipBcVersionCheck: false,
}
