package main

import (
	"cas-p2p/p2p"
	"log"
	"time"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO: OnPeer
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fsOpts := FileServerOpts{
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	fs := NewFileServer(fsOpts)

	go func() {
		time.Sleep(time.Second * 3)
		fs.Stop()
	}()

	if err := fs.Start(); err != nil {
		log.Fatal(err)
	}
}
