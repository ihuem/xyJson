package xyJson

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// JSONParser JSON解析器
type JSONParser struct {
	config ParserConfig
}

// NewJSONParser 创建新的JSON解析器
func NewJSONParser() IParser {
	return &JSONParser{
		config: GetGlobalConfig().Parser,
	}
}

// NewJSONParserWithConfig 使用指定配置创建JSON解析器
func NewJSONParserWithConfig(config ParserConfig) IParser {
	return &JSONParser{
		config: config,
	}
}

// Parse 解析JSON字节数组
func (p *JSONParser) Parse(data []byte) (IValue, error) {
	if len(data) == 0 {
		return nil, NewJSONError("empty input", 1, 1, 0)
	}
	
	// 使用Go标准库进行初步解析
	var raw interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, p.wrapJSONError(err, string(data))
	}
	
	// 转换为IValue
	return p.convertToIValue(raw)
}

// ParseString 解析JSON字符串
func (p *JSONParser) ParseString(s string) (IValue, error) {
	return p.Parse([]byte(s))
}

// convertToIValue 将interface{}转换为IValue
func (p *JSONParser) convertToIValue(raw interface{}) (IValue, error) {
	switch v := raw.(type) {
	case nil:
		return NewNullValue(), nil
	case bool:
		return NewBoolValue(v), nil
	case string:
		return NewStringValue(v), nil
	case float64:
		return NewNumberValue(v)
	case map[string]interface{}:
		return p.convertMapToObject(v)
	case []interface{}:
		return p.convertSliceToArray(v)
	default:
		return nil, NewTypeError("supported JSON type", fmt.Sprintf("%T", v), v)
	}
}

// convertMapToObject 将map转换为IObject
func (p *JSONParser) convertMapToObject(m map[string]interface{}) (IObject, error) {
	obj := NewObjectWithCapacity(len(m))
	for k, v := range m {
		value, err := p.convertToIValue(v)
		if err != nil {
			return nil, err
		}
		obj.Set(k, value)
	}
	return obj, nil
}

// convertSliceToArray 将slice转换为IArray
func (p *JSONParser) convertSliceToArray(s []interface{}) (IArray, error) {
	arr := NewArrayWithCapacity(len(s))
	for _, v := range s {
		value, err := p.convertToIValue(v)
		if err != nil {
			return nil, err
		}
		arr.Append(value)
	}
	return arr, nil
}

// wrapJSONError 包装JSON错误
func (p *JSONParser) wrapJSONError(err error, input string) error {
	if syntaxErr, ok := err.(*json.SyntaxError); ok {
		line, col := p.getLineColumn(input, syntaxErr.Offset)
		return NewJSONError(syntaxErr.Error(), line, col, syntaxErr.Offset)
	}
	if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
		line, col := p.getLineColumn(input, typeErr.Offset)
		return NewJSONError(typeErr.Error(), line, col, typeErr.Offset)
	}
	return NewJSONError(err.Error(), 1, 1, 0)
}

// getLineColumn 根据偏移量获取行列号
func (p *JSONParser) getLineColumn(input string, offset int64) (int, int) {
	line := 1
	col := 1
	for i := int64(0); i < offset && i < int64(len(input)); i++ {
		if input[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return line, col
}
