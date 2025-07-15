# xyJson API 参考文档 / API Reference

本文档提供了 xyJson 库的完整 API 参考，包括所有接口、函数和类型的详细说明。

This document provides a complete API reference for the xyJson library, including detailed descriptions of all interfaces, functions, and types.

## 目录 / Table of Contents

1. [核心接口](#核心接口--core-interfaces)
2. [全局函数](#全局函数--global-functions)
3. [工厂函数](#工厂函数--factory-functions)
4. [类型定义](#类型定义--type-definitions)
5. [配置选项](#配置选项--configuration-options)
6. [错误处理](#错误处理--error-handling)
7. [性能监控](#性能监控--performance-monitoring)

## 核心接口 / Core Interfaces

### IValue

所有JSON值的基础接口。

The fundamental interface for all JSON values.

```go
type IValue interface {
    Type() ValueType
    Raw() interface{}
    String() string
    IsNull() bool
    Clone() IValue
    Equals(other IValue) bool
}
```

#### 方法说明 / Method Descriptions

- **Type()** - 返回值的类型（null、string、number、boolean、object、array）
- **Raw()** - 返回对应的Go原生类型值，用于与标准库互操作
- **String()** - 返回JSON格式的字符串表示
- **IsNull()** - 检查是否为null值
- **Clone()** - 创建深拷贝
- **Equals(other IValue)** - 比较两个值是否相等

### IScalarValue

标量值接口，继承自IValue，用于字符串、数字、布尔值。

Scalar value interface, inherits from IValue, for strings, numbers, and booleans.

```go
type IScalarValue interface {
    IValue
    Int() (int, error)
    Int64() (int64, error)
    Float64() (float64, error)
    Bool() (bool, error)
    Time() (time.Time, error)
    Bytes() ([]byte, error)
}
```

#### 方法说明 / Method Descriptions

- **Int()** - 转换为int类型
- **Int64()** - 转换为int64类型
- **Float64()** - 转换为float64类型
- **Bool()** - 转换为bool类型
- **Time()** - 转换为time.Time类型（支持ISO 8601格式）
- **Bytes()** - 转换为字节数组

### IObject

JSON对象接口，提供键值对操作。

JSON object interface, providing key-value pair operations.

```go
type IObject interface {
    IValue
    Get(key string) IValue
    Set(key string, value interface{}) error
    Delete(key string) bool
    Has(key string) bool
    Keys() []string
    Size() int
    Clear()
    Range(fn func(key string, value IValue) bool)
}
```

#### 方法说明 / Method Descriptions

- **Get(key string)** - 根据键名获取值
- **Set(key string, value interface{})** - 设置键值对
- **Delete(key string)** - 删除指定键，返回是否存在
- **Has(key string)** - 检查是否包含指定键
- **Keys()** - 返回所有键名
- **Size()** - 返回键值对数量
- **Clear()** - 清空所有键值对
- **Range(fn func(key string, value IValue) bool)** - 遍历所有键值对

### IArray

JSON数组接口，提供索引操作。

JSON array interface, providing index-based operations.

```go
type IArray interface {
    IValue
    Get(index int) IValue
    Set(index int, value interface{}) error
    Append(value interface{}) error
    Insert(index int, value interface{}) error
    Delete(index int) error
    Length() int
    Clear()
    Range(fn func(index int, value IValue) bool)
}
```

#### 方法说明 / Method Descriptions

- **Get(index int)** - 根据索引获取值
- **Set(index int, value interface{})** - 设置指定索引的值
- **Append(value interface{})** - 追加值到数组末尾
- **Insert(index int, value interface{})** - 在指定位置插入值
- **Delete(index int)** - 删除指定索引的值
- **Length()** - 返回数组长度
- **Clear()** - 清空数组
- **Range(fn func(index int, value IValue) bool)** - 遍历数组元素

### IParser

JSON解析器接口。

JSON parser interface.

```go
type IParser interface {
    Parse(data []byte) (IValue, error)
    ParseString(jsonStr string) (IValue, error)
    SetMaxDepth(depth int)
    GetMaxDepth() int
}
```

#### 方法说明 / Method Descriptions

- **Parse(data []byte)** - 解析JSON字节数组
- **ParseString(jsonStr string)** - 解析JSON字符串
- **SetMaxDepth(depth int)** - 设置最大解析深度
- **GetMaxDepth()** - 获取最大解析深度

### ISerializer

JSON序列化器接口。

JSON serializer interface.

```go
type ISerializer interface {
    Serialize(value IValue) ([]byte, error)
    SerializeToString(value IValue) (string, error)
    SetOptions(options *SerializeOptions)
    GetOptions() *SerializeOptions
}
```

#### 方法说明 / Method Descriptions

- **Serialize(value IValue)** - 序列化为字节数组
- **SerializeToString(value IValue)** - 序列化为字符串
- **SetOptions(options *SerializeOptions)** - 设置序列化选项
- **GetOptions()** - 获取当前序列化选项

### IPathQuery

JSONPath查询接口。

JSONPath query interface.

```go
type IPathQuery interface {
    SelectAll(root IValue, path string) ([]IValue, error)
    SelectOne(root IValue, path string) (IValue, error)
    Set(root IValue, path string, value IValue) error
    Delete(root IValue, path string) error
    Exists(root IValue, path string) bool
    Count(root IValue, path string) int
}
```

#### 方法说明 / Method Descriptions

- **SelectAll(root IValue, path string)** - 查询所有匹配的值
- **SelectOne(root IValue, path string)** - 查询单个匹配的值
- **Set(root IValue, path string, value IValue)** - 根据路径设置值
- **Delete(root IValue, path string)** - 根据路径删除值
- **Exists(root IValue, path string)** - 检查路径是否存在
- **Count(root IValue, path string)** - 统计匹配路径的数量

## 全局函数 / Global Functions

### 解析函数 / Parsing Functions

```go
// 解析JSON字节数组
func Parse(data []byte) (IValue, error)

// 解析JSON字符串
func ParseString(data string) (IValue, error)

// 解析JSON，失败时panic
func MustParse(data []byte) IValue

// 解析JSON字符串，失败时panic
func MustParseString(data string) IValue
```

### 序列化函数 / Serialization Functions

```go
// 序列化为字节数组
func Serialize(value IValue) ([]byte, error)

// 序列化为字符串
func SerializeToString(value IValue) (string, error)

// 序列化，失败时panic
func MustSerialize(value IValue) []byte

// 序列化为字符串，失败时panic
func MustSerializeToString(value IValue) string

// 美化格式序列化
func Pretty(value IValue) (string, error)

// 美化格式序列化，失败时panic
func MustPretty(value IValue) string

// 紧凑格式序列化
func Compact(value IValue) (string, error)

// 紧凑格式序列化，失败时panic
func MustCompact(value IValue) string
```

### JSONPath查询函数 / JSONPath Query Functions

```go
// 根据路径获取值
func Get(root IValue, path string) (IValue, error)

// 根据路径获取值，失败时panic
func MustGet(root IValue, path string) IValue

// 根据路径获取所有匹配的值
func GetAll(root IValue, path string) ([]IValue, error)

// 根据路径设置值
func Set(root IValue, path string, value IValue) error

// 根据路径删除值
func Delete(root IValue, path string) error

// 检查路径是否存在
func Exists(root IValue, path string) bool

// 统计匹配路径的数量
func Count(root IValue, path string) int
```

### 类型转换函数 / Type Conversion Functions

```go
// 转换为字符串
func ToString(value IValue) (string, error)
func MustToString(value IValue) string

// 转换为整数
func ToInt(value IValue) (int, error)
func MustToInt(value IValue) int

// 转换为64位整数
func ToInt64(value IValue) (int64, error)
func MustToInt64(value IValue) int64

// 转换为64位浮点数
func ToFloat64(value IValue) (float64, error)
func MustToFloat64(value IValue) float64

// 转换为布尔值
func ToBool(value IValue) (bool, error)
func MustToBool(value IValue) bool

// 转换为时间
func ToTime(value IValue) (time.Time, error)
func MustToTime(value IValue) time.Time

// 转换为字节数组
func ToBytes(value IValue) ([]byte, error)
func MustToBytes(value IValue) []byte

// 转换为对象
func ToObject(value IValue) (IObject, error)
func MustToObject(value IValue) IObject

// 转换为数组
func ToArray(value IValue) (IArray, error)
func MustToArray(value IValue) IArray
```

## 工厂函数 / Factory Functions

### 值创建函数 / Value Creation Functions

```go
// 创建null值
func CreateNull() IValue

// 创建字符串值
func CreateString(value string) IValue

// 创建数字值
func CreateNumber(value interface{}) IValue
func MustCreateNumber(value interface{}) IValue

// 创建布尔值
func CreateBool(value bool) IValue

// 创建对象
func CreateObject() IObject
func CreateObjectWithCapacity(capacity int) IObject

// 创建数组
func CreateArray() IArray
func CreateArrayWithCapacity(capacity int) IArray

// 从Go原生类型创建值
func CreateFromRaw(value interface{}) (IValue, error)
func MustCreateFromRaw(value interface{}) IValue
```

### 构建器函数 / Builder Functions

```go
// 创建JSON构建器
func NewBuilder() *JSONBuilder
```

### 实例创建函数 / Instance Creation Functions

```go
// 创建解析器
func NewParser() IParser
func NewParserWithFactory(factory IValueFactory) IParser

// 创建序列化器
func NewSerializer() ISerializer
func NewSerializerWithOptions(options *SerializeOptions) ISerializer

// 创建路径查询器
func NewPathQuery() IPathQuery
func NewPathQueryWithFactory(factory IValueFactory) IPathQuery

// 创建值工厂
func NewValueFactory() IValueFactory
func NewValueFactoryWithPool(pool IObjectPool) IValueFactory

// 创建对象池
func NewObjectPool() IObjectPool
func NewObjectPoolWithOptions(options *ObjectPoolOptions) IObjectPool
```

### 便捷序列化器 / Convenience Serializers

```go
// 创建紧凑序列化器
func CompactSerializer() ISerializer

// 创建美化序列化器
func PrettySerializer(indent string) ISerializer
```

## 类型定义 / Type Definitions

### ValueType

JSON值类型枚举。

JSON value type enumeration.

```go
type ValueType int

const (
    NullType ValueType = iota
    StringType
    NumberType
    BooleanType
    ObjectType
    ArrayType
)
```

### SerializeOptions

序列化选项配置。

Serialization options configuration.

```go
type SerializeOptions struct {
    Indent     string // 缩进字符串
    Compact    bool   // 是否紧凑格式
    EscapeHTML bool   // 是否转义HTML字符
    SortKeys   bool   // 是否排序对象键
    MaxDepth   int    // 最大序列化深度
}
```

### ObjectPoolOptions

对象池选项配置。

Object pool options configuration.

```go
type ObjectPoolOptions struct {
    MaxValuePoolSize  int  // 值池最大大小
    MaxObjectPoolSize int  // 对象池最大大小
    MaxArrayPoolSize  int  // 数组池最大大小
    EnablePooling     bool // 是否启用池化
}
```

### PerformanceStats

性能统计信息。

Performance statistics.

```go
type PerformanceStats struct {
    ParseCount        int64         // 解析操作次数
    SerializeCount    int64         // 序列化操作次数
    TotalParseTime    time.Duration // 总解析时间
    TotalSerializeTime time.Duration // 总序列化时间
    PeakMemoryUsage   int64         // 峰值内存使用
    ErrorCount        int64         // 错误次数
}

// 计算平均解析时间
func (s *PerformanceStats) AverageParseTime() time.Duration

// 计算平均序列化时间
func (s *PerformanceStats) AverageSerializeTime() time.Duration
```

### PoolStats

对象池统计信息。

Object pool statistics.

```go
type PoolStats struct {
    TotalAllocated int64   // 总分配次数
    TotalReused    int64   // 总重用次数
    CurrentInUse   int64   // 当前使用中的对象数
    PoolHitRate    float64 // 池命中率
}
```

### MemorySnapshot

内存快照信息。

Memory snapshot information.

```go
type MemorySnapshot struct {
    Timestamp    time.Time // 快照时间戳
    TotalAlloc   uint64    // 总分配内存
    HeapAlloc    uint64    // 堆分配内存
    HeapObjects  uint64    // 堆对象数量
    NumGC        uint32    // GC次数
}
```

## 配置选项 / Configuration Options

### 默认配置 / Default Configurations

```go
// 默认序列化选项
var DefaultSerializeOptions = &SerializeOptions{
    Indent:     "",
    Compact:    false,
    EscapeHTML: true,
    SortKeys:   false,
    MaxDepth:   64,
}

// 默认对象池选项
var DefaultObjectPoolOptions = &ObjectPoolOptions{
    MaxValuePoolSize:  1000,
    MaxObjectPoolSize: 500,
    MaxArrayPoolSize:  500,
    EnablePooling:     true,
}

// 默认缩进
const DefaultIndent = "  "
```

### 全局配置函数 / Global Configuration Functions

```go
// 获取默认实例
func GetDefaultFactory() IValueFactory
func GetDefaultParser() IParser
func GetDefaultSerializer() ISerializer
func GetDefaultPathQuery() IPathQuery

// 设置默认实例
func SetDefaultFactory(factory IValueFactory)
func SetDefaultParser(parser IParser)
func SetDefaultSerializer(serializer ISerializer)
func SetDefaultPathQuery(pathQuery IPathQuery)
```

## 错误处理 / Error Handling

### 错误类型 / Error Types

```go
// 解析错误
type ParseError struct {
    Message string
    Line    int
    Column  int
    Offset  int
}

func (e *ParseError) Error() string

// 序列化错误
type SerializeError struct {
    Message string
    Path    string
}

func (e *SerializeError) Error() string

// JSONPath错误
type PathError struct {
    Message string
    Path    string
    Reason  string
}

func (e *PathError) Error() string
```

### 错误检查函数 / Error Checking Functions

```go
// 检查是否为解析错误
func IsParseError(err error) bool

// 检查是否为序列化错误
func IsSerializeError(err error) bool

// 检查是否为路径错误
func IsPathError(err error) bool
```

## 性能监控 / Performance Monitoring

### 监控控制函数 / Monitoring Control Functions

```go
// 启用性能监控
func EnablePerformanceMonitoring()

// 禁用性能监控
func DisablePerformanceMonitoring()

// 获取性能统计
func GetPerformanceStats() PerformanceStats

// 重置性能统计
func ResetPerformanceStats()
```

### 内存分析函数 / Memory Profiling Functions

```go
// 开始内存分析
func StartMemoryProfiling()

// 停止内存分析
func StopMemoryProfiling()

// 获取内存快照
func GetMemorySnapshots() []MemorySnapshot

// 获取最新内存快照
func GetLatestMemorySnapshot() *MemorySnapshot

// 获取内存趋势
func GetMemoryTrend() (trend string, growth float64)
```

### 计时器接口 / Timer Interface

```go
type Timer interface {
    End() time.Duration
}

// 获取全局监控器
func GetGlobalMonitor() *PerformanceMonitor

// 开始解析计时
func (m *PerformanceMonitor) StartParseTimer() Timer

// 开始序列化计时
func (m *PerformanceMonitor) StartSerializeTimer() Timer
```

## 版本信息 / Version Information

```go
// 版本常量
const (
    Version      = "1.0.0"
    VersionMajor = 1
    VersionMinor = 0
    VersionPatch = 0
)

// 获取版本信息
func GetVersion() string
```

## JSONPath语法支持 / JSONPath Syntax Support

### 基本语法 / Basic Syntax

- `$` - 根节点
- `.key` - 子节点
- `['key']` - 子节点（括号表示法）
- `[index]` - 数组索引
- `[start:end]` - 数组切片
- `*` - 通配符
- `..` - 递归下降
- `[?(@.key)]` - 过滤器表达式

### 支持的操作符 / Supported Operators

- `==` - 等于
- `!=` - 不等于
- `<` - 小于
- `<=` - 小于等于
- `>` - 大于
- `>=` - 大于等于
- `=~` - 正则匹配
- `in` - 包含

### 示例 / Examples

```go
// 基本路径
"$.store.book[0].title"

// 通配符
"$.store.book[*].author"

// 递归搜索
"$..price"

// 数组切片
"$.store.book[1:3]"

// 过滤器
"$.store.book[?(@.price < 30)]"

// 复杂过滤器
"$.store.book[?(@.author == 'John' && @.price > 20)]"
```

这个API参考文档涵盖了xyJson库的所有主要功能和接口。更多详细的使用示例，请参考examples目录中的示例代码。

This API reference document covers all major features and interfaces of the xyJson library. For more detailed usage examples, please refer to the example code in the examples directory.