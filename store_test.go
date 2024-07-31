package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

const (
	testCount = 50
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpictures"
	pathKey := CASPathTransformFunc(key)

	path := "52439/f5e49/ba33b/4b9b2/373dc/d6bfc/9e97c/afb7f"
	fileName := "52439f5e49ba33b4b9b2373dcd6bfc9e97cafb7f"

	if pathKey.PathName != path {
		t.Errorf("Expected: %s Received: %s", path, pathKey.PathName)
	}

	if pathKey.Filename != fileName {
		t.Errorf("Expected: %s Received: %s", fileName, pathKey.Filename)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t, s)

	for i := 0; i < testCount; i++ {
		key := fmt.Sprintf("foo_%d", i)
		data := []byte("random bytes") 
		
		// Write file
		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}
	
		// Verify if file exists
		if ok := s.Has(key); !ok {
			t.Errorf("Expected: true Received: %t", ok)
		}
	
		// Read file
		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}
	
		// Verify if file is the same
		b, err := io.ReadAll(r)
		if err != nil {
			t.Error(err)
		}
	
		if string(b) != string(data) {
			t.Errorf("Expected: %s Received: %s", data, b)
		}
	
		// Delete file
		if err := s.Delete(key); err != nil {
			t.Error(err)
		}
	}
}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store){
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}