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
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"zelonis/internal/flags"
	"zelonis/validator"
	"zelonis/wallet"
)

var (
	accountCommand = &cli.Command{
		Name:        "account",
		Usage:       "manage account",
		Description: "Manage accounts, list all existing accounts, import a private keys into a new\naccount, create a new account or update an existing account.",
		Subcommands: []*cli.Command{
			{
				Name: "list",

				Usage:  "Print summary of existing accounts",
				Action: accountList,
				Flags: []cli.Flag{
					&flags.DirectoryFlag{
						Name:  "datadir",
						Usage: "Data directory for the databases and keystore",
						Value: flags.DirectoryString(validator.DefaultDataDir()),
					},
				},

				Description: `Print a short summary of all accounts`,
			},
			{
				Name: "recover",

				Usage:  "Recover account",
				Action: accountRecover,
				Flags: []cli.Flag{
					&flags.DirectoryFlag{
						Name:  "datadir",
						Usage: "Data directory for the databases and keystore",
						Value: flags.DirectoryString(validator.DefaultDataDir()),
					},
				},

				Description: `Print a short summary of all accounts`,
			},
		},
	}
)

func accountList(c *cli.Context) error {

	wallet.ListAddress(validator.DefaultDataDir())
	return nil
}

func accountRecover(c *cli.Context) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Seed Phrase:")

	seedPhrase, _ := readLine(reader)

	fmt.Print("Enter Encrypted Key (Numeric):")
	input, _ := readLine(reader)

	//make sure its numeric

	inputInt, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		log.Fatal("Invalid Encrypted Key")
		return nil
	}
	key := input
	fmt.Println("Enter Encrypted Key (Numeric):", inputInt)
	fmt.Println("Seed Phrase entered", seedPhrase)
	fmt.Println("Please confirm (Y or N):")
	input, _ = readLine(reader)

	if input != "y" && input != "Y" {
		fmt.Println("Try Again...")
		return nil
	}
	//Recover wallet address
	rWallet := wallet.RecoverWallet(key, seedPhrase)
	rWallet.CreateKeyFile(validator.DefaultDataDir())

	fmt.Println("Wallet Address:", rWallet.Address)

	return nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return "", errors.WithStack(err)
	}

	return strings.TrimSpace(string(line)), nil
}
