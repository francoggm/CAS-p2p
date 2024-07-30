package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

const (
	blockSize = 5
)

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	Filename string
	RootPath string
}

func (pk PathKey) FullPath() string {
	return path.Join(pk.PathName, pk.Filename)
}

func DefaultPathTransformFunc(key string) PathKey {
	return PathKey{
		PathName: key,
		Filename: key,
		RootPath: key,
	}
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		Filename: hashStr,
		RootPath: paths[0],
	}
}

type StoreOpts struct {
	PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{opts}
}

func (s *Store) Has(key string) bool {
	pk := s.PathTransformFunc(key)

	_, err := os.Stat(pk.FullPath())
	return err != nil
}

func (s *Store) Delete(key string) error {
	pk := s.PathTransformFunc(key)
	return os.RemoveAll(pk.RootPath)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, f); err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	return os.Open(pathKey.FullPath())
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pk := s.PathTransformFunc(key)
	if err := os.MkdirAll(pk.PathName, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(pk.FullPath())
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Printf("Written %d bytes: %s", n, pk.FullPath())

	return nil
}
