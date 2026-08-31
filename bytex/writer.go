package bytex

import (
	"bytes"
	"encoding/binary"
)

// Writer 提供对 Bytes 的流式写入，支持类型化写入和位置控制。
// 基于标准库 bytes.Buffer 扩展。
type Writer struct {
	*bytes.Buffer
}

// NewWriter 创建一个新的 Writer
func NewWriter() *Writer {
	return &Writer{Buffer: &bytes.Buffer{}}
}

// NewWriterSize 创建指定初始容量的 Writer
func NewWriterSize(size int) *Writer {
	return &Writer{Buffer: bytes.NewBuffer(make([]byte, 0, size))}
}

// Bytes 返回已写入的字节数据
func (w *Writer) Bytes() Bytes {
	return Bytes(w.Buffer.Bytes())
}

// Len 返回已写入的字节长度
func (w *Writer) Len() int {
	return w.Buffer.Len()
}

// Reset 重置 Writer
func (w *Writer) Reset() {
	w.Buffer.Reset()
}

// WriteUint16 写入 uint16，支持字节序
func (w *Writer) WriteUint16(v uint16, order binary.ByteOrder) error {
	var buf [2]byte
	order.PutUint16(buf[:], v)
	_, err := w.Write(buf[:])
	return err
}

// WriteUint32 写入 uint32，支持字节序
func (w *Writer) WriteUint32(v uint32, order binary.ByteOrder) error {
	var buf [4]byte
	order.PutUint32(buf[:], v)
	_, err := w.Write(buf[:])
	return err
}

// WriteUint64 写入 uint64，支持字节序
func (w *Writer) WriteUint64(v uint64, order binary.ByteOrder) error {
	var buf [8]byte
	order.PutUint64(buf[:], v)
	_, err := w.Write(buf[:])
	return err
}

// WriteBytes 写入 Bytes 数据
func (w *Writer) WriteBytes(b Bytes) error {
	_, err := w.Write(b)
	return err
}

// WriteString 写入字符串
func (w *Writer) WriteString(s string) error {
	_, err := w.Buffer.WriteString(s)
	return err
}
