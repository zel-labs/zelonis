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
	http          *httpServer
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
		domain:        domain.NewDomain(cfg.DataDir),
	}

	if err := vn.openDataDir(); err != nil {
		return nil, err
	}
	vn.http = vn.newHTTPServer(5 * time.Second)

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

	gossipManager := gossip.NewManager(zeldb.NewDb("gossips", vn.cfg.DataDir), vn.cfg.PrivateKey, vn.domain)
	gossipManager.UpdateGossipManager(vn.domain, vn.cfg.GossipSeed, vn.cfg.GossipPort, vn.cfg.Validator, vn.cfg.Stake)
	go gossipManager.Start()

	vn.http.start()

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
	//Add gensis block to db
	//Add gensis transactions to db
	//Add gensis Account to db
	return nil
}
