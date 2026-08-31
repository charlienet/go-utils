package random

import (
	mrndv2 "math/rand/v2"
)

// RandGenerator 随机数生成器接口
type RandGenerator interface {
	Int() int
	Intn(int) int
	Int31() int32
	Int31n(int32) int32
	Int64() int64
	Int64n(int64) int64
}

// rangeGenerator 范围生成器包装
type rangeGenerator struct {
	g RandGenerator
}

func NewRangeGenerator(g RandGenerator) *rangeGenerator {
	return &rangeGenerator{g: g}
}

func (g rangeGenerator) Intr(min, max int) int {
	return min + g.g.Intn(max-min)
}

func (g rangeGenerator) Int31r(min, max int32) int32 {
	return min + g.g.Int31n(max-min)
}

func (g rangeGenerator) Int64r(min, max int64) int64 {
	return min + g.g.Int64n(max-min)
}

// mathRandV2Generator 基于 math/rand/v2 的伪随机数生成器
// 使用 math/rand/v2 顶层函数，内部已优化并发安全
type mathRandV2Generator struct{}

// NewRandGenerator 创建伪随机数生成器（基于 math/rand/v2）
func NewRandGenerator() *mathRandV2Generator {
	return &mathRandV2Generator{}
}

func (r *mathRandV2Generator) Int() int             { return mrndv2.Int() }
func (r *mathRandV2Generator) Intn(n int) int       { return mrndv2.IntN(n) }
func (r *mathRandV2Generator) Int31() int32         { return mrndv2.Int32() }
func (r *mathRandV2Generator) Int31n(n int32) int32 { return mrndv2.Int32N(n) }
func (r *mathRandV2Generator) Int64() int64         { return mrndv2.Int64() }
func (r *mathRandV2Generator) Int64n(n int64) int64 { return mrndv2.Int64N(n) }
