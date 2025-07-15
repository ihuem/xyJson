package xyJson

import (
	"time"
)

// IValue 是所有JSON值的基础接口，定义了JSON值的通用操作
// IValue is the fundamental interface for all JSON values, defining common operations for JSON values
//
// 所有JSON值类型（null、string、number、boolean、object、array）都实现此接口
// All JSON value types (null, string, number, boolean, object, array) implement this interface
//
// 使用示例 Usage Example:
//
//	value, _ := xyJson.ParseString(`"hello"`)
//	fmt.Println(value.Type())     // StringType
//	fmt.Println(value.String())   // "hello"
//	fmt.Println(value.IsNull())   // false
//	cloned := value.Clone()       // 创建深拷贝
//	fmt.Println(value.Equals(cloned)) // true
type IValue interface {
	// Type 返回JSON值的类型（null、string、number、boolean、object、array）
	// Type returns the type of the JSON value (null, string, number, boolean, object, array)
	//
	// 返回值 Returns:
	//   - ValueType: 值类型枚举 / Value type enumeration
	Type() ValueType

	// Raw 返回对应的Go原生类型值，主要用于与标准库json包的互操作
	// Raw returns the corresponding Go native type value, mainly for interoperability with standard json package
	//
	// 返回值 Returns:
	//   - interface{}: Go原生类型值（string、float64、bool、map[string]interface{}、[]interface{}、nil）
	//                  Go native type value (string, float64, bool, map[string]interface{}, []interface{}, nil)
	Raw() interface{}

	// String 返回JSON值的字符串表示形式
	// String returns the string representation of the JSON value
	//
	// 返回值 Returns:
	//   - string: JSON格式的字符串表示 / JSON format string representation
	String() string

	// IsNull 检查当前值是否为JSON null
	// IsNull checks if the current value is JSON null
	//
	// 返回值 Returns:
	//   - bool: 如果是null值返回true，否则返回false / Returns true if null, false otherwise
	IsNull() bool

	// Clone 创建当前值的深拷贝
	// Clone creates a deep copy of the current value
	//
	// 对于复杂类型（object、array），会递归复制所有子元素
	// For complex types (object, array), recursively copies all child elements
	//
	// 返回值 Returns:
	//   - IValue: 深拷贝的新值 / Deep copied new value
	Clone() IValue

	// Equals 比较两个JSON值是否相等
	// Equals compares if two JSON values are equal
	//
	// 参数 Parameters:
	//   - other: 要比较的另一个JSON值 / Another JSON value to compare
	//
	// 返回值 Returns:
	//   - bool: 如果两个值相等返回true / Returns true if the two values are equal
	Equals(other IValue) bool
}

// IScalarValue 标量值接口（字符串、数字、布尔值）
// IScalarValue represents scalar values (string, number, boolean)
type IScalarValue interface {
	IValue

	// Int 返回整数值
	// Int returns the integer value
	Int() (int, error)

	// Int64 返回64位整数值
	// Int64 returns the 64-bit integer value
	Int64() (int64, error)

	// Float64 返回64位浮点数值
	// Float64 returns the 64-bit float value
	Float64() (float64, error)

	// Bool 返回布尔值
	// Bool returns the boolean value
	Bool() (bool, error)

	// Time 返回时间值
	// Time returns the time value
	Time() (time.Time, error)

	// Bytes 返回字节数组
	// Bytes returns the byte array
	Bytes() ([]byte, error)
}

// IObject 表示JSON对象的接口，提供键值对操作功能
// IObject represents a JSON object interface, providing key-value pair operations
//
// JSON对象是键值对的无序集合，键必须是字符串，值可以是任意JSON类型
// A JSON object is an unordered collection of key-value pairs, keys must be strings, values can be any JSON type
//
// 使用示例 Usage Example:
//
//	obj := xyJson.CreateObject()
//	obj.Set("name", "Alice")
//	obj.Set("age", 25)
//	obj.Set("active", true)
//
//	fmt.Println(obj.Get("name").String()) // "Alice"
//	fmt.Println(obj.Has("age"))           // true
//	fmt.Println(obj.Size())               // 3
//
//	// 遍历所有键值对
//	obj.Range(func(key string, value IValue) bool {
//		fmt.Printf("%s: %s\n", key, value.String())
//		return true // 继续遍历
//	})
type IObject interface {
	IValue

	// Get 根据键名获取对应的JSON值
	// Get retrieves the JSON value associated with the given key
	//
	// 参数 Parameters:
	//   - key: 要查找的键名 / Key name to look up
	//
	// 返回值 Returns:
	//   - IValue: 对应的JSON值，如果键不存在则返回nil / Associated JSON value, nil if key doesn't exist
	Get(key string) IValue

	// Set 设置或更新指定键的值
	// Set sets or updates the value for the specified key
	//
	// 参数 Parameters:
	//   - key: 键名 / Key name
	//   - value: 要设置的值，可以是Go原生类型或IValue / Value to set, can be Go native type or IValue
	//
	// 返回值 Returns:
	//   - error: 设置失败时的错误信息 / Error information when setting fails
	Set(key string, value interface{}) error

	// Delete 删除指定的键值对
	// Delete removes the specified key-value pair
	//
	// 参数 Parameters:
	//   - key: 要删除的键名 / Key name to delete
	//
	// 返回值 Returns:
	//   - bool: 如果键存在并被删除返回true，否则返回false / Returns true if key existed and was deleted, false otherwise
	Delete(key string) bool

	// Has 检查对象是否包含指定的键
	// Has checks if the object contains the specified key
	//
	// 参数 Parameters:
	//   - key: 要检查的键名 / Key name to check
	//
	// 返回值 Returns:
	//   - bool: 如果键存在返回true / Returns true if key exists
	Has(key string) bool

	// Keys 返回对象中所有键名的切片
	// Keys returns a slice of all key names in the object
	//
	// 返回值 Returns:
	//   - []string: 所有键名的切片 / Slice of all key names
	Keys() []string

	// Size 返回对象中键值对的数量
	// Size returns the number of key-value pairs in the object
	//
	// 返回值 Returns:
	//   - int: 键值对数量 / Number of key-value pairs
	Size() int

	// Clear 清空对象中的所有键值对
	// Clear removes all key-value pairs from the object
	Clear()

	// Range 遍历对象中的所有键值对
	// Range iterates over all key-value pairs in the object
	//
	// 参数 Parameters:
	//   - fn: 遍历函数，接收键名和值作为参数，返回false可提前终止遍历
	//         Iteration function that receives key and value as parameters, return false to terminate early
	//
	// 示例 Example:
	//
	//	obj.Range(func(key string, value IValue) bool {
	//		if key == "target" {
	//			return false // 找到目标键，停止遍历
	//		}
	//		fmt.Printf("%s: %s\n", key, value.String())
	//		return true // 继续遍历
	//	})
	Range(fn func(key string, value IValue) bool)
}

// IArray JSON数组接口
// IArray represents a JSON array interface
type IArray interface {
	IValue

	// Get 根据索引获取值
	// Get retrieves a value by index
	Get(index int) IValue

	// Set 设置指定索引的值
	// Set sets the value at the specified index
	Set(index int, value interface{}) error

	// Append 追加值到数组末尾
	// Append adds a value to the end of the array
	Append(value interface{}) error

	// Insert 在指定位置插入值
	// Insert inserts a value at the specified position
	Insert(index int, value interface{}) error

	// Delete 删除指定索引的值
	// Delete removes the value at the specified index
	Delete(index int) error

	// Length 返回数组长度
	// Length returns the length of the array
	Length() int

	// Clear 清空数组
	// Clear removes all elements from the array
	Clear()

	// Range 遍历数组元素
	// Range iterates over array elements
	Range(fn func(index int, value IValue) bool)
}

// IParser JSON解析器接口
// IParser represents a JSON parser interface
type IParser interface {
	// Parse 解析JSON数据
	// Parse parses JSON data
	Parse(data []byte) (IValue, error)

	// ParseString 解析JSON字符串
	// ParseString parses a JSON string
	ParseString(jsonStr string) (IValue, error)

	// SetMaxDepth 设置最大解析深度
	// SetMaxDepth sets the maximum parsing depth
	SetMaxDepth(depth int)

	// GetMaxDepth 获取最大解析深度
	// GetMaxDepth gets the maximum parsing depth
	GetMaxDepth() int
}

// ISerializer JSON序列化器接口
// ISerializer represents a JSON serializer interface
type ISerializer interface {
	// Serialize 序列化JSON值
	// Serialize serializes a JSON value
	Serialize(value IValue) ([]byte, error)

	// SerializeToString 序列化为字符串
	// SerializeToString serializes to a string
	SerializeToString(value IValue) (string, error)

	// SetOptions 设置序列化选项
	// SetOptions sets serialization options
	SetOptions(options *SerializeOptions)

	// GetOptions 获取序列化选项
	// GetOptions gets serialization options
	GetOptions() *SerializeOptions
}

// IPathQuery JSONPath查询接口
// IPathQuery represents a JSONPath query interface
type IPathQuery interface {
	// SelectAll 根据路径查询多个值
	// SelectAll queries multiple values by path
	SelectAll(root IValue, path string) ([]IValue, error)

	// SelectOne 根据路径查询单个值
	// SelectOne queries a single value by path
	SelectOne(root IValue, path string) (IValue, error)

	// Set 根据路径设置值
	// Set sets a value by path
	Set(root IValue, path string, value IValue) error

	// Delete 根据路径删除值
	// Delete deletes a value by path
	Delete(root IValue, path string) error

	// Exists 检查路径是否存在
	// Exists checks if a path exists
	Exists(root IValue, path string) bool

	// Count 统计匹配路径的数量
	// Count counts the number of matching paths
	Count(root IValue, path string) int
}

// IValueFactory 值工厂接口
// IValueFactory represents a value factory interface
type IValueFactory interface {
	// CreateNull 创建null值
	// CreateNull creates a null value
	CreateNull() IValue

	// CreateString 创建字符串值
	// CreateString creates a string value
	CreateString(s string) IScalarValue

	// CreateNumber 创建数字值
	// CreateNumber creates a number value
	CreateNumber(n interface{}) IScalarValue

	// CreateBool 创建布尔值
	// CreateBool creates a boolean value
	CreateBool(b bool) IScalarValue

	// CreateObject 创建对象
	// CreateObject creates an object
	CreateObject() IObject

	// CreateArray 创建数组
	// CreateArray creates an array
	CreateArray() IArray

	// CreateFromRaw 从原始数据创建值
	// CreateFromRaw creates a value from raw data
	CreateFromRaw(data interface{}) (IValue, error)
}

// IObjectPool 对象池接口
// IObjectPool represents an object pool interface
type IObjectPool interface {
	// GetValue 从池中获取值对象
	// GetValue gets a value object from the pool
	GetValue() IValue

	// PutValue 将值对象放回池中
	// PutValue puts a value object back to the pool
	PutValue(value IValue)

	// GetObject 从池中获取对象
	// GetObject gets an object from the pool
	GetObject() IObject

	// PutObject 将对象放回池中
	// PutObject puts an object back to the pool
	PutObject(obj IObject)

	// GetArray 从池中获取数组
	// GetArray gets an array from the pool
	GetArray() IArray

	// PutArray 将数组放回池中
	// PutArray puts an array back to the pool
	PutArray(arr IArray)

	// GetStats 获取池统计信息
	// GetStats gets pool statistics
	GetStats() *PoolStats
}

// SerializeOptions 序列化选项
// SerializeOptions represents serialization options
type SerializeOptions struct {
	// Indent 缩进字符串
	// Indent is the indentation string
	Indent string

	// Compact 是否使用紧凑模式
	// Compact indicates whether to use compact mode
	Compact bool

	// EscapeHTML 是否转义HTML字符
	// EscapeHTML indicates whether to escape HTML characters
	EscapeHTML bool

	// EscapeUnicode 是否转义Unicode字符为\u格式
	// EscapeUnicode indicates whether to escape Unicode characters to \u format
	EscapeUnicode bool

	// SortKeys 是否对键名排序
	// SortKeys indicates whether to sort object keys
	SortKeys bool

	// MaxDepth 最大序列化深度
	// MaxDepth is the maximum serialization depth
	MaxDepth int
}

// PoolStats 对象池统计信息
// PoolStats represents object pool statistics
type PoolStats struct {
	// TotalAllocated 总分配对象数
	// TotalAllocated is the total number of allocated objects
	TotalAllocated int64

	// TotalReused 总复用对象数
	// TotalReused is the total number of reused objects
	TotalReused int64

	// CurrentInUse 当前使用中的对象数
	// CurrentInUse is the current number of objects in use
	CurrentInUse int64

	// PoolHitRate 池命中率
	// PoolHitRate is the pool hit rate
	PoolHitRate float64
}
