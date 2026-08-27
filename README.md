# go-utils

Go 通用工具库，提供日常开发中常用的实用工具包。

## 安装

```bash
go get github.com/charlienet/go-utils
```

## 包列表

| 包 | 说明 |
|---|------|
| [bytesconv](./bytesconv) | 零拷贝 string/[]byte 转换 |
| [compiledbuffer](./compiledbuffer) | 泛型编译结果缓存 |
| [currency](./currency) | 货币转换（元/分） |
| [expr](./expr) | 控制流语法糖（三元/If-Else/Switch） |
| [fs](./fs) | 文件系统工具 |
| [json](./json) | JSON 序列化（支持 jsoniter/命名风格转换） |
| [locker](./locker) | 多种锁实现 |
| [mathx](./mathx) | 数学扩展（舍入/费率计算） |
| [pool](./pool) | 泛型对象池 |
| [random](./random) | 随机数生成器 |
| [stringx](./stringx) | 命名风格转换（驼峰/蛇形/Pascal） |

## 使用示例

### bytesconv - 零拷贝转换

```go
import "github.com/charlienet/go-utils/bytesconv"

// 零拷贝转换（注意：返回的切片不应被修改）
b := bytesconv.StringToBytes("hello")
s := bytesconv.BytesToString(b)
```

### stringx - 命名风格转换

```go
import "github.com/charlienet/go-utils/stringx"

// 驼峰转蛇形
snake := stringx.CamelToSnake("userName")  // "user_name"

// 蛇形转驼峰
camel := stringx.SnakeToCamel("user_name")  // "UserName"
```

### currency - 货币转换

```go
import "github.com/charlienet/go-utils/currency"

// 元转分
fen := currency.YuanToFen("99.99")  // 9999

// 分转元
yuan := currency.FenToYuan(9999)    // "99.99"
```

### pool - 泛型对象池

```go
import "github.com/charlienet/go-utils/pool"

// 创建对象池
p := pool.New(10, func() *bytes.Buffer {
    return new(bytes.Buffer)
})

// 获取对象
buf := p.Get()
defer p.Put(buf)
```

## 依赖

- Go 1.21+
- github.com/shopspring/decimal（currency）
- github.com/json-iterator/go（json，可选）

## License

MIT License
