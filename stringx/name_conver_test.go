package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamel2Pascal(t *testing.T) {
	expected := []struct {
		actual string
		expect string
	}{
		{"updatedAt", "UpdatedAt"},
		{"name", "Name"},
		{"createdAt", "CreatedAt"},
		{"n", "N"},
		{"", ""},
	}

	for _, n := range expected {
		assert.Equal(t, n.expect, Camel2Pascal(n.actual))
	}
}

func TestPascal2Camel(t *testing.T) {
	expected := []struct {
		actual string
		expect string
	}{
		{"UpdatedAt", "updatedAt"},
		{"Name", "name"},
		{"CreatedAt", "createdAt"},
		{"N", "n"},
		{"", ""},
	}

	for _, n := range expected {
		assert.Equal(t, n.expect, Pascal2Camel(n.actual))
	}
}

func TestSnake2Pascal(t *testing.T) {
	expected := []struct {
		actual string
		expect string
	}{
		{"UPDATED_AT", "UpdatedAt"},
		{"Name", "Name"},
		{"created_at", "CreatedAt"},
		{"N", "N"},
		{"", ""},
	}

	for _, n := range expected {
		assert.Equal(t, n.expect, Snake2Pascal(n.actual))
	}
}

func TestSnake2Camel(t *testing.T) {
	expected := []struct {
		actual string
		expect string
	}{
		{"UPDATED_AT", "updatedAt"},
		{"Name", "name"},
		{"created_at", "createdAt"},
		{"N", "n"},
		{"", ""},
	}

	for _, n := range expected {
		assert.Equal(t, n.expect, Snake2Camel(n.actual))
	}
}

func TestPascal2Snake(t *testing.T) {
	expected := []struct {
		actual string
		expect string
	}{
		{"UpdatedAt", "Updated_At"},
		{"Name", "Name"},
		{"createdAt", "Created_At"},             // camelCase input now properly converted
		{"XMLParser", "XML_Parser"},             //缩写词处理
		{"HTTPSConnection", "HTTPS_Connection"}, //缩写词处理
		// 缩写词+小写词（P2 契约下有损，见 doc.go Known Limitations #3）
		{"HTTPSserver", "HTTP_Sserver"}, //缩写词+小写开头单词
		{"HTTPserver", "HTT_Pserver"},   //缩写词+小写开头单词
		{"XMLparser", "XM_Lparser"},     //缩写词+小写开头单词
		{"N", "N"},
		{"", ""},
	}

	for _, n := range expected {
		assert.Equal(t, n.expect, Pascal2Snake(n.actual))
	}
}

func TestPascal2UpperSnake(t *testing.T) {
	expected := []struct {
		actual string
		expect string
	}{
		{"UpdatedAt", "UPDATED_AT"},
		{"Name", "NAME"},
		{"createdAt", "CREATED_AT"},
		{"N", "N"},
		{"", ""},
	}

	for _, n := range expected {
		assert.Equal(t, n.expect, Pascal2UpperSnake(n.actual))
	}
}

func TestCamel2UpperSnake(t *testing.T) {
	expected := []struct {
		actual string
		expect string
	}{
		{"updatedAt", "UPDATED_AT"},
		{"name", "NAME"},
		{"createdAt", "CREATED_AT"},
		{"n", "N"},
		{"", ""},
	}

	for _, n := range expected {
		assert.Equal(t, n.expect, Camel2UpperSnake(n.actual))
	}
}

func TestCamel2Snake(t *testing.T) {
	expected := []struct {
		actual string
		expect string
	}{
		{"updatedAt", "updated_At"},
		{"name", "name"},
		{"createdAt", "created_At"},
		{"n", "n"},
		{"", ""},
	}

	for _, n := range expected {
		assert.Equal(t, n.expect, Camel2Snake(n.actual))
	}
}

func TestToSnake(t *testing.T) {
	expected := []struct {
		actual string
		expect string
	}{
		{"", ""},
		{"a", "a"},
		{"A", "a"},
		{"name", "name"},
		{"Name", "name"},
		{"UserName", "user_name"},
		{"updatedAt", "updated_at"},
		{"XMLParser", "xml_parser"},
		{"HTTPSConnection", "https_connection"},
		{"parseURL", "parse_url"},
		{"user123ID", "user123_id"},
		{"UpdatedAt123", "updated_at123"},
		{"HTTPSserver", "http_sserver"},
		{"HTTPserver", "htt_pserver"},
	}

	for _, n := range expected {
		assert.Equal(t, n.expect, ToSnake(n.actual))
	}
}

func TestUpper(t *testing.T) {
	expected := []struct {
		actual byte
		expect byte
	}{
		{'A', 'A'},
		{'Z', 'Z'},
		{'a', 'A'},
		{'z', 'Z'},
	}

	for _, n := range expected {
		assert.Equal(t, n.expect, toUpper(n.actual))
	}
}

func BenchmarkTransform(b *testing.B) {
	b.Run("Pascal2Camel", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Pascal2Camel("UpdatedAt")
		}
	})

	b.Run("Camel2Pascal", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Camel2Pascal("updatedAt")
		}
	})

	b.Run("Snake2Pascal", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Snake2Pascal("UPDATED_AT")
		}
	})

	b.Run("Snake2Camel", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Snake2Camel("UPDATED_AT")
		}

	})

	b.Run("Pascal2Snake", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Pascal2Snake("UpdatedAt")
		}
	})

	b.Run("Pascal2UpperSnake", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Pascal2UpperSnake("UpdatedAt")
		}
	})

	b.Run("Camel2Snake", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Camel2Snake("updatedAt")
		}
	})

	b.Run("Camel2UpperSnake", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Camel2UpperSnake("updatedAt")
		}
	})

	b.Run("Pascal2SnakeParallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				Pascal2Snake("UpdatedAt")
			}
		})
	})
}

func TestSplitByCapital2(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"abc", []string{"abc"}},
		{"UpdatedAt", []string{"Updated", "At"}},
		{"UpdatedByDayTuesday", []string{"Updated", "By", "Day", "Tuesday"}},
		{"Updated", []string{"Updated"}},
		{"", []string{""}},
		{"UserName", []string{"User", "Name"}},
		{"userNameIsAdmin", []string{"user", "Name", "Is", "Admin"}},
		// 缩写词 + 大写开头单词（正常场景）
		{"XMLParser", []string{"XML", "Parser"}},
		{"HTTPSConnection", []string{"HTTPS", "Connection"}},
		// 缩写词 + 小写词（P2 契约下有损，见 doc.go Known Limitations #3）
		{"HTTPSserver", []string{"HTTP", "Sserver"}},
		{"HTTPserver", []string{"HTT", "Pserver"}},
		{"XMLparser", []string{"XM", "Lparser"}},
		{"ABcdef", []string{"A", "Bcdef"}},
		// 其他缩写词场景
		{"parseURL", []string{"parse", "URL"}},
		{"getHTTPResponse", []string{"get", "HTTP", "Response"}},
		{"userID", []string{"user", "ID"}},
		{"simpleTest", []string{"simple", "Test"}},
		{"ABC", []string{"ABC"}},
		{"ABCDef", []string{"ABC", "Def"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := splitByCapital(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func BenchmarkSplit(b *testing.B) {
	b.Run("splitByCapital", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			splitByCapital("UpdatedByDayTuesday")
		}
	})
}

func TestEdgeCases(t *testing.T) {
	// 测试连续大写字母（缩写词处理）
	assert.Equal(t, "XML_Parser", Pascal2Snake("XMLParser"))
	assert.Equal(t, "Parse_URL", Pascal2Snake("ParseURL")) // Now correctly handles abbreviations
	assert.Equal(t, "HTTPS_Connection", Pascal2Snake("HTTPSConnection"))

	// 测试纯数字
	assert.Equal(t, "123", Pascal2Snake("123"))
	assert.Equal(t, "User123", Snake2Pascal("user123"))

	// 测试包含数字的混合
	assert.Equal(t, "Updated_At123", Pascal2Snake("UpdatedAt123"))

	// 测试空字符串
	assert.Equal(t, "", Pascal2Snake(""))
	assert.Equal(t, "", Camel2Snake(""))
	assert.Equal(t, "", Snake2Pascal(""))
	assert.Equal(t, "", Pascal2UpperSnake(""))
	assert.Equal(t, "", Camel2UpperSnake(""))
	assert.Equal(t, "", Snake2Camel(""))

	// 测试单字符
	assert.Equal(t, "A", Pascal2Snake("A"))
	assert.Equal(t, "A", Pascal2Snake("a")) // Now correctly converts single lowercase to Pascal then snake

	// 测试全小写（驼峰形式）
	assert.Equal(t, "Username", Snake2Pascal("username"))
	assert.Equal(t, "Username", Pascal2Snake("Username"))

	// 测试全大写
	assert.Equal(t, "USERNAME", Pascal2Snake("USERNAME")) // All caps remain as is

	// 测试Camel2Snake特殊情况
	assert.Equal(t, "xml_Parser", Camel2Snake("xmlParser")) // First word lowercase, then PascalCase
	assert.Equal(t, "updated_At", Camel2Snake("updatedAt")) // Main test case - now properly handled
	assert.Equal(t, "a", Camel2Snake("a"))                  // Single lowercase letter
	assert.Equal(t, "a", Camel2Snake("A"))                  // Single uppercase becomes Pascal then snake, then first char lowercased
	// F5: 大小写归一仅作用于 ASCII；多字节 UTF-8 序列字节完整、ASCII 词首归一化不变
	assert.Equal(t, "中文field", Snake2Camel("中文Field"))
}

func BenchmarkUcfirst(b *testing.B) {
	b.Run("Ucfirst", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Ucfirst("UCFIRST")
		}
	})
}

func TestUcfirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"HELLO", "Hello"},
		{"", ""},
		{"a", "A"},
		{"A", "A"},
		{"helloWorld", "Helloworld"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, Ucfirst(tt.input))
		})
	}
}

func TestLcfirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "hello"},
		{"hello", "hello"},
		{"HELLO", "hELLO"},
		{"", ""},
		{"A", "a"},
		{"a", "a"},
		{"HelloWorld", "helloWorld"}, // 只转首字母
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, Lcfirst(tt.input))
		})
	}
}
