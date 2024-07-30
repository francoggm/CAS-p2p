package main

import (
	"bytes"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpictures"
	pathKey := CASPathTransformFunc(key)

	expected := "52439/f5e49/ba33b/4b9b2/373dc/d6bfc/9e97c/afb7f"
	original := "52439f5e49ba33b4b9b2373dcd6bfc9e97cafb7f"

	if pathKey.PathName != expected {
		t.Errorf("Expected: %s Received: %s", expected, pathKey.PathName)
	}

	if pathKey.Filename != original {
		t.Errorf("Expected: %s Received: %s", expected, pathKey.Filename)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	key := "bestpics"
	data := []byte("some jpg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	if string(b) != string(data) {
		t.Errorf("Expected: %s Received: %s", data, b)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}
