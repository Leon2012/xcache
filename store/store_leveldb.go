package store

import (
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

type StoreLeveldb struct {
	db   *leveldb.DB
	path string
}

func NewStoreLeveldb(path string) (Store, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &StoreLeveldb{
		db:   db,
		path: path,
	}, nil
}

func (s *StoreLeveldb) Close() {
	s.db.Close()
}

func (s *StoreLeveldb) Set(key string, value []byte, flag, expire int) error {
	err := s.db.Put([]byte(key), value, nil)
	return err
}

func (s *StoreLeveldb) Get(key string) ([]byte, error) {
	data, err := s.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *StoreLeveldb) Del(key string) error {
	err := s.db.Delete([]byte(key), nil)
	return err
}

func (s *StoreLeveldb) Exist(key string) bool {
	b, err := s.db.Has([]byte(key), nil)
	if err != nil {
		return false
	} else {
		return b
	}
}

func (s *StoreLeveldb) Serialize() (map[string][]byte, error) {
	o := make(map[string][]byte)
	iter := s.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		o[string(key)] = value
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, err
	} else {
		return o, nil
	}
}

func (s *StoreLeveldb) Unserialize(c map[string][]byte) (Store, error) {
	s.Close()
	s.db = nil
	os.Remove(s.path)
	o, err := NewStoreLeveldb(s.path)
	if err != nil {
		return nil, err
	}
	return o, nil
}
