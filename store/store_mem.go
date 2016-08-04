package store

import (
	"fmt"
)

type StoreMem struct {
	cache map[string][]byte
}

func NewStoreMem() *StoreMem {
	return &StoreMem{
		cache: make(map[string][]byte),
	}
}

func (s *StoreMem) Set(key string, value []byte) error {
	s.cache[key] = value
	return nil
}

func (s *StoreMem) Get(key string) ([]byte, error) {
	val, ok := s.cache[key]
	if !ok {
		return nil, fmt.Errorf("not found!")
	}
	return val, nil
}

func (s *StoreMem) Del(key string) error {
	delete(s.cache, key)
	return nil
}

func (s *StoreMem) Serialize() (map[string][]byte, error) {
	o := make(map[string][]byte)
	for k, v := range s.cache {
		o[k] = v
	}
	return o, nil
}

func (s *StoreMem) Unserialize(c map[string][]byte) (Store, error) {
	o := NewStoreMem()
	o.cache = c
	return o, nil
}
