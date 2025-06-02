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

package wallet

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/term"
	"os"
	"os/signal"
	"syscall"
)

type PublicWalletAddress struct {
	Address []byte
}
type encryptedWalletInfo struct {
	PrivateKey    []*EncryptedInfo
	WalletAddress []*PublicWalletAddress
	EncryptionKey []*EncryptedInfo
}

func (w *walletInfo) recoverWithin() (*encryptedWalletInfo, error) {

	password := []byte(GetPassword("Enter password for the keys file:"))
	confirmPassword := []byte(GetPassword("Confirm password:"))
	if subtle.ConstantTimeCompare(password, confirmPassword) != 1 {
		return nil, errors.New("Passwords are not identical")
	}
	encryptedSeed, err := encryptInfo(w.Mnemonic, password)
	if err != nil {
		return nil, err
	}
	encryptedKey, err := encryptInfo(fmt.Sprintf("%v", w.EncryptKey), password)
	if err != nil {
		return nil, err
	}
	addr := &PublicWalletAddress{
		Address: []byte(w.Address),
	}
	return &encryptedWalletInfo{
		WalletAddress: []*PublicWalletAddress{addr},
		PrivateKey:    []*EncryptedInfo{encryptedSeed},
		EncryptionKey: []*EncryptedInfo{encryptedKey},
	}, nil
}
func (w *walletInfo) CreateKeyFile(dir string) error {
	walletInfo, err := w.recoverWithin()
	if err != nil {
		return err
	}
	file := File{
		Version:            LastVersion,
		EncryptedMnemonics: walletInfo.PrivateKey,
		ExtendedAddress:    walletInfo.WalletAddress,

		EncrytionKey: walletInfo.EncryptionKey,
	}
	err = file.SetPath(dir, true)
	if err != nil {
		return err
	}
	err = file.TryLock()
	if err != nil {
		return err
	}

	err = file.Save()
	if err != nil {
		return err
	}

	fmt.Printf("Wrote the keys into %s\n", file.Path())
	return nil
}

func GetPassword(prompt string) string {
	// Get the initial state of the terminal.
	initialTermState, e1 := term.GetState(int(syscall.Stdin))
	if e1 != nil {
		panic(e1)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		_ = term.Restore(int(syscall.Stdin), initialTermState)
		os.Exit(1)
	}()

	fmt.Print(prompt)
	p, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		panic(err)
	}

	signal.Stop(c)

	return string(p)
}

func (d *File) Save() error {
	if d.path == "" {
		return errors.New("cannot save a file with uninitialized path")
	}

	err := createFileDirectoryIfDoesntExist(d.path)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(d.path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(d.toJSON())
	if err != nil {
		return err
	}

	return nil
}

func (d *File) Path() string {
	return d.path
}
