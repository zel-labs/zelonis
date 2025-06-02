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
package main

import (
	"encoding/base64"
	"fmt"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"zelonis/params"
	"zelonis/utils/version"
	"zelonis/validator"
	"zelonis/wallet"
	"zelonis/zel/core"
)

const (
	clientIdentifier = "gzel" // Client identifier to advertise over the network

)

var app = cli.NewApp()

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func gzel(ctx *cli.Context) error {
	if args := ctx.Args().Slice(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}
	cfg := buildNode(ctx)

	vn, err := validator.New(&cfg.Node)
	if err != nil {
		return err
	}
	log.Println("Running in zelonis")
	vn.StartValidator()
	return nil
}

func buildNode(ctx *cli.Context) zelSetup {
	return zelSetup{
		Zel:  core.Defaults,
		Node: defaultNodeConfig(ctx),
	}
}

func defaultNodeConfig(ctx *cli.Context) validator.Config {

	keysFile, _ := wallet.ReadKeysFile(validator.DefaultDataDir())

	pass := wallet.GetPassword("Password:")

	mn, _ := wallet.DecryptedInfo(8, keysFile.EncryptedMnemonics[0], []byte(pass))
	key, _ := wallet.DecryptedInfo(8, keysFile.EncrytionKey[0], []byte(pass))
	recoveredWallet := wallet.RecoverWallet(key, mn)

	skey, _ := crypto.UnmarshalEd25519PrivateKey(recoveredWallet.PrivateKey)
	encoded, _ := crypto.MarshalPrivateKey(skey)

	git, _ := version.VCS()
	cfg := validator.DefaultSetup
	validatorInfo := ctx.Bool("v")
	if validatorInfo {
		cfg.Validator = true
		cfg.Stake = ctx.Float64("stake")
		if cfg.Stake < 100 {
			panic("Minimum stake amount is 100")
		}
	}
	cfg.PrivateKey = base64.StdEncoding.EncodeToString(encoded)
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit(git.Commit, git.Date)
	cfg.KeyStoreDir = validator.KeyStoreDirFlag.Name
	return cfg
}

func init() {
	app.Action = gzel
	app.Name = "gzelonis"
	app.Usage = "gzelonis [options]"
	app.Commands = []*cli.Command{
		accountCommand,
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "v",
			Usage: "is validator?",
			Value: false,
		},
		&cli.Float64Flag{
			Name:     "stake",
			Usage:    "Amount of stake in validator",
			Required: false,
			Value:    100,
		},
	}

}
