package memcache

import (
	"bufio"
	"bytes"
	"testing"
)

func TestRead(t *testing.T) {
	//b := []byte{115, 101, 116, 32, 107, 101, 121, 32, 48, 32, 51, 48, 32, 56, 118, 97, 108, 117, 101, 32}
	b := []byte{115, 101, 116, 32, 107, 101, 121, 32, 49, 32, 51, 48, 32, 52, 50, 97, 58, 51, 58, 123, 105, 58, 48, 59, 115, 58, 49, 58, 34, 97, 34, 59, 105, 58, 49, 59, 115, 58, 49, 58, 34, 98, 34, 59, 105, 58, 50, 59, 115, 58, 49, 58, 34, 99, 34, 59, 125}
	r := bufio.NewReader(bytes.NewReader(b))

	req, err := Read(r)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(req)
	}
}
