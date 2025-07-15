package xyJson

import (
	"time"
)

// ValueType 值类型枚举
type ValueType int

const (
	NullValueType ValueType = iota
	StringValueType
	NumberValueType
	BoolValueType
	ObjectValueType
	ArrayValueType
)

// String 返回值类型的字符串表示
func (vt ValueType) String() string {
	switch vt {
	case NullValueType:
		return "null"
	case StringValueType:
		return "string"
	case NumberValueType:
		return "number"
	case BoolValueType:
		return "boolean"
	case ObjectValueType:
		return "object"
	case ArrayValueType:
		return "array"
	default:
		return "unknown"
	}
}

// IValue JSON值接口
type IValue interface {
	// Type 获取值类型
	Type() ValueType

	// Raw 获取原始Go类型值
	Raw() interface{}

	// String 获取字符串表示
	String() string

	// IsNull 检查是否为null
	IsNull() bool

	// IsString 检查是否为字符串
	IsString() bool

	// IsNumber 检查是否为数字
	IsNumber() bool

	// IsBool 检查是否为布尔值
	IsBool() bool

	// IsObject 检查是否为对象
	IsObject() bool

	// IsArray 检查是否为数组
	IsArray() bool

	// Clone 深拷贝
	Clone() IValue

	// Equals 比较是否相等
	Equals(other IValue) bool

	// Int 转换为int
	Int() (int, error)

	// Int64 转换为int64
	Int64() (int64, error)

	// Float64 转换为float64
	Float64() (float64, error)

	// Bool 转换为bool
	Bool() (bool, error)

	// Time 转换为time.Time
	Time() (time.Time, error)

	// Bytes 转换为[]byte
	Bytes() ([]byte, error)
}

// IObject JSON对象接口
type IObject interface {
	IValue

	// Get 获取字段值
	Get(key string) (IValue, bool)

	// Set 设置字段值
	Set(key string, value IValue)

	// Delete 删除字段
	Delete(key string) bool

	// Has 检查字段是否存在
	Has(key string) bool

	// Keys 获取所有键
	Keys() []string

	// Size 获取字段数量
	Size() int

	// Clear 清空所有字段
	Clear()

	// Range 遍历所有字段
	Range(fn func(key string, value IValue) bool)

	// Merge 合并另一个对象
	Merge(other IObject)

	// SortedPairs 获取排序后的键值对
	SortedPairs() []KeyValuePair
}

// IArray JSON数组接口
type IArray interface {
	IValue

	// Get 获取指定索引的元素
	Get(index int) (IValue, bool)

	// Set 设置指定索引的元素
	Set(index int, value IValue) error

	// Append 追加元素
	Append(value IValue)

	// Insert 在指定位置插入元素
	Insert(index int, value IValue) error

	// Delete 删除指定索引的元素
	Delete(index int) error

	// Length 获取数组长度
	Length() int

	// Clear 清空数组
	Clear()

	// Range 遍历数组元素
	Range(fn func(index int, value IValue) bool)

	// AppendAll 批量追加元素
	AppendAll(values []IValue)

	// IndexOf 查找元素索引
	IndexOf(value IValue) int

	// Contains 检查是否包含元素
	Contains(value IValue) bool

	// RemoveValue 移除指定值的元素
	RemoveValue(value IValue) bool

	// Slice 获取子数组
	Slice(start, end int) IArray

	// Reverse 反转数组
	Reverse()

	// Filter 过滤数组元素
	Filter(predicate func(IValue) bool) IArray
}

// KeyValuePair 键值对
type KeyValuePair struct {
	Key   string
	Value IValue
}

// IParser JSON解析器接口
type IParser interface {
	// Parse 解析JSON字节数组
	Parse(data []byte) (IValue, error)

	// ParseString 解析JSON字符串
	ParseString(s string) (IValue, error)
}

// ISerializer JSON序列化器接口
type ISerializer interface {
	// Serialize 序列化为字节数组
	Serialize(value IValue) ([]byte, error)

	// SerializeToString 序列化为字符串
	SerializeToString(value IValue) (string, error)

	// Pretty 格式化输出
	Pretty(value IValue) (string, error)

	// Compact 压缩输出
	Compact(value IValue) (string, error)
}

// IPathQuery JSONPath查询接口
type IPathQuery interface {
	// SelectOne 选择单个值
	SelectOne(root IValue, path string) (IValue, error)

	// SelectAll 选择所有匹配的值
	SelectAll(root IValue, path string) ([]IValue, error)

	// Set 设置值
	Set(root IValue, path string, value IValue) error

	// Delete 删除值
	Delete(root IValue, path string) error

	// Exists 检查路径是否存在
	Exists(root IValue, path string) bool

	// Count 统计匹配数量
	Count(root IValue, path string) int
}

// IValueFactory 值工厂接口
type IValueFactory interface {
	// CreateNull 创建null值
	CreateNull() IValue

	// CreateString 创建字符串值
	CreateString(s string) IValue

	// CreateNumber 创建数字值
	CreateNumber(n interface{}) (IValue, error)

	// CreateBool 创建布尔值
	CreateBool(b bool) IValue

	// CreateObject 创建对象
	CreateObject() IObject

	// CreateObjectWithCapacity 创建指定容量的对象
	CreateObjectWithCapacity(capacity int) IObject

	// CreateArray 创建数组
	CreateArray() IArray

	// CreateArrayWithCapacity 创建指定容量的数组
	CreateArrayWithCapacity(capacity int) IArray

	// CreateFromRaw 从原始Go类型创建值
	CreateFromRaw(v interface{}) (IValue, error)
}

// IObjectPool 对象池接口
type IObjectPool interface {
	// GetObject 从池中获取对象
	GetObject() IObject

	// PutObject 将对象放回池中
	PutObject(obj IObject)

	// GetArray 从池中获取数组
	GetArray() IArray

	// PutArray 将数组放回池中
	PutArray(arr IArray)

	// GetStats 获取池统计信息
	GetStats() PoolStats

	// Clear 清空池
	Clear()
}

// PoolStats 对象池统计信息
type PoolStats struct {
	TotalAllocated uint64  // 总分配次数
	TotalReused    uint64  // 总复用次数
	CurrentInUse   uint64  // 当前使用中的对象数
	PoolHitRate    float64 // 池命中率
}

// IBuilder JSON构建器接口
type IBuilder interface {
	// SetString 设置字符串字段
	SetString(key, value string) IBuilder

	// SetInt 设置整数字段
	SetInt(key string, value int) IBuilder

	// SetInt64 设置64位整数字段
	SetInt64(key string, value int64) IBuilder

	// SetFloat64 设置浮点数字段
	SetFloat64(key string, value float64) IBuilder

	// SetBool 设置布尔字段
	SetBool(key string, value bool) IBuilder

	// SetNull 设置null字段
	SetNull(key string) IBuilder

	// SetTime 设置时间字段
	SetTime(key string, value time.Time) IBuilder

	// SetValue 设置任意值字段
	SetValue(key string, value IValue) IBuilder

	// AddString 向数组添加字符串
	AddString(value string) IBuilder

	// AddInt 向数组添加整数
	AddInt(value int) IBuilder

	// AddBool 向数组添加布尔值
	AddBool(value bool) IBuilder

	// AddNull 向数组添加null
	AddNull() IBuilder

	// AddValue 向数组添加任意值
	AddValue(value IValue) IBuilder

	// BeginObject 开始构建嵌套对象
	BeginObject(key string) IBuilder

	// BeginArray 开始构建嵌套数组
	BeginArray(key string) IBuilder

	// AddObject 向数组添加对象
	AddObject() IBuilder

	// AddArray 向数组添加数组
	AddArray() IBuilder

	// End 结束当前嵌套层级
	End() IBuilder

	// Build 构建最终值
	Build() (IValue, error)

	// MustBuild 构建最终值，失败时panic
	MustBuild() IValue

	// Error 获取构建过程中的错误
	Error() error

	// Reset 重置构建器
	Reset() IBuilder

	// ResetAsArray 重置为数组构建器
	ResetAsArray() IBuilder
}
