package main

import (
	"encoding/base64"
	"fmt"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"
	"log"
	"strconv"
	"time"
	ping "zelonis/capn"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"

	capnp "capnproto.org/go/capnp/v3"
	// Adjust to your module path
)

const protocolID = "/zel/1.0.0"

func main() {

	decoded, _ := base64.StdEncoding.DecodeString("CAESQLpm+VGlyR0uSl8P2q+CnqagCtjIyw9acxYKJjU1OgxI7cMzShfsG/B2ojnjRPv7KWfyyVGLAwmVw2xgyZUROYw=")
	privRestored, _ := crypto.UnmarshalPrivateKey(decoded)

	host, err := libp2p.New(libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0", // Listen on all interfaces on a random available port
	), libp2p.Identity(privRestored))
	if err != nil {
		log.Fatal(err)
	}

	host.SetStreamHandler(protocolID, handleStream)
	fmt.Println("Server peer ID:", host.ID(), host.Addrs())
	var naddr string
	for _, addr := range host.Addrs() {
		fmt.Printf("Listening on: %s/p2p/%s\n", addr, host.ID())
		naddr = fmt.Sprintf("%s/p2p/%s", addr, host.ID())
	}
	con := &connLogger{}

	log.Println("Listening on:", naddr)
	host.Network().Notify(con)

	go reciver(naddr)
	time.Sleep(5 * time.Second)
	go reciver(naddr)
	for {

	}
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

func handleStream(s network.Stream) {
	decoder := capnp.NewDecoder(s)
	go startReciver(decoder, s)
	encoder := capnp.NewEncoder(s)
	go startResponder(encoder, s)

}

func startResponder(encoder *capnp.Encoder, s network.Stream) {
	i := 0
	for {
		i++
		msg, seg, _ := capnp.NewMessage(capnp.SingleSegment(nil))

		nping, _ := ping.NewRootBlockInfo(seg)

		nping.SetMessage_([]byte(strconv.Itoa(i)))

		if err := encoder.Encode(msg); err != nil {
			log.Println(err)

		}
	}
}

func startReciver(decoder *capnp.Decoder, s network.Stream) {
	for {
		msg, err := decoder.Decode()
		if err != nil {

			log.Fatal(err)
		}

		block, err := ping.ReadRootBlockInfo(msg)
		if err != nil {
			log.Fatal(err)
		}

		text, _ := block.Message_()

		fmt.Printf("%s Received:%s\n", s.Conn().ID(), text)
	}
}

func sizeInKB(data []byte) float64 {
	return float64(len(data)) / 1024.0
}
