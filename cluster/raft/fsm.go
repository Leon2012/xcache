package raft

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	logger "github.com/Leon2012/xcache/log"
	"github.com/Leon2012/xcache/store"
	"github.com/hashicorp/raft"
)

type fsm struct {
	store store.Store
	mu    sync.Mutex
}

func NewFSM(s store.Store) *fsm {
	return &fsm{
		store: s,
	}
}

func (s *fsm) Apply(l *raft.Log) interface{} {
	var c command
	var err error
	if err = json.Unmarshal(l.Data, &c); err != nil {
		logger.Error("failed to unmarshal command: %s", err.Error())
		return err
	}
	err = nil
	switch c.Op {
	case "set":
		err = s.handleSet(c.Key, c.Value, c.Flag, c.Expire)
		break
	case "del":
		err = s.handleDel(c.Key)
		break
	default:
		err = fmt.Errorf("unrecognized command op: %s", c.Op)
	}
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

//生成快照
func (s *fsm) Snapshot() (raft.FSMSnapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, err := s.store.Serialize()
	if err != nil {
		return nil, err
	}
	return NewFSMSnapshot(m), nil
}

//重新读取
func (s *fsm) Restore(old io.ReadCloser) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	o := make(map[string][]byte)
	if err := json.NewDecoder(old).Decode(&o); err != nil {
		return err
	}

	t, err := s.store.Unserialize(o)
	if err != nil {
		return err
	}
	s.store = nil
	s.store = t
	return nil
}

func (s *fsm) Get(key string) ([]byte, error) {
	//s.mu.Lock()
	//defer s.mu.Unlock()
	return s.store.Get(key)
}

func (s *fsm) Has(key string) bool {
	return s.store.Exist(key)
}

func (s *fsm) handleSet(key string, value []byte, flag, expire int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store.Set(key, value, flag, expire)
}

func (s *fsm) handleDel(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store.Del(key)
}
