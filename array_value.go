package xyJson

import (
	"sync"
)

// arrayValue JSON数组实现
// arrayValue implements the IArray interface
type arrayValue struct {
	data []IValue
	mu   sync.RWMutex
}

// NewArray 创建新的JSON数组
// NewArray creates a new JSON array
func NewArray() IArray {
	return &arrayValue{
		data: make([]IValue, 0, DefaultArrayCapacity),
	}
}

// NewArrayWithCapacity 创建指定容量的JSON数组
// NewArrayWithCapacity creates a JSON array with specified capacity
func NewArrayWithCapacity(capacity int) IArray {
	if capacity <= 0 {
		capacity = DefaultArrayCapacity
	}
	return &arrayValue{
		data: make([]IValue, 0, capacity),
	}
}

// NewArrayFromSlice 从切片创建JSON数组
// NewArrayFromSlice creates a JSON array from a slice
func NewArrayFromSlice(slice []interface{}) (IArray, error) {
	arr := NewArrayWithCapacity(len(slice))
	factory := NewValueFactory()

	for _, item := range slice {
		value, err := factory.CreateFromRaw(item)
		if err != nil {
			return nil, err
		}
		if err := arr.Append(value); err != nil {
			return nil, err
		}
	}

	return arr, nil
}

// Type 返回值的类型
// Type returns the type of the value
func (av *arrayValue) Type() ValueType {
	return ArrayValueType
}

// Raw 返回原始Go类型值
// Raw returns the raw Go type value
func (av *arrayValue) Raw() interface{} {
	av.mu.RLock()
	defer av.mu.RUnlock()

	result := make([]interface{}, len(av.data))
	for i, value := range av.data {
		result[i] = value.Raw()
	}
	return result
}

// String 返回字符串表示
// String returns the string representation
func (av *arrayValue) String() string {
	// 数组的字符串表示通常是JSON格式，这里简化为类型名
	return "[object Array]"
}

// IsNull 检查是否为null值
// IsNull checks if the value is null
func (av *arrayValue) IsNull() bool {
	return false // 数组永远不为null
}

// Clone 创建值的深拷贝
// Clone creates a deep copy of the value
func (av *arrayValue) Clone() IValue {
	av.mu.RLock()
	defer av.mu.RUnlock()

	newArr := NewArrayWithCapacity(len(av.data))
	for _, value := range av.data {
		newArr.Append(value.Clone())
	}
	return newArr
}

// Equals 比较两个值是否相等
// Equals compares if two values are equal
func (av *arrayValue) Equals(other IValue) bool {
	if other == nil || other.Type() != ArrayValueType {
		return false
	}

	otherArr, ok := other.(IArray)
	if !ok {
		return false
	}

	av.mu.RLock()
	defer av.mu.RUnlock()

	if av.Length() != otherArr.Length() {
		return false
	}

	for i, value := range av.data {
		otherValue := otherArr.Get(i)
		if otherValue == nil || !value.Equals(otherValue) {
			return false
		}
	}

	return true
}

// Get 根据索引获取值
// Get retrieves a value by index
func (av *arrayValue) Get(index int) IValue {
	av.mu.RLock()
	defer av.mu.RUnlock()

	if index < 0 || index >= len(av.data) {
		return nil
	}

	return av.data[index]
}

// Set 设置指定索引的值
// Set sets the value at the specified index
func (av *arrayValue) Set(index int, value interface{}) error {
	av.mu.Lock()
	defer av.mu.Unlock()

	if index < 0 || index >= len(av.data) {
		return NewIndexOutOfRangeError(index, len(av.data), "")
	}

	var jsonValue IValue
	switch v := value.(type) {
	case IValue:
		jsonValue = v
	case nil:
		jsonValue = &scalarValue{valueType: NullValueType, rawData: nil}
	default:
		// 使用工厂创建值
		factory := NewValueFactory()
		var err error
		jsonValue, err = factory.CreateFromRaw(value)
		if err != nil {
			return err
		}
	}

	av.data[index] = jsonValue
	return nil
}

// Append 追加值到数组末尾
// Append adds a value to the end of the array
func (av *arrayValue) Append(value interface{}) error {
	var jsonValue IValue
	switch v := value.(type) {
	case IValue:
		jsonValue = v
	case nil:
		jsonValue = &scalarValue{valueType: NullValueType, rawData: nil}
	default:
		// 使用工厂创建值
		factory := NewValueFactory()
		var err error
		jsonValue, err = factory.CreateFromRaw(value)
		if err != nil {
			return err
		}
	}

	av.mu.Lock()
	defer av.mu.Unlock()

	av.data = append(av.data, jsonValue)
	return nil
}

// Insert 在指定位置插入值
// Insert inserts a value at the specified position
func (av *arrayValue) Insert(index int, value interface{}) error {
	av.mu.Lock()
	defer av.mu.Unlock()

	if index < 0 || index > len(av.data) {
		return NewIndexOutOfRangeError(index, len(av.data), "")
	}

	var jsonValue IValue
	switch v := value.(type) {
	case IValue:
		jsonValue = v
	case nil:
		jsonValue = &scalarValue{valueType: NullValueType, rawData: nil}
	default:
		// 使用工厂创建值
		factory := NewValueFactory()
		var err error
		jsonValue, err = factory.CreateFromRaw(value)
		if err != nil {
			return err
		}
	}

	// 扩展切片
	av.data = append(av.data, nil)
	// 移动元素
	copy(av.data[index+1:], av.data[index:])
	// 插入新值
	av.data[index] = jsonValue

	return nil
}

// Delete 删除指定索引的值
// Delete removes the value at the specified index
func (av *arrayValue) Delete(index int) error {
	av.mu.Lock()
	defer av.mu.Unlock()

	if index < 0 || index >= len(av.data) {
		return NewIndexOutOfRangeError(index, len(av.data), "")
	}

	// 移动元素并缩短切片
	copy(av.data[index:], av.data[index+1:])
	av.data = av.data[:len(av.data)-1]

	return nil
}

// Length 返回数组长度
// Length returns the length of the array
func (av *arrayValue) Length() int {
	av.mu.RLock()
	defer av.mu.RUnlock()

	return len(av.data)
}

// Clear 清空数组
// Clear removes all elements from the array
func (av *arrayValue) Clear() {
	av.mu.Lock()
	defer av.mu.Unlock()

	// 重置切片但保留容量
	av.data = av.data[:0]
}

// Range 遍历数组元素
// Range iterates over array elements
func (av *arrayValue) Range(fn func(index int, value IValue) bool) {
	if fn == nil {
		return
	}

	av.mu.RLock()
	// 创建数据的副本以避免在遍历时持有锁
	dataCopy := make([]IValue, len(av.data))
	copy(dataCopy, av.data)
	av.mu.RUnlock()

	for i, value := range dataCopy {
		if !fn(i, value) {
			break
		}
	}
}

// reset 重置数组状态（用于对象池）
// reset resets the array state (for object pool)
func (av *arrayValue) reset() {
	av.mu.Lock()
	defer av.mu.Unlock()

	// 清空数据但保留底层切片的容量
	av.data = av.data[:0]
}

// AppendAll 批量追加多个值
// AppendAll appends multiple values at once
func (av *arrayValue) AppendAll(values ...interface{}) error {
	factory := NewValueFactory()
	jsonValues := make([]IValue, 0, len(values))

	// 先转换所有值
	for _, value := range values {
		var jsonValue IValue
		switch v := value.(type) {
		case IValue:
			jsonValue = v
		case nil:
			jsonValue = &scalarValue{valueType: NullValueType, rawData: nil}
		default:
			var err error
			jsonValue, err = factory.CreateFromRaw(value)
			if err != nil {
				return err
			}
		}
		jsonValues = append(jsonValues, jsonValue)
	}

	// 批量追加
	av.mu.Lock()
	defer av.mu.Unlock()

	av.data = append(av.data, jsonValues...)
	return nil
}

// IndexOf 查找值的索引
// IndexOf finds the index of a value
func (av *arrayValue) IndexOf(value IValue) int {
	if value == nil {
		return -1
	}

	av.mu.RLock()
	defer av.mu.RUnlock()

	for i, item := range av.data {
		if item != nil && item.Equals(value) {
			return i
		}
	}

	return -1
}

// Contains 检查是否包含指定值
// Contains checks if the array contains a value
func (av *arrayValue) Contains(value IValue) bool {
	return av.IndexOf(value) >= 0
}

// RemoveValue 删除第一个匹配的值
// RemoveValue removes the first matching value
func (av *arrayValue) RemoveValue(value IValue) bool {
	index := av.IndexOf(value)
	if index >= 0 {
		return av.Delete(index) == nil
	}
	return false
}

// Slice 获取子数组
// Slice gets a sub-array
func (av *arrayValue) Slice(start, end int) (IArray, error) {
	av.mu.RLock()
	defer av.mu.RUnlock()

	length := len(av.data)

	// 处理负数索引
	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}

	// 边界检查
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start > end {
		return nil, NewInvalidOperationError("slice", "start index greater than end index")
	}

	// 创建新数组
	newArr := NewArrayWithCapacity(end - start)
	for i := start; i < end; i++ {
		newArr.Append(av.data[i])
	}

	return newArr, nil
}

// Reverse 反转数组
// Reverse reverses the array
func (av *arrayValue) Reverse() {
	av.mu.Lock()
	defer av.mu.Unlock()

	length := len(av.data)
	for i := 0; i < length/2; i++ {
		av.data[i], av.data[length-1-i] = av.data[length-1-i], av.data[i]
	}
}

// Filter 过滤数组元素
// Filter filters array elements
func (av *arrayValue) Filter(predicate func(index int, value IValue) bool) IArray {
	if predicate == nil {
		return NewArray()
	}

	av.mu.RLock()
	defer av.mu.RUnlock()

	result := NewArray()
	for i, value := range av.data {
		if predicate(i, value) {
			result.Append(value)
		}
	}

	return result
}
