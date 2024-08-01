package main

import (
	"cas-p2p/p2p"
	"fmt"
	"log"
)

type FileServerOpts struct {
	StorageRoot string
	PathTransformFunc

	Transport p2p.Transport
}

type FileServer struct {
	FileServerOpts
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
	}
}

func (fs *FileServer) Stop() {
	close(fs.quitChan)
}

func (fs *FileServer) loop() {
	defer func(){
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

func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}

	fs.loop()

	return nil
}
