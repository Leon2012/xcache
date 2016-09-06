package store

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
)

const defaultExpireSeconds = 12 * 30 * 24 * 60 * 60 //1 year

const (
	DEFAULT_FLAG    = 0
	COMPRESSED_FLAG = 2
)

var (
	ErrNotExist   = errors.New("not exist")
	ErrEmptyValue = errors.New("value is empty")
	ErrUnkown     = errors.New("unkown error")
	ErrExpire     = errors.New("expire")
	ErrKeyExist   = errors.New("key exist")
)

func doZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func doZlibUnComporess(src []byte) []byte {
	b := bytes.NewReader(src)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}
