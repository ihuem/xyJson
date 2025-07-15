package xyJson

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// scalarValue 标量值实现
type scalarValue struct {
	value interface{}
	vType ValueType
}

// NewScalarValue 创建标量值
func NewScalarValue(value interface{}, vType ValueType) IValue {
	return &scalarValue{
		value: value,
		vType: vType,
	}
}

// NewNullValue 创建null值
func NewNullValue() IValue {
	return &scalarValue{
		value: nil,
		vType: NullValueType,
	}
}

// NewStringValue 创建字符串值
func NewStringValue(s string) IValue {
	return &scalarValue{
		value: s,
		vType: StringValueType,
	}
}

// NewNumberValue 创建数字值
func NewNumberValue(n interface{}) (IValue, error) {
	switch v := n.(type) {
	case int:
		return &scalarValue{value: float64(v), vType: NumberValueType}, nil
	case int8:
		return &scalarValue{value: float64(v), vType: NumberValueType}, nil
	case int16:
		return &scalarValue{value: float64(v), vType: NumberValueType}, nil
	case int32:
		return &scalarValue{value: float64(v), vType: NumberValueType}, nil
	case int64:
		return &scalarValue{value: float64(v), vType: NumberValueType}, nil
	case uint:
		return &scalarValue{value: float64(v), vType: NumberValueType}, nil
	case uint8:
		return &scalarValue{value: float64(v), vType: NumberValueType}, nil
	case uint16:
		return &scalarValue{value: float64(v), vType: NumberValueType}, nil
	case uint32:
		return &scalarValue{value: float64(v), vType: NumberValueType}, nil
	case uint64:
		return &scalarValue{value: float64(v), vType: NumberValueType}, nil
	case float32:
		return &scalarValue{value: float64(v), vType: NumberValueType}, nil
	case float64:
		return &scalarValue{value: v, vType: NumberValueType}, nil
	default:
		return nil, NewTypeError("number", fmt.Sprintf("%T", n), n)
	}
}

// NewBoolValue 创建布尔值
func NewBoolValue(b bool) IValue {
	return &scalarValue{
		value: b,
		vType: BoolValueType,
	}
}

// Type 获取值类型
func (sv *scalarValue) Type() ValueType {
	return sv.vType
}

// Raw 获取原始Go类型值
func (sv *scalarValue) Raw() interface{} {
	return sv.value
}

// String 获取字符串表示
func (sv *scalarValue) String() string {
	switch sv.vType {
	case NullValueType:
		return "null"
	case StringValueType:
		return sv.value.(string)
	case NumberValueType:
		return fmt.Sprintf("%g", sv.value.(float64))
	case BoolValueType:
		if sv.value.(bool) {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", sv.value)
	}
}

// IsNull 检查是否为null
func (sv *scalarValue) IsNull() bool {
	return sv.vType == NullValueType
}

// IsString 检查是否为字符串
func (sv *scalarValue) IsString() bool {
	return sv.vType == StringValueType
}

// IsNumber 检查是否为数字
func (sv *scalarValue) IsNumber() bool {
	return sv.vType == NumberValueType
}

// IsBool 检查是否为布尔值
func (sv *scalarValue) IsBool() bool {
	return sv.vType == BoolValueType
}

// IsObject 检查是否为对象
func (sv *scalarValue) IsObject() bool {
	return false
}

// IsArray 检查是否为数组
func (sv *scalarValue) IsArray() bool {
	return false
}

// Clone 深拷贝
func (sv *scalarValue) Clone() IValue {
	return &scalarValue{
		value: sv.value,
		vType: sv.vType,
	}
}

// Equals 比较是否相等
func (sv *scalarValue) Equals(other IValue) bool {
	if other == nil || sv.Type() != other.Type() {
		return false
	}
	return sv.value == other.Raw()
}

const (
	maxInt = int(^uint(0) >> 1)
	minInt = -maxInt - 1
)

// Int 转换为int
func (sv *scalarValue) Int() (int, error) {
	switch sv.vType {
	case NullValueType:
		return 0, NewTypeError("int", "null", nil)
	case NumberValueType:
		f := sv.value.(float64)
		// 检查是否为整数
		if f != math.Trunc(f) {
			return 0, NewTypeError("int", "float with fractional part", f)
		}
		// 检查溢出
		if f > float64(maxInt) || f < float64(minInt) {
			return 0, NewTypeError("int", "number out of range", f)
		}
		return int(f), nil
	case StringValueType:
		return strconv.Atoi(sv.value.(string))
	case BoolValueType:
		if sv.value.(bool) {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, NewTypeError("int", sv.vType.String(), sv.value)
	}
}

// Int64 转换为int64
func (sv *scalarValue) Int64() (int64, error) {
	switch sv.vType {
	case NullValueType:
		return 0, NewTypeError("int64", "null", nil)
	case NumberValueType:
		f := sv.value.(float64)
		// 检查是否为整数
		if f != math.Trunc(f) {
			return 0, NewTypeError("int64", "float with fractional part", f)
		}
		// 检查溢出
		if f > float64(math.MaxInt64) || f < float64(math.MinInt64) {
			return 0, NewTypeError("int64", "number out of range", f)
		}
		return int64(f), nil
	case StringValueType:
		return strconv.ParseInt(sv.value.(string), 10, 64)
	case BoolValueType:
		if sv.value.(bool) {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, NewTypeError("int64", sv.vType.String(), sv.value)
	}
}

// Float64 转换为float64
func (sv *scalarValue) Float64() (float64, error) {
	switch sv.vType {
	case NullValueType:
		return 0, NewTypeError("float64", "null", nil)
	case NumberValueType:
		return sv.value.(float64), nil
	case StringValueType:
		return strconv.ParseFloat(sv.value.(string), 64)
	case BoolValueType:
		if sv.value.(bool) {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0, NewTypeError("float64", sv.vType.String(), sv.value)
	}
}

// Bool 转换为bool
func (sv *scalarValue) Bool() (bool, error) {
	switch sv.vType {
	case NullValueType:
		return false, nil
	case BoolValueType:
		return sv.value.(bool), nil
	case NumberValueType:
		return sv.value.(float64) != 0, nil
	case StringValueType:
		s := sv.value.(string)
		if s == "" {
			return false, nil
		}
		return strconv.ParseBool(s)
	default:
		return false, NewTypeError("bool", sv.vType.String(), sv.value)
	}
}

// Time 转换为time.Time
func (sv *scalarValue) Time() (time.Time, error) {
	switch sv.vType {
	case NullValueType:
		return time.Time{}, NewTypeError("time.Time", "null", nil)
	case StringValueType:
		s := sv.value.(string)
		// 尝试多种时间格式
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05Z",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, s); err == nil {
				return t, nil
			}
		}
		return time.Time{}, NewTypeError("time.Time", "invalid time format", s)
	case NumberValueType:
		// 假设是Unix时间戳
		timestamp := sv.value.(float64)
		return time.Unix(int64(timestamp), 0), nil
	default:
		return time.Time{}, NewTypeError("time.Time", sv.vType.String(), sv.value)
	}
}

// Bytes 转换为[]byte
func (sv *scalarValue) Bytes() ([]byte, error) {
	switch sv.vType {
	case NullValueType:
		return nil, nil
	case StringValueType:
		return []byte(sv.value.(string)), nil
	default:
		return []byte(sv.String()), nil
	}
}
