package store

import (
	"strings"
	"time"

	"github.com/petar/GoLLRB/llrb"
)

type pairItem struct {
	key    string
	value  []byte
	flag   int
	expire int64
}

func (pair *pairItem) Less(item llrb.Item) bool {
	switch x := item.(type) {
	case *pairItem:
		return strings.Compare(pair.key, x.key) < 0
	}
	return true
}

type rbTree struct {
	tree *llrb.LLRB
}

func NewRbTreeStore() Store {
	return &rbTree{tree: llrb.New()}
}

func (m *rbTree) Get(key string) ([]byte, error) {
	pair := m.tree.Get(&pairItem{key: key})
	if pair == nil {
		return nil, ErrNotExist
	}
	item, ok := pair.(*pairItem)
	if !ok {
		return nil, ErrUnkown
	}

	now := time.Now().Unix()
	if item.expire < now {
		return nil, ErrExpire
	}
	var v []byte
	if item.flag == COMPRESSED_FLAG {
		v = doZlibUnComporess(item.value)
	} else {
		v = item.value
	}
	return v, nil
}

func (m *rbTree) Set(key string, value []byte, flag, expire int) error {
	if len(value) == 0 || value == nil {
		return ErrEmptyValue
	}

	var v []byte
	if flag == COMPRESSED_FLAG { //zlib compressed value
		v = doZlibCompress(value)
	} else {
		v = value
	}
	var e int64
	if expire > 0 {
		e = time.Now().Unix() + int64(expire)
	} else {
		e = time.Now().Unix() + int64(defaultExpireSeconds)
	}
	m.tree.ReplaceOrInsert(&pairItem{
		key:    key,
		value:  v,
		flag:   flag,
		expire: e,
	})
	return nil
}

func (m *rbTree) Del(key string) error {
	m.tree.ReplaceOrInsert(&pairItem{
		key:    key,
		value:  nil,
		flag:   0,
		expire: 0,
	})
	return nil
}

func (m *rbTree) Exist(key string) bool {
	return m.tree.Has(&pairItem{key: key})
}

func (m *rbTree) Serialize() (map[string][]byte, error) {
	o := make(map[string][]byte)
	now := time.Now().Unix()
	m.tree.AscendGreaterOrEqual(&pairItem{key: ""}, func(item llrb.Item) bool {
		pair, ok := item.(*pairItem)
		if ok {
			if pair.expire >= now {
				var v []byte
				if pair.flag == COMPRESSED_FLAG {
					v = doZlibUnComporess(pair.value)
				} else {
					v = pair.value
				}
				o[pair.key] = v
			}
			return true
		} else {
			return false
		}
	})
	return o, nil
}

func (m *rbTree) Unserialize(c map[string][]byte) (Store, error) {
	store := NewRbTreeStore()
	for k, v := range c {
		store.Set(k, v, COMPRESSED_FLAG, defaultExpireSeconds)
	}
	return store, nil
}
