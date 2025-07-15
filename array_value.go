package xyJson

import (
	"sync"
	"time"
)

// arrayValue 数组值实现
type arrayValue struct {
	elements []IValue
	mu       sync.RWMutex
}

// NewArray 创建新的数组
func NewArray() IArray {
	return &arrayValue{
		elements: make([]IValue, 0),
	}
}

// NewArrayWithCapacity 创建指定容量的数组
func NewArrayWithCapacity(capacity int) IArray {
	return &arrayValue{
		elements: make([]IValue, 0, capacity),
	}
}

// NewArrayFromSlice 从切片创建数组
func NewArrayFromSlice(slice []interface{}) (IArray, error) {
	arr := NewArray()
	for _, v := range slice {
		value, err := createValueFromInterface(v)
		if err != nil {
			return nil, err
		}
		arr.Append(value)
	}
	return arr, nil
}

// Type 获取值类型
func (av *arrayValue) Type() ValueType {
	return ArrayValueType
}

// Raw 获取原始Go类型值
func (av *arrayValue) Raw() interface{} {
	av.mu.RLock()
	defer av.mu.RUnlock()
	
	result := make([]interface{}, len(av.elements))
	for i, v := range av.elements {
		result[i] = v.Raw()
	}
	return result
}

// String 获取字符串表示
func (av *arrayValue) String() string {
	serializer := NewJSONSerializer()
	result, _ := serializer.SerializeToString(av)
	return result
}

// IsNull 检查是否为null
func (av *arrayValue) IsNull() bool {
	return false
}

// IsString 检查是否为字符串
func (av *arrayValue) IsString() bool {
	return false
}

// IsNumber 检查是否为数字
func (av *arrayValue) IsNumber() bool {
	return false
}

// IsBool 检查是否为布尔值
func (av *arrayValue) IsBool() bool {
	return false
}

// IsObject 检查是否为对象
func (av *arrayValue) IsObject() bool {
	return false
}

// IsArray 检查是否为数组
func (av *arrayValue) IsArray() bool {
	return true
}

// Clone 深拷贝
func (av *arrayValue) Clone() IValue {
	av.mu.RLock()
	defer av.mu.RUnlock()
	
	newArr := NewArrayWithCapacity(len(av.elements))
	for _, v := range av.elements {
		newArr.Append(v.Clone())
	}
	return newArr
}

// Equals 比较是否相等
func (av *arrayValue) Equals(other IValue) bool {
	if other == nil || !other.IsArray() {
		return false
	}
	
	otherArr := other.(IArray)
	if av.Length() != otherArr.Length() {
		return false
	}
	
	av.mu.RLock()
	defer av.mu.RUnlock()
	
	for i, v := range av.elements {
		otherValue, _ := otherArr.Get(i)
		if !v.Equals(otherValue) {
			return false
		}
	}
	return true
}

// Int 转换为int
func (av *arrayValue) Int() (int, error) {
	return 0, NewTypeError("int", "array", av.Raw())
}

// Int64 转换为int64
func (av *arrayValue) Int64() (int64, error) {
	return 0, NewTypeError("int64", "array", av.Raw())
}

// Float64 转换为float64
func (av *arrayValue) Float64() (float64, error) {
	return 0, NewTypeError("float64", "array", av.Raw())
}

// Bool 转换为bool
func (av *arrayValue) Bool() (bool, error) {
	return av.Length() > 0, nil
}

// Time 转换为time.Time
func (av *arrayValue) Time() (time.Time, error) {
	return time.Time{}, NewTypeError("time.Time", "array", av.Raw())
}

// Bytes 转换为[]byte
func (av *arrayValue) Bytes() ([]byte, error) {
	serializer := NewJSONSerializer()
	return serializer.Serialize(av)
}

// Get 获取指定索引的元素
func (av *arrayValue) Get(index int) (IValue, bool) {
	av.mu.RLock()
	defer av.mu.RUnlock()
	
	if index < 0 || index >= len(av.elements) {
		return nil, false
	}
	return av.elements[index], true
}

// Set 设置指定索引的元素
func (av *arrayValue) Set(index int, value IValue) error {
	av.mu.Lock()
	defer av.mu.Unlock()
	
	if index < 0 || index >= len(av.elements) {
		return NewIndexError(index, len(av.elements))
	}
	av.elements[index] = value
	return nil
}

// Append 追加元素
func (av *arrayValue) Append(value IValue) {
	av.mu.Lock()
	defer av.mu.Unlock()
	
	av.elements = append(av.elements, value)
}

// Insert 在指定位置插入元素
func (av *arrayValue) Insert(index int, value IValue) error {
	av.mu.Lock()
	defer av.mu.Unlock()
	
	if index < 0 || index > len(av.elements) {
		return NewIndexError(index, len(av.elements))
	}
	
	// 扩展切片
	av.elements = append(av.elements, nil)
	// 移动元素
	copy(av.elements[index+1:], av.elements[index:])
	// 插入新元素
	av.elements[index] = value
	return nil
}

// Delete 删除指定索引的元素
func (av *arrayValue) Delete(index int) error {
	av.mu.Lock()
	defer av.mu.Unlock()
	
	if index < 0 || index >= len(av.elements) {
		return NewIndexError(index, len(av.elements))
	}
	
	// 移动元素
	copy(av.elements[index:], av.elements[index+1:])
	// 缩短切片
	av.elements = av.elements[:len(av.elements)-1]
	return nil
}

// Length 获取数组长度
func (av *arrayValue) Length() int {
	av.mu.RLock()
	defer av.mu.RUnlock()
	
	return len(av.elements)
}

// Clear 清空数组
func (av *arrayValue) Clear() {
	av.mu.Lock()
	defer av.mu.Unlock()
	
	av.elements = av.elements[:0]
}

// Range 遍历数组元素
func (av *arrayValue) Range(fn func(index int, value IValue) bool) {
	av.mu.RLock()
	elements := make([]IValue, len(av.elements))
	copy(elements, av.elements)
	av.mu.RUnlock()
	
	for i, v := range elements {
		if !fn(i, v) {
			break
		}
	}
}

// reset 重置数组（用于对象池）
func (av *arrayValue) reset() {
	av.mu.Lock()
	defer av.mu.Unlock()
	
	av.elements = av.elements[:0]
}

// AppendAll 批量追加元素
func (av *arrayValue) AppendAll(values []IValue) {
	av.mu.Lock()
	defer av.mu.Unlock()
	
	av.elements = append(av.elements, values...)
}

// IndexOf 查找元素索引
func (av *arrayValue) IndexOf(value IValue) int {
	av.mu.RLock()
	defer av.mu.RUnlock()
	
	for i, v := range av.elements {
		if v.Equals(value) {
			return i
		}
	}
	return -1
}

// Contains 检查是否包含元素
func (av *arrayValue) Contains(value IValue) bool {
	return av.IndexOf(value) >= 0
}

// RemoveValue 移除指定值的元素
func (av *arrayValue) RemoveValue(value IValue) bool {
	index := av.IndexOf(value)
	if index >= 0 {
		av.Delete(index)
		return true
	}
	return false
}

// Slice 获取子数组
func (av *arrayValue) Slice(start, end int) IArray {
	av.mu.RLock()
	defer av.mu.RUnlock()
	
	if start < 0 {
		start = 0
	}
	if end > len(av.elements) {
		end = len(av.elements)
	}
	if start >= end {
		return NewArray()
	}
	
	newArr := NewArrayWithCapacity(end - start)
	for i := start; i < end; i++ {
		newArr.Append(av.elements[i].Clone())
	}
	return newArr
}

// Reverse 反转数组
func (av *arrayValue) Reverse() {
	av.mu.Lock()
	defer av.mu.Unlock()
	
	for i, j := 0, len(av.elements)-1; i < j; i, j = i+1, j-1 {
		av.elements[i], av.elements[j] = av.elements[j], av.elements[i]
	}
}

// Filter 过滤数组元素
func (av *arrayValue) Filter(predicate func(IValue) bool) IArray {
	av.mu.RLock()
	defer av.mu.RUnlock()
	
	newArr := NewArray()
	for _, v := range av.elements {
		if predicate(v) {
			newArr.Append(v.Clone())
		}
	}
	return newArr
}
