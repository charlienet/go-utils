package mathx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFenToYuan(t *testing.T) {
	tests := []struct {
		fen      interface{}
		expected string
	}{
		{100, "1.00"},
		{1, "0.01"},
		{0, "0.00"},
		{999, "9.99"},
		{12345, "123.45"},
		{-100, "-1.00"},
		{int32(100), "1.00"},
		{int64(1), "0.01"},
		{uint(999), "9.99"},
		{uint64(12345), "123.45"},
	}

	for _, tt := range tests {
		switch v := tt.fen.(type) {
		case int:
			result := FenToYuan(v)
			assert.Equal(t, tt.expected, result)
		case int32:
			result := FenToYuan(v)
			assert.Equal(t, tt.expected, result)
		case int64:
			result := FenToYuan(v)
			assert.Equal(t, tt.expected, result)
		case uint:
			result := FenToYuan(v)
			assert.Equal(t, tt.expected, result)
		case uint64:
			result := FenToYuan(v)
			assert.Equal(t, tt.expected, result)
		}
	}
}

func TestYuanToFen(t *testing.T) {
	tests := []struct {
		yuan     string
		expected int64
		hasError bool
	}{
		{"1.00", 100, false},
		{"0.01", 1, false},
		{"0", 0, false},
		{"9.99", 999, false},
		{"123.45", 12345, false},
		{"-1.00", -100, false},
		// 错误场景
		{"abc", 0, true},
		{"1.234", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		result, err := YuanToFen(tt.yuan)
		if tt.hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		}
	}
}

func TestCurrencyRoundTrip(t *testing.T) {
	// 测试分转元再转分的往返一致性
	testCases := []int64{0, 1, 100, 999, 12345, -100, -12345}
	
	for _, original := range testCases {
		// 分转元
		yuanStr := FenToYuan(original)
		
		// 元转分
		backToFen, err := YuanToFen(yuanStr)
		assert.NoError(t, err)
		
		// 验证往返一致性
		assert.Equal(t, original, backToFen)
	}
}