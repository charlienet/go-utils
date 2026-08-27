package bytex

import (
	"encoding/binary"
	"fmt"
)

type endian int

const (
	BigEndian endian = iota + 1
	LittleEndian
)

// 实现 binary.ByteOrder 接口
func (e endian) Uint16(data []byte) uint16 {
	if e == BigEndian {
		return binary.BigEndian.Uint16(data)
	}
	return binary.LittleEndian.Uint16(data)
}

func (e endian) Uint32(data []byte) uint32 {
	if e == BigEndian {
		return binary.BigEndian.Uint32(data)
	}
	return binary.LittleEndian.Uint32(data)
}

func (e endian) Uint64(data []byte) uint64 {
	if e == BigEndian {
		return binary.BigEndian.Uint64(data)
	}
	return binary.LittleEndian.Uint64(data)
}

func (e endian) PutUint16(data []byte, v uint16) {
	if e == BigEndian {
		binary.BigEndian.PutUint16(data, v)
	} else {
		binary.LittleEndian.PutUint16(data, v)
	}
}

func (e endian) PutUint32(data []byte, v uint32) {
	if e == BigEndian {
		binary.BigEndian.PutUint32(data, v)
	} else {
		binary.LittleEndian.PutUint32(data, v)
	}
}

func (e endian) PutUint64(data []byte, v uint64) {
	if e == BigEndian {
		binary.BigEndian.PutUint64(data, v)
	} else {
		binary.LittleEndian.PutUint64(data, v)
	}
}

func (e endian) String() string {
	switch e {
	case BigEndian:
		return "BigEndian"
	case LittleEndian:
		return "LittleEndian"
	default:
		return fmt.Sprintf("endian(%d)", e)
	}
}

// BytesToUint64 将字节切片转换为 uint64。
// 支持 1-8 字节的输入，按指定字节序解析。
func (e endian) BytesToUint64(data []byte) (uint64, error) {
	if len(data) > 8 {
		return 0, fmt.Errorf("bytes to uint64, bytes length is too long, max 8 bytes, got %d", len(data))
	}

	// 创建一个8字节的缓冲区并填充数据
	buf := make([]byte, 8)
	
	if e == BigEndian {
		// 对于大端序，在前面补零
		copy(buf[8-len(data):], data)
		return binary.BigEndian.Uint64(buf), nil
	} else {
		// 对于小端序，在后面补零
		copy(buf, data)
		return binary.LittleEndian.Uint64(buf), nil
	}
}


