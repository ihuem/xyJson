package xyJson

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// NewObjectValue 创建新的对象值
// NewObjectValue creates a new object value
func NewObjectValue() IValue {
	return &objectValue{
		data: make(map[string]IValue),
	}
}

// NewArrayValue 创建新的数组值
// NewArrayValue creates a new array value
func NewArrayValue() IValue {
	return &arrayValue{
		data: make([]IValue, 0),
	}
}

// valueFactory 值工厂实现
// valueFactory implements the IValueFactory interface
type valueFactory struct {
	pool IObjectPool
}

// NewValueFactory 创建新的值工厂
// NewValueFactory creates a new value factory
func NewValueFactory() IValueFactory {
	return &valueFactory{
		pool: NewObjectPool(),
	}
}

// NewValueFactoryWithPool 使用指定对象池创建值工厂
// NewValueFactoryWithPool creates a value factory with the specified object pool
func NewValueFactoryWithPool(pool IObjectPool) IValueFactory {
	return &valueFactory{
		pool: pool,
	}
}

// CreateNull 创建null值
// CreateNull creates a null value
func (f *valueFactory) CreateNull() IValue {
	return &scalarValue{
		valueType: NullValueType,
		rawData:   nil,
	}
}

// CreateString 创建字符串值
// CreateString creates a string value
func (f *valueFactory) CreateString(s string) IScalarValue {
	return &scalarValue{
		valueType: StringValueType,
		rawData:   s,
	}
}

// CreateNumber 创建数字值
// CreateNumber creates a number value
func (f *valueFactory) CreateNumber(n interface{}) (IScalarValue, error) {
	if n == nil {
		return nil, NewInvalidOperationError("create number", "input cannot be nil")
	}

	switch v := n.(type) {
	case int:
		return &scalarValue{
			valueType: NumberValueType,
			rawData:   int64(v),
		}, nil
	case int8:
		return &scalarValue{
			valueType: NumberValueType,
			rawData:   int64(v),
		}, nil
	case int16:
		return &scalarValue{
			valueType: NumberValueType,
			rawData:   int64(v),
		}, nil
	case int32:
		return &scalarValue{
			valueType: NumberValueType,
			rawData:   int64(v),
		}, nil
	case int64:
		return &scalarValue{
			valueType: NumberValueType,
			rawData:   v,
		}, nil
	case uint:
		return &scalarValue{
			valueType: NumberValueType,
			rawData:   int64(v),
		}, nil
	case uint8:
		return &scalarValue{
			valueType: NumberValueType,
			rawData:   int64(v),
		}, nil
	case uint16:
		return &scalarValue{
			valueType: NumberValueType,
			rawData:   int64(v),
		}, nil
	case uint32:
		return &scalarValue{
			valueType: NumberValueType,
			rawData:   int64(v),
		}, nil
	case uint64:
		// 检查是否超出int64范围
		if v > 9223372036854775807 {
			return &scalarValue{
				valueType: NumberValueType,
				rawData:   float64(v),
			}, nil
		}
		return &scalarValue{
			valueType: NumberValueType,
			rawData:   int64(v),
		}, nil
	case float32:
		return &scalarValue{
			valueType: NumberValueType,
			rawData:   float64(v),
		}, nil
	case float64:
		return &scalarValue{
			valueType: NumberValueType,
			rawData:   v,
		}, nil
	case string:
		// 尝试解析字符串为数字
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return &scalarValue{
				valueType: NumberValueType,
				rawData:   i,
			}, nil
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return &scalarValue{
				valueType: NumberValueType,
				rawData:   f,
			}, nil
		}
		return nil, NewInvalidOperationError("create number", fmt.Sprintf("cannot parse '%s' as number", v))
	default:
		return nil, NewInvalidOperationError("create number", fmt.Sprintf("unsupported type: %T", n))
	}
}

// CreateBool 创建布尔值
// CreateBool creates a boolean value
func (f *valueFactory) CreateBool(b bool) IScalarValue {
	return &scalarValue{
		valueType: BoolValueType,
		rawData:   b,
	}
}

// CreateObject 创建对象
// CreateObject creates an object
func (f *valueFactory) CreateObject() IObject {
	if f.pool != nil {
		if obj := f.pool.GetObject(); obj != nil {
			return obj
		}
	}
	return NewObjectValue().(IObject)
}

// CreateArray 创建数组
// CreateArray creates an array
func (f *valueFactory) CreateArray() IArray {
	if f.pool != nil {
		if arr := f.pool.GetArray(); arr != nil {
			return arr
		}
	}
	return NewArrayValue().(IArray)
}

// CreateFromRaw 从原始数据创建值
// CreateFromRaw creates a value from raw data
func (f *valueFactory) CreateFromRaw(data interface{}) (IValue, error) {
	if data == nil {
		return f.CreateNull(), nil
	}

	switch v := data.(type) {
	case string:
		return f.CreateString(v), nil
	case bool:
		return f.CreateBool(v), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return f.CreateNumber(v)
	case time.Time:
		return f.CreateString(v.Format(time.RFC3339)), nil
	case []byte:
		return f.CreateString(string(v)), nil
	case map[string]interface{}:
		obj := f.CreateObject()
		for key, val := range v {
			childValue, err := f.CreateFromRaw(val)
			if err != nil {
				return nil, err
			}
			if err := obj.Set(key, childValue); err != nil {
				return nil, err
			}
		}
		return obj, nil
	case []interface{}:
		arr := f.CreateArray()
		for _, val := range v {
			childValue, err := f.CreateFromRaw(val)
			if err != nil {
				return nil, err
			}
			if err := arr.Append(childValue); err != nil {
				return nil, err
			}
		}
		return arr, nil
	case IValue:
		// 如果已经是IValue类型，直接返回
		return v, nil
	default:
		// 使用反射处理其他类型
		return f.createFromReflect(reflect.ValueOf(data))
	}
}

// createFromReflect 使用反射创建值
// createFromReflect creates a value using reflection
func (f *valueFactory) createFromReflect(rv reflect.Value) (IValue, error) {
	if !rv.IsValid() {
		return f.CreateNull(), nil
	}

	// 处理指针类型
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return f.CreateNull(), nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		return f.CreateString(rv.String()), nil
	case reflect.Bool:
		return f.CreateBool(rv.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.CreateNumber(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return f.CreateNumber(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return f.CreateNumber(rv.Float())
	case reflect.Slice, reflect.Array:
		arr := f.CreateArray()
		for i := 0; i < rv.Len(); i++ {
			elem, err := f.createFromReflect(rv.Index(i))
			if err != nil {
				return nil, err
			}
			if err := arr.Append(elem); err != nil {
				return nil, err
			}
		}
		return arr, nil
	case reflect.Map:
		obj := f.CreateObject()
		for _, key := range rv.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			val, err := f.createFromReflect(rv.MapIndex(key))
			if err != nil {
				return nil, err
			}
			if err := obj.Set(keyStr, val); err != nil {
				return nil, err
			}
		}
		return obj, nil
	case reflect.Struct:
		// 处理结构体类型
		obj := f.CreateObject()
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			field := rt.Field(i)
			if !field.IsExported() {
				continue
			}

			fieldName := field.Name
			// 检查json标签
			if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
				if idx := len(tag); idx > 0 {
					commaIdx := 0
					for commaIdx < len(tag) && tag[commaIdx] != ',' {
						commaIdx++
					}
					if commaIdx < len(tag) {
						fieldName = tag[:commaIdx]
					} else {
						fieldName = tag
					}
				}
			}

			val, err := f.createFromReflect(rv.Field(i))
			if err != nil {
				return nil, err
			}
			if err := obj.Set(fieldName, val); err != nil {
				return nil, err
			}
		}
		return obj, nil
	case reflect.Interface:
		if rv.IsNil() {
			return f.CreateNull(), nil
		}
		return f.createFromReflect(rv.Elem())
	default:
		// 对于其他类型，转换为字符串
		return f.CreateString(fmt.Sprintf("%v", rv.Interface())), nil
	}
}
