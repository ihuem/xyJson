package xyJson

import (
	"fmt"
)

// JSONError JSON解析错误
type JSONError struct {
	Message string
	Line    int
	Column  int
	Offset  int64
}

func (e *JSONError) Error() string {
	return fmt.Sprintf("JSON error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

// NewJSONError 创建JSON错误
func NewJSONError(message string, line, column int, offset int64) *JSONError {
	return &JSONError{
		Message: message,
		Line:    line,
		Column:  column,
		Offset:  offset,
	}
}

// TypeError 类型错误
type TypeError struct {
	Expected string
	Actual   string
	Value    interface{}
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("type error: expected %s, got %s (value: %v)", e.Expected, e.Actual, e.Value)
}

// NewTypeError 创建类型错误
func NewTypeError(expected, actual string, value interface{}) *TypeError {
	return &TypeError{
		Expected: expected,
		Actual:   actual,
		Value:    value,
	}
}

// PathError JSONPath错误
type PathError struct {
	Path    string
	Message string
}

func (e *PathError) Error() string {
	return fmt.Sprintf("JSONPath error in '%s': %s", e.Path, e.Message)
}

// NewPathError 创建路径错误
func NewPathError(path, message string) *PathError {
	return &PathError{
		Path:    path,
		Message: message,
	}
}

// IndexError 索引错误
type IndexError struct {
	Index  int
	Length int
}

func (e *IndexError) Error() string {
	return fmt.Sprintf("index error: index %d out of range [0, %d)", e.Index, e.Length)
}

// NewIndexError 创建索引错误
func NewIndexError(index, length int) *IndexError {
	return &IndexError{
		Index:  index,
		Length: length,
	}
}

// SerializationError 序列化错误
type SerializationError struct {
	Message string
	Type    string
}

func (e *SerializationError) Error() string {
	return fmt.Sprintf("serialization error: %s (type: %s)", e.Message, e.Type)
}

// NewSerializationError 创建序列化错误
func NewSerializationError(message, type_ string) *SerializationError {
	return &SerializationError{
		Message: message,
		Type:    type_,
	}
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
}

// NewValidationError 创建验证错误
func NewValidationError(field, message string, value interface{}) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}
