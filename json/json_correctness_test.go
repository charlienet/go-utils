//nolint:govet // 本文件包含故意设计的 struct tag 冲突测试场景
package json

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charlienet/go-utils/stringx"
	"github.com/stretchr/testify/assert"
)

// H1: omitempty 支持
func TestOmitempty_Support(t *testing.T) {
	type Person struct {
		Name   string `json:"name,omitempty"`
		Age    int    `json:"age,omitempty"`
		Email  string `json:"email,omitempty"`
		Hidden string `json:"hidden,omitempty"`
	}

	// 零值字段应被忽略
	p := Person{Name: "John", Age: 0, Email: "", Hidden: ""}
	result, err := Marshal(p)
	assert.NoError(t, err)
	assert.Contains(t, string(result), `"name"`)
	assert.NotContains(t, string(result), `"age"`)
	assert.NotContains(t, string(result), `"email"`)
	assert.NotContains(t, string(result), `"hidden"`)

	// 非零值字段应包含
	p2 := Person{Name: "Jane", Age: 25, Email: "jane@example.com", Hidden: "secret"}
	result2, err := Marshal(p2)
	assert.NoError(t, err)
	assert.Contains(t, string(result2), `"name"`)
	assert.Contains(t, string(result2), `"age"`)
	assert.Contains(t, string(result2), `"email"`)
	assert.Contains(t, string(result2), `"hidden"`)
}

// H2: []byte 的 base64 编码
func TestByteSlice_Base64(t *testing.T) {
	type Data struct {
		Content []byte `json:"content"`
	}

	d := Data{Content: []byte("hello world")}
	result, err := Marshal(d)
	assert.NoError(t, err)
	// base64 编码后的 "hello world"
	assert.Contains(t, string(result), `"aGVsbG8gd29ybGQ="`)
}

func TestByteSlice_Nil(t *testing.T) {
	type Data struct {
		Content []byte `json:"content"`
	}

	d := Data{Content: nil}
	result, err := Marshal(d)
	assert.NoError(t, err)
	assert.Contains(t, string(result), `"content":null`)
}

// H3: json.Marshaler 指针接收者测试
type ptrMarshalType struct{ V int }

func (m *ptrMarshalType) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `{"custom":%d}`, m.V), nil
}

func TestMarshaler_PtrReceiver(t *testing.T) {
	// 场景1：直接传指针 → 应调用 MarshalJSON
	m := &ptrMarshalType{V: 42}
	b, err := marshalJSON(m, stringx.Pascal2Camel)
	assert.NoError(t, err)
	assert.Equal(t, `{"custom":42}`, string(b))

	// 场景2：nil 指针 → 应输出 null
	var nilPtr *ptrMarshalType
	b, err = marshalJSON(nilPtr, stringx.Pascal2Camel)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(b))

	// 场景3：struct 字段（不可寻址）→ 不调用指针接收者
	type Container struct{ F ptrMarshalType }
	c := Container{F: ptrMarshalType{V: 99}}
	b, err = marshalJSON(c, stringx.Pascal2Camel)
	assert.NoError(t, err)
	// 不可寻址 → 按 struct 编码 → 键名被转换
	assert.Contains(t, string(b), `"f"`)
}

// H4: 命名转换后的键冲突检测
func TestKeyCollision_Detection(t *testing.T) {
	// 自定义转换函数：将所有字段名转为小写（会导致冲突）
	toLower := func(s string) string {
		return strings.ToLower(s)
	}

	type Conflict struct {
		NameA string
		Namea string // 转换后都变成 "namea"
	}

	c := Conflict{NameA: "john", Namea: "doe"}

	// 使用自定义转换器
	_, err := marshalJSON(c, toLower)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
}

func TestKeyCollision_Map(t *testing.T) {
	// 自定义转换函数：将所有键名转为小写
	toLower := func(s string) string {
		return strings.ToLower(s)
	}

	m := map[string]any{
		"keyA": "john",
		"keya": "doe", // 转换后都变成 "keya"
	}

	_, err := marshalJSON(m, toLower)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
}

// H5: 循环引用检测
func TestCycleDetection_Pointer(t *testing.T) {
	type Node struct {
		Value int
		Next  *Node
	}

	// 创建循环引用
	n1 := &Node{Value: 1}
	n2 := &Node{Value: 2}
	n1.Next = n2
	n2.Next = n1

	// 使用 Pascal2Camel 包装器来调用我们的 marshalJSON 函数
	_, err := Marshal(Pascal2Camel{n1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errCycle.Error())
}

func TestCycleDetection_Slice(t *testing.T) {
	type Node struct {
		Value int
		Next  []*Node
	}

	// 创建循环引用
	n1 := &Node{Value: 1}
	n2 := &Node{Value: 2}
	n1.Next = []*Node{n2}
	n2.Next = []*Node{n1}

	_, err := Marshal(Pascal2Camel{n1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errCycle.Error())
}

func TestCycleDetection_Map(t *testing.T) {
	type Node struct {
		Value int
		Next  map[string]*Node
	}

	// 创建循环引用
	n1 := &Node{Value: 1}
	n2 := &Node{Value: 2}
	n1.Next = map[string]*Node{"next": n2}
	n2.Next = map[string]*Node{"next": n1}

	_, err := Marshal(Pascal2Camel{n1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errCycle.Error())
}

// H5: 最大递归深度
func TestMaxDepth(t *testing.T) {
	// 创建深度嵌套的结构
	type Node struct {
		Value int
		Next  *Node
	}

	head := &Node{Value: 0}
	current := head
	for i := 1; i < 10050; i++ {
		current.Next = &Node{Value: i}
		current = current.Next
	}

	_, err := Marshal(Pascal2Camel{head})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errCycle.Error())
}

// 测试无循环引用时的正常处理
func TestNoCycle_Normal(t *testing.T) {
	type Node struct {
		Value int
		Next  *Node
	}

	// 创建非循环链表
	n1 := &Node{Value: 1}
	n2 := &Node{Value: 2}
	n3 := &Node{Value: 3}
	n1.Next = n2
	n2.Next = n3

	result, err := Marshal(n1)
	assert.NoError(t, err)
	assert.Contains(t, string(result), `"Value":1`)
	assert.Contains(t, string(result), `"Value":2`)
	assert.Contains(t, string(result), `"Value":3`)
}

// 测试 nil 切片和 map
func TestNilSliceAndMap(t *testing.T) {
	type Data struct {
		Items []string
		Map   map[string]int
	}

	d := Data{Items: nil, Map: nil}
	result, err := Marshal(d)
	assert.NoError(t, err)
	assert.Contains(t, string(result), `"Items":null`)
	assert.Contains(t, string(result), `"Map":null`)
}

// 测试空切片和 map
func TestEmptySliceAndMap(t *testing.T) {
	type Data struct {
		Items []string
		Map   map[string]int
	}

	d := Data{Items: []string{}, Map: map[string]int{}}
	result, err := Marshal(d)
	assert.NoError(t, err)
	assert.Contains(t, string(result), `"Items":[]`)
	assert.Contains(t, string(result), `"Map":{}`)
}

// M1: 嵌入字段冲突解决
func TestM1_TypeFields_ConflictResolution(t *testing.T) {
	// 场景1：同深度冲突（当前实现因 visited map bug，无法正确检测）
	// 标准：应丢弃冲突字段，输出 {}
	// 现状：输出 {"x":1}，因为只收集到一个 X 字段
	// TODO: 修复 json_conv.go typeFields 的 visited 逻辑
	type Inner struct{ X int }
	type A struct{ Inner }
	type B struct{ Inner }
	type Root1 struct {
		A
		B
	}
	r1 := Root1{A: A{Inner{1}}, B: B{Inner{2}}}
	b, err := marshalJSON(r1, stringx.Pascal2Camel)
	assert.NoError(t, err)
	// 修复后：标准库行为，同深度冲突字段被丢弃，输出 {}
	assert.Equal(t, `{}`, string(b))

	// 场景2：浅深度胜出
	type Outer struct {
		Name string
		Inner
	}
	o := Outer{Name: "outer", Inner: Inner{X: 99}}
	b, err = marshalJSON(o, stringx.Pascal2Camel)
	assert.NoError(t, err)
	assert.Contains(t, string(b), `"name"`)
	assert.Contains(t, string(b), `"x"`)

	// 场景3：同深度有 tag 唯一 → 胜出
	type A2 struct{ Inner }
	type B2 struct {
		X int `json:"X"`
	}
	type Root2 struct {
		A2
		B2
	}
	r2 := Root2{A2: A2{Inner{1}}, B2: B2{X: 2}}
	b, err = marshalJSON(r2, stringx.Pascal2Camel)
	assert.NoError(t, err)
	// 有 tag 的字段胜出，保留原始名称 "X"（fromTag=true）
	assert.Contains(t, string(b), `"X"`)

	// 场景4：同深度多个有 tag → 丢弃
	type A3 struct {
		X int `json:"X"`
	}
	type B3 struct {
		Y int `json:"X"` // 字段名不同但 tag 相同，触发冲突
	}
	type Root3 struct { //nolint:govet // 故意使用相同 tag 测试冲突检测
		A3
		B3
	}
	r3 := Root3{A3: A3{1}, B3: B3{2}}
	b, err = marshalJSON(r3, stringx.Pascal2Camel)
	assert.NoError(t, err)
	// 同名 tag 冲突，字段被丢弃，输出空对象
	assert.Equal(t, `{}`, string(b))
}

// M2: map 键排序
func TestM2_MapKeySort(t *testing.T) {
	m := map[string]any{"z": 1, "a": 2, "m": 3, "b": 4}
	b, err := marshalJSON(m, stringx.Pascal2Camel)
	assert.NoError(t, err)
	// 键应按字典序排列
	assert.Equal(t, `{"a":2,"b":4,"m":3,"z":1}`, string(b))

	// 多次调用应输出一致
	for range 10 {
		b2, _ := marshalJSON(m, stringx.Pascal2Camel)
		assert.Equal(t, string(b), string(b2))
	}
}
