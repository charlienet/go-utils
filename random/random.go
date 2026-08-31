package random

import (
	"crypto/rand"
	"encoding/binary"
	"io"

	"github.com/charlienet/go-utils/bytex"
	mrndv2 "math/rand/v2"
)

var (
	SecureGenerator = &secureRandGenerator{}
	NormalGenerator = NewRandGenerator()
)

type scopeConstraint interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Int[T scopeConstraint]() T {
	return T(NormalGenerator.Int())
}

func Intn[T scopeConstraint](max T) T {
	n := NormalGenerator.Int64n(int64(max))
	return T(n)
}

func IntRange[T scopeConstraint](min, max T) T {
	if min >= max {
		panic("random: invalid range, min must be less than max")
	}
	n := Intn(max - min)
	return T(n + min)
}

// RandBytes 生成密码学安全的随机字节
func RandBytes(length int) ([]byte, error) {
	r := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, r)
	return r, err
}

// FastRandBytes 生成快速伪随机字节（非密码学安全）
func FastRandBytes(length int) []byte {
	r := make([]byte, length)
	for i := 0; i < length; i += 8 {
		v := mrndv2.Uint64()
		end := i + 8
		if end > length {
			for j := i; j < length; j++ {
				r[j] = byte(v)
				v >>= 8
			}
		} else {
			binary.LittleEndian.PutUint64(r[i:], v)
		}
	}
	return r
}

// HexString 生成指定长度的快速十六进制字符串（小写）
// length 必须是偶数，因为每 2 个十六进制字符对应 1 个字节
func HexString(length int) string {
	if length%2 != 0 {
		panic("random: hex string length must be even")
	}
	bytes := FastRandBytes(length / 2)
	return bytex.Bytes(bytes).Hex()
}
