package bytex

import (
	"bytes"
	"encoding/binary"
	"io"
)

// Reader 提供对 Bytes 的流式读取，支持类型化读取和位置控制。
// 基于标准库 bytes.Reader 扩展，兼容 io.Reader、io.Seeker 等接口。
type Reader struct {
	*bytes.Reader
}

// NewReader 创建一个新的 Reader
func NewReader(data Bytes) *Reader {
	return &Reader{Reader: bytes.NewReader(data)}
}

// Position 返回当前读取位置
func (r *Reader) Position() int {
	pos, _ := r.Seek(0, io.SeekCurrent)
	return int(pos)
}

// Remaining 返回剩余可读取的字节数
func (r *Reader) Remaining() int {
	return r.Len()
}

// ReadUint16 读取 uint16，支持字节序
func (r *Reader) ReadUint16(order binary.ByteOrder) (uint16, error) {
	var buf [2]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}
	return order.Uint16(buf[:]), nil
}

// ReadUint32 读取 uint32，支持字节序
func (r *Reader) ReadUint32(order binary.ByteOrder) (uint32, error) {
	var buf [4]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}
	return order.Uint32(buf[:]), nil
}

// ReadUint64 读取 uint64，支持字节序
func (r *Reader) ReadUint64(order binary.ByteOrder) (uint64, error) {
	var buf [8]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}
	return order.Uint64(buf[:]), nil
}

// ReadBytes 读取 n 个字节，返回 Bytes 类型
func (r *Reader) ReadBytes(n int) (Bytes, error) {
	buf := make(Bytes, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// ReadString 读取 n 个字节并转换为字符串
func (r *Reader) ReadString(n int) (string, error) {
	buf := make(Bytes, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}
	return BytesToString(buf), nil
}

// Peek 预览接下来的 n 个字节，不移动位置
func (r *Reader) Peek(n int) (Bytes, error) {
	pos, _ := r.Seek(0, io.SeekCurrent)
	if int(pos)+n > r.Len()+int(pos) {
		return nil, io.EOF
	}
	buf := make(Bytes, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	// 恢复位置
	_, _ = r.Seek(pos, io.SeekStart)
	return buf, nil
}

// ResetTo 重置读取位置到指定偏移
func (r *Reader) ResetTo(offset int64) {
	_, _ = r.Seek(offset, io.SeekStart)
}

// Reset 重置读取位置到开头
func (r *Reader) Reset() {
	_, _ = r.Seek(0, io.SeekStart)
}
