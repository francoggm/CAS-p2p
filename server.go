package main

import (
	"cas-p2p/p2p"
	"fmt"
	"log"
	"net"
	"sync"
)

type FileServerOpts struct {
	StorageRoot string
	PathTransformFunc

	Transport      p2p.Transport
	BootstrapNodes []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[net.Addr]p2p.Peer

	store    *Store
	quitChan chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		StoreRoot:         opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitChan:       make(chan struct{}),
		peers:          make(map[net.Addr]p2p.Peer),
	}
}

func (fs *FileServer) Stop() {
	close(fs.quitChan)
}

func (fs *FileServer) OnPeer(p p2p.Peer) error{
	fs.peerLock.Lock()
	defer fs.peerLock.Unlock()

	fs.peers[p.RemoteAddr()] = p
	log.Printf("Peer connected: %s", p.RemoteAddr())
	return nil
}

func (fs *FileServer) loop() {
	defer func() {
		log.Println("File server stopped")
		fs.Transport.Close()
	}()

	for {
		select {
		case msg := <-fs.Transport.Consume():
			fmt.Println(msg)
		case <-fs.quitChan:
			return
		}
	}
}

func (fs *FileServer) bootstrapNetwork() error {
	for _, addr := range fs.BootstrapNodes {
		if addr == "" {
			continue
		}

		go func(addr string) {
			if err := fs.Transport.Dial(addr); err != nil {
				log.Println("Dial error: ", err)
			}
		}(addr)
	}

	return nil
}

func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}

	fs.bootstrapNetwork()

	fs.loop()

	return nil
}
