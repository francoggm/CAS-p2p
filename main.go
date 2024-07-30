package main

import (
	"cas-p2p/p2p"
	"log"
)

func OnPeer(p p2p.Peer) error {
	log.Printf("New peer: %+v\n", p)
	return nil
}

func main() {
	opts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		Decoder:       p2p.DefaultDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
		OnPeer:        OnPeer,
	}

	tr := p2p.NewTCPTransport(opts)

	go func() {
		for rpc := range tr.Consume() {
			log.Printf("Received RPC: %s\n", rpc)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
