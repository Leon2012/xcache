package raft

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	logger "github.com/Leon2012/xcache/log"
	"github.com/Leon2012/xcache/store"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

type command struct {
	Op     string `json:"op,omitempty"`
	Key    string `json:"key,omitempty"`
	Value  []byte `json:"value,omitempty"`
	Flag   int    `json:"flag,omitempty"`
	Expire int    `json:"expire,omitempty"`
}

const (
	DEFAULT_RAFT_TIMEOUT = 10 * time.Second
	DEFAULT_TCP_TIMEOUT  = 10 * time.Second
	DEFAULT_MAX_POOL     = 3
)

func init() {
	logger.SetModule("RAFT")
}

type RaftImpl struct {
	RaftDir  string
	RaftBind string
	//mu       sync.Mutex
	raft    *raft.Raft
	fsmImpl *fsm
}

func NewRaft(dir, bind string) *RaftImpl {
	return &RaftImpl{
		RaftDir:  dir,
		RaftBind: bind,
	}
}

func (r *RaftImpl) Init(enableSingle bool, s store.Store) error {

	if s == nil {
		return fmt.Errorf("Please set store before open")
	}

	//设置raft config
	config := raft.DefaultConfig()
	config.LogOutput = os.Stdout

	//检查是否存在peers，并读取
	peers, err := readPeersJSON(filepath.Join(r.RaftDir, "peers.json"))
	if err != nil {
		return err
	}

	//设置是否单节点
	if enableSingle && len(peers) <= 1 {
		logger.Info("enabling signle-node mode")
		config.EnableSingleNode = true
		config.DisableBootstrapAfterElect = false
	}

	//设置raft实例
	addr, err := net.ResolveTCPAddr("tcp", r.RaftBind)
	if err != nil {
		return err
	}

	trans, err := raft.NewTCPTransport(r.RaftBind, addr, DEFAULT_MAX_POOL, DEFAULT_TCP_TIMEOUT, os.Stderr)
	if err != nil {
		return err
	}

	//创建peer storage
	peerStore := raft.NewJSONPeers(r.RaftDir, trans)

	//创建snapshot Storage
	snapshotStore, err := raft.NewFileSnapshotStore(r.RaftDir, 2, os.Stderr)
	if err != nil {
		return err
	}

	//创建raft.Log storage
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(r.RaftDir, "raft.db"))
	if err != nil {
		logStore.Close()
		return err
	}

	//创建log cache
	cacheStore, err := raft.NewLogCache(512, logStore)
	if err != nil {
		logStore.Close()
		return err
	}

	//new raft
	fsm := NewFSM(s)
	r.fsmImpl = fsm

	ra, err := raft.NewRaft(config, r.fsmImpl, cacheStore, logStore, snapshotStore, peerStore, trans)
	if err != nil {
		logStore.Close()
		trans.Close()
		return err
	}

	r.raft = ra

	go r.monitorLeadership()

	return nil
}

func (r *RaftImpl) Join(addr string) error {
	// r.mu.Lock()
	// defer r.mu.Unlock()

	logger.Info("received join request for remote node as %s", addr)
	f := r.raft.AddPeer(addr)
	if f.Error() != nil {
		return f.Error()
	}
	logger.Info("node at %s joined successfully", addr)
	return nil
}

func (r *RaftImpl) Has(key string) bool {
	// r.mu.Lock()
	// defer r.mu.Unlock()

	return r.fsmImpl.Has(key)
}

func (r *RaftImpl) Incr(key string, offset int64) (int64, error) {
	if !r.Has(key) {
		return 0, fmt.Errorf("key not exist")
	}
	data, err := r.Get(key)
	if err != nil {
		return 0, err
	}
	logger.Info("value1 : %s", string(data))
	value, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0, err
	}
	//value := bytesToInt64(data)
	logger.Info("value : %d", value)
	value = value + offset
	str := strconv.Itoa(int(value))
	logger.Info("value2 : %s", str)
	//data = int64ToBytes(value)
	//data = []byte(str)
	err = r.Set(key, []byte(str), 0, 3600)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (r *RaftImpl) Decr(key string, offset int64) (int64, error) {
	if !r.Has(key) {
		return 0, fmt.Errorf("key not exist")
	}
	data, err := r.Get(key)
	if err != nil {
		return 0, err
	}
	logger.Info("value1 : %s", string(data))
	value, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0, err
	}
	//value := bytesToInt64(data)
	logger.Info("value : %d", value)
	value = value - offset
	str := strconv.Itoa(int(value))
	//data = int64ToBytes(value)
	data = []byte(str)
	err = r.Set(key, data, 0, 3600)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (r *RaftImpl) Get(key string) ([]byte, error) {
	// r.mu.Lock()
	// defer r.mu.Unlock()
	return r.fsmImpl.Get(key)
}

func (r *RaftImpl) Add(key string, value []byte, flag, expire int) error {
	// r.mu.Lock()
	// defer r.mu.Unlock()

	if r.Has(key) {
		return fmt.Errorf("key exist")
	}
	return r.Set(key, value, flag, expire)
}

func (r *RaftImpl) Replace(key string, value []byte, flag, expire int) error {
	// r.mu.Lock()
	// defer r.mu.Unlock()
	if !r.Has(key) {
		return fmt.Errorf("key not exist")
	}
	return r.Set(key, value, flag, expire)
}

func (r *RaftImpl) Set(key string, value []byte, flag, expire int) error {
	// r.mu.Lock()
	// defer r.mu.Unlock()

	//return r.fsmImpl.handleSet(key, value, flag, expire)
	if r.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}
	c := &command{
		Op:     "set",
		Key:    key,
		Value:  value,
		Flag:   flag,
		Expire: expire,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f := r.raft.Apply(b, DEFAULT_RAFT_TIMEOUT)
	if err, ok := f.(error); ok { //判断是否返回error
		return err
	}
	logger.Info("fsm set success!!!!")
	return nil
}

func (r *RaftImpl) Del(key string) error {
	// r.mu.Lock()
	// defer r.mu.Unlock()
	if r.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}
	c := &command{
		Op:  "del",
		Key: key,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f := r.raft.Apply(b, DEFAULT_RAFT_TIMEOUT)
	if err, ok := f.(error); ok { //判断是否返回error
		return err
	}
	return nil
}

func (r *RaftImpl) Name() string {
	return r.RaftBind
}

func (r *RaftImpl) IsLeader() bool {
	return r.raft.State() == raft.Leader
}

func (r *RaftImpl) monitorLeadership() {
	leaderCh := r.raft.LeaderCh()
	for {
		select {
		case isLeader := <-leaderCh:
			if isLeader {
				logger.Warning("============================= cluster %s leadership acquired ===============================", r.Name())
			}
		}
	}
}

//读取peers.json文件
func readPeersJSON(path string) ([]string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	//文件为空
	if len(b) == 0 {
		return nil, nil
	}

	var peers []string
	decoder := json.NewDecoder(bytes.NewReader(b))
	if err := decoder.Decode(&peers); err != nil {
		return nil, err
	}
	return peers, nil
}

func int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.PutUvarint(buf, uint64(i))
	return buf
}

func bytesToInt64(b []byte) int64 {
	buf := bytes.NewBuffer(b)
	i, err := binary.ReadUvarint(buf)
	if err != nil {
		return 0
	} else {
		return int64(i)
	}
}
