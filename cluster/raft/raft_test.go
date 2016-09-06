package raft

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func Test_bytesToInt32(t *testing.T) {
	//a := []byte{0x01}
	// i := bytesToInt32(a)
	// t.Logf("i: %d", i)
	//b := []byte{123, 0, 0, 0, 0, 0, 0, 0}
	b := []byte("100")
	buf := bytes.NewBuffer(b)
	i, err := binary.ReadVarint(buf)
	if err != nil {
		t.Log(0)
	} else {
		t.Logf("%d", int64(i))
	}

}

func Test_int32ToBytes(t *testing.T) {
	i := 123
	var buf = make([]byte, 8)
	binary.PutUvarint(buf, uint64(i))
	t.Log(buf)
}
