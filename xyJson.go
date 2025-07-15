package xyJson

import (
	"time"
)

// 包级别的全局实例
// Package-level global instances
var (
	defaultFactory    IValueFactory
	defaultParser     IParser
	defaultSerializer ISerializer
	defaultPathQuery  IPathQuery
)

// init 初始化默认实例
// init initializes default instances
func init() {
	pool := NewObjectPool()
	defaultFactory = NewValueFactoryWithPool(pool)
	defaultParser = NewParserWithFactory(defaultFactory)
	defaultSerializer = NewSerializer()
	defaultPathQuery = NewPathQueryWithFactory(defaultFactory)
}

// Parse 解析JSON字节数组
// Parse parses JSON byte array
func Parse(data []byte) (IValue, error) {
	timer := GetGlobalMonitor().StartParseTimer()
	defer timer.End()

	return defaultParser.Parse(data)
}

// ParseString 解析JSON字符串
// ParseString parses JSON string
func ParseString(data string) (IValue, error) {
	timer := GetGlobalMonitor().StartParseTimer()
	defer timer.End()

	return defaultParser.ParseString(data)
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

// Serialize 序列化JSON值到字节数组
// Serialize serializes JSON value to byte array
func Serialize(value IValue) ([]byte, error) {
	timer := GetGlobalMonitor().StartSerializeTimer()
	defer timer.End()

	return defaultSerializer.Serialize(value)
}

// SerializeToString 序列化JSON值到字符串
// SerializeToString serializes JSON value to string
func SerializeToString(value IValue) (string, error) {
	timer := GetGlobalMonitor().StartSerializeTimer()
	defer timer.End()

	return defaultSerializer.SerializeToString(value)
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

// Get 根据路径获取值
// Get gets value by path
func Get(root IValue, path string) (IValue, error) {
	return defaultPathQuery.SelectOne(root, path)
}

// MustGet 根据路径获取值，如果失败则panic
// MustGet gets value by path, panics on failure
func MustGet(root IValue, path string) IValue {
	result, err := Get(root, path)
	if err != nil {
		panic(err)
	}
	return result
}

// GetAll 根据路径获取所有匹配的值
// GetAll gets all matching values by path
func GetAll(root IValue, path string) ([]IValue, error) {
	return defaultPathQuery.SelectAll(root, path)
}

// Set 根据路径设置值
// Set sets value by path
func Set(root IValue, path string, value IValue) error {
	return defaultPathQuery.Set(root, path, value)
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

// CreateNull 创建null值
// CreateNull creates a null value
func CreateNull() IValue {
	return defaultFactory.CreateNull()
}

// CreateString 创建字符串值
// CreateString creates a string value
func CreateString(value string) IValue {
	return defaultFactory.CreateString(value)
}

// CreateNumber 创建数字值
// CreateNumber creates a number value
func CreateNumber(value interface{}) (IValue, error) {
	return defaultFactory.CreateNumber(value)
}

// MustCreateNumber 创建数字值，如果失败则panic
// MustCreateNumber creates a number value, panics on failure
func MustCreateNumber(value interface{}) IValue {
	result, err := CreateNumber(value)
	if err != nil {
		panic(err)
	}
	return result
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
