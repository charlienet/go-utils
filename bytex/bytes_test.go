package bytex

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytes_Equal(t *testing.T) {
	tests := []struct {
		name     string
		a        Bytes
		b        Bytes
		expected bool
	}{
		{"相同内容", Bytes([]byte("hello")), Bytes([]byte("hello")), true},
		{"不同内容", Bytes([]byte("hello")), Bytes([]byte("world")), false},
		{"长度不同", Bytes([]byte("hello")), Bytes([]byte("hi")), false},
		{"空切片", Bytes([]byte{}), Bytes([]byte{}), true},
		{"nil 和空", nil, Bytes([]byte{}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Equal(tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBytes_Compare(t *testing.T) {
	tests := []struct {
		name     string
		a        Bytes
		b        Bytes
		expected int
	}{
		{"相等", Bytes([]byte("hello")), Bytes([]byte("hello")), 0},
		{"小于", Bytes([]byte("abc")), Bytes([]byte("def")), -1},
		{"大于", Bytes([]byte("def")), Bytes([]byte("abc")), 1},
		{"前缀", Bytes([]byte("hello")), Bytes([]byte("helloworld")), -1},
		{"空切片", Bytes([]byte{}), Bytes([]byte{}), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Compare(tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBytes_Contains(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		sub      []byte
		expected bool
	}{
		{"包含", Bytes([]byte("hello world")), []byte("world"), true},
		{"不包含", Bytes([]byte("hello world")), []byte("xyz"), false},
		{"空子序列", Bytes([]byte("hello")), []byte{}, true},
		{"完全匹配", Bytes([]byte("hello")), []byte("hello"), true},
		{"空数据", Bytes([]byte{}), []byte("test"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.data.Contains(tt.sub)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBytes_Index(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		sub      []byte
		expected int
	}{
		{"找到", Bytes([]byte("hello world")), []byte("world"), 6},
		{"未找到", Bytes([]byte("hello world")), []byte("xyz"), -1},
		{"开头", Bytes([]byte("hello")), []byte("hello"), 0},
		{"结尾", Bytes([]byte("hello")), []byte("lo"), 3},
		{"空子序列", Bytes([]byte("hello")), []byte{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.data.Index(tt.sub)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBytes_HasPrefix(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		prefix   []byte
		expected bool
	}{
		{"有前缀", Bytes([]byte("hello world")), []byte("hello"), true},
		{"无前缀", Bytes([]byte("hello world")), []byte("world"), false},
		{"空前缀", Bytes([]byte("hello")), []byte{}, true},
		{"完全匹配", Bytes([]byte("hello")), []byte("hello"), true},
		{"前缀更长", Bytes([]byte("hi")), []byte("hello"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.data.HasPrefix(tt.prefix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBytes_HasSuffix(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		suffix   []byte
		expected bool
	}{
		{"有后缀", Bytes([]byte("hello world")), []byte("world"), true},
		{"无后缀", Bytes([]byte("hello world")), []byte("hello"), false},
		{"空后缀", Bytes([]byte("hello")), []byte{}, true},
		{"完全匹配", Bytes([]byte("hello")), []byte("hello"), true},
		{"后缀更长", Bytes([]byte("hi")), []byte("hello"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.data.HasSuffix(tt.suffix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBytes_ToUint64(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		endian   binary.ByteOrder
		expected uint64
		wantErr  bool
	}{
		{"大端序", Bytes([]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}), binary.BigEndian, 0x0123456789ABCDEF, false},
		{"小端序", Bytes([]byte{0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01}), binary.LittleEndian, 0x0123456789ABCDEF, false},
		{"长度错误", Bytes([]byte{0x01, 0x02, 0x03}), binary.BigEndian, 0, true},
		{"空数据", Bytes([]byte{}), binary.BigEndian, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.data.ToUint64(tt.endian)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBytes_ToInt64(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		endian   binary.ByteOrder
		expected int64
		wantErr  bool
	}{
		{"正数", Bytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7F}), binary.BigEndian, 127, false},
		{"负数", Bytes([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x81}), binary.BigEndian, -127, false},
		{"长度错误", Bytes([]byte{0x01, 0x02}), binary.BigEndian, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.data.ToInt64(tt.endian)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBytes_ToUint32(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		endian   binary.ByteOrder
		expected uint32
		wantErr  bool
	}{
		{"大端序", Bytes([]byte{0x12, 0x34, 0x56, 0x78}), binary.BigEndian, 0x12345678, false},
		{"小端序", Bytes([]byte{0x78, 0x56, 0x34, 0x12}), binary.LittleEndian, 0x12345678, false},
		{"长度错误", Bytes([]byte{0x01, 0x02, 0x03}), binary.BigEndian, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.data.ToUint32(tt.endian)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBytes_ToInt32(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		endian   binary.ByteOrder
		expected int32
		wantErr  bool
	}{
		{"正数", Bytes([]byte{0x00, 0x00, 0x00, 0x7F}), binary.BigEndian, 127, false},
		{"负数", Bytes([]byte{0xFF, 0xFF, 0xFF, 0x81}), binary.BigEndian, -127, false},
		{"长度错误", Bytes([]byte{0x01}), binary.BigEndian, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.data.ToInt32(tt.endian)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBytes_ToUint16(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		endian   binary.ByteOrder
		expected uint16
		wantErr  bool
	}{
		{"大端序", Bytes([]byte{0x12, 0x34}), binary.BigEndian, 0x1234, false},
		{"小端序", Bytes([]byte{0x34, 0x12}), binary.LittleEndian, 0x1234, false},
		{"长度错误", Bytes([]byte{0x01}), binary.BigEndian, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.data.ToUint16(tt.endian)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBytes_ToInt16(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		endian   binary.ByteOrder
		expected int16
		wantErr  bool
	}{
		{"正数", Bytes([]byte{0x00, 0x7F}), binary.BigEndian, 127, false},
		{"负数", Bytes([]byte{0xFF, 0x81}), binary.BigEndian, -127, false},
		{"长度错误", Bytes([]byte{0x01}), binary.BigEndian, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.data.ToInt16(tt.endian)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBytes_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		expected string
	}{
		{"普通字符串", Bytes([]byte("hello")), `"aGVsbG8="`},
		{"空数据", Bytes([]byte{}), `""`},
		{"nil", nil, `null`},
		{"二进制数据", Bytes([]byte{0x00, 0xff, 0xab}), `"AP+r"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.data.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestBytes_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Bytes
		wantErr  bool
	}{
		{"普通字符串", `"aGVsbG8="`, Bytes([]byte("hello")), false},
		{"空数据", `""`, Bytes([]byte{}), false},
		{"二进制数据", `"AP+r"`, Bytes([]byte{0x00, 0xff, 0xab}), false},
		{"无效 JSON", `invalid`, nil, true},
		{"无效 Base64", `"not-valid-base64!!!"`, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Bytes
			err := result.UnmarshalJSON([]byte(tt.input))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBytes_MarshalText(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		expected string
	}{
		{"普通字符串", Bytes([]byte("hello")), "68656c6c6f"},
		{"空数据", Bytes([]byte{}), ""},
		{"二进制数据", Bytes([]byte{0x00, 0xff, 0xab}), "00ffab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.data.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestBytes_UnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Bytes
		wantErr  bool
	}{
		{"普通字符串", "68656c6c6f", Bytes([]byte("hello")), false},
		{"空数据", "", Bytes([]byte{}), false},
		{"二进制数据", "00ffab", Bytes([]byte{0x00, 0xff, 0xab}), false},
		{"大写 Hex", "68656C6C6F", Bytes([]byte("hello")), false},
		{"无效 Hex", "xyz", nil, true},
		{"奇数长度", "abc", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Bytes
			err := result.UnmarshalText([]byte(tt.input))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBytes_Serialization_RoundTrip(t *testing.T) {
	original := Bytes([]byte("hello world"))

	// JSON 往返
	t.Run("JSON", func(t *testing.T) {
		data, err := original.MarshalJSON()
		assert.NoError(t, err)

		var result Bytes
		err = result.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Equal(t, original, result)
	})

	// Text 往返
	t.Run("Text", func(t *testing.T) {
		data, err := original.MarshalText()
		assert.NoError(t, err)

		var result Bytes
		err = result.UnmarshalText(data)
		assert.NoError(t, err)
		assert.Equal(t, original, result)
	})
}

func TestBytes_HexDump(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		expected string
	}{
		{
			"空数据",
			Bytes([]byte{}),
			"",
		},
		{
			"单字节",
			Bytes([]byte{0x00}),
			"00",
		},
		{
			"一行完整",
			Bytes([]byte{0x00, 0x32, 0x23, 0x64, 0x89, 0x32, 0x32, 0x32, 0x00, 0x32, 0x23, 0x64, 0x89, 0x32, 0x32, 0x32}),
			"00 32 23 64 89 32 32 32 00 32 23 64 89 32 32 32",
		},
		{
			"两行",
			Bytes([]byte{0x00, 0x32, 0x23, 0x64, 0x89, 0x32, 0x32, 0x32, 0x00, 0x32, 0x23, 0x64, 0x89, 0x32, 0x32, 0x32, 0xab, 0xcd}),
			"00 32 23 64 89 32 32 32 00 32 23 64 89 32 32 32\nab cd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.data.HexDump()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBytes_HexDumpWidth(t *testing.T) {
	tests := []struct {
		name     string
		data     Bytes
		width    int
		expected string
	}{
		{
			"4字节宽",
			Bytes([]byte{0x00, 0x32, 0x23, 0x64, 0x89, 0x32, 0x32, 0x32}),
			4,
			"00 32 23 64\n89 32 32 32",
		},
		{
			"8字节宽",
			Bytes([]byte{0x00, 0x32, 0x23, 0x64, 0x89, 0x32, 0x32, 0x32, 0xab}),
			8,
			"00 32 23 64 89 32 32 32\nab",
		},
		{
			"无效宽度回退到16",
			Bytes([]byte{0x00, 0x32, 0x23, 0x64}),
			0,
			"00 32 23 64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.data.HexDumpWidth(tt.width)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBytes_RoundTrip(t *testing.T) {
	// 测试 From* 和 To* 的往返转换
	t.Run("Uint64", func(t *testing.T) {
		original := uint64(0x0123456789ABCDEF)
		b := FromUint64(original, binary.BigEndian)
		result, err := b.ToUint64(binary.BigEndian)
		assert.NoError(t, err)
		assert.Equal(t, original, result)
	})

	t.Run("Int64", func(t *testing.T) {
		original := int64(-1234567890)
		b := FromInt64(original, binary.BigEndian)
		result, err := b.ToInt64(binary.BigEndian)
		assert.NoError(t, err)
		assert.Equal(t, original, result)
	})

	t.Run("Uint32", func(t *testing.T) {
		original := uint32(0x12345678)
		b := FromUint32(original, binary.BigEndian)
		result, err := b.ToUint32(binary.BigEndian)
		assert.NoError(t, err)
		assert.Equal(t, original, result)
	})

	t.Run("Int32", func(t *testing.T) {
		original := int32(-123456)
		b := FromInt32(original, binary.BigEndian)
		result, err := b.ToInt32(binary.BigEndian)
		assert.NoError(t, err)
		assert.Equal(t, original, result)
	})

	t.Run("Uint16", func(t *testing.T) {
		original := uint16(0x1234)
		b := FromUint16(original, binary.BigEndian)
		result, err := b.ToUint16(binary.BigEndian)
		assert.NoError(t, err)
		assert.Equal(t, original, result)
	})

	t.Run("Int16", func(t *testing.T) {
		original := int16(-1234)
		b := FromInt16(original, binary.BigEndian)
		result, err := b.ToInt16(binary.BigEndian)
		assert.NoError(t, err)
		assert.Equal(t, original, result)
	})
}
