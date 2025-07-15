package xyJson

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

// scalarValue 标量值实现（字符串、数字、布尔值、null）
// scalarValue implements scalar values (string, number, boolean, null)
type scalarValue struct {
	valueType ValueType
	rawData   interface{}
}

// Type 返回值的类型
// Type returns the type of the value
func (sv *scalarValue) Type() ValueType {
	return sv.valueType
}

// Raw 返回原始Go类型值
// Raw returns the raw Go type value
func (sv *scalarValue) Raw() interface{} {
	return sv.rawData
}

// String 返回字符串表示
// String returns the string representation
func (sv *scalarValue) String() string {
	if sv.IsNull() {
		return ""
	}

	switch sv.valueType {
	case StringValueType:
		if str, ok := sv.rawData.(string); ok {
			return str
		}
		return ""
	case NumberValueType:
		return sv.numberToString()
	case BoolValueType:
		if b, ok := sv.rawData.(bool); ok {
			if b {
				return "true"
			}
			return "false"
		}
		return "false"
	default:
		return ""
	}
}

// IsNull 检查是否为null值
// IsNull checks if the value is null
func (sv *scalarValue) IsNull() bool {
	return sv.valueType == NullValueType || sv.rawData == nil
}

// Clone 创建值的深拷贝
// Clone creates a deep copy of the value
func (sv *scalarValue) Clone() IValue {
	return &scalarValue{
		valueType: sv.valueType,
		rawData:   sv.rawData,
	}
}

// Equals 比较两个值是否相等
// Equals compares if two values are equal
func (sv *scalarValue) Equals(other IValue) bool {
	if other == nil {
		return false
	}

	if sv.Type() != other.Type() {
		return false
	}

	if sv.IsNull() && other.IsNull() {
		return true
	}

	if sv.IsNull() || other.IsNull() {
		return false
	}

	// 比较原始数据
	return sv.rawData == other.Raw()
}

// Int 返回整数值
// Int returns the integer value
func (sv *scalarValue) Int() (int, error) {
	if sv.IsNull() {
		return 0, NewTypeMismatchError(NumberValueType, NullValueType, "")
	}

	switch sv.valueType {
	case NumberValueType:
		switch v := sv.rawData.(type) {
		case int64:
			// 检查int64到int的溢出
			const maxInt = int64(^uint(0) >> 1)
			const minInt = -maxInt - 1
			if v > maxInt || v < minInt {
				return 0, NewInvalidOperationError("int conversion", fmt.Sprintf("value %d overflows int", v))
			}
			return int(v), nil
		case float64:
			// 检查浮点数精度丢失和溢出
			if v != float64(int64(v)) {
				return 0, NewInvalidOperationError("int conversion", fmt.Sprintf("float64 %g has fractional part", v))
			}
			intVal := int64(v)
			const maxInt = int64(^uint(0) >> 1)
			const minInt = -maxInt - 1
			if intVal > maxInt || intVal < minInt {
				return 0, NewInvalidOperationError("int conversion", fmt.Sprintf("value %g overflows int", v))
			}
			return int(intVal), nil
		default:
			return 0, NewInvalidOperationError("int conversion", fmt.Sprintf("unexpected number type: %T", v))
		}
	case StringValueType:
		if str, ok := sv.rawData.(string); ok {
			if i, err := strconv.Atoi(str); err == nil {
				return i, nil
			}
			return 0, NewInvalidOperationError("int conversion", fmt.Sprintf("cannot parse '%s' as int", str))
		}
		return 0, NewTypeMismatchError(NumberValueType, StringValueType, "")
	case BoolValueType:
		if b, ok := sv.rawData.(bool); ok {
			if b {
				return 1, nil
			}
			return 0, nil
		}
		return 0, NewTypeMismatchError(NumberValueType, BoolValueType, "")
	default:
		return 0, NewTypeMismatchError(NumberValueType, sv.valueType, "")
	}
}

// Int64 返回64位整数值
// Int64 returns the 64-bit integer value
func (sv *scalarValue) Int64() (int64, error) {
	if sv.IsNull() {
		return 0, NewTypeMismatchError(NumberValueType, NullValueType, "")
	}

	switch sv.valueType {
	case NumberValueType:
		switch v := sv.rawData.(type) {
		case int64:
			return v, nil
		case float64:
			// 检查浮点数精度丢失
			if v != float64(int64(v)) {
				return 0, NewInvalidOperationError("int64 conversion", fmt.Sprintf("float64 %g has fractional part", v))
			}
			// 检查溢出
			if v > 9223372036854775807 || v < -9223372036854775808 {
				return 0, NewInvalidOperationError("int64 conversion", fmt.Sprintf("value %g overflows int64", v))
			}
			return int64(v), nil
		default:
			return 0, NewInvalidOperationError("int64 conversion", fmt.Sprintf("unexpected number type: %T", v))
		}
	case StringValueType:
		if str, ok := sv.rawData.(string); ok {
			if i, err := strconv.ParseInt(str, 10, 64); err == nil {
				return i, nil
			}
			return 0, NewInvalidOperationError("int64 conversion", fmt.Sprintf("cannot parse '%s' as int64", str))
		}
		return 0, NewTypeMismatchError(NumberValueType, StringValueType, "")
	case BoolValueType:
		if b, ok := sv.rawData.(bool); ok {
			if b {
				return 1, nil
			}
			return 0, nil
		}
		return 0, NewTypeMismatchError(NumberValueType, BoolValueType, "")
	default:
		return 0, NewTypeMismatchError(NumberValueType, sv.valueType, "")
	}
}

// Float64 返回64位浮点数值
// Float64 returns the 64-bit float value
func (sv *scalarValue) Float64() (float64, error) {
	if sv.IsNull() {
		return 0, NewTypeMismatchError(NumberValueType, NullValueType, "")
	}

	switch sv.valueType {
	case NumberValueType:
		switch v := sv.rawData.(type) {
		case int64:
			return float64(v), nil
		case float64:
			return v, nil
		default:
			return 0, NewInvalidOperationError("float64 conversion", fmt.Sprintf("unexpected number type: %T", v))
		}
	case StringValueType:
		if str, ok := sv.rawData.(string); ok {
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				return f, nil
			}
			return 0, NewInvalidOperationError("float64 conversion", fmt.Sprintf("cannot parse '%s' as float64", str))
		}
		return 0, NewTypeMismatchError(NumberValueType, StringValueType, "")
	case BoolValueType:
		if b, ok := sv.rawData.(bool); ok {
			if b {
				return 1.0, nil
			}
			return 0.0, nil
		}
		return 0, NewTypeMismatchError(NumberValueType, BoolValueType, "")
	default:
		return 0, NewTypeMismatchError(NumberValueType, sv.valueType, "")
	}
}

// Bool 返回布尔值
// Bool returns the boolean value
func (sv *scalarValue) Bool() (bool, error) {
	if sv.IsNull() {
		return false, NewTypeMismatchError(BoolValueType, NullValueType, "")
	}

	switch sv.valueType {
	case BoolValueType:
		if b, ok := sv.rawData.(bool); ok {
			return b, nil
		}
		return false, NewInvalidOperationError("bool conversion", "invalid bool data")
	case NumberValueType:
		switch v := sv.rawData.(type) {
		case int64:
			return v != 0, nil
		case float64:
			return v != 0.0, nil
		default:
			return false, NewInvalidOperationError("bool conversion", fmt.Sprintf("unexpected number type: %T", v))
		}
	case StringValueType:
		if str, ok := sv.rawData.(string); ok {
			if b, err := strconv.ParseBool(str); err == nil {
				return b, nil
			}
			// 对于非标准布尔字符串，使用长度判断
			return len(str) > 0, nil
		}
		return false, NewTypeMismatchError(BoolValueType, StringValueType, "")
	default:
		return false, NewTypeMismatchError(BoolValueType, sv.valueType, "")
	}
}

// Time 返回时间值
// Time returns the time value
func (sv *scalarValue) Time() (time.Time, error) {
	if sv.IsNull() {
		return time.Time{}, NewTypeMismatchError(StringValueType, NullValueType, "")
	}

	if sv.valueType != StringValueType {
		return time.Time{}, NewTypeMismatchError(StringValueType, sv.valueType, "")
	}

	str, ok := sv.rawData.(string)
	if !ok {
		return time.Time{}, NewInvalidOperationError("time conversion", "invalid string data")
	}

	// 尝试多种时间格式
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			return t, nil
		}
	}

	return time.Time{}, NewInvalidOperationError("time conversion", fmt.Sprintf("cannot parse '%s' as time", str))
}

// Bytes 返回字节数组
// Bytes returns the byte array
func (sv *scalarValue) Bytes() ([]byte, error) {
	if sv.IsNull() {
		return nil, NewTypeMismatchError(StringValueType, NullValueType, "")
	}

	if sv.valueType != StringValueType {
		return nil, NewTypeMismatchError(StringValueType, sv.valueType, "")
	}

	str, ok := sv.rawData.(string)
	if !ok {
		return nil, NewInvalidOperationError("bytes conversion", "invalid string data")
	}

	// 尝试base64解码
	if data, err := base64.StdEncoding.DecodeString(str); err == nil {
		return data, nil
	}

	// 如果不是base64，直接返回字符串的字节表示
	return []byte(str), nil
}

// numberToString 将数字转换为字符串
// numberToString converts a number to string
func (sv *scalarValue) numberToString() string {
	switch v := sv.rawData.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		// 使用-1精度让Go自动选择最短表示
		return strconv.FormatFloat(v, 'g', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}
