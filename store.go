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
	blockSize             = 5
	defaultRootFolderName = "store"
)

type PathKey struct {
	PathName    string
	Filename    string
	RootPathKey string
}

func (pk PathKey) FullPath() string {
	return path.Join(pk.PathName, pk.Filename)
}

type PathTransformFunc func(string) PathKey

func DefaultPathTransformFunc(key string) PathKey {
	return PathKey{
		PathName:    key,
		Filename:    key,
		RootPathKey: key,
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
		PathName:    strings.Join(paths, "/"),
		Filename:    hashStr,
		RootPathKey: paths[0],
	}
}

type StoreOpts struct {
	StoreRoot string
	PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if opts.StoreRoot == "" {
		opts.StoreRoot = defaultRootFolderName
	}

	return &Store{opts}
}

func (s *Store) Has(key string) bool {
	pk := s.PathTransformFunc(key)

	_, err := os.Stat(path.Join(s.StoreRoot, pk.FullPath()))
	return err == nil
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.StoreRoot)
}

func (s *Store) Delete(key string) error {
	pk := s.PathTransformFunc(key)
	return os.RemoveAll(path.Join(s.StoreRoot, pk.RootPathKey))
}

func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)
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
	pk := s.PathTransformFunc(key)
	return os.Open(path.Join(s.StoreRoot, pk.FullPath()))
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pk := s.PathTransformFunc(key)
	pathNameWithRoot := path.Join(s.StoreRoot, pk.PathName)

	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	fullPath := path.Join(s.StoreRoot, pk.FullPath())

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Printf("Written %d bytes: %s", n, fullPath)

	return nil
}
