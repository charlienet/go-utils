package bytex_test

import (
	"encoding/binary"
	"testing"

	"github.com/charlienet/go-utils/bytex"
)

var (
	sinkString string
	sinkBytes  []byte
)

func BenchmarkStringToBytes(b *testing.B) {
	s := "hello world"
	by := []byte(s)

	b.Run("native_[]byte", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkBytes = []byte(s)
		}
	})

	b.Run("native_string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkString = string(by)
		}
	})

	b.Run("StringToBytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkBytes = bytex.StringToBytes(s)
		}
	})

	b.Run("BytesToString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkString = bytex.BytesToString(by)
		}
	})
}

var sinkUint64 uint64

func BenchmarkHex(b *testing.B) {
	data := bytex.Bytes([]byte("hello world, this is a benchmark test"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = data.Hex()
	}
}

func BenchmarkUpperHex(b *testing.B) {
	data := bytex.Bytes([]byte("hello world, this is a benchmark test"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = data.UpperHex()
	}
}

func BenchmarkFromUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkBytes = bytex.FromUint64(0x1234567890ABCDEF, binary.BigEndian)
	}
}

func BenchmarkFromInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkBytes = bytex.FromInt64(-1234567890, binary.BigEndian)
	}
}

func BenchmarkBase64(b *testing.B) {
	data := bytex.Bytes([]byte("hello world, this is a benchmark test for base64 encoding"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = data.Base64()
	}
}

func BenchmarkString(b *testing.B) {
	data := bytex.Bytes([]byte("hello world, this is a benchmark test"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = data.String()
	}
}

func BenchmarkClone(b *testing.B) {
	data := bytex.Bytes([]byte("hello world, this is a benchmark test"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBytes = data.Clone()
	}
}

func BenchmarkFromString(b *testing.B) {
	s := "hello world, this is a benchmark test"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBytes = bytex.FromString(s)
	}
}

func BenchmarkBytesToUint64(b *testing.B) {
	data := []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkUint64, _ = bytex.BigEndian.BytesToUint64(data)
	}
}
