package store

import (
	"strconv"
	"testing"
)

func Benchmark_rbtree(b *testing.B) {
	store := NewRbTreeStore()
	storeBenchmark(store, b)
}

func Benchmark_mem(b *testing.B) {
	store := NewStoreMem()
	storeBenchmark(store, b)
}

func Benchmark_leveldb(b *testing.B) {
	store, err := NewStoreLeveldb("/home/vagrant/db111.db")
	if err != nil {
		b.Error(err)
	} else {
		storeBenchmark(store, b)
	}
}

func storeBenchmark(s Store, b *testing.B) {
	for i := 0; i < b.N; i++ {
		key := "key_" + strconv.Itoa(i)
		value := "value_" + strconv.Itoa(i)
		s.Set(key, []byte(value), DEFAULT_FLAG, 60)
	}
}
