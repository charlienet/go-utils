package mathx

import (
    "fmt"
    "github.com/shopspring/decimal"
)

// integer 是整数类型约束
type integer interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
        ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// FenToYuan 分转元，支持任意整数类型输入，返回固定2位小数的字符串。
// 例：FenToYuan(123) → "1.23"，FenToYuan(int64(1)) → "0.01"
func FenToYuan[T integer](fen T) string {
    f := int64(fen)
    sign := ""
    if f < 0 {
        sign = "-"
        f = -f
    }
    return fmt.Sprintf("%s%d.%02d", sign, f/100, f%100)
}

// YuanToFen 元转分，接受字符串输入（如 "1.23"），返回 int64 分和 error。
// 使用 decimal 库保证精度，避免 float64 的精度丢失问题。
// 例：YuanToFen("1.23") → (123, nil)
func YuanToFen(yuan string) (int64, error) {
    d, err := decimal.NewFromString(yuan)
    if err != nil {
        return 0, fmt.Errorf("invalid yuan string %q: %w", yuan, err)
    }
    // 乘以 100 转分
    cents := d.Mul(decimal.NewFromInt(100))
    if !cents.IsInteger() {
        return 0, fmt.Errorf("yuan string %q has more than 2 decimal places", yuan)
    }
    return cents.IntPart(), nil
}