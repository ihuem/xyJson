package xyJson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// 性能优化相关常量
// Performance optimization constants
const (
	// MaxStructDepth 最大结构体嵌套深度
	// MaxStructDepth maximum struct nesting depth
	MaxStructDepth = 100

	// StructCacheSize 结构体信息缓存大小
	// StructCacheSize struct info cache size
	StructCacheSize = 1000
	// ReflectValuePoolSize reflect.Value对象池大小
	// ReflectValuePoolSize is the size of reflect.Value object pool
	ReflectValuePoolSize = 100
	// VisitedMapPoolSize visited map对象池大小
	// VisitedMapPoolSize is the size of visited map object pool
	VisitedMapPoolSize = 50
)

// jsonTag JSON标签信息结构体
// jsonTag represents JSON tag information
type jsonTag struct {
	// Name 字段名称
	// Name field name
	Name string

	// OmitEmpty 是否忽略空值
	// OmitEmpty whether to omit empty values
	OmitEmpty bool

	// Skip 是否跳过该字段
	// Skip whether to skip this field
	Skip bool

	// AsString 是否强制转换为字符串
	// AsString whether to force convert to string
	AsString bool
}

// fieldInfo 字段信息结构体
// fieldInfo represents field information
type fieldInfo struct {
	// Index 字段索引
	// Index field index
	Index int

	// Name 字段名称
	// Name field name
	Name string

	// Type 字段类型
	// Type field type
	Type reflect.Type

	// Tag JSON标签信息
	// Tag JSON tag information
	Tag jsonTag

	// IsPtr 是否为指针类型
	// IsPtr whether it's a pointer type
	IsPtr bool
}

// structInfo 结构体信息缓存
// structInfo represents cached struct information
type structInfo struct {
	// Fields 字段映射表
	// Fields field mapping
	Fields map[string]*fieldInfo
}

// 全局结构体信息缓存
// Global struct info cache
var (
	structCache = make(map[reflect.Type]*structInfo)
	cacheMutex  sync.RWMutex

	// 对象池用于复用常用对象
	// Object pools for reusing common objects
	visitedMapPool = sync.Pool{
		New: func() interface{} {
			return make(map[IValue]bool, 16)
		},
	}

	// 预分配的reflect.Type缓存
	// Pre-allocated reflect.Type cache
	timeType    = reflect.TypeOf(time.Time{})
	stringType  = reflect.TypeOf("")
	boolType    = reflect.TypeOf(true)
	int64Type   = reflect.TypeOf(int64(0))
	float64Type = reflect.TypeOf(float64(0))
)

// serializer JSON序列化器实现
// serializer implements the JSON serializer
type serializer struct {
	options *SerializeOptions
}

// NewSerializer 创建新的JSON序列化器
// NewSerializer creates a new JSON serializer
func NewSerializer() ISerializer {
	return &serializer{
		options: &SerializeOptions{
			Indent:        "",
			EscapeHTML:    true,
			EscapeUnicode: false,
			SortKeys:      false,
			Compact:       false,
			MaxDepth:      DefaultMaxDepth,
		},
	}
}

// NewSerializerWithOptions 使用指定选项创建JSON序列化器
// NewSerializerWithOptions creates a JSON serializer with specified options
func NewSerializerWithOptions(options *SerializeOptions) ISerializer {
	if options == nil {
		options = &SerializeOptions{
			Indent:        "",
			EscapeHTML:    true,
			EscapeUnicode: false,
			SortKeys:      false,
			Compact:       false,
			MaxDepth:      DefaultMaxDepth,
		}
	}
	return &serializer{
		options: options,
	}
}

// Serialize 序列化JSON值到字节数组
// Serialize serializes JSON value to byte array
func (s *serializer) Serialize(value IValue) ([]byte, error) {
	if value == nil {
		return nil, NewInvalidJSONError("cannot serialize nil value", nil)
	}

	var buf bytes.Buffer
	visited := make(map[IValue]bool)
	err := s.serializeValue(value, &buf, 0, visited)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// SerializeToString 序列化JSON值到字符串
// SerializeToString serializes JSON value to string
func (s *serializer) SerializeToString(value IValue) (string, error) {
	data, err := s.Serialize(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetOptions 设置序列化选项
// SetOptions sets serialization options
func (s *serializer) SetOptions(options *SerializeOptions) {
	if options != nil {
		s.options = options
	}
}

// GetOptions 获取序列化选项
// GetOptions gets serialization options
func (s *serializer) GetOptions() *SerializeOptions {
	return s.options
}

// SerializeToStruct 将JSON值序列化到结构体
// SerializeToStruct serializes JSON value to struct
func (s *serializer) SerializeToStruct(value IValue, target interface{}) error {
	if value == nil {
		return NewInvalidJSONError("cannot serialize nil value to struct", nil)
	}

	if target == nil {
		return NewNullPointerError("target cannot be nil")
	}

	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr {
		return NewJSONError(ErrInvalidOperation, "target must be a pointer", nil)
	}

	if rv.IsNil() {
		return NewNullPointerError("target pointer cannot be nil")
	}

	elem := rv.Elem()
	if elem.Kind() != reflect.Struct {
		return NewJSONError(ErrInvalidOperation, "target must be a pointer to struct", nil)
	}

	// 从对象池获取visited map
	visited := visitedMapPool.Get().(map[IValue]bool)
	defer func() {
		// 清空并归还到对象池
		for k := range visited {
			delete(visited, k)
		}
		visitedMapPool.Put(visited)
	}()

	return s.mapValueToStruct(value, elem, visited, 0)
}

// MustSerializeToStruct 将JSON值序列化到结构体，失败时panic
// MustSerializeToStruct serializes JSON value to struct, panics on failure
func (s *serializer) MustSerializeToStruct(value IValue, target interface{}) {
	if err := s.SerializeToStruct(value, target); err != nil {
		panic(err)
	}
}

// UnmarshalToStructFast 快速解析JSON字节数组到Go结构体（跳过IValue中间表示）
// UnmarshalToStructFast fast parses JSON byte array to Go struct (skips IValue intermediate representation)
func (s *serializer) UnmarshalToStructFast(data []byte, target interface{}) error {
	if target == nil {
		return NewNullPointerError("target cannot be nil")
	}

	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr {
		return NewInvalidJSONError("target must be a pointer", nil)
	}

	rv = rv.Elem()
	if !rv.CanSet() {
		return NewInvalidJSONError("target must be settable", nil)
	}

	// 使用官方json包进行快速解析
	// Use official json package for fast parsing
	return json.Unmarshal(data, target)
}

// UnmarshalToStructCustom 使用自定义解析器解析JSON到结构体（不依赖官方包）
// UnmarshalToStructCustom unmarshal JSON to struct using custom parser (no official package dependency)
func (s *serializer) UnmarshalToStructCustom(data []byte, target interface{}) error {
	parser := NewCustomParser()
	return parser.UnmarshalDirect(data, target)
}

// UnmarshalStringToStructCustom 使用自定义解析器解析JSON字符串到结构体
// UnmarshalStringToStructCustom unmarshal JSON string to struct using custom parser
func (s *serializer) UnmarshalStringToStructCustom(data string, target interface{}) error {
	parser := NewCustomParser()
	return parser.UnmarshalDirectString(data, target)
}

// MustUnmarshalToStructCustom 使用自定义解析器解析JSON到结构体（panic版本）
// MustUnmarshalToStructCustom unmarshal JSON to struct using custom parser (panic version)
func (s *serializer) MustUnmarshalToStructCustom(data []byte, target interface{}) {
	if err := s.UnmarshalToStructCustom(data, target); err != nil {
		panic(err)
	}
}

// MustUnmarshalStringToStructCustom 使用自定义解析器解析JSON字符串到结构体（panic版本）
// MustUnmarshalStringToStructCustom unmarshal JSON string to struct using custom parser (panic version)
func (s *serializer) MustUnmarshalStringToStructCustom(data string, target interface{}) {
	if err := s.UnmarshalStringToStructCustom(data, target); err != nil {
		panic(err)
	}
}

// serializeValue 序列化值的内部实现
// serializeValue internal implementation for serializing values
func (s *serializer) serializeValue(value IValue, buf *bytes.Buffer, depth int, visited map[IValue]bool) error {
	if value == nil {
		buf.WriteString("null")
		return nil
	}

	// 检查最大深度
	if depth > s.options.MaxDepth {
		return NewInvalidJSONError("maximum serialization depth exceeded", nil)
	}

	// 检查循环引用
	if visited[value] {
		return NewInvalidJSONError("circular reference detected", nil)
	}
	visited[value] = true
	defer delete(visited, value)

	switch value.Type() {
	case NullValueType:
		buf.WriteString("null")
	case StringValueType:
		return s.serializeString(value.String(), buf)
	case NumberValueType:
		return s.serializeNumber(value, buf)
	case BoolValueType:
		if scalar, ok := value.(IScalarValue); ok {
			boolVal, err := scalar.Bool()
			if err != nil {
				return err
			}
			if boolVal {
				buf.WriteString("true")
			} else {
				buf.WriteString("false")
			}
		} else {
			return NewTypeMismatchError(BoolValueType, value.Type(), "")
		}
	case ObjectValueType:
		return s.serializeObject(value.(IObject), buf, depth, visited)
	case ArrayValueType:
		return s.serializeArray(value.(IArray), buf, depth, visited)
	default:
		return NewInvalidJSONError("unknown value type", nil)
	}

	return nil
}

// serializeString 序列化字符串
// serializeString serializes a string
func (s *serializer) serializeString(str string, buf *bytes.Buffer) error {
	buf.WriteByte('"')

	// 标准转义
	for _, r := range str {
		switch r {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\b':
			buf.WriteString(`\b`)
		case '\f':
			buf.WriteString(`\f`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		case '<', '>', '&':
			if s.options.EscapeHTML {
				buf.WriteString(fmt.Sprintf("\\u%04x", r))
			} else {
				buf.WriteRune(r)
			}
		default:
			if r < 0x20 || r == 0x7f {
				// 控制字符需要转义
				buf.WriteString(fmt.Sprintf("\\u%04x", r))
			} else if r > 0x7f && s.options.EscapeUnicode {
				// 非ASCII字符在Unicode转义模式下需要转义
				buf.WriteString(fmt.Sprintf("\\u%04x", r))
			} else {
				buf.WriteRune(r)
			}
		}
	}

	buf.WriteByte('"')
	return nil
}

// serializeNumber 序列化数字
// serializeNumber serializes a number
func (s *serializer) serializeNumber(value IValue, buf *bytes.Buffer) error {
	scalar, ok := value.(IScalarValue)
	if !ok {
		return NewTypeMismatchError(NumberValueType, value.Type(), "")
	}

	// 尝试获取整数
	if intVal, err := scalar.Int64(); err == nil {
		buf.WriteString(strconv.FormatInt(intVal, 10))
		return nil
	}

	// 获取浮点数
	floatVal, err := scalar.Float64()
	if err != nil {
		return err
	}

	// 检查特殊值
	if math.IsNaN(floatVal) { // NaN
		buf.WriteString("null")
		return nil
	}
	if math.IsInf(floatVal, 0) { // Infinity
		buf.WriteString("null")
		return nil
	}

	buf.WriteString(strconv.FormatFloat(floatVal, 'g', -1, 64))

	return nil
}

// serializeObject 序列化对象
// serializeObject serializes an object
func (s *serializer) serializeObject(obj IObject, buf *bytes.Buffer, depth int, visited map[IValue]bool) error {
	buf.WriteByte('{')

	keys := obj.Keys()
	if len(keys) == 0 {
		buf.WriteByte('}')
		return nil
	}

	// 排序键
	if s.options.SortKeys {
		sort.Strings(keys)
	}

	first := true
	for _, key := range keys {
		value := obj.Get(key)
		if value == nil {
			continue
		}

		if !first {
			buf.WriteByte(',')
		}
		first = false

		// 添加缩进
		if s.options.Indent != "" && !s.options.Compact {
			buf.WriteByte('\n')
			for i := 0; i <= depth; i++ {
				buf.WriteString(s.options.Indent)
			}
		}

		// 序列化键
		if err := s.serializeString(key, buf); err != nil {
			return err
		}

		buf.WriteByte(':')
		if s.options.Indent != "" && !s.options.Compact {
			buf.WriteByte(' ')
		}

		// 序列化值
		if err := s.serializeValue(value, buf, depth+1, visited); err != nil {
			return err
		}
	}

	// 添加结束缩进
	if s.options.Indent != "" && !s.options.Compact && !first {
		buf.WriteByte('\n')
		for i := 0; i < depth; i++ {
			buf.WriteString(s.options.Indent)
		}
	}

	buf.WriteByte('}')
	return nil
}

// serializeArray 序列化数组
// serializeArray serializes an array
func (s *serializer) serializeArray(arr IArray, buf *bytes.Buffer, depth int, visited map[IValue]bool) error {
	buf.WriteByte('[')

	length := arr.Length()
	if length == 0 {
		buf.WriteByte(']')
		return nil
	}

	first := true
	for i := 0; i < length; i++ {
		value := arr.Get(i)
		if value == nil {
			continue
		}

		if !first {
			buf.WriteByte(',')
		}
		first = false

		// 添加缩进
		if s.options.Indent != "" && !s.options.Compact {
			buf.WriteByte('\n')
			for j := 0; j <= depth; j++ {
				buf.WriteString(s.options.Indent)
			}
		}

		// 序列化值
		if err := s.serializeValue(value, buf, depth+1, visited); err != nil {
			return err
		}
	}

	// 添加结束缩进
	if s.options.Indent != "" && !s.options.Compact && !first {
		buf.WriteByte('\n')
		for i := 0; i < depth; i++ {
			buf.WriteString(s.options.Indent)
		}
	}

	buf.WriteByte(']')
	return nil
}

// CompactSerializer 创建紧凑序列化器
// CompactSerializer creates a compact serializer
func CompactSerializer() ISerializer {
	return NewSerializerWithOptions(&SerializeOptions{
		Indent:        "",
		EscapeHTML:    true,
		EscapeUnicode: false,
		SortKeys:      false,
		Compact:       true,
		MaxDepth:      DefaultMaxDepth,
	})
}

// PrettySerializer 创建格式化序列化器
// PrettySerializer creates a pretty serializer
func PrettySerializer(indent string) ISerializer {
	if indent == "" {
		indent = DefaultIndent
	}
	return NewSerializerWithOptions(&SerializeOptions{
		Indent:        indent,
		EscapeHTML:    true,
		EscapeUnicode: false,
		SortKeys:      true,
		Compact:       false,
		MaxDepth:      DefaultMaxDepth,
	})
}

// HTMLSafeSerializer 创建HTML安全序列化器
// HTMLSafeSerializer creates an HTML-safe serializer
func HTMLSafeSerializer() ISerializer {
	return NewSerializerWithOptions(&SerializeOptions{
		Indent:        "",
		EscapeHTML:    true,
		EscapeUnicode: true,
		SortKeys:      false,
		Compact:       false,
		MaxDepth:      DefaultMaxDepth,
	})
}

// MinimalSerializer 创建最小化序列化器
// MinimalSerializer creates a minimal serializer
func MinimalSerializer() ISerializer {
	return NewSerializerWithOptions(&SerializeOptions{
		Indent:        "",
		EscapeHTML:    false,
		EscapeUnicode: false,
		SortKeys:      false,
		Compact:       true,
		MaxDepth:      DefaultMaxDepth,
	})
}

// escapeStringForHTML HTML转义字符串
// escapeStringForHTML escapes string for HTML
func escapeStringForHTML(s string) string {
	var buf strings.Builder
	for _, r := range s {
		switch r {
		case '<':
			buf.WriteString("\\u003c")
		case '>':
			buf.WriteString("\\u003e")
		case '&':
			buf.WriteString("\\u0026")
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\b':
			buf.WriteString(`\b`)
		case '\f':
			buf.WriteString(`\f`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		default:
			if r < 0x20 || r == 0x7f {
				buf.WriteString(fmt.Sprintf("\\u%04x", r))
			} else {
				buf.WriteRune(r)
			}
		}
	}
	return buf.String()
}

// isValidUTF8 检查字符串是否为有效的UTF-8
// isValidUTF8 checks if string is valid UTF-8
func isValidUTF8(s string) bool {
	return utf8.ValidString(s)
}

// mapValueToStruct 将IValue映射到结构体
// mapValueToStruct maps IValue to struct
func (s *serializer) mapValueToStruct(value IValue, rv reflect.Value, visited map[IValue]bool, depth int) error {
	if depth > MaxStructDepth {
		return NewInvalidJSONError("maximum struct depth exceeded", nil)
	}

	if visited[value] {
		return NewInvalidJSONError("circular reference detected", nil)
	}
	visited[value] = true
	defer delete(visited, value)

	switch value.Type() {
	case ObjectValueType:
		return s.mapObjectToStruct(value.(IObject), rv, visited, depth)
	case ArrayValueType:
		return s.mapArrayToStruct(value.(IArray), rv, visited, depth)
	default:
		return s.mapScalarToStruct(value, rv)
	}
}

// mapObjectToStruct 将IObject映射到结构体
// mapObjectToStruct maps IObject to struct
func (s *serializer) mapObjectToStruct(obj IObject, rv reflect.Value, visited map[IValue]bool, depth int) error {
	if rv.Kind() != reflect.Struct {
		return NewTypeMismatchError(ObjectValueType, ValueType(rv.Kind()), "")
	}

	structInfo := getStructInfo(rv.Type())

	// 遍历JSON对象的所有字段
	var lastErr error
	obj.Range(func(key string, value IValue) bool {
		fieldInfo, exists := structInfo.Fields[key]
		if !exists || fieldInfo.Tag.Skip {
			return true // 继续遍历
		}

		fieldValue := rv.Field(fieldInfo.Index)
		if !fieldValue.CanSet() {
			return true // 跳过不可设置的字段
		}

		if err := s.setFieldValue(fieldValue, value, fieldInfo, visited, depth+1); err != nil {
			lastErr = err
			return false // 停止遍历
		}

		return true
	})

	return lastErr
}

// mapArrayToStruct 将IArray映射到结构体（通常不支持，除非是特殊情况）
// mapArrayToStruct maps IArray to struct (usually not supported)
func (s *serializer) mapArrayToStruct(arr IArray, rv reflect.Value, visited map[IValue]bool, depth int) error {
	switch rv.Kind() {
	case reflect.Slice:
		return s.mapArrayToSlice(arr, rv, visited, depth)
	case reflect.Array:
		return s.mapArrayToArray(arr, rv, visited, depth)
	default:
		return NewTypeMismatchError(ArrayValueType, ValueType(rv.Kind()), "")
	}
}

// mapArrayToSlice 将IArray映射到切片
// mapArrayToSlice maps IArray to slice
func (s *serializer) mapArrayToSlice(arr IArray, rv reflect.Value, visited map[IValue]bool, depth int) error {
	elemType := rv.Type().Elem()
	length := arr.Length()

	// 创建新的切片
	slice := reflect.MakeSlice(rv.Type(), length, length)

	for i := 0; i < length; i++ {
		value := arr.Get(i)
		if value == nil {
			continue
		}

		elemValue := slice.Index(i)
		if err := s.setValueByType(elemValue, value, elemType, visited, depth+1); err != nil {
			return err
		}
	}

	rv.Set(slice)
	return nil
}

// mapArrayToArray 将IArray映射到数组
// mapArrayToArray maps IArray to array
func (s *serializer) mapArrayToArray(arr IArray, rv reflect.Value, visited map[IValue]bool, depth int) error {
	elemType := rv.Type().Elem()
	length := arr.Length()
	arrayLen := rv.Len()

	// 确保不超出数组长度
	if length > arrayLen {
		length = arrayLen
	}

	for i := 0; i < length; i++ {
		value := arr.Get(i)
		if value == nil {
			continue
		}

		elemValue := rv.Index(i)
		if err := s.setValueByType(elemValue, value, elemType, visited, depth+1); err != nil {
			return err
		}
	}

	return nil
}

// mapScalarToStruct 将标量值映射到结构体（通常不支持）
// mapScalarToStruct maps scalar value to struct (usually not supported)
func (s *serializer) mapScalarToStruct(value IValue, rv reflect.Value) error {
	return NewTypeMismatchError(value.Type(), ValueType(rv.Kind()), "")
}

// setFieldValue 设置字段值
// setFieldValue sets field value
func (s *serializer) setFieldValue(fieldValue reflect.Value, value IValue, fieldInfo *fieldInfo, visited map[IValue]bool, depth int) error {
	// 处理指针类型
	if fieldInfo.IsPtr {
		if value.IsNull() {
			fieldValue.Set(reflect.Zero(fieldValue.Type()))
			return nil
		}

		// 创建新的指针值
		newPtr := reflect.New(fieldValue.Type().Elem())
		if err := s.setValueByType(newPtr.Elem(), value, fieldValue.Type().Elem(), visited, depth); err != nil {
			return err
		}
		fieldValue.Set(newPtr)
		return nil
	}

	return s.setValueByType(fieldValue, value, fieldInfo.Type, visited, depth)
}

// setValueByType 根据类型设置值
// setValueByType sets value by type
func (s *serializer) setValueByType(rv reflect.Value, value IValue, targetType reflect.Type, visited map[IValue]bool, depth int) error {
	if value.IsNull() {
		rv.Set(reflect.Zero(targetType))
		return nil
	}

	// 快速路径：使用预缓存的类型比较
	if targetType == timeType {
		return s.setTimeValue(rv, value)
	}

	kind := targetType.Kind()
	valueType := value.Type()

	// 优化的类型匹配
	switch kind {
	case reflect.String:
		if valueType == StringValueType {
			rv.SetString(value.AsString())
			return nil
		}
		// 回退到通用方法
		return s.setStringValue(rv, value)

	case reflect.Bool:
		if valueType == BoolValueType {
			rv.SetBool(value.AsBool())
			return nil
		}
		return NewTypeMismatchError(valueType, BoolValueType, "")

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if valueType == NumberValueType {
			return s.setIntValueFast(rv, value, kind)
		}
		return NewTypeMismatchError(valueType, NumberValueType, "")

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if valueType == NumberValueType {
			return s.setUintValueFast(rv, value, kind)
		}
		return NewTypeMismatchError(valueType, NumberValueType, "")

	case reflect.Float32, reflect.Float64:
		if valueType == NumberValueType {
			rv.SetFloat(value.AsFloat64())
			return nil
		}
		return NewTypeMismatchError(valueType, NumberValueType, "")

	case reflect.Slice:
		if valueType == ArrayValueType {
			return s.mapArrayToSlice(value.(IArray), rv, visited, depth)
		}
		return NewTypeMismatchError(valueType, ArrayValueType, "")

	case reflect.Array:
		if valueType == ArrayValueType {
			return s.mapArrayToArray(value.(IArray), rv, visited, depth)
		}
		return NewTypeMismatchError(valueType, ArrayValueType, "")

	case reflect.Map:
		return s.setMapValue(rv, value, targetType, visited, depth)

	case reflect.Struct:
		if valueType == ObjectValueType {
			return s.mapObjectToStruct(value.(IObject), rv, visited, depth)
		}
		return NewTypeMismatchError(valueType, ObjectValueType, "")

	case reflect.Ptr:
		if value.IsNull() {
			rv.Set(reflect.Zero(targetType))
			return nil
		}
		newPtr := reflect.New(targetType.Elem())
		if err := s.setValueByType(newPtr.Elem(), value, targetType.Elem(), visited, depth); err != nil {
			return err
		}
		rv.Set(newPtr)
		return nil

	case reflect.Interface:
		// 设置为原始值
		rv.Set(reflect.ValueOf(value.Raw()))
		return nil

	default:
		return NewJSONError(ErrTypeMismatch, fmt.Sprintf("unsupported type: %s", kind), nil)
	}
}

// setStringValue 设置字符串值
// setStringValue sets string value
func (s *serializer) setStringValue(rv reflect.Value, value IValue) error {
	if value.Type() == StringValueType {
		rv.SetString(value.AsString())
		return nil
	}
	// 尝试转换其他类型到字符串
	rv.SetString(value.AsString())
	return nil
}

// setBoolValue 设置布尔值
// setBoolValue sets boolean value
func (s *serializer) setBoolValue(rv reflect.Value, value IValue) error {
	if value.Type() == BoolValueType {
		rv.SetBool(value.AsBool())
		return nil
	}
	return NewTypeMismatchError(value.Type(), BoolValueType, "")
}

// setIntValue 设置整数值
// setIntValue sets integer value
func (s *serializer) setIntValue(rv reflect.Value, value IValue, targetType reflect.Type) error {
	if value.Type() == NumberValueType {
		intVal := value.AsInt64()
		// 检查范围
		switch targetType.Kind() {
		case reflect.Int8:
			if intVal < -128 || intVal > 127 {
				return NewJSONError(ErrTypeMismatch, "value out of int8 range", nil)
			}
		case reflect.Int16:
			if intVal < -32768 || intVal > 32767 {
				return NewJSONError(ErrTypeMismatch, "value out of int16 range", nil)
			}
		case reflect.Int32:
			if intVal < -2147483648 || intVal > 2147483647 {
				return NewJSONError(ErrTypeMismatch, "value out of int32 range", nil)
			}
		}
		rv.SetInt(intVal)
		return nil
	}
	return NewTypeMismatchError(value.Type(), NumberValueType, "")
}

// setIntValueFast 快速设置整数值（已知类型匹配）
// setIntValueFast sets integer value fast (type already matched)
func (s *serializer) setIntValueFast(rv reflect.Value, value IValue, kind reflect.Kind) error {
	intVal := value.AsInt64()
	// 优化的范围检查
	switch kind {
	case reflect.Int8:
		if intVal < -128 || intVal > 127 {
			return NewJSONError(ErrTypeMismatch, "value out of int8 range", nil)
		}
	case reflect.Int16:
		if intVal < -32768 || intVal > 32767 {
			return NewJSONError(ErrTypeMismatch, "value out of int16 range", nil)
		}
	case reflect.Int32:
		if intVal < -2147483648 || intVal > 2147483647 {
			return NewJSONError(ErrTypeMismatch, "value out of int32 range", nil)
		}
	}
	rv.SetInt(intVal)
	return nil
}

// setUintValue 设置无符号整数值
// setUintValue sets unsigned integer value
func (s *serializer) setUintValue(rv reflect.Value, value IValue, targetType reflect.Type) error {
	if value.Type() == NumberValueType {
		intVal := value.AsInt64()
		if intVal < 0 {
			return NewJSONError(ErrTypeMismatch, "negative value for unsigned integer", nil)
		}
		uintVal := uint64(intVal)
		// 检查范围
		switch targetType.Kind() {
		case reflect.Uint8:
			if uintVal > 255 {
				return NewJSONError(ErrTypeMismatch, "value out of uint8 range", nil)
			}
		case reflect.Uint16:
			if uintVal > 65535 {
				return NewJSONError(ErrTypeMismatch, "value out of uint16 range", nil)
			}
		case reflect.Uint32:
			if uintVal > 4294967295 {
				return NewJSONError(ErrTypeMismatch, "value out of uint32 range", nil)
			}
		}
		rv.SetUint(uintVal)
		return nil
	}
	return NewTypeMismatchError(value.Type(), NumberValueType, "")
}

// setUintValueFast 快速设置无符号整数值（已知类型匹配）
// setUintValueFast sets unsigned integer value fast (type already matched)
func (s *serializer) setUintValueFast(rv reflect.Value, value IValue, kind reflect.Kind) error {
	intVal := value.AsInt64()
	if intVal < 0 {
		return NewJSONError(ErrTypeMismatch, "negative value for unsigned integer", nil)
	}
	uintVal := uint64(intVal)
	// 优化的范围检查
	switch kind {
	case reflect.Uint8:
		if uintVal > 255 {
			return NewJSONError(ErrTypeMismatch, "value out of uint8 range", nil)
		}
	case reflect.Uint16:
		if uintVal > 65535 {
			return NewJSONError(ErrTypeMismatch, "value out of uint16 range", nil)
		}
	case reflect.Uint32:
		if uintVal > 4294967295 {
			return NewJSONError(ErrTypeMismatch, "value out of uint32 range", nil)
		}
	}
	rv.SetUint(uintVal)
	return nil
}

// setFloatValue 设置浮点数值
// setFloatValue sets float value
func (s *serializer) setFloatValue(rv reflect.Value, value IValue, targetType reflect.Type) error {
	if value.Type() == NumberValueType {
		floatVal := value.AsFloat64()
		rv.SetFloat(floatVal)
		return nil
	}
	return NewTypeMismatchError(value.Type(), NumberValueType, "")
}

// setTimeValue 设置时间值
// setTimeValue sets time value
func (s *serializer) setTimeValue(rv reflect.Value, value IValue) error {
	if value.Type() == StringValueType {
		timeStr := value.AsString()
		// 尝试多种时间格式
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05Z",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}

		for _, format := range formats {
			if t, err := time.Parse(format, timeStr); err == nil {
				rv.Set(reflect.ValueOf(t))
				return nil
			}
		}
		return NewJSONError(ErrTypeMismatch, "invalid time format", nil)
	}
	return NewTypeMismatchError(value.Type(), StringValueType, "")
}

// setMapValue 设置Map值
// setMapValue sets map value
func (s *serializer) setMapValue(rv reflect.Value, value IValue, targetType reflect.Type, visited map[IValue]bool, depth int) error {
	if value.Type() != ObjectValueType {
		return NewTypeMismatchError(value.Type(), ObjectValueType, "")
	}

	obj := value.(IObject)
	keyType := targetType.Key()
	valueType := targetType.Elem()

	// 只支持string类型的key
	if keyType.Kind() != reflect.String {
		return NewJSONError(ErrTypeMismatch, "map key must be string", nil)
	}

	// 创建新的map
	newMap := reflect.MakeMap(targetType)

	var lastErr error
	obj.Range(func(key string, val IValue) bool {
		mapValue := reflect.New(valueType).Elem()
		if err := s.setValueByType(mapValue, val, valueType, visited, depth+1); err != nil {
			lastErr = err
			return false
		}
		newMap.SetMapIndex(reflect.ValueOf(key), mapValue)
		return true
	})

	if lastErr != nil {
		return lastErr
	}

	rv.Set(newMap)
	return nil
}

// parseJSONTag 解析JSON标签
// parseJSONTag parses JSON tag
func parseJSONTag(tag string) jsonTag {
	if tag == "" {
		return jsonTag{}
	}

	if tag == "-" {
		return jsonTag{Skip: true}
	}

	parts := strings.Split(tag, ",")
	result := jsonTag{Name: parts[0]}

	for _, opt := range parts[1:] {
		switch strings.TrimSpace(opt) {
		case "omitempty":
			result.OmitEmpty = true
		case "string":
			result.AsString = true
		}
	}

	return result
}

// getStructInfo 获取或创建结构体信息
// getStructInfo gets or creates struct info
func getStructInfo(t reflect.Type) *structInfo {
	cacheMutex.RLock()
	if info, exists := structCache[t]; exists {
		cacheMutex.RUnlock()
		return info
	}
	cacheMutex.RUnlock()

	// 创建新的结构体信息
	return createStructInfo(t)
}

// createStructInfo 创建结构体信息
// createStructInfo creates struct info
func createStructInfo(t reflect.Type) *structInfo {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// 双重检查
	if info, exists := structCache[t]; exists {
		return info
	}

	info := &structInfo{
		Fields: make(map[string]*fieldInfo),
	}

	// 遍历所有字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 跳过非导出字段
		if !field.IsExported() {
			continue
		}

		// 解析JSON标签
		tag := parseJSONTag(field.Tag.Get("json"))

		// 跳过标记为忽略的字段
		if tag.Skip {
			continue
		}

		// 确定字段名
		fieldName := field.Name
		if tag.Name != "" {
			fieldName = tag.Name
		}

		// 检查是否为指针类型
		isPtr := field.Type.Kind() == reflect.Ptr
		fieldType := field.Type
		if isPtr {
			fieldType = field.Type.Elem()
		}

		info.Fields[fieldName] = &fieldInfo{
			Index: i,
			Name:  fieldName,
			Type:  fieldType,
			Tag:   tag,
			IsPtr: isPtr,
		}
	}

	// 缓存结构体信息
	if len(structCache) < StructCacheSize {
		structCache[t] = info
	}

	return info
}
