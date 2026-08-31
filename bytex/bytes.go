package bytex

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	hexTable      = "0123456789ABCDEF"
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

// Base32 返回 Base32 编码字符串
func (r Bytes) Base32() string {
	return base32.StdEncoding.EncodeToString(r)
}

// Base58 返回 Base58 编码字符串（比特币地址常用）
// Base58 字母表：123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz
func (r Bytes) Base58() string {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	if len(r) == 0 {
		return ""
	}

	// 计算前导零的数量
	zeros := 0
	for zeros < len(r) && r[zeros] == 0 {
		zeros++
	}

	// 将字节转换为 base58 数字
	b58 := make([]int, len(r)*138/100+1) // log(256)/log(58), rounded up

	for _, digit := range r {
		carry := int(digit)
		for j := len(b58) - 1; j >= 0; j-- {
			carry += b58[j] << 8 // b58[j] * 256
			b58[j] = carry % 58
			carry /= 58
		}
	}

	// 找到第一个非零数字
	start := 0
	for start < len(b58) && b58[start] == 0 {
		start++
	}

	// 构建结果字符串
	var result strings.Builder
	result.Grow(zeros + len(b58) - start)

	// 添加前导 '1' (代表零)
	for i := 0; i < zeros; i++ {
		result.WriteByte(alphabet[0])
	}

	// 添加实际的 base58 编码
	for i := start; i < len(b58); i++ {
		result.WriteByte(alphabet[b58[i]])
	}

	return result.String()
}

// Base64URL 返回 URL 安全的 Base64 编码字符串
func (r Bytes) Base64URL() string {
	return base64.URLEncoding.EncodeToString(r)
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

// Equal 比较两个字节切片是否相等
func Equal(a, b []byte) bool {
	return bytes.Equal(a, b)
}

// Compare 比较两个字节切片的大小，返回 -1（小于）、0（等于）或 1（大于）
func Compare(a, b []byte) int {
	return bytes.Compare(a, b)
}

// Contains 检查 b 是否包含 subslice
func Contains(b, subslice []byte) bool {
	return bytes.Contains(b, subslice)
}

// Index 查找 sep 在 s 中第一次出现的位置，未找到返回 -1
func Index(s, sep []byte) int {
	return bytes.Index(s, sep)
}

// HasPrefix 检查 s 是否以 prefix 开头
func HasPrefix(s, prefix []byte) bool {
	return bytes.HasPrefix(s, prefix)
}

// HasSuffix 检查 s 是否以 suffix 结尾
func HasSuffix(s, suffix []byte) bool {
	return bytes.HasSuffix(s, suffix)
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

// Join 拼接多个字节切片，使用 sep 分隔
func Join(sep []byte, items ...[]byte) Bytes {
	result := bytes.Join(items, sep)
	return Bytes(result)
}

// Split 分割字节切片
func Split(b, sep []byte) []Bytes {
	split := bytes.Split(b, sep)
	result := make([]Bytes, len(split))
	for i, s := range split {
		result[i] = Bytes(s)
	}
	return result
}

// Repeat 重复字节切片 count 次
func Repeat(b []byte, count int) Bytes {
	repeated := bytes.Repeat(b, count)
	return Bytes(repeated)
}

// HexDump 返回格式化的十六进制输出，每行 16 字节。
// 适合调试和日志输出。
func (r Bytes) HexDump() string {
	return r.HexDumpWidth(16)
}

// HexDumpWidth 返回格式化的十六进制输出，每行 width 字节。
// 适合调试和日志输出。
func (r Bytes) HexDumpWidth(width int) string {
	if len(r) == 0 {
		return ""
	}
	if width <= 0 {
		width = 16
	}

	var buf bytes.Buffer
	for i := 0; i < len(r); i += width {
		end := min(i+width, len(r))
		for j := i; j < end; j++ {
			if j > i {
				buf.WriteByte(' ')
			}
			buf.WriteString(lowerHexTable[r[j]>>4 : r[j]>>4+1])
			buf.WriteString(lowerHexTable[r[j]&0x0f : r[j]&0x0f+1])
		}
		if end < len(r) {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

// Reverse 返回反转后的字节切片副本
func (r Bytes) Reverse() Bytes {
	result := make(Bytes, len(r))
	for i, j := 0, len(r)-1; i < len(r); i, j = i+1, j-1 {
		result[i] = r[j]
	}
	return result
}

// Slice 返回字节切片的子切片 [start:end]
// 如果 start 或 end 超出范围，会自动调整
func (r Bytes) Slice(start, end int) Bytes {
	// 调整 start
	if start < 0 {
		start = 0
	} else if start > len(r) {
		start = len(r)
	}

	// 调整 end
	if end < start {
		end = start
	} else if end > len(r) {
		end = len(r)
	}

	return r[start:end]
}

// Trim 返回去除指定字节集合后的副本
func (r Bytes) Trim(cutset []byte) Bytes {
	return Bytes(bytes.Trim(r, string(cutset)))
}

// TrimLeft 返回去除左侧指定字节集合后的副本
func (r Bytes) TrimLeft(cutset []byte) Bytes {
	return Bytes(bytes.TrimLeft(r, string(cutset)))
}

// TrimRight 返回去除右侧指定字节集合后的副本
func (r Bytes) TrimRight(cutset []byte) Bytes {
	return Bytes(bytes.TrimRight(r, string(cutset)))
}
