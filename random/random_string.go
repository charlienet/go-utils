package random

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/charlienet/go-utils/bytex"
	"io"
	mrndv2 "math/rand/v2"
	"strings"
)

const (
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	digit     = "0123456789"
	nomix     = "BCDFGHJKMPQRTVWXY2346789"
	letter    = uppercase + lowercase
	allChars  = uppercase + lowercase + digit
	hex       = digit + "ABCDEF"
)

type charScope struct {
	bytes  []byte
	length int
	bits   int
	mask   int
}

var (
	Uppercase = StringScope(uppercase) // 大写字母
	Lowercase = StringScope(lowercase) // 小写字母
	Digit     = StringScope(digit)     // 数字
	Nomix     = StringScope(nomix)     // 不混淆字符
	Letter    = StringScope(letter)    // 字母
	Hex       = StringScope(hex)       // 十六进制字符
	AllChars  = StringScope(allChars)  // 所有字符
)

func StringScope(str string) *charScope {
	len := len(str)

	scope := &charScope{
		bytes:  bytex.StringToBytes(str),
		length: len,
		bits:   1,
	}

	for scope.mask < len {
		scope.bits++
		scope.mask = 1<<scope.bits - 1
	}

	return scope
}

func (scope *charScope) allocRet(length int, prefix ...string) []byte {
	if len(prefix) > 0 {
		pre := strings.Join(prefix, "")
		ret := make([]byte, len(pre), length+len(pre))
		copy(ret, pre)
		return ret
	}
	return make([]byte, 0, length)
}

// Generate 使用快速伪随机生成器生成指定长度的随机字符串
func (scope *charScope) Generate(length int, prefix ...string) string {
	ret := scope.allocRet(length, prefix...)

	n := length
	var cache uint64
	bitsAvailable := 0

	for i := n - 1; i >= 0; {
		if bitsAvailable < scope.bits {
			cache = mrndv2.Uint64()
			bitsAvailable = 64
		}
		idx := int(cache & uint64(scope.mask))
		cache >>= uint(scope.bits)
		bitsAvailable -= scope.bits
		if idx < scope.length {
			ret = append(ret, scope.bytes[idx])
			i--
		}
	}
	return bytex.BytesToString(ret)
}

// GenerateSecure 使用密码学安全随机生成器生成指定长度的随机字符串
func (scope *charScope) GenerateSecure(length int, prefix ...string) string {
	ret := scope.allocRet(length, prefix...)

	n := length
	var cache uint64
	bitsAvailable := 0
	var buf [8]byte

	for i := n - 1; i >= 0; {
		if bitsAvailable < scope.bits {
			_, err := io.ReadFull(rand.Reader, buf[:])
			if err != nil {
				panic(err)
			}
			cache = binary.LittleEndian.Uint64(buf[:])
			bitsAvailable = 64
		}
		idx := int(cache & uint64(scope.mask))
		cache >>= uint(scope.bits)
		bitsAvailable -= scope.bits
		if idx < scope.length {
			ret = append(ret, scope.bytes[idx])
			i--
		}
	}
	return bytex.BytesToString(ret)
}
