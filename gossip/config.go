package gossip

import (
	"encoding/base64"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
	"log"
	"sync"
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
}
type gossipLister struct {
	server         *host.Host
	activeOutGoing int
	activeIncoming int
	minConnection  int
	domain         *domain.Domain
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

func (m *Manager) UpdateGossipManager(domain *domain.Domain, gossipSeed []string, gossipPort int) {
	m.gossipSeed = gossipSeed
	m.domain = domain
	m.gossipPort = gossipPort
}

func (m *Manager) Start() error {
	log.Println("Running Gossip Manager")
	go m.startListner()
	m.startFlow()
	return nil
}

const (
	defaultListen     = "0.0.0.0"
	defaultIpType     = "ip4"
	defaultListenType = "tcp"
	defaultKey        = "12D3KooWRpVRrcc2qDM9nv4rcy4Nynb6NW1rWwG7oTVrEX1PUMvb"
	protocolID        = "/zel/1.0.0"
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
		activeOutGoing: 0,
		activeIncoming: 0,
		minConnection:  8,
		domain:         m.domain,
	}

	host.SetStreamHandler(protocolID, m.gossipLister.handleStream)

	var naddr string
	for _, addr := range host.Addrs() {

		naddr = fmt.Sprintf("%s/p2p/%s", addr, host.ID())
	}
	con := &connLogger{}
	host.Network().Notify(con)

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

func (m *gossipLister) handleStream(s network.Stream) {

	m.handShake(s)

}
