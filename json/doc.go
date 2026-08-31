/*
Package json 提供增强的 JSON 序列化功能，支持标准库和 jsoniter 库的切换以及命名风格转换。

主要功能：
- 提供标准 encoding/json 包的兼容接口
- 支持通过 build tag 在标准库和 jsoniter 之间切换
- 提供多种命名风格转换功能（蛇形转驼峰、帕斯卡转蛇形等）
- 支持结构体字段命名风格的自动转换

导出函数和类型：
- Marshal/Unmarshal/MarshalIndent/NewDecoder/NewEncoder: 标准 JSON 操作函数
- RegisterFuzzyDecoders(): 注册模糊解码器（仅在 jsoniter 模式下生效）
- Snake2Camel/Pascal2Snake/Camel2Pascal 等类型: 命名风格转换包装器
  - Snake2Camel: 蛇形转驼峰
  - Snake2Pascal: 蛇形转帕斯卡
  - Pascal2Camel: 帕斯卡转驼峰
  - Camel2Pascal: 驼峰转帕斯卡
  - Pascal2Snake: 帕斯卡转蛇形
  - Pascal2UpperSnake: 帕斯卡转大写蛇形

构建标签：
- 默认使用标准 encoding/json 库
- 使用 `jsoniter` 标签编译时切换到 github.com/json-iterator/go 库

使用示例：

	// 示例1: 基本 JSON 序列化/反序列化
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	p := Person{Name: "John", Age: 30}
	data, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	var p2 Person
	err = json.Unmarshal(data, &p2)

	// 示例2: 使用命名风格转换（帕斯卡转驼峰）
	type User struct {
		UserName string
		UserAge  int
		Email    string
	}

	user := User{
		UserName: "john_doe",
		UserAge:  30,
		Email:    "john@example.com",
	}

	// 自动将字段名从帕斯卡转为驼峰
	jsonData := json.MustStruct2Json(json.Pascal2Camel{user})
	// 输出: {"userName":"john_doe","userAge":30,"email":"john@example.com"}

	// 示例3: 蛇形转驼峰转换
	m := map[string]interface{}{
		"user_name": "jane",
		"email_addr": "jane@example.com",
	}

	jsonData = json.MustStruct2Json(json.Snake2Camel{m})
	// 输出: {"userName":"jane","emailAddr":"jane@example.com"}

	// 示例4: 使用 build tag 切换底层实现
	// go build -tags=jsoniter ./... # 使用 jsoniter
	// go build ./...               # 使用标准库（默认）

注意事项：
- 循环引用检测：会检测并拒绝循环引用的序列化（返回 errCycle）
- 最大递归深度：限制为 10000 层，防止栈溢出
- 命名风格转换后键冲突：转换后会检测键冲突，冲突时返回错误
- omitempty 支持：支持 omitempty tag，零值字段会被忽略
- []byte 特殊处理：[]byte 类型会被编码为 base64 字符串
*/
package json
