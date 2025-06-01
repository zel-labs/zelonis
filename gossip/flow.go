package gossip

import (
	"capnproto.org/go/capnp/v3"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/mr-tron/base58"
	ma "github.com/multiformats/go-multiaddr"
	"log"
	"net"
	"time"
	flowv1 "zelonis/gossip/flow"
	"zelonis/gossip/flow/appMsg"
)

type seeders []*seeder
type seeder struct {
	seed string
	ip   net.IP
	addr string
}

func (m *Manager) startFlow() {

	//get addr from database if database empty use gossip seed
	for {
		nseeder := m.gossipAddr(m.gossipSeeder())
		//Convert seeder to libp2p address

		err := m.createGossip(nseeder.gossipAddr())
		if err != nil {
			log.Println(err)
		}

		for _, gossipM := range m.gossipManager {
			if gossipM.active {

				continue
			}

			gossipM.startGossip()
		}
		log.Println("Gossip Manager is active ", len(m.gossipManager))
		time.Sleep(10 * time.Second)
	}

}

type gossipFlowManager []*gossip

type gossip struct {
	gossipAddr string

	host host.Host
	*zelPeer
	cancel     context.CancelFunc
	lastPing   time.Time
	isIncoming bool
	isOutgoing bool
	isSyncing  bool
	active     bool
}

func (g *gossip) startGossip() error {
	g.active = true

	hashshake := appMsg.NewHandshakeWithInfo()

	msg, err := hashshake.Encode()
	if err != nil {
		return err
	}
	err = g.zelPeer.sendMsg(msg)
	if err != nil {
		return err
	}

	msg, err = g.zelPeer.getMsg()
	if err != nil {
		return err
	}
	msgHandshake := appMsg.NewHandshake()
	if err = msgHandshake.Decode(msg); err != nil {
		return err
	}

	if err = msgHandshake.Verify(hashshake); err != nil {
		log.Println(hashshake.Version)
		log.Println(msgHandshake.Version)
		log.Println(err)
		return err
	}

	flow := flowv1.CreateFollow(g.zelPeer.encoder, g.zelPeer.decoder, g.zelPeer.conn, g.zelPeer.domain, g.zelPeer.validator, g.zelPeer.stake, g.zelPeer.NodeStatus)
	flow.Start(1)
	return nil
}

func (m *Manager) createGossip(s seeders) error {

	gossipFlows := make(gossipFlowManager, 0)

	for _, val := range s {
		//Check if peer is already being used
		if cGossip, status := m.checkIfGossipActive(val.addr); status {

			gossipFlows = append(gossipFlows, cGossip)
			continue
		}
		peerAddr := m.addr

		ctx, cancel := context.WithCancel(context.Background())
		host, err := libp2p.New()
		if err != nil {
			return err
		}

		maddr, err := ma.NewMultiaddr(peerAddr)
		if err != nil {
			return err
		}
		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			return err
		}
		if err := host.Connect(ctx, *info); err != nil {
			return err
		}

		stream, err := host.NewStream(ctx, info.ID, protocolID)
		if err != nil {
			return err
		}
		encoder := capnp.NewEncoder(stream)
		decoder := capnp.NewDecoder(stream)
		z := &zelPeer{
			conn:       stream.Conn(),
			encoder:    encoder,
			decoder:    decoder,
			handshake:  false,
			domain:     m.domain,
			validator:  m.validator,
			stake:      m.stake,
			NodeStatus: m.NodeStatus,
		}

		gossipP2P := &gossip{
			gossipAddr: val.addr,
			cancel:     cancel,
			host:       host,
			zelPeer:    z,
			lastPing:   time.Now(),
			active:     false,
		}
		gossipFlows = append(gossipFlows, gossipP2P)
	}

	m.gossipManager = gossipFlows
	return nil
}

func (m *Manager) checkIfGossipActive(addr string) (*gossip, bool) {
	currentGossip := m.gossipManager
	for _, gossipP2P := range currentGossip {
		if gossipP2P.gossipAddr == addr && gossipP2P.active {
			return gossipP2P, true
		}
	}
	return nil, false
}

func (s seeders) gossipAddr() seeders {

	for key, val := range s {

		s[key].addr = fmt.Sprintf("/ip4/%s/tcp/30331/p2p/%s", val.ip.String(), defaultKey)
	}
	return s
}

func (m *Manager) gossipAddr(seed []string) seeders {
	nseeders := make(seeders, 0)
	for _, addr := range seed {

		ips, err := net.LookupIP(addr)
		if err != nil {
			continue
		}
		nseeder := &seeder{
			seed: addr,
			ip:   ips[0],
		}
		nseeders = append(nseeders, nseeder)

	}
	return nseeders
}

func (m *Manager) gossipSeeder() []string {
	seedStr := make([]string, 0)
	for _, seed := range m.gossipSeed {
		nseed := m.decryptSeed(seed)

		seedStr = append(seedStr, nseed)
	}
	return seedStr
}

func (m *Manager) discoverSeed(seed []byte) string {
	return fmt.Sprintf("%sco%s.%s", "dis", "ver", seed)
}

func (m *Manager) decryptSeed(str string) string {

	decoded, _ := hex.DecodeString(str)
	decoded, _ = base58.Decode(string(decoded))
	return m.discoverSeed(decoded)
}
