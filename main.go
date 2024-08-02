package main

import (
	"cas-p2p/p2p"
	"log"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fsOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	fs := NewFileServer(fsOpts) 
	tcpTransport.OnPeer = fs.OnPeer

	return fs
}

func main() {
	// Mocking two file servers
	
	fs1 := makeServer(":3000")
	go func() {
		log.Fatal(fs1.Start())
	}()

	fs2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(fs2.Start())
	}()

	select{}
}
