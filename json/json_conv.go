package json

import (
	"bytes"
	"encoding/base64"
	ejson "encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/charlienet/go-utils/stringx"
)

var (
	fieldCache sync.Map // map[reflect.Type][]field
)

// json.Marshaler 接口的反射类型
var jsonMarshalerType = reflect.TypeFor[ejson.Marshaler]()

const maxDepth = 10000 // 最大递归深度

type encodeState struct {
	convert  func(string) string
	keyCache map[string][]byte    // M3: origName → JSON 转义后的键（含引号）
	visited  map[uintptr]struct{} // H5: 循环引用检测
	depth    uint                 // H5: 递归深度
}

func (s *encodeState) fieldKeyJSON(f *field) []byte {
	if f.fromTag {
		cacheKey := "\x00" + f.name
		if cached, ok := s.keyCache[cacheKey]; ok {
			return cached
		}
		b, _ := Marshal(f.name)
		s.keyCache[cacheKey] = b
		return b
	}
	if cached, ok := s.keyCache[f.origName]; ok {
		return cached
	}
	b, _ := Marshal(s.convert(f.origName))
	s.keyCache[f.origName] = b
	return b
}

func (s *encodeState) mapKeyJSON(k string) []byte {
	if cached, ok := s.keyCache[k]; ok {
		return cached
	}
	b, _ := Marshal(s.convert(k))
	s.keyCache[k] = b
	return b
}

var errCycle = errors.New("json: unsupported value: cycle detected")

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	}
	return false
}

func marshalerOf(val reflect.Value) (ejson.Marshaler, bool) {
	if val.Type().Implements(jsonMarshalerType) {
		if val.Kind() == reflect.Pointer && val.IsNil() {
			return nil, false
		}
		return val.Interface().(ejson.Marshaler), true
	}
	if val.CanAddr() && reflect.PointerTo(val.Type()).Implements(jsonMarshalerType) {
		return val.Addr().Interface().(ejson.Marshaler), true
	}
	return nil, false
}

// field 表示结构体的一个导出字段（经冲突解决后的最终字段）
type field struct {
	name      string // 输出的键名（若来自 tag 则为 tag 名，否则待转换）
	origName  string // 原始字段名（用于无 tag 时转换）
	index     []int  // 字段索引路径
	typ       reflect.Type
	omitempty bool
	fromTag   bool // true 表示 name 直接来自 tag，不再转换
}

// ====================== 命名风格转换器 ======================

type Snake2Camel struct{ Value any }

func (c Snake2Camel) MarshalJSON() ([]byte, error) {
	return marshalJSON(c.Value, stringx.Snake2Camel)
}

type Snake2Pascal struct{ Value any }

func (c Snake2Pascal) MarshalJSON() ([]byte, error) {
	return marshalJSON(c.Value, stringx.Snake2Pascal)
}

type Pascal2Camel struct{ Value any }

func (c Pascal2Camel) MarshalJSON() ([]byte, error) {
	return marshalJSON(c.Value, stringx.Pascal2Camel)
}

type Camel2Pascal struct{ Value any }

func (c Camel2Pascal) MarshalJSON() ([]byte, error) {
	return marshalJSON(c.Value, stringx.Camel2Pascal)
}

type Pascal2Snake struct{ Value any }

func (c Pascal2Snake) MarshalJSON() ([]byte, error) {
	return marshalJSON(c.Value, stringx.Pascal2Snake)
}

type Pascal2UpperSnake struct{ Value any }

func (c Pascal2UpperSnake) MarshalJSON() ([]byte, error) {
	return marshalJSON(c.Value, stringx.Pascal2UpperSnake)
}

// ====================== 核心递归编码函数 ======================

// marshalJSON 是公共入口，调用 encodeValue 递归处理值
func marshalJSON(v any, convert func(string) string) ([]byte, error) {
	state := &encodeState{
		convert:  convert,
		keyCache: make(map[string][]byte),
		visited:  make(map[uintptr]struct{}),
		depth:    0,
	}
	return encodeValue(reflect.ValueOf(v), state)
}

// encodeValue 递归处理任意值，返回其 JSON 片段
func encodeValue(val reflect.Value, state *encodeState) ([]byte, error) {
	if state.depth > maxDepth {
		return nil, errCycle
	}
	if !val.IsValid() {
		return []byte("null"), nil
	}

	// H3: 解引用前检查 Marshaler
	if m, ok := marshalerOf(val); ok {
		return m.MarshalJSON()
	}

	switch val.Kind() {
	case reflect.Pointer:
		if val.IsNil() {
			return []byte("null"), nil
		}
		ptr := val.Pointer()
		if _, seen := state.visited[ptr]; seen {
			return nil, errCycle
		}
		state.visited[ptr] = struct{}{}
		state.depth++
		b, err := encodeValue(val.Elem(), state)
		state.depth--
		delete(state.visited, ptr)
		return b, err

	case reflect.Interface:
		if val.IsNil() {
			return []byte("null"), nil
		}
		return encodeValue(val.Elem(), state)

	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			if val.IsNil() {
				return []byte("null"), nil
			}
			ptr := val.Pointer()
			if _, seen := state.visited[ptr]; seen {
				return nil, errCycle
			}
			state.visited[ptr] = struct{}{}
			encoded := base64.StdEncoding.EncodeToString(val.Bytes())
			b, err := Marshal(encoded)
			delete(state.visited, ptr)
			return b, err
		}
		if val.IsNil() {
			return []byte("null"), nil
		}
		ptr := val.Pointer()
		if _, seen := state.visited[ptr]; seen {
			return nil, errCycle
		}
		state.visited[ptr] = struct{}{}
		state.depth++
		b, err := encodeSlice(val, state)
		state.depth--
		delete(state.visited, ptr)
		return b, err

	case reflect.Map:
		if val.IsNil() {
			return []byte("null"), nil
		}
		ptr := val.Pointer()
		if _, seen := state.visited[ptr]; seen {
			return nil, errCycle
		}
		state.visited[ptr] = struct{}{}
		state.depth++
		b, err := encodeMap(val, state)
		state.depth--
		delete(state.visited, ptr)
		return b, err

	case reflect.Struct:
		state.depth++
		b, err := encodeStruct(val, state)
		state.depth--
		return b, err

	case reflect.Array:
		state.depth++
		b, err := encodeSlice(val, state)
		state.depth--
		return b, err

	default:
		return Marshal(val.Interface())
	}
}

// encodeStruct 编码结构体，键名优先使用 tag，否则由 convert 转换
func encodeStruct(val reflect.Value, state *encodeState) ([]byte, error) {
	typ := val.Type()
	fields := cachedTypeFields(typ)

	buf := &bytes.Buffer{}
	buf.WriteByte('{')
	first := true
	seen := make(map[string]bool, len(fields))

	for _, f := range fields {
		fieldVal := val.FieldByIndex(f.index)
		if !fieldVal.IsValid() {
			continue
		}

		// H1: omitempty
		if f.omitempty && isEmptyValue(fieldVal) {
			continue
		}

		var finalKey string
		if f.fromTag {
			finalKey = f.name
		} else {
			finalKey = state.convert(f.origName)
		}

		// H4: 碰撞检测
		if seen[finalKey] {
			return nil, fmt.Errorf("json: duplicate key %q after name conversion", finalKey)
		}
		seen[finalKey] = true

		elemBytes, err := encodeValue(fieldVal, state)
		if err != nil {
			return nil, err
		}

		if !first {
			buf.WriteByte(',')
		}
		first = false

		keyBytes := state.fieldKeyJSON(&f)
		buf.Write(keyBytes)
		buf.WriteByte(':')
		buf.Write(elemBytes)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// encodeMap 编码 map，键名经过 convert 转换（不处理 tag），值递归处理
func encodeMap(val reflect.Value, state *encodeState) ([]byte, error) {
	// 非字符串键的 map 在 JSON 中不支持
	if val.Type().Key().Kind() != reflect.String {
		return nil, &ejson.UnsupportedTypeError{Type: val.Type()}
	}

	keys := make([]string, 0, val.Len())
	iter := val.MapRange()
	for iter.Next() {
		keys = append(keys, iter.Key().String())
	}
	sort.Strings(keys)

	buf := &bytes.Buffer{}
	buf.WriteByte('{')
	seen := make(map[string]bool, len(keys))
	first := true

	for _, k := range keys {
		v := val.MapIndex(reflect.ValueOf(k))
		convertedKey := state.convert(k)

		if seen[convertedKey] {
			return nil, fmt.Errorf("json: duplicate key %q after name conversion in map", convertedKey)
		}
		seen[convertedKey] = true

		valueBytes, err := encodeValue(v, state)
		if err != nil {
			return nil, err
		}

		if !first {
			buf.WriteByte(',')
		}
		first = false

		keyBytes := state.mapKeyJSON(k)
		buf.Write(keyBytes)
		buf.WriteByte(':')
		buf.Write(valueBytes)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// encodeSlice 编码切片或数组
func encodeSlice(val reflect.Value, state *encodeState) ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteByte('[')
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		elemBytes, err := encodeValue(elem, state)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.Write(elemBytes)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

// ====================== 结构体字段缓存与解析 ======================

// cachedTypeFields 获取类型的字段列表（带缓存）
func cachedTypeFields(typ reflect.Type) []field {
	if cached, ok := fieldCache.Load(typ); ok {
		return cached.([]field)
	}
	fields := typeFields(typ)
	fieldCache.Store(typ, fields)
	return fields
}

// typeFields 解析结构体类型的所有导出字段，严格遵循标准库冲突处理规则
func typeFields(typ reflect.Type) []field {
	var fields []field

	// 递归收集所有字段（包括匿名结构体提升）
	// 移除 visited 检查，允许同一嵌入类型通过不同路径被多次访问
	// 这样在冲突解决阶段可以正确检测同深度冲突并丢弃字段
	var visit func(t reflect.Type, baseIndex []int)
	visit = func(t reflect.Type, baseIndex []int) {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			// 忽略非导出字段
			if f.PkgPath != "" && !f.Anonymous {
				continue
			}

			// 当前字段的索引路径
			index := make([]int, len(baseIndex)+1)
			copy(index, baseIndex)
			index[len(baseIndex)] = i

			// 处理匿名字段
			if f.Anonymous {
				// 检查匿名字段自身的 json tag
				tag := f.Tag.Get("json")
				if tag != "" && tag != "-" {
					// 有 tag，将其作为普通字段处理（不展开）
					name, omitempty := parseTag(tag)
					if name == "" {
						name = f.Name // 实际上标准库中，若 tag 存在但为空，会使用字段名，但这里 name 已由 parseTag 处理
					}
					fields = append(fields, field{
						name:      name,
						origName:  f.Name,
						index:     index,
						typ:       f.Type,
						omitempty: omitempty,
						fromTag:   true, // 因为 tag 非空
					})
					continue
				}
				// 无 tag，展开匿名字段
				ft := f.Type
				if ft.Kind() == reflect.Pointer {
					ft = ft.Elem()
				}
				if ft.Kind() == reflect.Struct {
					visit(ft, index)
				} else {
					// 非 struct 匿名内嵌：将其作为普通字段处理
					fields = append(fields, field{
						name:      f.Name,
						origName:  f.Name,
						index:     index,
						typ:       f.Type,
						omitempty: false,
						fromTag:   false,
					})
				}
				continue
			}

			// 非匿名字段，解析 json tag
			tag := f.Tag.Get("json")
			if tag == "-" {
				continue
			}
			name, omitempty := parseTag(tag)
			fromTag := false
			if name == "" {
				name = f.Name
			} else {
				fromTag = true
			}
			fields = append(fields, field{
				name:      name,
				origName:  f.Name,
				index:     index,
				typ:       f.Type,
				omitempty: omitempty,
				fromTag:   fromTag,
			})
		}
	}

	visit(typ, nil)

	// 按名称分组，解决冲突
	byName := make(map[string][]field)
	for _, f := range fields {
		byName[f.name] = append(byName[f.name], f)
	}

	// 为每个名称选择胜出的字段（标准库规则）
	result := make([]field, 0, len(byName))
	for _, cands := range byName {
		if len(cands) == 1 {
			result = append(result, cands[0])
			continue
		}

		// 找最小深度
		minDepth := len(cands[0].index)
		for _, c := range cands[1:] {
			if len(c.index) < minDepth {
				minDepth = len(c.index)
			}
		}

		// 最小深度的候选
		minDepthCands := make([]field, 0)
		for _, c := range cands {
			if len(c.index) == minDepth {
				minDepthCands = append(minDepthCands, c)
			}
		}

		// 如果最小深度只有一个，选中
		if len(minDepthCands) == 1 {
			result = append(result, minDepthCands[0])
			continue
		}

		// 最小深度有多个，统计有 tag 的
		taggedCands := make([]field, 0)
		for _, c := range minDepthCands {
			if c.fromTag {
				taggedCands = append(taggedCands, c)
			}
		}

		// 如果有 tag 的唯一，选中
		if len(taggedCands) == 1 {
			result = append(result, taggedCands[0])
			continue
		}

		// 否则丢弃（不添加到 result）
	}

	// 按深度和索引排序，使输出顺序与标准库一致
	sort.Slice(result, func(i, j int) bool {
		if len(result[i].index) != len(result[j].index) {
			return len(result[i].index) < len(result[j].index)
		}
		return lessIndex(result[i].index, result[j].index)
	})

	return result
}

// parseTag 解析 json tag，返回名称和 omitempty 标志
func parseTag(tag string) (name string, omitempty bool) {
	if tag == "" {
		return "", false
	}
	parts := strings.Split(tag, ",")
	name = parts[0]
	for _, opt := range parts[1:] {
		if opt == "omitempty" {
			omitempty = true
		}
	}
	return
}

// lessIndex 比较两个索引切片，按字典序
func lessIndex(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}
