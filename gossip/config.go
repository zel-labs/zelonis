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
package gossip

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multiaddr"
	"log"
	"math/big"
	"reflect"
	"sync"
	"time"
	"zelonis/gossip/flow/appMsg"

	"zelonis/external"

	"zelonis/validator/domain"
	"zelonis/zeldb"
)

type Manager struct {
	db            *zeldb.ZelDB
	gossipSeed    []string
	domain        *domain.Domain
	wg            *sync.WaitGroup
	gossipPort    int
	gossipLister  *gossipLister
	privateKey    string
	gossipManager gossipFlowManager
	addr          string
	validator     bool
	stake         float64
	*external.NodeStatus
}
type gossipLister struct {
	server         *host.Host
	activeOutGoing int
	activeIncoming int
	minConnection  int
	domain         *domain.Domain
	validator      bool
	stake          float64
	privateKey     string
	*external.NodeStatus
}

func NewManager(db *zeldb.ZelDB, privateKey string, domain *domain.Domain) *Manager {
	return &Manager{
		db:            db,
		wg:            &sync.WaitGroup{},
		privateKey:    privateKey,
		gossipManager: make(gossipFlowManager, 0),
		domain:        domain,
	}
}

func (m *Manager) UpdateGossipManager(domain *domain.Domain, gossipSeed []string, gossipPort int, validator bool, stake float64) {
	m.gossipSeed = gossipSeed
	m.domain = domain
	m.gossipPort = gossipPort
	m.stake = stake
	m.validator = validator
	m.NodeStatus = external.NewNodeStatus()
}

func (m *Manager) Start() error {

	go m.startListner()
	m.startFlow()

	return nil
}

const (
	defaultListen     = "0.0.0.0"
	defaultIpType     = "ip4"
	defaultListenType = "tcp"
	defaultKey        = "12D3KooWRpVRrcc2qDM9nv4rcy4Nynb6NW1rWwG7oTVrEX1PUMvb"
	protocolID        = "/zel/0.0.3"
)

func (m *Manager) startListner() {
	//build listen address

	decoded, _ := Ed25519StringToPrivateKey(m.privateKey)

	listenAddr := fmt.Sprintf("/%s/%s/%s/%d", defaultIpType, defaultListen, defaultListenType, m.gossipPort)
	host, err := libp2p.New(libp2p.ListenAddrStrings(
		listenAddr,
	), libp2p.Identity(decoded))
	if err != nil {
		log.Fatal(err)
	}
	m.gossipLister = &gossipLister{
		server:         &host,
		privateKey:     m.privateKey,
		activeOutGoing: 0,
		activeIncoming: 0,
		minConnection:  8,
		domain:         m.domain,
		validator:      m.validator,
		stake:          m.stake,
		NodeStatus:     m.NodeStatus,
	}

	host.SetStreamHandler(protocolID, m.gossipLister.handleStream)

	var naddr string
	for _, addr := range host.Addrs() {

		naddr = fmt.Sprintf("%s/p2p/%s", addr, host.ID())
	}
	con := &connLogger{}
	host.Network().Notify(con)
	go m.checkNodeStatus()
	m.addr = naddr

	select {}
}

type connLogger struct{}

func (cl *connLogger) Listen(net network.Network, addr multiaddr.Multiaddr) {
	log.Println("Listening on:", net, addr)
}
func (cl *connLogger) ListenClose(net network.Network, addr multiaddr.Multiaddr) {}
func (cl *connLogger) Connected(net network.Network, conn network.Conn) {
	fmt.Println("✅ New peer connected:", conn.RemotePeer())
	fmt.Println("🔗 From address:", conn.RemoteMultiaddr())

}
func (cl *connLogger) Disconnected(net network.Network, conn network.Conn) {
	fmt.Println("❌ Peer disconnected:", conn.RemotePeer())
}

func Ed25519StringToPrivateKey(b64Key string) (crypto.PrivKey, error) {
	// Step 1: Decode the base64 string
	privBytes, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Step 2: Unmarshal into a libp2p private key
	priv, err := crypto.UnmarshalPrivateKey(privBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private key: %w", err)
	}

	return priv, nil
}

func toStdLibEd25519(priv crypto.PrivKey) (ed25519.PrivateKey, error) {
	raw, err := priv.Raw()
	if err != nil {
		return nil, err
	}
	return ed25519.PrivateKey(raw), nil
}

func (m *gossipLister) handleStream(s network.Stream) {

	m.handShake(s)

}

func (m *Manager) privKey() ed25519.PrivateKey {
	privKey, err := Ed25519StringToPrivateKey(m.privateKey)
	if err != nil {
		panic(err)
	}
	priKey, err := toStdLibEd25519(privKey)
	if err != nil {
		panic(err)
	}
	return priKey
}

func (m *Manager) getWalletAddress() []byte {

	pubKey := m.privKey().Public().(ed25519.PublicKey)
	return []byte(base58.Encode(pubKey))

}

var isValidatorOldRunning = false
var isValidatorRunning = false
var isValidatorEnabled = false
var txadded = false

func (m *Manager) checkNodeStatus() {
	//get wallet address from private key

	time.Sleep(3 * time.Second)
	m.checkIfValidStake()
	if !isValidatorOldRunning && !isValidatorEnabled {
		return
	}

	for {
		log.Println("Trying to start validator node")
		m.checkIfValidStake()

		if m.NodeStatus.Synced && time.Since(m.NodeStatus.SyncedTime).Seconds() >= 1000 && m.NodeStatus.LastHeight == 0 && (isValidatorOldRunning || isValidatorEnabled) && !isValidatorRunning {
			//Start validator for gensis
			//Check account balance

			isValidatorRunning = true
		} else if m.NodeStatus.Synced {
			if isValidatorOldRunning {

				m.validatorActive()

			} else if isValidatorEnabled && !txadded {
				//Add transaction to mempool
				wallet := m.getWalletAddress()
				stakeAmount := []byte(fmt.Sprintf("%v", m.stake))
				tx := m.domain.TxManager().BuildTxFromType(wallet, wallet, stakeAmount, m.LastBlockHash, external.TxStakingSend)
				sig := m.domain.TxManager().SignTxAndVerify(tx, m.privKey())
				tx.Signature = sig
				log.Printf("%x", tx.TxHash)
				m.domain.TxManager().Mempool().AddTxToMempool(tx)
				//Share tx hash to all the users

				txadded = true
			} else if !isValidatorRunning && isValidatorEnabled {
				//start validator

				isValidatorRunning = true
				m.validatorActive()
				//Propose Block
			}
			//m.domain.GetAccountBalance()

		} else {
			log.Println("Validator node is still running")
			log.Println(m.NodeStatus)

		}
		time.Sleep(3 * time.Second)
	}
}

func (m *Manager) BroadcastTransaction(tx *external.Transaction) {
	m.shareTx(tx)
}

func (m *Manager) shareTx(tx *external.Transaction) {
	for _, gossipFlow := range m.gossipManager {

		gossipFlow.zelPeer.encodeAndSend(tx, appMsg.SendInviTransaction)
	}
}

func (m *Manager) validatorActive() {
	for {
		//get order of blocks creation
		time.Sleep(335 * time.Millisecond)
		proposedBlock := m.domain.StartValidatorMode(m.privKey(), m.getWalletAddress())
		m.shareProposedBlock(proposedBlock)

	}
}

func (m *Manager) shareProposedBlock(block *external.Block) {
	for _, gossipFlow := range m.gossipManager {

		gossipFlow.zelPeer.encodeAndSend(block, appMsg.SendProposeBlock)
	}
	log.Printf("Propose Block shared %x", block.Header.BlockHash)
}

func (m *Manager) checkIfValidStake() {
	if m.validator {
		wallet := m.getWalletAddress()
		account, status := m.domain.GetAccountBalance(wallet)
		if !status {
			log.Println("Invalid account balance:", wallet)
		}
		if !reflect.DeepEqual(account.Stake, []byte("0")) {
			//Check active stake val
			isValidatorOldRunning = true
			isValidatorEnabled = true
			return
		}
		balStr := fmt.Sprintf("%s", account.Balance)
		stakeVal := big.NewFloat(m.stake)
		bal, _ := new(big.Float).SetString(balStr)

		if bal.Cmp(stakeVal) == 1 {

			isValidatorEnabled = true
			return
		}

	}
	return
}
