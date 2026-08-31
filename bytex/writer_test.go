package bytex

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestNewWriter(t *testing.T) {
	writer := NewWriter()
	if writer == nil {
		t.Fatal("NewWriter() returned nil")
	}
	if writer.Len() != 0 {
		t.Errorf("Expected length 0, got %d", writer.Len())
	}
}

func TestNewWriterSize(t *testing.T) {
	writer := NewWriterSize(100)
	if writer == nil {
		t.Fatal("NewWriterSize(100) returned nil")
	}
	if writer.Len() != 0 {
		t.Errorf("Expected length 0, got %d", writer.Len())
	}
}

func TestWriter_WriteUint16(t *testing.T) {
	writer := NewWriter()

	err := writer.WriteUint16(0x1234, binary.BigEndian)
	if err != nil {
		t.Fatalf("WriteUint16 failed: %v", err)
	}

	expected := []byte{0x12, 0x34}
	actual := []byte(writer.Bytes())
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestWriter_WriteUint32(t *testing.T) {
	writer := NewWriter()

	err := writer.WriteUint32(0x12345678, binary.BigEndian)
	if err != nil {
		t.Fatalf("WriteUint32 failed: %v", err)
	}

	expected := []byte{0x12, 0x34, 0x56, 0x78}
	actual := []byte(writer.Bytes())
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestWriter_WriteUint64(t *testing.T) {
	writer := NewWriter()

	err := writer.WriteUint64(0x1234567890ABCDEF, binary.BigEndian)
	if err != nil {
		t.Fatalf("WriteUint64 failed: %v", err)
	}

	expected := []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF}
	actual := []byte(writer.Bytes())
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestWriter_WriteBytes(t *testing.T) {
	writer := NewWriter()

	data := []byte("hello")
	err := writer.WriteBytes(data)
	if err != nil {
		t.Fatalf("WriteBytes failed: %v", err)
	}

	expected := []byte("hello")
	actual := []byte(writer.Bytes())
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestWriter_WriteString(t *testing.T) {
	writer := NewWriter()

	err := writer.WriteString("world")
	if err != nil {
		t.Fatalf("WriteString failed: %v", err)
	}

	expected := []byte("world")
	actual := []byte(writer.Bytes())
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestWriter_MultipleWrites(t *testing.T) {
	writer := NewWriter()

	// Write multiple values
	_ = writer.WriteUint16(0x1234, binary.BigEndian)
	_ = writer.WriteUint32(0x567890AB, binary.BigEndian)
	_ = writer.WriteString("test")

	expected := []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 't', 'e', 's', 't'}
	actual := []byte(writer.Bytes())
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestWriter_Len(t *testing.T) {
	writer := NewWriter()

	initialLen := writer.Len()
	if initialLen != 0 {
		t.Errorf("Expected initial length 0, got %d", initialLen)
	}

	_ = writer.WriteString("hello")
	newLen := writer.Len()
	if newLen != 5 {
		t.Errorf("Expected length 5 after writing 'hello', got %d", newLen)
	}
}

func TestWriter_Reset(t *testing.T) {
	writer := NewWriter()
	_ = writer.WriteString("hello")

	if writer.Len() == 0 {
		t.Error("Expected non-zero length after writing")
	}

	writer.Reset()

	if writer.Len() != 0 {
		t.Errorf("Expected length 0 after reset, got %d", writer.Len())
	}
}
