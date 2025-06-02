/*
Copyright (C) 2025 Zelonis Contributors

This file is part of Zelonis.

Zelonis is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Zelonis is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Zelonis. If not, see <https://www.gnu.org/licenses/>.
*/

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
