package store

import (
	"testing"
	"time"
)

func TestRBtreeSet(t *testing.T) {
	store := NewRbTreeStore()
	err := store.Set("key", []byte("value"), COMPRESSED_FLAG, 60)
	if err != nil {
		t.Error(err)
	}
}

func TestRBtreeGet(t *testing.T) {
	store := NewRbTreeStore()
	err := store.Set("key", []byte("value"), COMPRESSED_FLAG, 1)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(2 * time.Second)

	value, err := store.Get("key")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(value))
	}
}

func TestRBtreeSerialize(t *testing.T) {
	store := NewRbTreeStore()
	err := store.Set("key", []byte("value"), COMPRESSED_FLAG, 60)
	if err != nil {
		t.Error(err)
	}
	o, err := store.Serialize()
	if err != nil {
		t.Error(err)
	}
	t.Log(o)
}
