package xyJson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// JSONSerializer JSON序列化器
type JSONSerializer struct {
	config SerializerConfig
}

// NewJSONSerializer 创建新的JSON序列化器
func NewJSONSerializer() ISerializer {
	return &JSONSerializer{
		config: GetGlobalConfig().Serializer,
	}
}

// NewJSONSerializerWithConfig 使用指定配置创建JSON序列化器
func NewJSONSerializerWithConfig(config SerializerConfig) ISerializer {
	return &JSONSerializer{
		config: config,
	}
}

// Serialize 序列化为字节数组
func (s *JSONSerializer) Serialize(value IValue) ([]byte, error) {
	if value == nil {
		return []byte("null"), nil
	}
	
	var buf bytes.Buffer
	err := s.serializeValue(value, &buf, 0)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// SerializeToString 序列化为字符串
func (s *JSONSerializer) SerializeToString(value IValue) (string, error) {
	data, err := s.Serialize(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Pretty 格式化输出
func (s *JSONSerializer) Pretty(value IValue) (string, error) {
	originalIndent := s.config.Indent
	originalCompact := s.config.CompactOutput
	
	s.config.Indent = "  "
	s.config.CompactOutput = false
	
	result, err := s.SerializeToString(value)
	
	s.config.Indent = originalIndent
	s.config.CompactOutput = originalCompact
	
	return result, err
}

// Compact 压缩输出
func (s *JSONSerializer) Compact(value IValue) (string, error) {
	originalIndent := s.config.Indent
	originalCompact := s.config.CompactOutput
	
	s.config.Indent = ""
	s.config.CompactOutput = true
	
	result, err := s.SerializeToString(value)
	
	s.config.Indent = originalIndent
	s.config.CompactOutput = originalCompact
	
	return result, err
}

// serializeValue 序列化值
func (s *JSONSerializer) serializeValue(value IValue, buf *bytes.Buffer, depth int) error {
	if value == nil {
		buf.WriteString("null")
		return nil
	}
	
	switch value.Type() {
	case NullValueType:
		buf.WriteString("null")
	case BoolValueType:
		if value.Raw().(bool) {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case StringValueType:
		s.serializeString(value.Raw().(string), buf)
	case NumberValueType:
		s.serializeNumber(value.Raw().(float64), buf)
	case ObjectValueType:
		return s.serializeObject(value.(IObject), buf, depth)
	case ArrayValueType:
		return s.serializeArray(value.(IArray), buf, depth)
	default:
		return NewSerializationError("unsupported value type", value.Type().String())
	}
	return nil
}

// serializeString 序列化字符串
func (s *JSONSerializer) serializeString(str string, buf *bytes.Buffer) {
	buf.WriteByte('"')
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
		default:
			if r < 32 || (s.config.EscapeHTML && (r == '<' || r == '>' || r == '&')) {
				buf.WriteString(fmt.Sprintf(`\u%04x`, r))
			} else {
				buf.WriteRune(r)
			}
		}
	}
	buf.WriteByte('"')
}

// serializeNumber 序列化数字
func (s *JSONSerializer) serializeNumber(num float64, buf *bytes.Buffer) {
	buf.WriteString(strconv.FormatFloat(num, 'g', -1, 64))
}

// serializeObject 序列化对象
func (s *JSONSerializer) serializeObject(obj IObject, buf *bytes.Buffer, depth int) error {
	buf.WriteByte('{')
	
	var pairs []KeyValuePair
	if s.config.SortKeys {
		pairs = obj.SortedPairs()
	} else {
		pairs = make([]KeyValuePair, 0, obj.Size())
		obj.Range(func(key string, value IValue) bool {
			pairs = append(pairs, KeyValuePair{Key: key, Value: value})
			return true
		})
	}
	
	for i, pair := range pairs {
		if i > 0 {
			buf.WriteByte(',')
		}
		
		if !s.config.CompactOutput && s.config.Indent != "" {
			buf.WriteByte('\n')
			buf.WriteString(strings.Repeat(s.config.Indent, depth+1))
		}
		
		s.serializeString(pair.Key, buf)
		buf.WriteByte(':')
		
		if !s.config.CompactOutput {
			buf.WriteByte(' ')
		}
		
		err := s.serializeValue(pair.Value, buf, depth+1)
		if err != nil {
			return err
		}
	}
	
	if len(pairs) > 0 && !s.config.CompactOutput && s.config.Indent != "" {
		buf.WriteByte('\n')
		buf.WriteString(strings.Repeat(s.config.Indent, depth))
	}
	
	buf.WriteByte('}')
	return nil
}

// serializeArray 序列化数组
func (s *JSONSerializer) serializeArray(arr IArray, buf *bytes.Buffer, depth int) error {
	buf.WriteByte('[')
	
	length := arr.Length()
	for i := 0; i < length; i++ {
		if i > 0 {
			buf.WriteByte(',')
		}
		
		if !s.config.CompactOutput && s.config.Indent != "" {
			buf.WriteByte('\n')
			buf.WriteString(strings.Repeat(s.config.Indent, depth+1))
		}
		
		value, _ := arr.Get(i)
		err := s.serializeValue(value, buf, depth+1)
		if err != nil {
			return err
		}
	}
	
	if length > 0 && !s.config.CompactOutput && s.config.Indent != "" {
		buf.WriteByte('\n')
		buf.WriteString(strings.Repeat(s.config.Indent, depth))
	}
	
	buf.WriteByte(']')
	return nil
}
