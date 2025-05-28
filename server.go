package main

import (
	_ "bufio"
	"context"
	"fmt"
	"log"
	"time"
	ping "zelonis/capn"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"

	capnp "capnproto.org/go/capnp/v3"
)

func reciver(addr string) {
	time.Sleep(5 * time.Second)
	peerAddr := addr

	ctx := context.Background()
	host, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}

	maddr, err := ma.NewMultiaddr(peerAddr)
	if err != nil {
		log.Fatal(err)
	}
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatal(err)
	}
	if err := host.Connect(ctx, *info); err != nil {
		log.Fatal(err)
	}

	stream, err := host.NewStream(ctx, info.ID, "/zel/1.0.0")
	if err != nil {
		log.Fatal(err)
	}

	encoder := capnp.NewEncoder(stream)
	decoder := capnp.NewDecoder(stream)

	for {

		msg, seg, _ := capnp.NewMessage(capnp.SingleSegment(nil))
		nping, _ := ping.NewRootBlockInfo(seg)

		nping.SetMessage_([]byte(fmt.Sprintf("Sent Info To User")))

		if err := encoder.Encode(msg); err != nil {
			log.Println(err)
			break
		}

		msg, err := decoder.Decode()
		if err != nil {

			log.Fatal(err)
		}

		block, err := ping.ReadRootBlockInfo(msg)
		if err != nil {
			log.Fatal(err)
		}
		text, _ := block.Message_()

		fmt.Printf("Received:%s\n", text)
	}

}
