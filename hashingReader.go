package bencoding

import (
	"bufio"
	"crypto/sha1"
	"hash"
)

type hashingRreader struct {
	b          *bufio.Reader
	hash       hash.Hash
	shouldHash bool
}

func (hr *hashingRreader) Peek(n int) ([]byte, error) {
	b, e := hr.b.Peek(n)
	return b, e
}

func (hr *hashingRreader) ReadByte() (byte, error) {
	b, e := hr.b.ReadByte()
	if hr.shouldHash && e == nil {
		hr.hash.Write([]byte{b})
	}
	return b, e
}

func (hr *hashingRreader) ReadBytes(delim byte) ([]byte, error) {
	b, e := hr.b.ReadBytes(delim)
	if hr.shouldHash && e == nil {
		hr.hash.Write(b)
	}
	return b, e
}

func (hr *hashingRreader) StartHasing() {
	if !hr.shouldHash {
		hr.hash = sha1.New()
		hr.shouldHash = true
	}
}

func (hr *hashingRreader) StopHashing() {
	hr.shouldHash = false
}
