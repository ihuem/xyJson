package xyJson

import (
	"sort"
	"sync"
	"time"
)

// objectValue 对象值实现
type objectValue struct {
	fields map[string]IValue
	mu     sync.RWMutex
}

// NewObject 创建新的对象
func NewObject() IObject {
	return &objectValue{
		fields: make(map[string]IValue),
	}
}

// NewObjectWithCapacity 创建指定容量的对象
func NewObjectWithCapacity(capacity int) IObject {
	return &objectValue{
		fields: make(map[string]IValue, capacity),
	}
}

// NewObjectFromMap 从map创建对象
func NewObjectFromMap(m map[string]interface{}) (IObject, error) {
	obj := NewObject()
	for k, v := range m {
		value, err := createValueFromInterface(v)
		if err != nil {
			return nil, err
		}
		obj.Set(k, value)
	}
	return obj, nil
}

// Type 获取值类型
func (ov *objectValue) Type() ValueType {
	return ObjectValueType
}

// Raw 获取原始Go类型值
func (ov *objectValue) Raw() interface{} {
	ov.mu.RLock()
	defer ov.mu.RUnlock()
	
	result := make(map[string]interface{}, len(ov.fields))
	for k, v := range ov.fields {
		result[k] = v.Raw()
	}
	return result
}

// String 获取字符串表示
func (ov *objectValue) String() string {
	serializer := NewJSONSerializer()
	result, _ := serializer.SerializeToString(ov)
	return result
}

// IsNull 检查是否为null
func (ov *objectValue) IsNull() bool {
	return false
}

// IsString 检查是否为字符串
func (ov *objectValue) IsString() bool {
	return false
}

// IsNumber 检查是否为数字
func (ov *objectValue) IsNumber() bool {
	return false
}

// IsBool 检查是否为布尔值
func (ov *objectValue) IsBool() bool {
	return false
}

// IsObject 检查是否为对象
func (ov *objectValue) IsObject() bool {
	return true
}

// IsArray 检查是否为数组
func (ov *objectValue) IsArray() bool {
	return false
}

// Clone 深拷贝
func (ov *objectValue) Clone() IValue {
	ov.mu.RLock()
	defer ov.mu.RUnlock()
	
	newObj := NewObject()
	for k, v := range ov.fields {
		newObj.Set(k, v.Clone())
	}
	return newObj
}

// Equals 比较是否相等
func (ov *objectValue) Equals(other IValue) bool {
	if other == nil || !other.IsObject() {
		return false
	}
	
	otherObj := other.(IObject)
	if ov.Size() != otherObj.Size() {
		return false
	}
	
	ov.mu.RLock()
	defer ov.mu.RUnlock()
	
	for k, v := range ov.fields {
		otherValue, exists := otherObj.Get(k)
		if !exists || !v.Equals(otherValue) {
			return false
		}
	}
	return true
}

// Int 转换为int
func (ov *objectValue) Int() (int, error) {
	return 0, NewTypeError("int", "object", ov.Raw())
}

// Int64 转换为int64
func (ov *objectValue) Int64() (int64, error) {
	return 0, NewTypeError("int64", "object", ov.Raw())
}

// Float64 转换为float64
func (ov *objectValue) Float64() (float64, error) {
	return 0, NewTypeError("float64", "object", ov.Raw())
}

// Bool 转换为bool
func (ov *objectValue) Bool() (bool, error) {
	return ov.Size() > 0, nil
}

// Time 转换为time.Time
func (ov *objectValue) Time() (time.Time, error) {
	return time.Time{}, NewTypeError("time.Time", "object", ov.Raw())
}

// Bytes 转换为[]byte
func (ov *objectValue) Bytes() ([]byte, error) {
	serializer := NewJSONSerializer()
	return serializer.Serialize(ov)
}

// Get 获取字段值
func (ov *objectValue) Get(key string) (IValue, bool) {
	ov.mu.RLock()
	defer ov.mu.RUnlock()
	
	value, exists := ov.fields[key]
	return value, exists
}

// Set 设置字段值
func (ov *objectValue) Set(key string, value IValue) {
	ov.mu.Lock()
	defer ov.mu.Unlock()
	
	ov.fields[key] = value
}

// Delete 删除字段
func (ov *objectValue) Delete(key string) bool {
	ov.mu.Lock()
	defer ov.mu.Unlock()
	
	if _, exists := ov.fields[key]; exists {
		delete(ov.fields, key)
		return true
	}
	return false
}

// Has 检查字段是否存在
func (ov *objectValue) Has(key string) bool {
	ov.mu.RLock()
	defer ov.mu.RUnlock()
	
	_, exists := ov.fields[key]
	return exists
}

// Keys 获取所有键
func (ov *objectValue) Keys() []string {
	ov.mu.RLock()
	defer ov.mu.RUnlock()
	
	keys := make([]string, 0, len(ov.fields))
	for k := range ov.fields {
		keys = append(keys, k)
	}
	return keys
}

// Size 获取字段数量
func (ov *objectValue) Size() int {
	ov.mu.RLock()
	defer ov.mu.RUnlock()
	
	return len(ov.fields)
}

// Clear 清空所有字段
func (ov *objectValue) Clear() {
	ov.mu.Lock()
	defer ov.mu.Unlock()
	
	ov.fields = make(map[string]IValue)
}

// Range 遍历所有字段
func (ov *objectValue) Range(fn func(key string, value IValue) bool) {
	ov.mu.RLock()
	fields := make(map[string]IValue, len(ov.fields))
	for k, v := range ov.fields {
		fields[k] = v
	}
	ov.mu.RUnlock()
	
	for k, v := range fields {
		if !fn(k, v) {
			break
		}
	}
}

// Merge 合并另一个对象
func (ov *objectValue) Merge(other IObject) {
	other.Range(func(key string, value IValue) bool {
		ov.Set(key, value.Clone())
		return true
	})
}

// SortedPairs 获取排序后的键值对
func (ov *objectValue) SortedPairs() []KeyValuePair {
	ov.mu.RLock()
	defer ov.mu.RUnlock()
	
	pairs := make([]KeyValuePair, 0, len(ov.fields))
	for k, v := range ov.fields {
		pairs = append(pairs, KeyValuePair{Key: k, Value: v})
	}
	
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key < pairs[j].Key
	})
	
	return pairs
}

// reset 重置对象（用于对象池）
func (ov *objectValue) reset() {
	ov.mu.Lock()
	defer ov.mu.Unlock()
	
	// 清空字段但保留底层map
	for k := range ov.fields {
		delete(ov.fields, k)
	}
}

// createValueFromInterface 从interface{}创建IValue
func createValueFromInterface(v interface{}) (IValue, error) {
	switch val := v.(type) {
	case nil:
		return NewNullValue(), nil
	case bool:
		return NewBoolValue(val), nil
	case string:
		return NewStringValue(val), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return NewNumberValue(val)
	case map[string]interface{}:
		return NewObjectFromMap(val)
	case []interface{}:
		return NewArrayFromSlice(val)
	default:
		return nil, NewTypeError("supported type", fmt.Sprintf("%T", v), v)
	}
}
