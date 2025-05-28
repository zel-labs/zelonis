package main

import (
	"zelonis/validator"
	"zelonis/zel/core"
)

type zelSetup struct {
	Zel      core.Config
	Node     validator.Config
	Zelstats statsSetup
}

type statsSetup struct {
	URL string `toml:",omitempty"`
}
