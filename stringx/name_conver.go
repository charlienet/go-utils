package stringx

import "strings"

// CamelCase 		userName
// PascalCase		UserName
// SnakeCase		user_name
// PascalSnakeCase	User_Name
// UpperSnakeCase	USER_NAME

// Pascal转换为驼峰
func Pascal2Camel(name string) string {
	if len(name) == 0 {
		return name
	}

	data := []byte(name)
	data[0] = toLower(data[0])

	return string(data)
}

// 驼峰转换为Pascal
func Camel2Pascal(name string) string {
	if len(name) == 0 {
		return name
	}

	data := []byte(name)
	data[0] = toUpper(data[0])

	return string(data)
}

func Pascal2Snake(name string) string {
	if len(name) == 0 {
		return name
	}

	// 如果是 camelCase，先转 PascalCase
	if name[0] >= 'a' && name[0] <= 'z' {
		name = Camel2Pascal(name)
	}

	return strings.Join(splitByCapital(name), "_")
}

func Pascal2UpperSnake(name string) string {
	if len(name) == 0 {
		return name
	}

	// 如果是 camelCase，先转 PascalCase
	if name[0] >= 'a' && name[0] <= 'z' {
		name = Camel2Pascal(name)
	}

	joined := []byte(strings.Join(splitByCapital(name), "_"))
	for i := range joined {
		joined[i] = toUpper(joined[i])
	}

	return string(joined)
}

func Camel2UpperSnake(name string) string {
	if len(name) == 0 {
		return name
	}

	// 将驼峰转为帕斯卡，然后转换为大写蛇形
	pascalCase := Camel2Pascal(name)
	return Pascal2UpperSnake(pascalCase)
}

func Camel2Snake(name string) string {
	if len(name) == 0 {
		return name
	}

	// 先转 PascalCase
	pascalName := Camel2Pascal(name)
	// 再转 snake_case
	snakeResult := Pascal2Snake(pascalName)

	// 将首字母转为小写，其余字符保持不变
	data := []byte(snakeResult)
	data[0] = toLower(data[0])
	return string(data)
}

// ToSnake 将 PascalCase/camelCase 转换为全小写的 snake_case（如 "UserName" -> "user_name"）。
// 与 Pascal2Snake 遵循相同的切分契约（含缩写词有损行为，见 doc.go Known Limitations #3）。
// 输入若已是小写 snake_case 则幂等返回。
func ToSnake(name string) string {
	if len(name) == 0 {
		return name
	}

	// 如果是 camelCase，先转 PascalCase
	if name[0] >= 'a' && name[0] <= 'z' {
		name = Camel2Pascal(name)
	}

	joined := []byte(strings.Join(splitByCapital(name), "_"))
	for i := range joined {
		joined[i] = toLower(joined[i])
	}

	return string(joined)
}

func Snake2Pascal(name string) string {
	if len(name) == 0 {
		return name
	}

	names := strings.Split(name, "_")
	for i, n := range names {
		names[i] = Ucfirst(n)
	}

	return strings.Join(names, "")
}

func Snake2Camel(name string) string {
	if len(name) == 0 {
		return name
	}

	names := strings.Split(name, "_")
	// 对于camelCase，第一个单词应全部小写
	names[0] = lowerASCII(names[0])
	for i := 1; i < len(names); i++ {
		names[i] = Ucfirst(names[i])
	}

	return strings.Join(names, "")
}

// Ucfirst 首字母大写，其余字母全部小写（破坏性：会丢失其余字符的大小写信息）。
// 如仅需首字母变换请使用 Lcfirst。Snake2Pascal 依赖本函数的归一化行为。
func Ucfirst(str string) string {
	if len(str) == 0 {
		return str
	}

	data := []byte(str)
	data[0] = toUpper(data[0])
	for i := 1; i < len(data); i++ {
		data[i] = toLower(data[i])
	}

	return string(data)
}

// Lcfirst 首字母小写，其余保持不变
func Lcfirst(str string) string {
	if len(str) == 0 {
		return str
	}
	data := []byte(str)
	data[0] = toLower(data[0])
	return string(data)
}

func toUpper(c byte) byte {
	if 'a' <= c && c <= 'z' {
		c -= 'a' - 'A'
	}

	return c
}

func toLower(s byte) byte {
	if 'A' <= s && s <= 'Z' {
		s += 'a' - 'A'
	}
	return s
}

// lowerASCII 仅将 ASCII 大写字母转为小写，非 ASCII 字节保持不变。
// 无大写字母时零分配直接返回原串。
func lowerASCII(s string) string {
	for i := 0; i < len(s); i++ {
		if 'A' <= s[i] && s[i] <= 'Z' {
			d := []byte(s)
			for j := i; j < len(d); j++ {
				d[j] = toLower(d[j])
			}
			return string(d)
		}
	}
	return s
}

// splitByCapital 按 PascalCase/camelCase 词边界切分，遵循社区标准：
//
//	规则1：小写或数字之后的大写字母开新词（userName、user123ID、parseURL）
//	规则2：大写之后的首个小写字母，使前一个大写字母归本词
//	       （连续大写段的末位大写归后词：XMLParser -> XML|Parser；
//	        HTTPSserver -> HTTP|Sserver 为声明的有损行为，见 doc.go）
//	数字不单独成词，跟随当前词（UpdatedAt123 -> Updated|At123）
func splitByCapital(s string) []string {
	if len(s) == 0 {
		return []string{s}
	}

	a := make([]string, 0, 8)
	last := 0

	for i := 1; i < len(s); i++ {
		cur, prev := s[i], s[i-1]

		if isUpper(cur) && !isUpper(prev) {
			// 规则1：切点 i
			a = append(a, s[last:i])
			last = i
		} else if isLower(cur) && isUpper(prev) && i-1 > last {
			// 规则2：切点 i-1；i-1 == last 说明该大写已是词首（ParseURL 的 P），跳过防空段
			a = append(a, s[last:i-1])
			last = i - 1
		}
		// 大写跟大写（后续统一由规则2回切）、数字、小写跟小写：不切
	}

	if last < len(s) {
		a = append(a, s[last:])
	}

	return a
}

func isUpper(c byte) bool { return c >= 'A' && c <= 'Z' }
func isLower(c byte) bool { return c >= 'a' && c <= 'z' }
