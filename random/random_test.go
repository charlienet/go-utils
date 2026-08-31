package random

import (
	stdhex "encoding/hex"
	"sync"
	"testing"

	"github.com/charlienet/go-utils/bytex"
	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	// 生成多个随机数，确保不全部相同
	results := make(map[int]bool)
	for range 100 {
		results[Int[int]()] = true
	}
	assert.Greater(t, len(results), 1, "should generate different random numbers")
}

func TestIntn(t *testing.T) {
	max := 100
	for range 1000 {
		n := Intn(max)
		assert.GreaterOrEqual(t, n, 0)
		assert.Less(t, n, max)
	}
}

func TestIntRange(t *testing.T) {
	min, max := 10, 20
	for range 1000 {
		n := IntRange(min, max)
		assert.GreaterOrEqual(t, n, min)
		assert.Less(t, n, max)
	}
}

func TestRandBytes(t *testing.T) {
	b, err := RandBytes(32)
	assert.NoError(t, err)
	assert.Len(t, b, 32)

	// 确保不是全零
	allZero := true
	for _, v := range b {
		if v != 0 {
			allZero = false
			break
		}
	}
	assert.False(t, allZero, "random bytes should not be all zeros")
}

func TestFastRandBytes(t *testing.T) {
	// 测试不同长度
	for _, length := range []int{0, 1, 7, 8, 9, 15, 16, 17, 32, 64, 100} {
		b := FastRandBytes(length)
		assert.Len(t, b, length, "length %d", length)
	}

	// 确保不是全零（32 字节）
	b := FastRandBytes(32)
	allZero := true
	for _, v := range b {
		if v != 0 {
			allZero = false
			break
		}
	}
	assert.False(t, allZero, "fast random bytes should not be all zeros")
}

func TestHexString(t *testing.T) {
	// 测试不同长度
	for _, length := range []int{0, 2, 4, 8, 16, 32, 64} {
		s := HexString(length)
		assert.Len(t, s, length, "length %d", length)
		// 验证是有效的十六进制字符
		for _, c := range s {
			assert.True(t, (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f'),
				"should only contain lowercase hex characters, got %c", c)
		}
	}

	// 测试奇数长度 panic
	assert.Panics(t, func() { HexString(3) }, "odd length should panic")
	assert.Panics(t, func() { HexString(7) }, "odd length should panic")
}

func TestSecureGenerator(t *testing.T) {
	// 测试 SecureGenerator
	for range 100 {
		n := SecureGenerator.Intn(100)
		assert.GreaterOrEqual(t, n, 0)
		assert.Less(t, n, 100)
	}
}

func TestNormalGenerator(t *testing.T) {
	// 测试 NormalGenerator
	for range 100 {
		n := NormalGenerator.Intn(100)
		assert.GreaterOrEqual(t, n, 0)
		assert.Less(t, n, 100)
	}
}

func TestStringScope_Generate(t *testing.T) {
	// 测试生成随机字符串
	s := Uppercase.Generate(10)
	assert.Len(t, s, 10)
	for _, c := range s {
		assert.True(t, c >= 'A' && c <= 'Z', "should only contain uppercase letters")
	}

	s = Digit.Generate(8)
	assert.Len(t, s, 8)
	for _, c := range s {
		assert.True(t, c >= '0' && c <= '9', "should only contain digits")
	}

	s = Letter.Generate(20)
	assert.Len(t, s, 20)

	// 测试带前缀
	s = Digit.Generate(10, "ID_")
	assert.Len(t, s, 13) // 3 (prefix) + 10
	assert.Equal(t, "ID_", s[:3])
}

func TestStringScope_AllowConsecutiveDuplicates(t *testing.T) {
	// 移除去重逻辑后，允许连续重复字符（这是随机性的自然表现）
	s := AllChars.Generate(100)
	assert.Len(t, s, 100)

	// 验证生成的字符串只包含有效字符
	for _, c := range s {
		assert.True(t,
			(c >= 'A' && c <= 'Z') ||
				(c >= 'a' && c <= 'z') ||
				(c >= '0' && c <= '9'),
			"should only contain alphanumeric characters")
	}
}

func TestEdgeCases(t *testing.T) {
	// 测试 Intn 边界情况
	t.Run("Intn panics with zero or negative values", func(t *testing.T) {
		assert.Panics(t, func() { Intn(0) }, "Intn(0) should panic")
		assert.Panics(t, func() { Intn(-1) }, "Intn(-1) should panic")
		assert.Panics(t, func() { Intn(-100) }, "Intn(-100) should panic")
	})

	// 测试 IntRange 边界情况
	t.Run("IntRange panics with invalid ranges", func(t *testing.T) {
		assert.Panics(t, func() { IntRange(10, 10) }, "IntRange(10, 10) should panic (equal values)")
		assert.Panics(t, func() { IntRange(20, 10) }, "IntRange(20, 10) should panic (min > max)")
	})

	// 测试 Generate(0) 返回空字符串
	t.Run("Generate zero length returns empty string", func(t *testing.T) {
		result := AllChars.Generate(0)
		assert.Equal(t, "", result, "Generate(0) should return empty string")

		result = AllChars.Generate(0, "prefix")
		assert.Equal(t, "prefix", result, "Generate(0) with prefix should return just the prefix")
	})
}

func TestConcurrentSafety(t *testing.T) {
	// 测试并发安全性
	const numGoroutines = 10
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(goroutineID int) {
			defer wg.Done()

			for range iterations {
				// 测试各种随机函数的并发调用
				_ = Int[int]()
				_ = Intn(1000)
				_ = IntRange(0, 100)

				// 生成随机字符串
				_ = AllChars.Generate(10)
				_ = Digit.Generate(5)

				// 使用不同的生成器
				_ = SecureGenerator.Intn(100)
				_ = NormalGenerator.Intn(100)
			}
		}(i)
	}

	wg.Wait()
	// 如果没有panic或死锁，测试就成功了
}

func TestLargeNumbers(t *testing.T) {
	// 测试大数值的 Intn
	t.Run("Intn with large numbers", func(t *testing.T) {
		// 测试接近 MaxInt64 的值
		maxInt32 := int64(^uint32(0) >> 1) // 2^31 - 1
		result := Intn(int(maxInt32))
		assert.GreaterOrEqual(t, result, 0)
		assert.Less(t, int64(result), maxInt32)

		// 测试一个大的值
		bigNum := int(1e9) // 1 billion
		result = Intn(bigNum)
		assert.GreaterOrEqual(t, result, 0)
		assert.Less(t, result, bigNum)
	})
}

// ============ Benchmark 基线测试 ============

// 泛型函数基准测试
func BenchmarkInt(b *testing.B) {
	for b.Loop() {
		_ = Int[int]()
	}
}

func BenchmarkIntn(b *testing.B) {
	for b.Loop() {
		_ = Intn(100)
	}
}

func BenchmarkIntRange(b *testing.B) {
	for b.Loop() {
		_ = IntRange(10, 100)
	}
}

// 生成器基准测试
func BenchmarkSecureRandIntn(b *testing.B) {
	for b.Loop() {
		_ = SecureGenerator.Intn(100)
	}
}

func BenchmarkNormalRandIntn(b *testing.B) {
	for b.Loop() {
		_ = NormalGenerator.Intn(100)
	}
}

// RandBytes 基准测试
func BenchmarkRandBytes(b *testing.B) {
	for b.Loop() {
		_, _ = RandBytes(32)
	}
}

func BenchmarkFastRandBytes(b *testing.B) {
	for b.Loop() {
		_ = FastRandBytes(32)
	}
}

// 不同长度对比
func BenchmarkRandBytes_Len8(b *testing.B) {
	for b.Loop() {
		_, _ = RandBytes(8)
	}
}

func BenchmarkFastRandBytes_Len8(b *testing.B) {
	for b.Loop() {
		_ = FastRandBytes(8)
	}
}

func BenchmarkRandBytes_Len64(b *testing.B) {
	for b.Loop() {
		_, _ = RandBytes(64)
	}
}

func BenchmarkFastRandBytes_Len64(b *testing.B) {
	for b.Loop() {
		_ = FastRandBytes(64)
	}
}

func BenchmarkRandBytes_Len256(b *testing.B) {
	for b.Loop() {
		_, _ = RandBytes(256)
	}
}

func BenchmarkFastRandBytes_Len256(b *testing.B) {
	for b.Loop() {
		_ = FastRandBytes(256)
	}
}

func BenchmarkRandBytes_Len1024(b *testing.B) {
	for b.Loop() {
		_, _ = RandBytes(1024)
	}
}

func BenchmarkFastRandBytes_Len1024(b *testing.B) {
	for b.Loop() {
		_ = FastRandBytes(1024)
	}
}

// 不同字符集基准测试（长度固定 32）
func BenchmarkStringGenerate_Digit(b *testing.B) {
	for b.Loop() {
		_ = Digit.Generate(32)
	}
}

func BenchmarkStringGenerate_Hex(b *testing.B) {
	for b.Loop() {
		_ = Hex.Generate(32)
	}
}

func BenchmarkStringGenerate_Uppercase(b *testing.B) {
	for b.Loop() {
		_ = Uppercase.Generate(32)
	}
}

func BenchmarkStringGenerate_AllChars(b *testing.B) {
	for b.Loop() {
		_ = AllChars.Generate(32)
	}
}

// 不同长度基准测试（字符集固定 AllChars）
func BenchmarkStringGenerate_Len8(b *testing.B) {
	for b.Loop() {
		_ = AllChars.Generate(8)
	}
}

func BenchmarkStringGenerate_Len16(b *testing.B) {
	for b.Loop() {
		_ = AllChars.Generate(16)
	}
}

func BenchmarkStringGenerate_Len32(b *testing.B) {
	for b.Loop() {
		_ = AllChars.Generate(32)
	}
}

func BenchmarkStringGenerate_Len64(b *testing.B) {
	for b.Loop() {
		_ = AllChars.Generate(64)
	}
}

func BenchmarkStringGenerate_Len128(b *testing.B) {
	for b.Loop() {
		_ = AllChars.Generate(128)
	}
}

// 带前缀基准测试
func BenchmarkStringGenerate_WithPrefix(b *testing.B) {
	for b.Loop() {
		_ = AllChars.Generate(32, "ID-")
	}
}

// 并发基准测试
func BenchmarkStringGenerate_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = AllChars.Generate(32)
		}
	})
}

func BenchmarkNormalRandIntn_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = NormalGenerator.Intn(100)
		}
	})
}

func BenchmarkSecureRandIntn_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = SecureGenerator.Intn(100)
		}
	})
}

// 十六进制字符串生成方式对比
// 方法 1: Hex.Generate(32) - 字符集索引
func BenchmarkHexString_Generate(b *testing.B) {
	for b.Loop() {
		_ = Hex.Generate(32)
	}
}

// 方法 2: FastRandBytes(16) + hex.EncodeToString - 字节转十六进制
func BenchmarkHexString_FastRandBytes(b *testing.B) {
	for b.Loop() {
		bytes := FastRandBytes(16) // 16 字节 = 32 个十六进制字符
		_ = stdhex.EncodeToString(bytes)
	}
}

// 方法 3: FastRandBytes(16) + bytex.Bytes.Hex() - bytex 零拷贝转换
func BenchmarkHexString_BytexHex(b *testing.B) {
	for b.Loop() {
		bytes := FastRandBytes(16)
		_ = bytex.Bytes(bytes).Hex()
	}
}

// 不同长度的十六进制对比
func BenchmarkHexString_Generate_Len8(b *testing.B) {
	for b.Loop() {
		_ = Hex.Generate(8)
	}
}

func BenchmarkHexString_FastRandBytes_Len8(b *testing.B) {
	for b.Loop() {
		bytes := FastRandBytes(4)
		_ = stdhex.EncodeToString(bytes)
	}
}

func BenchmarkHexString_Generate_Len64(b *testing.B) {
	for b.Loop() {
		_ = Hex.Generate(64)
	}
}

func BenchmarkHexString_FastRandBytes_Len64(b *testing.B) {
	for b.Loop() {
		bytes := FastRandBytes(32)
		_ = stdhex.EncodeToString(bytes)
	}
}

func BenchmarkHexString_Generate_Len128(b *testing.B) {
	for b.Loop() {
		_ = Hex.Generate(128)
	}
}

func BenchmarkHexString_FastRandBytes_Len128(b *testing.B) {
	for b.Loop() {
		bytes := FastRandBytes(64)
		_ = stdhex.EncodeToString(bytes)
	}
}
