package xyJson

import (
	"fmt"
	"reflect"
	"time"
)

// ValueFactory 值工厂实现
type ValueFactory struct{}

// NewValueFactory 创建新的值工厂
func NewValueFactory() IValueFactory {
	return &ValueFactory{}
}

// CreateNull 创建null值
func (f *ValueFactory) CreateNull() IValue {
	return NewNullValue()
}

// CreateString 创建字符串值
func (f *ValueFactory) CreateString(s string) IValue {
	return NewStringValue(s)
}

// CreateNumber 创建数字值
func (f *ValueFactory) CreateNumber(n interface{}) (IValue, error) {
	return NewNumberValue(n)
}

// CreateBool 创建布尔值
func (f *ValueFactory) CreateBool(b bool) IValue {
	return NewBoolValue(b)
}

// CreateObject 创建对象
func (f *ValueFactory) CreateObject() IObject {
	return NewObject()
}

// CreateObjectWithCapacity 创建指定容量的对象
func (f *ValueFactory) CreateObjectWithCapacity(capacity int) IObject {
	return NewObjectWithCapacity(capacity)
}

// CreateArray 创建数组
func (f *ValueFactory) CreateArray() IArray {
	return NewArray()
}

// CreateArrayWithCapacity 创建指定容量的数组
func (f *ValueFactory) CreateArrayWithCapacity(capacity int) IArray {
	return NewArrayWithCapacity(capacity)
}

// CreateFromRaw 从原始Go类型创建值
func (f *ValueFactory) CreateFromRaw(v interface{}) (IValue, error) {
	if v == nil {
		return f.CreateNull(), nil
	}
	
	switch val := v.(type) {
	case bool:
		return f.CreateBool(val), nil
	case string:
		return f.CreateString(val), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return f.CreateNumber(val)
	case time.Time:
		return f.CreateString(val.Format(time.RFC3339)), nil
	case []byte:
		return f.CreateString(string(val)), nil
	case map[string]interface{}:
		return f.createObjectFromMap(val)
	case []interface{}:
		return f.createArrayFromSlice(val)
	default:
		// 使用反射处理其他类型
		return f.createFromReflect(reflect.ValueOf(v))
	}
}

// createObjectFromMap 从map创建对象
func (f *ValueFactory) createObjectFromMap(m map[string]interface{}) (IObject, error) {
	obj := f.CreateObjectWithCapacity(len(m))
	for k, v := range m {
		value, err := f.CreateFromRaw(v)
		if err != nil {
			return nil, err
		}
		obj.Set(k, value)
	}
	return obj, nil
}

// createArrayFromSlice 从切片创建数组
func (f *ValueFactory) createArrayFromSlice(s []interface{}) (IArray, error) {
	arr := f.CreateArrayWithCapacity(len(s))
	for _, v := range s {
		value, err := f.CreateFromRaw(v)
		if err != nil {
			return nil, err
		}
		arr.Append(value)
	}
	return arr, nil
}

// createFromReflect 使用反射创建值
func (f *ValueFactory) createFromReflect(rv reflect.Value) (IValue, error) {
	if !rv.IsValid() {
		return f.CreateNull(), nil
	}
	
	// 处理指针
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return f.CreateNull(), nil
		}
		rv = rv.Elem()
	}
	
	switch rv.Kind() {
	case reflect.Bool:
		return f.CreateBool(rv.Bool()), nil
	
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.CreateNumber(rv.Int())
	
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return f.CreateNumber(rv.Uint())
	
	case reflect.Float32, reflect.Float64:
		return f.CreateNumber(rv.Float())
	
	case reflect.String:
		return f.CreateString(rv.String()), nil
	
	case reflect.Slice, reflect.Array:
		return f.createArrayFromReflect(rv)
	
	case reflect.Map:
		return f.createObjectFromReflect(rv)
	
	case reflect.Struct:
		return f.createObjectFromStruct(rv)
	
	case reflect.Interface:
		if rv.IsNil() {
			return f.CreateNull(), nil
		}
		return f.createFromReflect(rv.Elem())
	
	default:
		return nil, NewTypeError("supported type", rv.Kind().String(), rv.Interface())
	}
}

// createArrayFromReflect 使用反射创建数组
func (f *ValueFactory) createArrayFromReflect(rv reflect.Value) (IArray, error) {
	length := rv.Len()
	arr := f.CreateArrayWithCapacity(length)
	
	for i := 0; i < length; i++ {
		elem := rv.Index(i)
		value, err := f.createFromReflect(elem)
		if err != nil {
			return nil, err
		}
		arr.Append(value)
	}
	
	return arr, nil
}

// createObjectFromReflect 使用反射创建对象
func (f *ValueFactory) createObjectFromReflect(rv reflect.Value) (IObject, error) {
	if rv.Kind() != reflect.Map {
		return nil, NewTypeError("map", rv.Kind().String(), rv.Interface())
	}
	
	obj := f.CreateObjectWithCapacity(rv.Len())
	
	for _, key := range rv.MapKeys() {
		if key.Kind() != reflect.String {
			return nil, NewTypeError("string key", key.Kind().String(), key.Interface())
		}
		
		value, err := f.createFromReflect(rv.MapIndex(key))
		if err != nil {
			return nil, err
		}
		
		obj.Set(key.String(), value)
	}
	
	return obj, nil
}

// createObjectFromStruct 从结构体创建对象
func (f *ValueFactory) createObjectFromStruct(rv reflect.Value) (IObject, error) {
	rt := rv.Type()
	obj := f.CreateObjectWithCapacity(rt.NumField())
	
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.Field(i)
		
		// 跳过未导出的字段
		if !fieldValue.CanInterface() {
			continue
		}
		
		// 获取JSON标签
		jsonTag := field.Tag.Get("json")
		fieldName := field.Name
		
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] == "-" {
				continue // 跳过标记为忽略的字段
			}
			if parts[0] != "" {
				fieldName = parts[0]
			}
			
			// 检查omitempty选项
			if len(parts) > 1 && contains(parts[1:], "omitempty") {
				if isEmptyValue(fieldValue) {
					continue
				}
			}
		}
		
		value, err := f.createFromReflect(fieldValue)
		if err != nil {
			return nil, err
		}
		
		obj.Set(fieldName, value)
	}
	
	return obj, nil
}

// contains 检查字符串切片是否包含指定字符串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isEmptyValue 检查值是否为空
func isEmptyValue(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return rv.IsNil()
	}
	return false
}
