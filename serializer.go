package xyJson

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
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
