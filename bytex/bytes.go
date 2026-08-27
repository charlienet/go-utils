package bytex

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

const (
	hexTable    = "0123456789ABCDEF"
	lowerHexTable = "0123456789abcdef"
)

type Bytes []byte

// FromString 将字符串转换为 Bytes。
// 注意：返回的 Bytes 与原始字符串共享内存，不可修改。
// 修改返回的 Bytes 会导致未定义行为。
func FromString(s string) Bytes {
	return Bytes(StringToBytes(s))
}

// FromBytes 将字节切片转换为 Bytes 类型。
// 注意：返回的 Bytes 与传入的切片共享底层数组，修改任一方都会影响另一方。
// 如需独立副本，请使用 FromBytes(b).Clone()。
func FromBytes(b []byte) Bytes {
	return Bytes(b)
}

// FromHexString 从十六进制字符串创建 Bytes。
// s 必须是偶数长度的十六进制字符串（如 "68656c6c6f"），否则返回错误。
func FromHexString(s string) (Bytes, error) {
	b, err := hex.DecodeString(s)
	return Bytes(b), err
}

// FromBase64String 从 Base64 字符串创建 Bytes。
// s 必须是标准 Base64 编码的字符串，否则返回错误。
func FromBase64String(s string) (Bytes, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	return Bytes(b), err
}


// FromUint64 将 uint64 转换为字节切片。
// endian 参数指定字节序：binary.BigEndian 或 binary.LittleEndian。
func FromUint64(v uint64, endian binary.ByteOrder) Bytes {
	b := make(Bytes, 8)
	endian.PutUint64(b, v)
	return b
}

// FromInt64 将 int64 转换为字节切片。
// endian 参数指定字节序：binary.BigEndian 或 binary.LittleEndian。
func FromInt64(v int64, endian binary.ByteOrder) Bytes {
	b := make(Bytes, 8)
	endian.PutUint64(b, uint64(v))
	return b
}

// FromUint32 将 uint32 转换为字节切片。
// endian 参数指定字节序：binary.BigEndian 或 binary.LittleEndian。
func FromUint32(v uint32, endian binary.ByteOrder) Bytes {
	b := make(Bytes, 4)
	endian.PutUint32(b, v)
	return b
}

// FromInt32 将 int32 转换为字节切片。
// endian 参数指定字节序：binary.BigEndian 或 binary.LittleEndian。
func FromInt32(v int32, endian binary.ByteOrder) Bytes {
	b := make(Bytes, 4)
	endian.PutUint32(b, uint32(v))
	return b
}

// FromUint16 将 uint16 转换为字节切片。
// endian 参数指定字节序：binary.BigEndian 或 binary.LittleEndian。
func FromUint16(v uint16, endian binary.ByteOrder) Bytes {
	b := make(Bytes, 2)
	endian.PutUint16(b, v)
	return b
}

// FromInt16 将 int16 转换为字节切片。
// endian 参数指定字节序：binary.BigEndian 或 binary.LittleEndian。
func FromInt16(v int16, endian binary.ByteOrder) Bytes {
	b := make(Bytes, 2)
	endian.PutUint16(b, uint16(v))
	return b
}

// Hex 返回小写十六进制字符串表示。
func (r Bytes) Hex() string {
	return r.encodeHex(lowerHexTable)
}

// UpperHex 返回大写十六进制字符串表示。
func (r Bytes) UpperHex() string {
	return r.encodeHex(hexTable)
}

// encodeHex 使用指定的十六进制表编码字节为字符串。
func (r Bytes) encodeHex(table string) string {
	dst := make([]byte, hex.EncodedLen(len(r)))
	j := 0
	for _, v := range r {
		dst[j] = table[v>>4]
		dst[j+1] = table[v&0x0f]
		j += 2
	}
	return BytesToString(dst)
}

func (r Bytes) Base64() string {
	return base64.StdEncoding.EncodeToString(r)
}

// Bytes 返回底层字节切片引用。
// 注意：修改返回的切片会影响原 Bytes 数据。如需独立副本，请使用 Clone()。
func (r Bytes) Bytes() []byte {
	return r
}

// String 返回安全的可打印字符串表示。
// 将字节数据转换为带转义的字符串格式（如 "hello\x00\xff"），
// 确保所有字节都能安全打印，不会输出异常字符。
func (r Bytes) String() string {
	return strconv.Quote(BytesToString(r))
}

// Clone 返回字节数据的深拷贝。
// 修改返回的 Bytes 不会影响原数据。
// 对于 nil receiver，返回空 Bytes（非 nil）。
func (r Bytes) Clone() Bytes {
	c := make(Bytes, len(r))
	copy(c, r)
	return c
}

// Len 返回字节长度。
func (r Bytes) Len() int {
	return len(r)
}

func (r Bytes) Open() io.Reader {
	return bytes.NewReader(r)
}

// Equal 比较两个 Bytes 是否相等
func (r Bytes) Equal(other Bytes) bool {
	return bytes.Equal(r, other)
}

// Compare 比较两个 Bytes 的大小，返回 -1（小于）、0（等于）或 1（大于）
func (r Bytes) Compare(other Bytes) int {
	return bytes.Compare(r, other)
}

// Contains 检查是否包含指定的子序列
func (r Bytes) Contains(sub []byte) bool {
	return bytes.Contains(r, sub)
}

// Index 查找子序列第一次出现的位置，未找到返回 -1
func (r Bytes) Index(sub []byte) int {
	return bytes.Index(r, sub)
}

// HasPrefix 检查是否以指定的前缀开头
func (r Bytes) HasPrefix(prefix []byte) bool {
	return bytes.HasPrefix(r, prefix)
}

// HasSuffix 检查是否以指定的后缀结尾
func (r Bytes) HasSuffix(suffix []byte) bool {
	return bytes.HasSuffix(r, suffix)
}

// ToUint64 从字节解析 uint64
func (r Bytes) ToUint64(endian binary.ByteOrder) (uint64, error) {
	if len(r) != 8 {
		return 0, fmt.Errorf("invalid length %d, expected 8", len(r))
	}
	return endian.Uint64(r), nil
}

// ToInt64 从字节解析 int64
func (r Bytes) ToInt64(endian binary.ByteOrder) (int64, error) {
	v, err := r.ToUint64(endian)
	return int64(v), err
}

// ToUint32 从字节解析 uint32
func (r Bytes) ToUint32(endian binary.ByteOrder) (uint32, error) {
	if len(r) != 4 {
		return 0, fmt.Errorf("invalid length %d, expected 4", len(r))
	}
	return endian.Uint32(r), nil
}

// ToInt32 从字节解析 int32
func (r Bytes) ToInt32(endian binary.ByteOrder) (int32, error) {
	v, err := r.ToUint32(endian)
	return int32(v), err
}

// ToUint16 从字节解析 uint16
func (r Bytes) ToUint16(endian binary.ByteOrder) (uint16, error) {
	if len(r) != 2 {
		return 0, fmt.Errorf("invalid length %d, expected 2", len(r))
	}
	return endian.Uint16(r), nil
}

// ToInt16 从字节解析 int16
func (r Bytes) ToInt16(endian binary.ByteOrder) (int16, error) {
	v, err := r.ToUint16(endian)
	return int16(v), err
}

// MarshalJSON 实现 json.Marshaler 接口。
// 使用 Base64 编码，与标准库对 []byte 的 JSON 序列化行为一致。
func (r Bytes) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}
	encoded := base64.StdEncoding.EncodeToString(r)
	return []byte(`"` + encoded + `"`), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口。
// 使用 Base64 解码，与标准库对 []byte 的 JSON 反序列化行为一致。
func (r *Bytes) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid JSON string: %s", string(data))
	}
	decoded, err := base64.StdEncoding.DecodeString(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	*r = Bytes(decoded)
	return nil
}

// MarshalText 实现 encoding.TextMarshaler 接口。
// 使用 Hex 编码，便于日志输出和调试。
func (r Bytes) MarshalText() ([]byte, error) {
	return []byte(r.Hex()), nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler 接口。
// 使用 Hex 解码。
func (r *Bytes) UnmarshalText(data []byte) error {
	decoded, err := hex.DecodeString(string(data))
	if err != nil {
		return err
	}
	*r = Bytes(decoded)
	return nil
}
