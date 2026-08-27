package bytex

import (
	"encoding/binary"
	"io"
	"testing"
)

func TestNewReader(t *testing.T) {
	data := FromString("hello world")
	reader := NewReader(data)
	
	if reader.Position() != 0 {
		t.Errorf("Expected position 0, got %d", reader.Position())
	}
	
	if reader.Remaining() != len(data) {
		t.Errorf("Expected remaining %d, got %d", len(data), reader.Remaining())
	}
}

func TestPositionAndRemaining(t *testing.T) {
	data := FromString("hello")
	reader := NewReader(data)
	
	if reader.Position() != 0 {
		t.Errorf("Expected position 0, got %d", reader.Position())
	}
	
	if reader.Remaining() != 5 {
		t.Errorf("Expected remaining 5, got %d", reader.Remaining())
	}
	
	// Move position using Seek
	reader.Seek(2, io.SeekStart)
	if reader.Position() != 2 {
		t.Errorf("Expected position 2, got %d", reader.Position())
	}
	
	if reader.Remaining() != 3 {
		t.Errorf("Expected remaining 3, got %d", reader.Remaining())
	}
}

func TestSeek(t *testing.T) {
	data := FromString("hello world")
	reader := NewReader(data)
	
	// Test seeking from beginning
	pos, err := reader.Seek(3, 0)
	if err != nil || pos != 3 {
		t.Errorf("Seek from beginning failed: %v, pos: %d", err, pos)
	}
	
	// Test seeking from current position
	pos, err = reader.Seek(2, 1)
	if err != nil || pos != 5 {
		t.Errorf("Seek from current position failed: %v, pos: %d", err, pos)
	}
	
	// Test seeking from end
	pos, err = reader.Seek(-2, 2)
	if err != nil || pos != 9 {
		t.Errorf("Seek from end failed: %v, pos: %d", err, pos)
	}
	
	// Test invalid whence
	_, err = reader.Seek(0, 3)
	if err == nil {
		t.Error("Expected error for invalid whence")
	}
	
	// Test seeking beyond end (bytes.Reader allows this)
	pos, err = reader.Seek(100, 0)
	if err != nil || pos != 100 {
		t.Errorf("Seek beyond end should succeed: pos=%d, err=%v", pos, err)
	}
}

func TestReadByte(t *testing.T) {
	data := FromString("abc")
	reader := NewReader(data)
	
	b, err := reader.ReadByte()
	if err != nil || b != 'a' {
		t.Errorf("First read failed: %v, byte: %c", err, b)
	}
	
	b, err = reader.ReadByte()
	if err != nil || b != 'b' {
		t.Errorf("Second read failed: %v, byte: %c", err, b)
	}
	
	b, err = reader.ReadByte()
	if err != nil || b != 'c' {
		t.Errorf("Third read failed: %v, byte: %c", err, b)
	}
	
	// Should return EOF now
	_, err = reader.ReadByte()
	if err != io.EOF {
		t.Errorf("Expected EOF, got: %v", err)
	}
}

func TestReadUint16(t *testing.T) {
	// Test big endian
	data := FromBytes([]byte{0x12, 0x34})
	reader := NewReader(data)
	
	value, err := reader.ReadUint16(binary.BigEndian)
	if err != nil || value != 0x1234 {
		t.Errorf("Big endian read failed: %v, value: %x", err, value)
	}
	
	// Test little endian
	data = FromBytes([]byte{0x34, 0x12})
	reader = NewReader(data)
	
	value, err = reader.ReadUint16(binary.LittleEndian)
	if err != nil || value != 0x1234 {
		t.Errorf("Little endian read failed: %v, value: %x", err, value)
	}
	
	// Test EOF condition
	data = FromBytes([]byte{0x12}) // Only one byte
	reader = NewReader(data)
	
	_, err = reader.ReadUint16(binary.BigEndian)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("Expected ErrUnexpectedEOF for insufficient data, got: %v", err)
	}
}

func TestReadUint32(t *testing.T) {
	// Test big endian
	data := FromBytes([]byte{0x12, 0x34, 0x56, 0x78})
	reader := NewReader(data)
	
	value, err := reader.ReadUint32(binary.BigEndian)
	if err != nil || value != 0x12345678 {
		t.Errorf("Big endian read failed: %v, value: %x", err, value)
	}
	
	// Test little endian
	data = FromBytes([]byte{0x78, 0x56, 0x34, 0x12})
	reader = NewReader(data)
	
	value, err = reader.ReadUint32(binary.LittleEndian)
	if err != nil || value != 0x12345678 {
		t.Errorf("Little endian read failed: %v, value: %x", err, value)
	}
	
	// Test EOF condition
	data = FromBytes([]byte{0x12, 0x34, 0x56}) // Only three bytes
	reader = NewReader(data)
	
	_, err = reader.ReadUint32(binary.BigEndian)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("Expected ErrUnexpectedEOF for insufficient data, got: %v", err)
	}
}

func TestReadUint64(t *testing.T) {
	// Test big endian
	data := FromBytes([]byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0})
	reader := NewReader(data)
	
	value, err := reader.ReadUint64(binary.BigEndian)
	if err != nil || value != 0x123456789ABCDEF0 {
		t.Errorf("Big endian read failed: %v, value: %x", err, value)
	}
	
	// Test little endian
	data = FromBytes([]byte{0xF0, 0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12})
	reader = NewReader(data)
	
	value, err = reader.ReadUint64(binary.LittleEndian)
	if err != nil || value != 0x123456789ABCDEF0 {
		t.Errorf("Little endian read failed: %v, value: %x", err, value)
	}
	
	// Test EOF condition
	data = FromBytes([]byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE}) // Only seven bytes
	reader = NewReader(data)
	
	_, err = reader.ReadUint64(binary.BigEndian)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("Expected ErrUnexpectedEOF for insufficient data, got: %v", err)
	}
}

func TestReadBytes(t *testing.T) {
	data := FromString("hello world")
	reader := NewReader(data)
	
	bytes, err := reader.ReadBytes(5)
	if err != nil || !Equal(bytes, FromString("hello")) {
		t.Errorf("ReadBytes failed: %v, bytes: %s", err, bytes.String())
	}
	
	bytes, err = reader.ReadBytes(1)
	if err != nil || !Equal(bytes, FromString(" ")) {
		t.Errorf("ReadBytes failed: %v, bytes: %s", err, bytes.String())
	}
	
	// Test EOF condition
	_, err = reader.ReadBytes(100)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("Expected ErrUnexpectedEOF for insufficient data, got: %v", err)
	}
}

func TestReadString(t *testing.T) {
	data := FromString("hello world")
	reader := NewReader(data)
	
	str, err := reader.ReadString(5)
	if err != nil || str != "hello" {
		t.Errorf("ReadString failed: %v, string: %s", err, str)
	}
	
	str, err = reader.ReadString(6)
	if err != nil || str != " world" {
		t.Errorf("ReadString failed: %v, string: %s", err, str)
	}
	
	// Test EOF condition
	_, err = reader.ReadString(1)
	if err != io.EOF {
		t.Errorf("Expected EOF for insufficient data, got: %v", err)
	}
}

func TestPeek(t *testing.T) {
	data := FromString("hello")
	reader := NewReader(data)
	
	peeked, err := reader.Peek(3)
	if err != nil || !Equal(peeked, FromString("hel")) {
		t.Errorf("Peek failed: %v, peeked: %s", err, peeked.String())
	}
	
	// Position should not have changed
	if reader.Position() != 0 {
		t.Errorf("Peek changed position, expected 0, got %d", reader.Position())
	}
	
	// Test EOF condition
	_, err = reader.Peek(10)
	if err != io.EOF {
		t.Errorf("Expected EOF for insufficient data in peek, got: %v", err)
	}
}

func TestReset(t *testing.T) {
	data := FromString("hello")
	reader := NewReader(data)
	
	// Move position
	reader.ReadByte()
	reader.ReadByte()
	
	if reader.Position() != 2 {
		t.Errorf("Expected position 2 after reads, got %d", reader.Position())
	}
	
	reader.Reset()
	
	if reader.Position() != 0 {
		t.Errorf("Expected position 0 after reset, got %d", reader.Position())
	}
	
	if reader.Remaining() != 5 {
		t.Errorf("Expected remaining 5 after reset, got %d", reader.Remaining())
	}
}

func TestIOReaderInterface(t *testing.T) {
	data := FromString("hello world")
	reader := NewReader(data)
	
	buffer := make([]byte, 5)
	n, err := reader.Read(buffer)
	if err != nil || n != 5 || string(buffer[:n]) != "hello" {
		t.Errorf("Read failed: n=%d, err=%v, buffer=%s", n, err, string(buffer[:n]))
	}
	
	// Read the rest
	buffer = make([]byte, 10)
	n, err = reader.Read(buffer)
	if err != nil || n != 6 || string(buffer[:n]) != " world" {  // Fixed: should be " world" not " worl"
		t.Errorf("Read failed: n=%d, err=%v, buffer=%s", n, err, string(buffer[:n]))
	}
	
	// Should return EOF now
	buffer = make([]byte, 1)
	n, err = reader.Read(buffer)
	if err != io.EOF || n != 0 {
		t.Errorf("Expected EOF with 0 bytes, got n=%d, err=%v", n, err)
	}
}

func TestProtocolParsingScenario(t *testing.T) {
	// Simulate parsing a simple protocol: [length:uint32][data:string]
	data := make([]byte, 0)
	data = append(data, FromUint32(5, binary.BigEndian)...) // Length: 5
	data = append(data, FromString("hello")...)             // Data: "hello"
	data = append(data, FromUint32(5, binary.BigEndian)...) // Length: 5
	data = append(data, FromString("world")...)             // Data: "world"
	
	reader := NewReader(FromBytes(data))
	
	// Read first message
	length, err := reader.ReadUint32(binary.BigEndian)
	if err != nil || length != 5 {
		t.Fatalf("Failed to read length: %v, length: %d", err, length)
	}
	
	msg, err := reader.ReadString(int(length))
	if err != nil || msg != "hello" {
		t.Fatalf("Failed to read message: %v, msg: %s", err, msg)
	}
	
	// Read second message
	length, err = reader.ReadUint32(binary.BigEndian)
	if err != nil || length != 5 {
		t.Fatalf("Failed to read length: %v, length: %d", err, length)
	}
	
	msg, err = reader.ReadString(int(length))
	if err != nil || msg != "world" {
		t.Fatalf("Failed to read message: %v, msg: %s", err, msg)
	}
	
	// Should be at end
	if reader.Remaining() != 0 {
		t.Errorf("Expected 0 remaining bytes, got %d", reader.Remaining())
	}
}

func TestEmptyData(t *testing.T) {
	reader := NewReader(nil)
	
	_, err := reader.ReadByte()
	if err != io.EOF {
		t.Errorf("Expected EOF for empty data, got: %v", err)
	}
	
	_, err = reader.ReadUint16(binary.BigEndian)
	if err != io.EOF {
		t.Errorf("Expected EOF for empty data, got: %v", err)
	}
	
	_, err = reader.ReadUint32(binary.BigEndian)
	if err != io.EOF {
		t.Errorf("Expected EOF for empty data, got: %v", err)
	}
	
	_, err = reader.ReadUint64(binary.BigEndian)
	if err != io.EOF {
		t.Errorf("Expected EOF for empty data, got: %v", err)
	}
	
	_, err = reader.ReadBytes(1)
	if err != io.EOF {
		t.Errorf("Expected EOF for empty data, got: %v", err)
	}
	
	_, err = reader.ReadString(1)
	if err != io.EOF {
		t.Errorf("Expected EOF for empty data, got: %v", err)
	}
	
	_, err = reader.Peek(1)
	if err != io.EOF {
		t.Errorf("Expected EOF for empty data, got: %v", err)
	}
}