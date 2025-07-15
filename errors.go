package xyJson

import (
	"fmt"
)

// ErrorCode 错误码枚举
// ErrorCode represents different types of JSON processing errors
type ErrorCode int

const (
	// ErrInvalidJSON 无效JSON格式
	// ErrInvalidJSON indicates invalid JSON format
	ErrInvalidJSON ErrorCode = iota + 1000
	// ErrPathNotFound 路径不存在
	// ErrPathNotFound indicates the specified path does not exist
	ErrPathNotFound
	// ErrTypeMismatch 类型不匹配
	// ErrTypeMismatch indicates a type mismatch error
	ErrTypeMismatch
	// ErrIndexOutOfRange 索引超出范围
	// ErrIndexOutOfRange indicates array index is out of range
	ErrIndexOutOfRange
	// ErrKeyNotFound 键名不存在
	// ErrKeyNotFound indicates the specified key does not exist
	ErrKeyNotFound
	// ErrCircularReference 循环引用
	// ErrCircularReference indicates a circular reference in the data structure
	ErrCircularReference
	// ErrMaxDepthExceeded 超过最大嵌套深度
	// ErrMaxDepthExceeded indicates maximum nesting depth exceeded
	ErrMaxDepthExceeded
	// ErrInvalidPath 无效路径表达式
	// ErrInvalidPath indicates invalid path expression
	ErrInvalidPath
	// ErrNullPointer 空指针错误
	// ErrNullPointer indicates null pointer error
	ErrNullPointer
	// ErrInvalidOperation 无效操作
	// ErrInvalidOperation indicates invalid operation
	ErrInvalidOperation
)

// String 返回错误码的字符串表示
// String returns the string representation of the error code
func (ec ErrorCode) String() string {
	switch ec {
	case ErrInvalidJSON:
		return "INVALID_JSON"
	case ErrPathNotFound:
		return "PATH_NOT_FOUND"
	case ErrTypeMismatch:
		return "TYPE_MISMATCH"
	case ErrIndexOutOfRange:
		return "INDEX_OUT_OF_RANGE"
	case ErrKeyNotFound:
		return "KEY_NOT_FOUND"
	case ErrCircularReference:
		return "CIRCULAR_REFERENCE"
	case ErrMaxDepthExceeded:
		return "MAX_DEPTH_EXCEEDED"
	case ErrInvalidPath:
		return "INVALID_PATH"
	case ErrNullPointer:
		return "NULL_POINTER"
	case ErrInvalidOperation:
		return "INVALID_OPERATION"
	default:
		return "UNKNOWN_ERROR"
	}
}

// JSONError 自定义JSON错误类型
// JSONError represents a custom JSON processing error
type JSONError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Cause   error     `json:"cause,omitempty"`
	Path    string    `json:"path,omitempty"`
	Line    int       `json:"line,omitempty"`
	Column  int       `json:"column,omitempty"`
	Context string    `json:"context,omitempty"`
}

// Error 实现error接口
// Error implements the error interface
func (je *JSONError) Error() string {
	if je.Path != "" {
		return fmt.Sprintf("[%s] %s at path '%s'", je.Code.String(), je.Message, je.Path)
	}
	if je.Line > 0 && je.Column > 0 {
		return fmt.Sprintf("[%s] %s at line %d, column %d", je.Code.String(), je.Message, je.Line, je.Column)
	}
	return fmt.Sprintf("[%s] %s", je.Code.String(), je.Message)
}

// Unwrap 返回底层错误
// Unwrap returns the underlying error
func (je *JSONError) Unwrap() error {
	return je.Cause
}

// WithPath 添加路径信息
// WithPath adds path information to the error
func (je *JSONError) WithPath(path string) *JSONError {
	je.Path = path
	return je
}

// WithPosition 添加位置信息
// WithPosition adds position information to the error
func (je *JSONError) WithPosition(line, column int) *JSONError {
	je.Line = line
	je.Column = column
	return je
}

// WithContext 添加上下文信息
// WithContext adds context information to the error
func (je *JSONError) WithContext(context string) *JSONError {
	je.Context = context
	return je
}

// NewJSONError 创建新的JSON错误
// NewJSONError creates a new JSON error
func NewJSONError(code ErrorCode, message string, cause error) *JSONError {
	return &JSONError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NewInvalidJSONError 创建无效JSON错误
// NewInvalidJSONError creates an invalid JSON error
func NewInvalidJSONError(message string, cause error) *JSONError {
	return NewJSONError(ErrInvalidJSON, message, cause)
}

// NewPathNotFoundError 创建路径不存在错误
// NewPathNotFoundError creates a path not found error
func NewPathNotFoundError(path string) *JSONError {
	return NewJSONError(ErrPathNotFound, fmt.Sprintf("path '%s' not found", path), nil).WithPath(path)
}

// NewTypeMismatchError 创建类型不匹配错误
// NewTypeMismatchError creates a type mismatch error
func NewTypeMismatchError(expected, actual ValueType, path string) *JSONError {
	message := fmt.Sprintf("expected %s but got %s", expected.String(), actual.String())
	return NewJSONError(ErrTypeMismatch, message, nil).WithPath(path)
}

// NewIndexOutOfRangeError 创建索引超出范围错误
// NewIndexOutOfRangeError creates an index out of range error
func NewIndexOutOfRangeError(index, length int, path string) *JSONError {
	message := fmt.Sprintf("index %d out of range [0, %d)", index, length)
	return NewJSONError(ErrIndexOutOfRange, message, nil).WithPath(path)
}

// NewKeyNotFoundError 创建键名不存在错误
// NewKeyNotFoundError creates a key not found error
func NewKeyNotFoundError(key, path string) *JSONError {
	message := fmt.Sprintf("key '%s' not found", key)
	return NewJSONError(ErrKeyNotFound, message, nil).WithPath(path)
}

// NewCircularReferenceError 创建循环引用错误
// NewCircularReferenceError creates a circular reference error
func NewCircularReferenceError(path string) *JSONError {
	message := "circular reference detected"
	return NewJSONError(ErrCircularReference, message, nil).WithPath(path)
}

// NewMaxDepthExceededError 创建超过最大深度错误
// NewMaxDepthExceededError creates a max depth exceeded error
func NewMaxDepthExceededError(depth int) *JSONError {
	message := fmt.Sprintf("maximum nesting depth %d exceeded", depth)
	return NewJSONError(ErrMaxDepthExceeded, message, nil)
}

// NewInvalidPathError 创建无效路径错误
// NewInvalidPathError creates an invalid path error
func NewInvalidPathError(path string, cause error) *JSONError {
	message := fmt.Sprintf("invalid path expression: %s", path)
	return NewJSONError(ErrInvalidPath, message, cause).WithPath(path)
}

// NewNullPointerError 创建空指针错误
// NewNullPointerError creates a null pointer error
func NewNullPointerError(context string) *JSONError {
	message := "null pointer access"
	return NewJSONError(ErrNullPointer, message, nil).WithContext(context)
}

// NewInvalidOperationError 创建无效操作错误
// NewInvalidOperationError creates an invalid operation error
func NewInvalidOperationError(operation, context string) *JSONError {
	message := fmt.Sprintf("invalid operation '%s'", operation)
	return NewJSONError(ErrInvalidOperation, message, nil).WithContext(context)
}
