package json

import (
	"testing"
)

// 基准测试：结构体序列化
func BenchmarkMarshal_Struct(b *testing.B) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	p := Person{Name: "John", Age: 30}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Marshal(p)
	}
}

// 基准测试：嵌套结构体序列化
func BenchmarkMarshal_NestedStruct(b *testing.B) {
	type Address struct {
		City    string `json:"city"`
		Country string `json:"country"`
	}
	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}

	p := Person{
		Name: "John",
		Age:  30,
		Address: Address{
			City:    "Beijing",
			Country: "China",
		},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Marshal(p)
	}
}

// 基准测试：切片序列化
func BenchmarkMarshal_Slice(b *testing.B) {
	items := make([]string, 100)
	for i := range 100 {
		items[i] = "item"
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Marshal(items)
	}
}

// 基准测试：map 序列化
func BenchmarkMarshal_Map(b *testing.B) {
	m := make(map[string]int, 100)
	for i := range 100 {
		m["key"] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Marshal(m)
	}
}

// 基准测试：命名风格转换
func BenchmarkMarshal_NameConversion(b *testing.B) {
	type Person struct {
		UserName string
		UserAge  int
		Email    string
	}

	p := Person{
		UserName: "john_doe",
		UserAge:  30,
		Email:    "john@example.com",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Marshal(Pascal2Camel{p})
	}
}

// 基准测试：键缓存效果
func BenchmarkMarshal_KeyCache(b *testing.B) {
	type Person struct {
		UserName string
		UserAge  int
		Email    string
	}

	p := Person{
		UserName: "john_doe",
		UserAge:  30,
		Email:    "john@example.com",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Marshal(Pascal2Camel{p})
	}
}

// 基准测试：指针序列化
func BenchmarkMarshal_Pointer(b *testing.B) {
	type Person struct {
		Name *string
		Age  *int
	}

	name := "John"
	age := 30
	p := Person{Name: &name, Age: &age}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Marshal(p)
	}
}

// 基准测试：[]byte 序列化（base64）
func BenchmarkMarshal_ByteSlice(b *testing.B) {
	type Data struct {
		Content []byte `json:"content"`
	}

	data := Data{Content: make([]byte, 1024)}
	for i := range data.Content {
		data.Content[i] = byte(i % 256)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Marshal(data)
	}
}

// 基准测试：omitempty 字段
func BenchmarkMarshal_Omitempty(b *testing.B) {
	type Person struct {
		Name   string `json:"name,omitempty"`
		Age    int    `json:"age,omitempty"`
		Email  string `json:"email,omitempty"`
		Hidden string `json:"hidden,omitempty"`
	}

	p := Person{Name: "John", Age: 0, Email: "", Hidden: ""}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Marshal(p)
	}
}

// 基准测试：大结构体序列化
func BenchmarkMarshal_LargeStruct(b *testing.B) {
	type Item struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	type LargeStruct struct {
		Items []Item `json:"items"`
	}

	large := LargeStruct{
		Items: make([]Item, 1000),
	}
	for i := range 1000 {
		large.Items[i] = Item{
			ID:    i,
			Name:  "item",
			Value: "value",
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Marshal(large)
	}
}
