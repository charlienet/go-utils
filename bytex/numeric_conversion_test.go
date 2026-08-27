package bytex

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromUint64(t *testing.T) {
	tests := []struct {
		name     string
		value    uint64
		endian   binary.ByteOrder
		expected []byte
	}{
		{
			name:     "BigEndian zero",
			value:    0,
			endian:   binary.BigEndian,
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "BigEndian max",
			value:    ^uint64(0), // All bits set
			endian:   binary.BigEndian,
			expected: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name:     "BigEndian specific value",
			value:    0x123456789ABCDEF0,
			endian:   binary.BigEndian,
			expected: []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0},
		},
		{
			name:     "LittleEndian specific value",
			value:    0x123456789ABCDEF0,
			endian:   binary.LittleEndian,
			expected: []byte{0xF0, 0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromUint64(tt.value, tt.endian)
			assert.Equal(t, tt.expected, []byte(result))
			
			// Test round trip
			roundTripValue := tt.endian.Uint64(result)
			assert.Equal(t, tt.value, roundTripValue)
		})
	}
}

func TestFromInt64(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		endian   binary.ByteOrder
		expected []byte
	}{
		{
			name:     "BigEndian positive",
			value:    0x123456789ABCDEF0,
			endian:   binary.BigEndian,
			expected: []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0},
		},
		{
			name:     "LittleEndian positive",
			value:    0x123456789ABCDEF0,
			endian:   binary.LittleEndian,
			expected: []byte{0xF0, 0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12},
		},
		{
			name:     "BigEndian negative",
			value:    -1,
			endian:   binary.BigEndian,
			expected: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name:     "BigEndian_min",
			value:    math.MinInt64,
			endian:   binary.BigEndian,
			expected: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:     "LittleEndian_min",
			value:    math.MinInt64,
			endian:   binary.LittleEndian,
			expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromInt64(tt.value, tt.endian)
			assert.Equal(t, tt.expected, []byte(result))
			
			// Test round trip (reading back as uint64 for comparison)
			roundTripValue := int64(tt.endian.Uint64(result))
			assert.Equal(t, tt.value, roundTripValue)
		})
	}
}

func TestFromUint32(t *testing.T) {
	tests := []struct {
		name     string
		value    uint32
		endian   binary.ByteOrder
		expected []byte
	}{
		{
			name:     "BigEndian zero",
			value:    0,
			endian:   binary.BigEndian,
			expected: []byte{0, 0, 0, 0},
		},
		{
			name:     "BigEndian max",
			value:    ^uint32(0), // All bits set
			endian:   binary.BigEndian,
			expected: []byte{0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name:     "BigEndian specific value",
			value:    0x12345678,
			endian:   binary.BigEndian,
			expected: []byte{0x12, 0x34, 0x56, 0x78},
		},
		{
			name:     "LittleEndian specific value",
			value:    0x12345678,
			endian:   binary.LittleEndian,
			expected: []byte{0x78, 0x56, 0x34, 0x12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromUint32(tt.value, tt.endian)
			assert.Equal(t, tt.expected, []byte(result))
			
			// Test round trip
			roundTripValue := tt.endian.Uint32(result)
			assert.Equal(t, tt.value, roundTripValue)
		})
	}
}

func TestFromInt32(t *testing.T) {
	tests := []struct {
		name     string
		value    int32
		endian   binary.ByteOrder
		expected []byte
	}{
		{
			name:     "BigEndian positive",
			value:    0x12345678,
			endian:   binary.BigEndian,
			expected: []byte{0x12, 0x34, 0x56, 0x78},
		},
		{
			name:     "LittleEndian positive",
			value:    0x12345678,
			endian:   binary.LittleEndian,
			expected: []byte{0x78, 0x56, 0x34, 0x12},
		},
		{
			name:     "BigEndian negative",
			value:    -1,
			endian:   binary.BigEndian,
			expected: []byte{0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name:     "BigEndian_min",
			value:    math.MinInt32,
			endian:   binary.BigEndian,
			expected: []byte{0x80, 0x00, 0x00, 0x00},
		},
		{
			name:     "LittleEndian_min",
			value:    math.MinInt32,
			endian:   binary.LittleEndian,
			expected: []byte{0x00, 0x00, 0x00, 0x80},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromInt32(tt.value, tt.endian)
			assert.Equal(t, tt.expected, []byte(result))
			
			// Test round trip (reading back as uint32 for comparison)
			roundTripValue := int32(tt.endian.Uint32(result))
			assert.Equal(t, tt.value, roundTripValue)
		})
	}
}

func TestFromUint16(t *testing.T) {
	tests := []struct {
		name     string
		value    uint16
		endian   binary.ByteOrder
		expected []byte
	}{
		{
			name:     "BigEndian zero",
			value:    0,
			endian:   binary.BigEndian,
			expected: []byte{0, 0},
		},
		{
			name:     "BigEndian max",
			value:    ^uint16(0), // All bits set
			endian:   binary.BigEndian,
			expected: []byte{0xFF, 0xFF},
		},
		{
			name:     "BigEndian specific value",
			value:    0x1234,
			endian:   binary.BigEndian,
			expected: []byte{0x12, 0x34},
		},
		{
			name:     "LittleEndian specific value",
			value:    0x1234,
			endian:   binary.LittleEndian,
			expected: []byte{0x34, 0x12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromUint16(tt.value, tt.endian)
			assert.Equal(t, tt.expected, []byte(result))
			
			// Test round trip
			roundTripValue := tt.endian.Uint16(result)
			assert.Equal(t, tt.value, roundTripValue)
		})
	}
}

func TestFromInt16(t *testing.T) {
	tests := []struct {
		name     string
		value    int16
		endian   binary.ByteOrder
		expected []byte
	}{
		{
			name:     "BigEndian positive",
			value:    0x1234,
			endian:   binary.BigEndian,
			expected: []byte{0x12, 0x34},
		},
		{
			name:     "LittleEndian positive",
			value:    0x1234,
			endian:   binary.LittleEndian,
			expected: []byte{0x34, 0x12},
		},
		{
			name:     "BigEndian negative",
			value:    -1,
			endian:   binary.BigEndian,
			expected: []byte{0xFF, 0xFF},
		},
		{
			name:     "BigEndian_min",
			value:    math.MinInt16,
			endian:   binary.BigEndian,
			expected: []byte{0x80, 0x00},
		},
		{
			name:     "LittleEndian_min",
			value:    math.MinInt16,
			endian:   binary.LittleEndian,
			expected: []byte{0x00, 0x80},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromInt16(tt.value, tt.endian)
			assert.Equal(t, tt.expected, []byte(result))
			
			// Test round trip (reading back as uint16 for comparison)
			roundTripValue := int16(tt.endian.Uint16(result))
			assert.Equal(t, tt.value, roundTripValue)
		})
	}
}

func TestEndianImplementsBinaryByteOrder(t *testing.T) {
	// Verify that our custom endian types implement the binary.ByteOrder interface
	var _ binary.ByteOrder = BigEndian
	var _ binary.ByteOrder = LittleEndian
	
	// Test that the implementations work correctly
	data := make([]byte, 8)
	value := uint64(0x123456789ABCDEF0)
	
	// Test BigEndian
	BigEndian.PutUint64(data, value)
	readBack := BigEndian.Uint64(data)
	assert.Equal(t, value, readBack)
	
	// Reset data
	data = make([]byte, 8)
	
	// Test LittleEndian
	LittleEndian.PutUint64(data, value)
	readBack = LittleEndian.Uint64(data)
	assert.Equal(t, value, readBack)
}

func TestZeroCopyFromString(t *testing.T) {
	original := "hello world"
	result := FromString(original)
	
	// The content should match
	assert.Equal(t, original, string(result))
	
	// Length should match
	assert.Len(t, result, len(original))
}