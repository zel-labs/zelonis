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

package validator

import (
	"errors"
	"github.com/dgraph-io/badger/v4"
	"github.com/gofrs/flock"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"zelonis/gossip"
	"zelonis/utils/logger"
	"zelonis/validator/domain"
	"zelonis/zeldb"
)

type Validator struct {
	cfg *Config

	stop          chan struct{}
	rpc           *RpcServer
	startStopLock sync.Mutex
	log           logger.Logger
	dirLock       *flock.Flock

	domain *domain.Domain
}

type DbTrackers struct {
	name string
	*badger.DB
}

func New(cfg *Config) (*Validator, error) {
	cfgCopy := *cfg
	cfg = &cfgCopy
	if cfg.DataDir != "" {
		absdatadir, err := filepath.Abs(cfg.DataDir)
		log.Println(absdatadir)

		if err != nil {
			return nil, err
		}
		cfg.DataDir = absdatadir
	}

	if strings.ContainsAny(cfg.Name, `/\`) {
		return nil, errors.New(`Config.Name must not contain '/' or '\'`)
	}
	if cfg.Name == datadirDefaultKeyStore {
		return nil, errors.New(`Config.Name cannot be "` + datadirDefaultKeyStore + `"`)
	}
	if strings.HasSuffix(cfg.Name, ".ipc") {
		return nil, errors.New(`Config.Name cannot end in ".ipc"`)
	}

	vn := &Validator{
		cfg:           cfg,
		stop:          make(chan struct{}),
		startStopLock: sync.Mutex{},
		domain:        domain.NewDomain(cfg.DataDir, cfg.ReVerify),
	}

	if err := vn.openDataDir(); err != nil {
		return nil, err
	}
	vn.rpc = NewHTTPServer(5*time.Second, vn.domain)

	return vn, nil
}

func (vn *Validator) StartValidator() error {

	err := vn.startValidator(false)
	if err != nil {
		return err
	}
	return nil
}

func (vn *Validator) startValidator(isConsole bool) error {
	//Get Nodes
	if err := vn.Start(); err != nil {
		logger.Crit("Error starting protocol stack: %v", err)
	}
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		shutdown := func() {
			logger.Info("Got interrupt, shutting down...")
			go vn.Close()
			for i := 10; i > 0; i-- {
				<-sigc
				if i > 1 {
					logger.Warn("Already shutting down, interrupt more to panic.", "times", i-1)
				}
			}

		}
		if isConsole {
			// In JS console mode, SIGINT is ignored because it's handled by the console.
			// However, SIGTERM still shuts down the node.
			for {
				sig := <-sigc
				if sig == syscall.SIGTERM {
					shutdown()
					return
				}
			}
		} else {
			<-sigc
			shutdown()
		}

	}()
	return nil
}

func (vn *Validator) Close() error {
	vn.startStopLock.Lock()
	defer vn.startStopLock.Unlock()

	return ErrNodeStopped

}

func (vn *Validator) Start() error {
	vn.startStopLock.Lock()
	defer vn.startStopLock.Unlock()
	vn.verifyGenesis()

	gossipManager := gossip.NewManager(zeldb.NewDb("gossips", vn.cfg.DataDir, false), vn.cfg.PrivateKey, vn.domain)
	gossipManager.UpdateGossipManager(vn.domain, vn.cfg.GossipSeed, vn.cfg.GossipPort, vn.cfg.Validator, vn.cfg.Stake)
	go gossipManager.Start()

	vn.rpc.Start(gossipManager.FlowManager())

	return nil
}

func (vn *Validator) verifyGenesis() error {
	//Create gensis Block from info
	log.Println("Verifying genesis block...")
	//Build Gensis Header
	genesisBlock := vn.generateGenesis()
	status, err := vn.domain.VerifyInsertBlockAndTransaction(genesisBlock)
	if err != nil {
		log.Println(err)
		return err
	}
	if !status {
		panic("genesis block could not be verified")
	}
	if vn.cfg.ReVerify {
		vn.domain.ReverifyTx()
	}
	return nil
}
