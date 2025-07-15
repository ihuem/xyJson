package xyJson

import (
	"sort"
	"sync"
	"time"
)

// objectValue JSON对象实现
// objectValue implements the IObject interface
type objectValue struct {
	data map[string]IValue
	mu   sync.RWMutex
}

// NewObject 创建新的JSON对象
// NewObject creates a new JSON object
func NewObject() IObject {
	return &objectValue{
		data: make(map[string]IValue, DefaultMapCapacity),
	}
}

// NewObjectWithCapacity 创建指定容量的JSON对象
// NewObjectWithCapacity creates a JSON object with specified capacity
func NewObjectWithCapacity(capacity int) IObject {
	if capacity <= 0 {
		capacity = DefaultMapCapacity
	}
	return &objectValue{
		data: make(map[string]IValue, capacity),
	}
}

// Type 返回值的类型
// Type returns the type of the value
func (ov *objectValue) Type() ValueType {
	return ObjectValueType
}

// Raw 返回原始Go类型值
// Raw returns the raw Go type value
func (ov *objectValue) Raw() interface{} {
	ov.mu.RLock()
	defer ov.mu.RUnlock()

	result := make(map[string]interface{}, len(ov.data))
	for key, value := range ov.data {
		result[key] = value.Raw()
	}
	return result
}

// String 返回字符串表示
// String returns the string representation
func (ov *objectValue) String() string {
	// 对象的字符串表示通常是JSON格式，这里简化为类型名
	return "[object Object]"
}

// IsNull 检查是否为null值
// IsNull checks if the value is null
func (ov *objectValue) IsNull() bool {
	return false // 对象永远不为null
}

// Clone 创建值的深拷贝
// Clone creates a deep copy of the value
func (ov *objectValue) Clone() IValue {
	ov.mu.RLock()
	defer ov.mu.RUnlock()

	newObj := NewObjectWithCapacity(len(ov.data))
	for key, value := range ov.data {
		newObj.Set(key, value.Clone())
	}
	return newObj
}

// Equals 比较两个值是否相等
// Equals compares if two values are equal
func (ov *objectValue) Equals(other IValue) bool {
	if other == nil || other.Type() != ObjectValueType {
		return false
	}

	otherObj, ok := other.(IObject)
	if !ok {
		return false
	}

	ov.mu.RLock()
	defer ov.mu.RUnlock()

	if ov.Size() != otherObj.Size() {
		return false
	}

	for key, value := range ov.data {
		otherValue := otherObj.Get(key)
		if otherValue == nil || !value.Equals(otherValue) {
			return false
		}
	}

	return true
}

// Get 根据键名获取值
// Get retrieves a value by key
func (ov *objectValue) Get(key string) IValue {
	ov.mu.RLock()
	defer ov.mu.RUnlock()

	return ov.data[key]
}

// Set 设置键值对
// Set sets a key-value pair
func (ov *objectValue) Set(key string, value interface{}) error {
	if key == "" {
		return NewInvalidOperationError("set object key", "key cannot be empty")
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

	ov.mu.Lock()
	defer ov.mu.Unlock()

	ov.data[key] = jsonValue
	return nil
}

// Delete 删除指定键
// Delete removes the specified key
func (ov *objectValue) Delete(key string) bool {
	ov.mu.Lock()
	defer ov.mu.Unlock()

	if _, exists := ov.data[key]; exists {
		delete(ov.data, key)
		return true
	}
	return false
}

// Has 检查是否包含指定键
// Has checks if the object contains the specified key
func (ov *objectValue) Has(key string) bool {
	ov.mu.RLock()
	defer ov.mu.RUnlock()

	_, exists := ov.data[key]
	return exists
}

// Keys 返回所有键名
// Keys returns all key names
func (ov *objectValue) Keys() []string {
	ov.mu.RLock()
	defer ov.mu.RUnlock()

	keys := make([]string, 0, len(ov.data))
	for key := range ov.data {
		keys = append(keys, key)
	}

	// 对键名进行排序，确保结果的一致性
	sort.Strings(keys)
	return keys
}

// Size 返回键值对数量
// Size returns the number of key-value pairs
func (ov *objectValue) Size() int {
	ov.mu.RLock()
	defer ov.mu.RUnlock()

	return len(ov.data)
}

// Clear 清空所有键值对
// Clear removes all key-value pairs
func (ov *objectValue) Clear() {
	ov.mu.Lock()
	defer ov.mu.Unlock()

	// 创建新的map而不是逐个删除，更高效
	ov.data = make(map[string]IValue, DefaultMapCapacity)
}

// Range 遍历所有键值对
// Range iterates over all key-value pairs
func (ov *objectValue) Range(fn func(key string, value IValue) bool) {
	if fn == nil {
		return
	}

	ov.mu.RLock()
	// 创建键的副本以避免在遍历时持有锁
	keys := make([]string, 0, len(ov.data))
	for key := range ov.data {
		keys = append(keys, key)
	}
	ov.mu.RUnlock()

	// 对键进行排序以确保遍历顺序的一致性
	sort.Strings(keys)

	for _, key := range keys {
		ov.mu.RLock()
		value, exists := ov.data[key]
		ov.mu.RUnlock()

		if exists {
			if !fn(key, value) {
				break
			}
		}
	}
}

// reset 重置对象状态（用于对象池）
// reset resets the object state (for object pool)
func (ov *objectValue) reset() {
	ov.mu.Lock()
	defer ov.mu.Unlock()

	// 清空数据但保留底层map的容量
	for key := range ov.data {
		delete(ov.data, key)
	}
}

// GetSorted 按键名排序返回所有键值对
// GetSorted returns all key-value pairs sorted by key
func (ov *objectValue) GetSorted() []struct {
	Key   string
	Value IValue
} {
	ov.mu.RLock()
	defer ov.mu.RUnlock()

	result := make([]struct {
		Key   string
		Value IValue
	}, 0, len(ov.data))

	keys := make([]string, 0, len(ov.data))
	for key := range ov.data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		result = append(result, struct {
			Key   string
			Value IValue
		}{
			Key:   key,
			Value: ov.data[key],
		})
	}

	return result
}

// Merge 合并另一个对象的键值对
// Merge merges key-value pairs from another object
func (ov *objectValue) Merge(other IObject) error {
	if other == nil {
		return NewNullPointerError("merge object")
	}

	other.Range(func(key string, value IValue) bool {
		ov.Set(key, value)
		return true
	})

	return nil
}

// AsString 将值转换为字符串，对象类型返回空字符串
// AsString converts the value to string, returns empty string for object type
func (ov *objectValue) AsString() string {
	return ""
}

// AsInt 将值转换为整数，对象类型返回0
// AsInt converts the value to integer, returns 0 for object type
func (ov *objectValue) AsInt() int {
	return 0
}

// AsInt64 将值转换为64位整数，对象类型返回0
// AsInt64 converts the value to 64-bit integer, returns 0 for object type
func (ov *objectValue) AsInt64() int64 {
	return 0
}

// AsFloat64 将值转换为64位浮点数，对象类型返回0.0
// AsFloat64 converts the value to 64-bit float, returns 0.0 for object type
func (ov *objectValue) AsFloat64() float64 {
	return 0.0
}

// AsBool 将值转换为布尔值，对象类型返回false
// AsBool converts the value to boolean, returns false for object type
func (ov *objectValue) AsBool() bool {
	return false
}

// AsBytes 将值转换为字节数组，对象类型返回nil
// AsBytes converts the value to byte array, returns nil for object type
func (ov *objectValue) AsBytes() []byte {
	return nil
}

// AsTime 将值转换为时间，对象类型返回零时间
// AsTime converts the value to time, returns zero time for object type
func (ov *objectValue) AsTime() time.Time {
	return time.Time{}
}

// AsObject 将值转换为对象，对象类型返回自身
// AsObject converts the value to object, returns self for object type
func (ov *objectValue) AsObject() IObject {
	return ov
}

// AsArray 将值转换为数组，对象类型返回nil
// AsArray converts the value to array, returns nil for object type
func (ov *objectValue) AsArray() IArray {
	return nil
}
