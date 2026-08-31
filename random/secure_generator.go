package random

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"math/big"
	"sync"
)

var bigIntPool = sync.Pool{
	New: func() any { return new(big.Int) },
}

type secureRandGenerator struct{}

func (secureRandGenerator) Int() int {
	var buf [8]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return int(binary.LittleEndian.Uint64(buf[:]) << 1 >> 1)
}

func (s secureRandGenerator) Intn(max int) int {
	return int(s.Int64n(int64(max)))
}

func (secureRandGenerator) Int31() int32 {
	var buf [4]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return int32(binary.LittleEndian.Uint32(buf[:]) << 1 >> 1)
}

func (s secureRandGenerator) Int31n(max int32) int32 {
	return int32(s.Int64n(int64(max)))
}

func (secureRandGenerator) Int64() int64 {
	var buf [8]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return int64(binary.LittleEndian.Uint64(buf[:]))
}

func (secureRandGenerator) Int64n(max int64) int64 {
	b := bigIntPool.Get().(*big.Int)
	defer bigIntPool.Put(b)
	n, err := rand.Int(rand.Reader, b.SetInt64(max))
	if err != nil {
		panic(err) // 密码学场景下静默失败比panic更危险
	}
	return n.Int64()
}
