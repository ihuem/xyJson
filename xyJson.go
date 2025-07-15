// Package xyJson provides a high-performance JSON processing library for Go.
// It offers comprehensive JSON parsing, serialization, and JSONPath query capabilities
// with built-in performance monitoring and memory optimization features.
//
// xyJson 是一个高性能的 Go JSON 处理库，提供全面的 JSON 解析、序列化和 JSONPath 查询功能，
// 内置性能监控和内存优化特性。
//
// Key Features / 主要特性:
//   - High-performance JSON parsing and serialization / 高性能 JSON 解析和序列化
//   - JSONPath query support / JSONPath 查询支持
//   - Built-in performance monitoring / 内置性能监控
//   - Memory pool optimization / 内存池优化
//   - Thread-safe operations / 线程安全操作
//   - Comprehensive error handling / 完善的错误处理
//
// Basic Usage / 基本用法:
//
//	// Parse JSON / 解析 JSON
//	value, err := xyJson.Parse([]byte(`{"name":"John","age":30}`))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Access values / 访问值
//	name := xyJson.MustGet(value, "$.name")
//	fmt.Println(name.String()) // "John"
//
//	// Serialize JSON / 序列化 JSON
//	jsonStr, err := xyJson.SerializeToString(value)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Performance Monitoring / 性能监控:
//
//	// Enable monitoring / 启用监控
//	xyJson.EnablePerformanceMonitoring()
//
//	// Get statistics / 获取统计信息
//	stats := xyJson.GetPerformanceStats()
//	fmt.Printf("Parse operations: %d\n", stats.ParseCount)
package xyJson

import (
	"sync"
	"time"
)

// 包级别的全局实例，提供便捷的默认功能访问
// Package-level global instances providing convenient access to default functionality
var (
	// defaultFactory 默认值工厂，用于创建各种JSON值类型
	// defaultFactory is the default value factory for creating various JSON value types
	defaultFactory IValueFactory

	// defaultParser 默认解析器，用于解析JSON数据
	// defaultParser is the default parser for parsing JSON data
	defaultParser IParser

	// defaultSerializer 默认序列化器，用于序列化JSON值
	// defaultSerializer is the default serializer for serializing JSON values
	defaultSerializer ISerializer

	// defaultPathQuery 默认路径查询器，用于JSONPath查询
	// defaultPathQuery is the default path query for JSONPath operations
	defaultPathQuery IPathQuery

	// parserPool parser对象池，用于重用parser实例以提高性能和减少内存分配
	// parserPool is a parser object pool for reusing parser instances to improve performance and reduce memory allocation
	parserPool sync.Pool
)

// init 初始化默认实例
// init initializes default instances
func init() {
	pool := NewObjectPool()
	defaultFactory = NewValueFactoryWithPool(pool)
	defaultParser = NewParserWithFactory(defaultFactory)
	defaultSerializer = NewSerializer()
	defaultPathQuery = NewPathQueryWithFactory(defaultFactory)

	// 初始化parser对象池
	// Initialize parser object pool
	parserPool.New = func() interface{} {
		return NewParserWithFactory(defaultFactory)
	}
}

// Parse 解析JSON字节数组为IValue接口
// Parse parses a JSON byte array into an IValue interface
//
// 参数 Parameters:
//   - data: 要解析的JSON字节数组 / JSON byte array to parse
//
// 返回值 Returns:
//   - IValue: 解析后的JSON值，可以是对象、数组或标量值 / Parsed JSON value (object, array, or scalar)
//   - error: 解析错误，如果JSON格式无效 / Parse error if JSON format is invalid
//
// 示例 Example:
//
//	data := []byte(`{"name":"Alice","age":25,"active":true}`)
//	value, err := xyJson.Parse(data)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// 访问对象字段 / Access object fields
//	obj := value.(xyJson.IObject)
//	name := obj.Get("name").String() // "Alice"
func Parse(data []byte) (IValue, error) {
	timer := GetGlobalMonitor().StartParseTimer()
	var hasError bool
	defer func() {
		if hasError {
			timer.EndWithError()
		} else {
			timer.End()
		}
	}()

	// 从对象池获取parser实例以提高性能
	// Get parser instance from object pool for better performance
	parser := parserPool.Get().(IParser)
	defer parserPool.Put(parser)

	result, err := parser.Parse(data)
	if err != nil {
		hasError = true
	}
	return result, err
}

// ParseString 解析JSON字符串为IValue接口
// ParseString parses a JSON string into an IValue interface
//
// 参数 Parameters:
//   - data: 要解析的JSON字符串 / JSON string to parse
//
// 返回值 Returns:
//   - IValue: 解析后的JSON值 / Parsed JSON value
//   - error: 解析错误 / Parse error
//
// 示例 Example:
//
//	jsonStr := `[1,2,3,{"key":"value"}]`
//	value, err := xyJson.ParseString(jsonStr)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// 访问数组元素 / Access array elements
//	arr := value.(xyJson.IArray)
//	firstItem := arr.Get(0) // 1
func ParseString(data string) (IValue, error) {
	timer := GetGlobalMonitor().StartParseTimer()
	var hasError bool
	defer func() {
		if hasError {
			timer.EndWithError()
		} else {
			timer.End()
		}
	}()

	// 从对象池获取parser实例以提高性能
	// Get parser instance from object pool for better performance
	parser := parserPool.Get().(IParser)
	defer parserPool.Put(parser)

	result, err := parser.ParseString(data)
	if err != nil {
		hasError = true
	}
	return result, err
}

// MustParse 解析JSON，如果失败则panic
// MustParse parses JSON, panics on failure
func MustParse(data []byte) IValue {
	result, err := Parse(data)
	if err != nil {
		panic(err)
	}
	return result
}

// MustParseString 解析JSON字符串，如果失败则panic
// MustParseString parses JSON string, panics on failure
func MustParseString(data string) IValue {
	result, err := ParseString(data)
	if err != nil {
		panic(err)
	}
	return result
}

// Serialize 将JSON值序列化为字节数组
// Serialize serializes a JSON value to a byte array
//
// 参数 Parameters:
//   - value: 要序列化的JSON值 / JSON value to serialize
//
// 返回值 Returns:
//   - []byte: 序列化后的JSON字节数组 / Serialized JSON byte array
//   - error: 序列化错误 / Serialization error
//
// 示例 Example:
//
//	obj := xyJson.CreateObject()
//	obj.Set("name", "Bob")
//	obj.Set("age", 30)
//	data, err := xyJson.Serialize(obj)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(string(data)) // {"name":"Bob","age":30}
func Serialize(value IValue) ([]byte, error) {
	timer := GetGlobalMonitor().StartSerializeTimer()
	var hasError bool
	defer func() {
		if hasError {
			timer.EndWithError()
		} else {
			timer.End()
		}
	}()

	result, err := defaultSerializer.Serialize(value)
	if err != nil {
		hasError = true
	}
	return result, err
}

// SerializeToString 将JSON值序列化为字符串
// SerializeToString serializes a JSON value to a string
//
// 参数 Parameters:
//   - value: 要序列化的JSON值 / JSON value to serialize
//
// 返回值 Returns:
//   - string: 序列化后的JSON字符串 / Serialized JSON string
//   - error: 序列化错误 / Serialization error
//
// 示例 Example:
//
//	arr := xyJson.CreateArray()
//	arr.Append(1)
//	arr.Append("hello")
//	arr.Append(true)
//	jsonStr, err := xyJson.SerializeToString(arr)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(jsonStr) // [1,"hello",true]
func SerializeToString(value IValue) (string, error) {
	timer := GetGlobalMonitor().StartSerializeTimer()
	var hasError bool
	defer func() {
		if hasError {
			timer.EndWithError()
		} else {
			timer.End()
		}
	}()

	result, err := defaultSerializer.SerializeToString(value)
	if err != nil {
		hasError = true
	}
	return result, err
}

// MustSerialize 序列化JSON值，如果失败则panic
// MustSerialize serializes JSON value, panics on failure
func MustSerialize(value IValue) []byte {
	result, err := Serialize(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustSerializeToString 序列化JSON值到字符串，如果失败则panic
// MustSerializeToString serializes JSON value to string, panics on failure
func MustSerializeToString(value IValue) string {
	result, err := SerializeToString(value)
	if err != nil {
		panic(err)
	}
	return result
}

// Pretty 格式化JSON值
// Pretty formats JSON value with indentation
func Pretty(value IValue) (string, error) {
	prettySerializer := PrettySerializer(DefaultIndent)
	return prettySerializer.SerializeToString(value)
}

// MustPretty 格式化JSON值，如果失败则panic
// MustPretty formats JSON value, panics on failure
func MustPretty(value IValue) string {
	result, err := Pretty(value)
	if err != nil {
		panic(err)
	}
	return result
}

// Compact 压缩JSON值
// Compact compacts JSON value
func Compact(value IValue) (string, error) {
	compactSerializer := CompactSerializer()
	return compactSerializer.SerializeToString(value)
}

// MustCompact 压缩JSON值，如果失败则panic
// MustCompact compacts JSON value, panics on failure
func MustCompact(value IValue) string {
	result, err := Compact(value)
	if err != nil {
		panic(err)
	}
	return result
}

// Get 使用JSONPath表达式从根值中获取单个匹配的值
// Get retrieves a single matching value from the root using a JSONPath expression
//
// 参数 Parameters:
//   - root: 根JSON值，作为查询的起点 / Root JSON value as the starting point for the query
//   - path: JSONPath表达式字符串 / JSONPath expression string
//
// 返回值 Returns:
//   - IValue: 匹配的JSON值，如果没有匹配则返回nil / Matching JSON value, nil if no match
//   - error: 查询错误或路径格式错误 / Query error or path format error
//
// 支持的JSONPath语法 Supported JSONPath syntax:
//   - $.key - 获取对象的key字段 / Get object's key field
//   - $[0] - 获取数组的第一个元素 / Get first element of array
//   - $.*.key - 获取所有对象的key字段 / Get key field from all objects
//   - $..key - 递归搜索key字段 / Recursively search for key field
//
// 示例 Example:
//
//	data := `{"users":[{"name":"Alice","age":25},{"name":"Bob","age":30}]}`
//	root, _ := xyJson.ParseString(data)
//	// 获取第一个用户的名字 / Get first user's name
//	name, err := xyJson.Get(root, "$.users[0].name")
//	if err == nil {
//		fmt.Println(name.String()) // "Alice"
//	}
func Get(root IValue, path string) (IValue, error) {
	return defaultPathQuery.SelectOne(root, path)
}

// MustGet 使用JSONPath获取值，如果失败则panic
// MustGet gets value using JSONPath, panics on failure
//
// 这是Get函数的便捷版本，适用于确信路径存在的场景
// This is a convenience version of Get, suitable for scenarios where the path is certain to exist
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - IValue: 匹配的JSON值 / Matching JSON value
//
// 注意 Note: 如果路径不存在或查询失败，此函数会panic / This function panics if path doesn't exist or query fails
func MustGet(root IValue, path string) IValue {
	result, err := Get(root, path)
	if err != nil {
		panic(err)
	}
	return result
}

// GetAll 使用JSONPath表达式获取所有匹配的值
// GetAll retrieves all matching values using a JSONPath expression
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - []IValue: 所有匹配的JSON值数组 / Array of all matching JSON values
//   - error: 查询错误 / Query error
//
// 示例 Example:
//
//	data := `{"products":[{"price":10},{"price":20},{"price":15}]}`
//	root, _ := xyJson.ParseString(data)
//	// 获取所有产品的价格 / Get all product prices
//	prices, err := xyJson.GetAll(root, "$.products[*].price")
//	if err == nil {
//		for _, price := range prices {
//			fmt.Println(price.String()) // 10, 20, 15
//		}
//	}
func GetAll(root IValue, path string) ([]IValue, error) {
	return defaultPathQuery.SelectAll(root, path)
}

// Set 根据路径设置值
// Set sets value by path
func Set(root IValue, path string, value any) error {
	v, err := defaultFactory.CreateFromRaw(value)
	if err != nil {
		return err
	}
	return defaultPathQuery.Set(root, path, v)
}

// Delete 根据路径删除值
// Delete deletes value by path
func Delete(root IValue, path string) error {
	return defaultPathQuery.Delete(root, path)
}

// Exists 检查路径是否存在
// Exists checks if path exists
func Exists(root IValue, path string) bool {
	return defaultPathQuery.Exists(root, path)
}

// Count 统计匹配路径的数量
// Count counts matching paths
func Count(root IValue, path string) int {
	return defaultPathQuery.Count(root, path)
}

// BatchResult 批量操作结果结构体
// BatchResult represents the result of a batch operation
type BatchResult struct {
	// Path 操作的JSONPath路径
	// Path is the JSONPath that was operated on
	Path string
	// Value 获取到的值（仅用于GetBatch）
	// Value is the retrieved value (only for GetBatch)
	Value IValue
	// Error 操作过程中的错误
	// Error is any error that occurred during the operation
	Error error
}

// BatchSetOperation 批量设置操作结构体
// BatchSetOperation represents a single set operation in a batch
type BatchSetOperation struct {
	// Path 要设置的JSONPath路径
	// Path is the JSONPath to set the value at
	Path string
	// Value 要设置的值
	// Value is the value to set
	Value interface{}
}

// BatchSetResult 批量设置操作结果结构体
// BatchSetResult represents the result of a batch set operation
type BatchSetResult struct {
	// Path 操作的JSONPath路径
	// Path is the JSONPath that was operated on
	Path string
	// Error 操作过程中的错误
	// Error is any error that occurred during the operation
	Error error
}

// GetBatch 批量获取多个路径的值
// GetBatch retrieves values for multiple paths in a single operation
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - paths: JSONPath表达式数组 / Array of JSONPath expressions
//
// 返回值 Returns:
//   - []BatchResult: 批量操作结果数组，结果顺序与输入路径顺序一致 / Array of batch results, order matches input paths
//
// 示例 Example:
//
//	data := `{"user":{"name":"Alice","age":25},"settings":{"theme":"dark"}}`
//	root, _ := xyJson.ParseString(data)
//	paths := []string{"$.user.name", "$.user.age", "$.settings.theme", "$.nonexistent"}
//	results := xyJson.GetBatch(root, paths)
//	for _, result := range results {
//		if result.Error != nil {
//			fmt.Printf("Error getting %s: %v\n", result.Path, result.Error)
//		} else {
//			fmt.Printf("%s = %v\n", result.Path, result.Value)
//		}
//	}
func GetBatch(root IValue, paths []string) []BatchResult {
	results := make([]BatchResult, len(paths))
	for i, path := range paths {
		value, err := Get(root, path)
		results[i] = BatchResult{
			Path:  path,
			Value: value,
			Error: err,
		}
	}
	return results
}

// SetBatch 批量设置多个路径的值
// SetBatch sets values for multiple paths in a single operation
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - operations: 批量设置操作数组 / Array of batch set operations
//
// 返回值 Returns:
//   - []BatchSetResult: 批量设置结果数组，结果顺序与输入操作顺序一致 / Array of batch set results, order matches input operations
//
// 示例 Example:
//
//	root := xyJson.CreateObject()
//	operations := []xyJson.BatchSetOperation{
//		{Path: "$.user.name", Value: "Bob"},
//		{Path: "$.user.age", Value: 30},
//		{Path: "$.settings.theme", Value: "light"},
//	}
//	results := xyJson.SetBatch(root, operations)
//	for _, result := range results {
//		if result.Error != nil {
//			fmt.Printf("Error setting %s: %v\n", result.Path, result.Error)
//		} else {
//			fmt.Printf("Successfully set %s\n", result.Path)
//		}
//	}
func SetBatch(root IValue, operations []BatchSetOperation) []BatchSetResult {
	results := make([]BatchSetResult, len(operations))
	for i, op := range operations {
		err := Set(root, op.Path, op.Value)
		results[i] = BatchSetResult{
			Path:  op.Path,
			Error: err,
		}
	}
	return results
}

// GetString 使用JSONPath获取字符串值
// GetString gets string value using JSONPath
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - string: 字符串值 / String value
//   - error: 查询或转换错误 / Query or conversion error
//
// 示例 Example:
//
//	data := `{"user":{"name":"Alice"}}`
//	root, _ := xyJson.ParseString(data)
//	name, err := xyJson.GetString(root, "$.user.name")
//	if err == nil {
//		fmt.Println(name) // "Alice"
//	}
func GetString(root IValue, path string) (string, error) {
	value, err := Get(root, path)
	if err != nil {
		return "", err
	}
	return ToString(value)
}

// GetInt 使用JSONPath获取整数值
// GetInt gets integer value using JSONPath
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - int: 整数值 / Integer value
//   - error: 查询或转换错误 / Query or conversion error
func GetInt(root IValue, path string) (int, error) {
	value, err := Get(root, path)
	if err != nil {
		return 0, err
	}
	return ToInt(value)
}

// GetInt64 使用JSONPath获取64位整数值
// GetInt64 gets 64-bit integer value using JSONPath
func GetInt64(root IValue, path string) (int64, error) {
	value, err := Get(root, path)
	if err != nil {
		return 0, err
	}
	return ToInt64(value)
}

// GetFloat64 使用JSONPath获取浮点数值
// GetFloat64 gets float64 value using JSONPath
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - float64: 浮点数值 / Float64 value
//   - error: 查询或转换错误 / Query or conversion error
//
// 示例 Example:
//
//	data := `{"product":{"price":29.99}}`
//	root, _ := xyJson.ParseString(data)
//	price, err := xyJson.GetFloat64(root, "$.product.price")
//	if err == nil {
//		fmt.Printf("Price: %.2f\n", price) // "Price: 29.99"
//	}
func GetFloat64(root IValue, path string) (float64, error) {
	value, err := Get(root, path)
	if err != nil {
		return 0, err
	}
	return ToFloat64(value)
}

// GetBool 使用JSONPath获取布尔值
// GetBool gets boolean value using JSONPath
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - bool: 布尔值 / Boolean value
//   - error: 查询或转换错误 / Query or conversion error
func GetBool(root IValue, path string) (bool, error) {
	value, err := Get(root, path)
	if err != nil {
		return false, err
	}
	return ToBool(value)
}

// GetObject 使用JSONPath获取对象值
// GetObject gets object value using JSONPath
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - IObject: 对象值 / Object value
//   - error: 查询或转换错误 / Query or conversion error
func GetObject(root IValue, path string) (IObject, error) {
	value, err := Get(root, path)
	if err != nil {
		return nil, err
	}
	return ToObject(value)
}

// GetArray 使用JSONPath获取数组值
// GetArray gets array value using JSONPath
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - IArray: 数组值 / Array value
//   - error: 查询或转换错误 / Query or conversion error
func GetArray(root IValue, path string) (IArray, error) {
	value, err := Get(root, path)
	if err != nil {
		return nil, err
	}
	return ToArray(value)
}

// MustGetString 使用JSONPath获取字符串值，如果失败则panic
// MustGetString gets string value using JSONPath, panics on failure
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - string: 字符串值 / String value
//
// 警告 Warning: 此方法在失败时会panic，仅在确信数据正确时使用
// This method panics on failure, use only when you're certain the data is correct
// 推荐使用TryGetString作为更安全的替代方案 / Consider using TryGetString as a safer alternative
func MustGetString(root IValue, path string) string {
	result, err := GetString(root, path)
	if err != nil {
		panic(err)
	}
	return result
}

// MustGetInt 使用JSONPath获取整数值，如果失败则panic
// MustGetInt gets integer value using JSONPath, panics on failure
//
// 警告 Warning: 此方法在失败时会panic，仅在确信数据正确时使用
// This method panics on failure, use only when you're certain the data is correct
// 推荐使用TryGetInt作为更安全的替代方案 / Consider using TryGetInt as a safer alternative
func MustGetInt(root IValue, path string) int {
	result, err := GetInt(root, path)
	if err != nil {
		panic(err)
	}
	return result
}

// MustGetInt64 使用JSONPath获取64位整数值，如果失败则panic
// MustGetInt64 gets 64-bit integer value using JSONPath, panics on failure
//
// 警告 Warning: 此方法在失败时会panic，仅在确信数据正确时使用
// This method panics on failure, use only when you're certain the data is correct
// 推荐使用TryGetInt64作为更安全的替代方案 / Consider using TryGetInt64 as a safer alternative
func MustGetInt64(root IValue, path string) int64 {
	result, err := GetInt64(root, path)
	if err != nil {
		panic(err)
	}
	return result
}

// MustGetFloat64 使用JSONPath获取浮点数值，如果失败则panic
// MustGetFloat64 gets float64 value using JSONPath, panics on failure
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - float64: 浮点数值 / Float64 value
//
// 警告 Warning: 此方法在失败时会panic，仅在确信数据正确时使用
// This method panics on failure, use only when you're certain the data is correct
// 推荐使用TryGetFloat64作为更安全的替代方案 / Consider using TryGetFloat64 as a safer alternative
func MustGetFloat64(root IValue, path string) float64 {
	result, err := GetFloat64(root, path)
	if err != nil {
		panic(err)
	}
	return result
}

// MustGetBool 使用JSONPath获取布尔值，如果失败则panic
// MustGetBool gets boolean value using JSONPath, panics on failure
//
// 警告 Warning: 此方法在失败时会panic，仅在确信数据正确时使用
// This method panics on failure, use only when you're certain the data is correct
// 推荐使用TryGetBool作为更安全的替代方案 / Consider using TryGetBool as a safer alternative
func MustGetBool(root IValue, path string) bool {
	result, err := GetBool(root, path)
	if err != nil {
		panic(err)
	}
	return result
}

// MustGetObject 使用JSONPath获取对象值，如果失败则panic
// MustGetObject gets object value using JSONPath, panics on failure
//
// 警告 Warning: 此方法在失败时会panic，仅在确信数据正确时使用
// This method panics on failure, use only when you're certain the data is correct
// 推荐使用TryGetObject作为更安全的替代方案 / Consider using TryGetObject as a safer alternative
func MustGetObject(root IValue, path string) IObject {
	result, err := GetObject(root, path)
	if err != nil {
		panic(err)
	}
	return result
}

// MustGetArray 使用JSONPath获取数组值，如果失败则panic
// MustGetArray gets array value using JSONPath, panics on failure
//
// 警告 Warning: 此方法在失败时会panic，仅在确信数据正确时使用
// This method panics on failure, use only when you're certain the data is correct
func MustGetArray(root IValue, path string) IArray {
	result, err := GetArray(root, path)
	if err != nil {
		panic(err)
	}
	return result
}

// TryGetString 使用JSONPath尝试获取字符串值
// TryGetString attempts to get string value using JSONPath
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - string: 字符串值，失败时返回空字符串 / String value, empty string on failure
//   - bool: 是否成功获取 / Whether the operation succeeded
//
// 示例 Example:
//
//	data := `{"user":{"name":"Alice"}}`
//	root, _ := xyJson.ParseString(data)
//	if name, ok := xyJson.TryGetString(root, "$.user.name"); ok {
//		fmt.Println("姓名:", name) // "姓名: Alice"
//	} else {
//		fmt.Println("获取姓名失败")
//	}
func TryGetString(root IValue, path string) (string, bool) {
	result, err := GetString(root, path)
	if err != nil {
		return "", false
	}
	return result, true
}

// TryGetInt 使用JSONPath尝试获取整数值
// TryGetInt attempts to get integer value using JSONPath
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - int: 整数值，失败时返回0 / Integer value, 0 on failure
//   - bool: 是否成功获取 / Whether the operation succeeded
func TryGetInt(root IValue, path string) (int, bool) {
	result, err := GetInt(root, path)
	if err != nil {
		return 0, false
	}
	return result, true
}

// TryGetInt64 使用JSONPath尝试获取64位整数值
// TryGetInt64 attempts to get 64-bit integer value using JSONPath
func TryGetInt64(root IValue, path string) (int64, bool) {
	result, err := GetInt64(root, path)
	if err != nil {
		return 0, false
	}
	return result, true
}

// TryGetFloat64 使用JSONPath尝试获取浮点数值
// TryGetFloat64 attempts to get float64 value using JSONPath
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - float64: 浮点数值，失败时返回0.0 / Float64 value, 0.0 on failure
//   - bool: 是否成功获取 / Whether the operation succeeded
//
// 示例 Example:
//
//	data := `{"product":{"price":29.99}}`
//	root, _ := xyJson.ParseString(data)
//	if price, ok := xyJson.TryGetFloat64(root, "$.product.price"); ok {
//		fmt.Printf("价格: %.2f\n", price) // "价格: 29.99"
//	} else {
//		fmt.Println("获取价格失败")
//	}
func TryGetFloat64(root IValue, path string) (float64, bool) {
	result, err := GetFloat64(root, path)
	if err != nil {
		return 0.0, false
	}
	return result, true
}

// TryGetBool 使用JSONPath尝试获取布尔值
// TryGetBool attempts to get boolean value using JSONPath
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - bool: 布尔值，失败时返回false / Boolean value, false on failure
//   - bool: 是否成功获取 / Whether the operation succeeded
func TryGetBool(root IValue, path string) (bool, bool) {
	result, err := GetBool(root, path)
	if err != nil {
		return false, false
	}
	return result, true
}

// TryGetObject 使用JSONPath尝试获取对象值
// TryGetObject attempts to get object value using JSONPath
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - IObject: 对象值，失败时返回nil / Object value, nil on failure
//   - bool: 是否成功获取 / Whether the operation succeeded
func TryGetObject(root IValue, path string) (IObject, bool) {
	result, err := GetObject(root, path)
	if err != nil {
		return nil, false
	}
	return result, true
}

// TryGetArray 使用JSONPath尝试获取数组值
// TryGetArray attempts to get array value using JSONPath
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//
// 返回值 Returns:
//   - IArray: 数组值，失败时返回nil / Array value, nil on failure
//   - bool: 是否成功获取 / Whether the operation succeeded
func TryGetArray(root IValue, path string) (IArray, bool) {
	result, err := GetArray(root, path)
	if err != nil {
		return nil, false
	}
	return result, true
}

// GetStringWithDefault 使用JSONPath获取字符串值，失败时返回默认值
// GetStringWithDefault gets string value using JSONPath, returns default value on failure
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//   - defaultValue: 默认值 / Default value
//
// 返回值 Returns:
//   - string: 字符串值或默认值 / String value or default value
//
// 示例 Example:
//
//	data := `{"user":{"name":"Alice"}}`
//	root, _ := xyJson.ParseString(data)
//	name := xyJson.GetStringWithDefault(root, "$.user.name", "Unknown") // "Alice"
//	city := xyJson.GetStringWithDefault(root, "$.user.city", "Unknown") // "Unknown"
func GetStringWithDefault(root IValue, path string, defaultValue string) string {
	if result, ok := TryGetString(root, path); ok {
		return result
	}
	return defaultValue
}

// GetIntWithDefault 使用JSONPath获取整数值，失败时返回默认值
// GetIntWithDefault gets integer value using JSONPath, returns default value on failure
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//   - defaultValue: 默认值 / Default value
//
// 返回值 Returns:
//   - int: 整数值或默认值 / Integer value or default value
func GetIntWithDefault(root IValue, path string, defaultValue int) int {
	if result, ok := TryGetInt(root, path); ok {
		return result
	}
	return defaultValue
}

// GetInt64WithDefault 使用JSONPath获取64位整数值，失败时返回默认值
// GetInt64WithDefault gets 64-bit integer value using JSONPath, returns default value on failure
func GetInt64WithDefault(root IValue, path string, defaultValue int64) int64 {
	if result, ok := TryGetInt64(root, path); ok {
		return result
	}
	return defaultValue
}

// GetFloat64WithDefault 使用JSONPath获取浮点数值，失败时返回默认值
// GetFloat64WithDefault gets float64 value using JSONPath, returns default value on failure
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//   - defaultValue: 默认值 / Default value
//
// 返回值 Returns:
//   - float64: 浮点数值或默认值 / Float64 value or default value
//
// 示例 Example:
//
//	data := `{"product":{"price":29.99}}`
//	root, _ := xyJson.ParseString(data)
//	price := xyJson.GetFloat64WithDefault(root, "$.product.price", 0.0) // 29.99
//	discount := xyJson.GetFloat64WithDefault(root, "$.product.discount", 0.0) // 0.0
func GetFloat64WithDefault(root IValue, path string, defaultValue float64) float64 {
	if result, ok := TryGetFloat64(root, path); ok {
		return result
	}
	return defaultValue
}

// GetBoolWithDefault 使用JSONPath获取布尔值，失败时返回默认值
// GetBoolWithDefault gets boolean value using JSONPath, returns default value on failure
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//   - defaultValue: 默认值 / Default value
//
// 返回值 Returns:
//   - bool: 布尔值或默认值 / Boolean value or default value
func GetBoolWithDefault(root IValue, path string, defaultValue bool) bool {
	if result, ok := TryGetBool(root, path); ok {
		return result
	}
	return defaultValue
}

// GetObjectWithDefault 使用JSONPath获取对象值，失败时返回默认值
// GetObjectWithDefault gets object value using JSONPath, returns default value on failure
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//   - defaultValue: 默认值，如果为nil则返回空对象 / Default value, returns empty object if nil
//
// 返回值 Returns:
//   - IObject: 对象值或默认值 / Object value or default value
//
// 注意 Note:
//
//	当defaultValue为nil时，返回一个空的对象，而不是Go的nil
//	When defaultValue is nil, returns an empty object instead of Go nil
func GetObjectWithDefault(root IValue, path string, defaultValue IObject) IObject {
	if result, ok := TryGetObject(root, path); ok {
		return result
	}
	// 如果默认值为nil，返回空对象而不是nil
	// If default value is nil, return empty object instead of nil
	if defaultValue == nil {
		return CreateObject()
	}
	return defaultValue
}

// GetArrayWithDefault 使用JSONPath获取数组值，失败时返回默认值
// GetArrayWithDefault gets array value using JSONPath, returns default value on failure
//
// 参数 Parameters:
//   - root: 根JSON值 / Root JSON value
//   - path: JSONPath表达式 / JSONPath expression
//   - defaultValue: 默认值，如果为nil则返回空数组 / Default value, returns empty array if nil
//
// 返回值 Returns:
//   - IArray: 数组值或默认值 / Array value or default value
//
// 注意 Note:
//
//	当defaultValue为nil时，返回一个空的数组，而不是Go的nil
//	When defaultValue is nil, returns an empty array instead of Go nil
func GetArrayWithDefault(root IValue, path string, defaultValue IArray) IArray {
	if result, ok := TryGetArray(root, path); ok {
		return result
	}
	// 如果默认值为nil，返回空数组而不是nil
	// If default value is nil, return empty array instead of nil
	if defaultValue == nil {
		return CreateArray()
	}
	return defaultValue
}

// CreateNull 创建一个表示JSON null的值
// CreateNull creates a value representing JSON null
//
// 返回值 Returns:
//   - IValue: JSON null值 / JSON null value
//
// 示例 Example:
//
//	nullValue := xyJson.CreateNull()
//	fmt.Println(nullValue.IsNull()) // true
//	fmt.Println(nullValue.String()) // "null"
func CreateNull() IValue {
	return defaultFactory.CreateNull()
}

// CreateString 创建一个JSON字符串值
// CreateString creates a JSON string value
//
// 参数 Parameters:
//   - value: 字符串内容 / String content
//
// 返回值 Returns:
//   - IValue: JSON字符串值 / JSON string value
//
// 示例 Example:
//
//	strValue := xyJson.CreateString("Hello, World!")
//	fmt.Println(strValue.String()) // "Hello, World!"
//	fmt.Println(strValue.Type()) // StringType
func CreateString(value string) IValue {
	return defaultFactory.CreateString(value)
}

// CreateNumber 创建一个JSON数字值
// CreateNumber creates a JSON number value
//
// 参数 Parameters:
//   - value: 数字值，支持int、int64、float64等数字类型 / Number value, supports int, int64, float64, etc.
//
// 返回值 Returns:
//   - IValue: JSON数字值 / JSON number value
//
// 示例 Example:
//
//	intValue := xyJson.CreateNumber(42)
//	floatValue := xyJson.CreateNumber(3.14159)
//	fmt.Println(intValue.String()) // "42"
//	fmt.Println(floatValue.String()) // "3.14159"
func CreateNumber(value interface{}) IValue {
	return defaultFactory.CreateNumber(value)
}

// MustCreateNumber 创建数字值，如果失败则panic
// MustCreateNumber creates a number value, panics on failure
func MustCreateNumber(value interface{}) IValue {
	return CreateNumber(value)
}

// CreateBool 创建布尔值
// CreateBool creates a boolean value
func CreateBool(value bool) IValue {
	return defaultFactory.CreateBool(value)
}

// CreateObject 创建对象
// CreateObject creates an object
func CreateObject() IObject {
	return defaultFactory.CreateObject()
}

// CreateObjectWithCapacity 创建指定容量的对象
// CreateObjectWithCapacity creates an object with specified capacity
func CreateObjectWithCapacity(capacity int) IObject {
	obj := NewObjectWithCapacity(capacity)
	return obj
}

// CreateArray 创建数组
// CreateArray creates an array
func CreateArray() IArray {
	return defaultFactory.CreateArray()
}

// CreateArrayWithCapacity 创建指定容量的数组
// CreateArrayWithCapacity creates an array with specified capacity
func CreateArrayWithCapacity(capacity int) IArray {
	arr := NewArrayWithCapacity(capacity)
	return arr
}

// CreateFromRaw 从原始数据创建JSON值
// CreateFromRaw creates JSON value from raw data
func CreateFromRaw(value interface{}) (IValue, error) {
	return defaultFactory.CreateFromRaw(value)
}

// MustCreateFromRaw 从原始数据创建JSON值，如果失败则panic
// MustCreateFromRaw creates JSON value from raw data, panics on failure
func MustCreateFromRaw(value interface{}) IValue {
	result, err := CreateFromRaw(value)
	if err != nil {
		panic(err)
	}
	return result
}

// NewBuilder 创建JSON构建器
// NewBuilder creates a JSON builder
func NewBuilder() *JSONBuilder {
	return NewJSONBuilderWithFactory(defaultFactory)
}

// GetDefaultFactory 获取默认工厂
// GetDefaultFactory gets the default factory
func GetDefaultFactory() IValueFactory {
	return defaultFactory
}

// GetDefaultParser 获取默认解析器
// GetDefaultParser gets the default parser
func GetDefaultParser() IParser {
	return defaultParser
}

// GetDefaultSerializer 获取默认序列化器
// GetDefaultSerializer gets the default serializer
func GetDefaultSerializer() ISerializer {
	return defaultSerializer
}

// GetDefaultPathQuery 获取默认路径查询器
// GetDefaultPathQuery gets the default path query
func GetDefaultPathQuery() IPathQuery {
	return defaultPathQuery
}

// SetDefaultFactory 设置默认工厂
// SetDefaultFactory sets the default factory
func SetDefaultFactory(factory IValueFactory) {
	if factory != nil {
		defaultFactory = factory
		defaultParser = NewParserWithFactory(factory)
		defaultPathQuery = NewPathQueryWithFactory(factory)
	}
}

// SetDefaultParser 设置默认解析器
// SetDefaultParser sets the default parser
func SetDefaultParser(parser IParser) {
	if parser != nil {
		defaultParser = parser
	}
}

// SetDefaultSerializer 设置默认序列化器
// SetDefaultSerializer sets the default serializer
func SetDefaultSerializer(serializer ISerializer) {
	if serializer != nil {
		defaultSerializer = serializer
	}
}

// SetDefaultPathQuery 设置默认路径查询器
// SetDefaultPathQuery sets the default path query
func SetDefaultPathQuery(pathQuery IPathQuery) {
	if pathQuery != nil {
		defaultPathQuery = pathQuery
	}
}

// EnablePerformanceMonitoring 启用性能监控
// EnablePerformanceMonitoring enables performance monitoring
func EnablePerformanceMonitoring() {
	GetGlobalMonitor().Enable()
}

// DisablePerformanceMonitoring 禁用性能监控
// DisablePerformanceMonitoring disables performance monitoring
func DisablePerformanceMonitoring() {
	GetGlobalMonitor().Disable()
}

// GetPerformanceStats 获取性能统计
// GetPerformanceStats gets performance statistics
func GetPerformanceStats() PerformanceStats {
	return GetGlobalMonitor().GetStats()
}

// ResetPerformanceStats 重置性能统计
// ResetPerformanceStats resets performance statistics
func ResetPerformanceStats() {
	GetGlobalMonitor().Reset()
}

// StartMemoryProfiling 开始内存分析
// StartMemoryProfiling starts memory profiling
func StartMemoryProfiling() {
	GetGlobalProfiler().Start()
}

// StopMemoryProfiling 停止内存分析
// StopMemoryProfiling stops memory profiling
func StopMemoryProfiling() {
	GetGlobalProfiler().Stop()
}

// GetMemorySnapshots 获取内存快照
// GetMemorySnapshots gets memory snapshots
func GetMemorySnapshots() []MemorySnapshot {
	return GetGlobalProfiler().GetSnapshots()
}

// GetLatestMemorySnapshot 获取最新内存快照
// GetLatestMemorySnapshot gets the latest memory snapshot
func GetLatestMemorySnapshot() *MemorySnapshot {
	return GetGlobalProfiler().GetLatestSnapshot()
}

// GetMemoryTrend 获取内存趋势
// GetMemoryTrend gets memory trend
func GetMemoryTrend() (trend string, growth float64) {
	return GetGlobalProfiler().GetMemoryTrend()
}

// 便捷的类型转换函数
// Convenient type conversion functions

// ToString 转换为字符串
// ToString converts to string
func ToString(value IValue) (string, error) {
	if value == nil {
		return "", NewTypeMismatchError(StringValueType, NullValueType, "")
	}
	if value.Type() == StringValueType {
		return value.String(), nil
	}
	if scalar, ok := value.(IScalarValue); ok {
		return scalar.String(), nil
	}
	return value.String(), nil
}

// ToInt 转换为整数
// ToInt converts to integer
func ToInt(value IValue) (int, error) {
	if scalar, ok := value.(IScalarValue); ok {
		return scalar.Int()
	}
	return 0, NewTypeMismatchError(NumberValueType, value.Type(), "")
}

// ToInt64 转换为64位整数
// ToInt64 converts to 64-bit integer
func ToInt64(value IValue) (int64, error) {
	if scalar, ok := value.(IScalarValue); ok {
		return scalar.Int64()
	}
	return 0, NewTypeMismatchError(NumberValueType, value.Type(), "")
}

// ToFloat64 转换为64位浮点数
// ToFloat64 converts to 64-bit float
func ToFloat64(value IValue) (float64, error) {
	if scalar, ok := value.(IScalarValue); ok {
		return scalar.Float64()
	}
	return 0, NewTypeMismatchError(NumberValueType, value.Type(), "")
}

// ToBool 转换为布尔值
// ToBool converts to boolean
func ToBool(value IValue) (bool, error) {
	if scalar, ok := value.(IScalarValue); ok {
		return scalar.Bool()
	}
	return false, NewTypeMismatchError(BoolValueType, value.Type(), "")
}

// ToTime 转换为时间
// ToTime converts to time
func ToTime(value IValue) (time.Time, error) {
	if scalar, ok := value.(IScalarValue); ok {
		return scalar.Time()
	}
	return time.Time{}, NewTypeMismatchError(StringValueType, value.Type(), "")
}

// ToBytes 转换为字节数组
// ToBytes converts to byte array
func ToBytes(value IValue) ([]byte, error) {
	if scalar, ok := value.(IScalarValue); ok {
		return scalar.Bytes()
	}
	return nil, NewTypeMismatchError(StringValueType, value.Type(), "")
}

// ToObject 转换为对象
// ToObject converts to object
func ToObject(value IValue) (IObject, error) {
	if obj, ok := value.(IObject); ok {
		return obj, nil
	}
	return nil, NewTypeMismatchError(ObjectValueType, value.Type(), "")
}

// ToArray 转换为数组
// ToArray converts to array
func ToArray(value IValue) (IArray, error) {
	if arr, ok := value.(IArray); ok {
		return arr, nil
	}
	return nil, NewTypeMismatchError(ArrayValueType, value.Type(), "")
}

// 便捷的Must版本转换函数
// Convenient Must version conversion functions

// MustToString 转换为字符串，失败则panic
// MustToString converts to string, panics on failure
func MustToString(value IValue) string {
	result, err := ToString(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustToInt 转换为整数，失败则panic
// MustToInt converts to integer, panics on failure
func MustToInt(value IValue) int {
	result, err := ToInt(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustToInt64 转换为64位整数，失败则panic
// MustToInt64 converts to 64-bit integer, panics on failure
func MustToInt64(value IValue) int64 {
	result, err := ToInt64(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustToFloat64 转换为64位浮点数，失败则panic
// MustToFloat64 converts to 64-bit float, panics on failure
func MustToFloat64(value IValue) float64 {
	result, err := ToFloat64(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustToBool 转换为布尔值，失败则panic
// MustToBool converts to boolean, panics on failure
func MustToBool(value IValue) bool {
	result, err := ToBool(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustToTime 转换为时间，失败则panic
// MustToTime converts to time, panics on failure
func MustToTime(value IValue) time.Time {
	result, err := ToTime(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustToBytes 转换为字节数组，失败则panic
// MustToBytes converts to byte array, panics on failure
func MustToBytes(value IValue) []byte {
	result, err := ToBytes(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustToObject 转换为对象，失败则panic
// MustToObject converts to object, panics on failure
func MustToObject(value IValue) IObject {
	result, err := ToObject(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustToArray 转换为数组，失败则panic
// MustToArray converts to array, panics on failure
func MustToArray(value IValue) IArray {
	result, err := ToArray(value)
	if err != nil {
		panic(err)
	}
	return result
}

// 版本信息
// Version information
const (
	Version      = "1.0.0"
	VersionMajor = 1
	VersionMinor = 0
	VersionPatch = 0
)

// GetVersion 获取版本信息
// GetVersion gets version information
func GetVersion() string {
	return Version
}
