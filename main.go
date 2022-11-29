package main

import (
	"encoding/hex"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func FromHex(s string) []byte {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte { return h[:] }

func main() {
    hash := HexToHash("e47c27b80db796fd333ab4fe523417984731766e52f1db44ebcd565d0a40a025")

    db, err := leveldb.OpenFile("../state-root/chaindata", nil)
    if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
        db, err = leveldb.RecoverFile("../state-root/chaindata", nil)
    }
    if err != nil {
        fmt.Println("Error open!", err)
    }

    data, err := db.Get(hash.Bytes(), nil)
    if err != nil {
        fmt.Println("Error get!", err)
    }

    fmt.Println("Hello, World!", data)
}
