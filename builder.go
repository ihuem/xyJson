package xyJson

import (
	"time"
)

// JSONBuilder JSON构建器，支持链式调用
// JSONBuilder is a JSON builder that supports method chaining
type JSONBuilder struct {
	factory IValueFactory
	root    IValue
	current IValue
	path    []string
	err     error
}

// NewJSONBuilder 创建新的JSON构建器
// NewJSONBuilder creates a new JSON builder
func NewJSONBuilder() *JSONBuilder {
	factory := NewValueFactory()
	root := factory.CreateObject()
	return &JSONBuilder{
		factory: factory,
		root:    root,
		current: root,
		path:    make([]string, 0),
	}
}

// NewJSONBuilderWithFactory 使用指定工厂创建JSON构建器
// NewJSONBuilderWithFactory creates a JSON builder with the specified factory
func NewJSONBuilderWithFactory(factory IValueFactory) *JSONBuilder {
	if factory == nil {
		factory = NewValueFactory()
	}
	root := factory.CreateObject()
	return &JSONBuilder{
		factory: factory,
		root:    root,
		current: root,
		path:    make([]string, 0),
	}
}

// NewArrayBuilder 创建数组构建器
// NewArrayBuilder creates an array builder
func NewArrayBuilder() *JSONBuilder {
	factory := NewValueFactory()
	root := factory.CreateArray()
	return &JSONBuilder{
		factory: factory,
		root:    root,
		current: root,
		path:    make([]string, 0),
	}
}

// SetString 设置字符串值（链式调用）
// SetString sets a string value (method chaining)
func (b *JSONBuilder) SetString(key, value string) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if obj, ok := b.current.(IObject); ok {
		if err := obj.Set(key, b.factory.CreateString(value)); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ObjectValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// SetInt 设置整数值（链式调用）
// SetInt sets an integer value (method chaining)
func (b *JSONBuilder) SetInt(key string, value int) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if obj, ok := b.current.(IObject); ok {
		numValue, err := b.factory.CreateNumber(value)
		if err != nil {
			b.err = err
			return b
		}
		if err := obj.Set(key, numValue); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ObjectValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// SetInt64 设置64位整数值（链式调用）
// SetInt64 sets a 64-bit integer value (method chaining)
func (b *JSONBuilder) SetInt64(key string, value int64) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if obj, ok := b.current.(IObject); ok {
		numValue, err := b.factory.CreateNumber(value)
		if err != nil {
			b.err = err
			return b
		}
		if err := obj.Set(key, numValue); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ObjectValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// SetFloat64 设置64位浮点数值（链式调用）
// SetFloat64 sets a 64-bit float value (method chaining)
func (b *JSONBuilder) SetFloat64(key string, value float64) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if obj, ok := b.current.(IObject); ok {
		numValue, err := b.factory.CreateNumber(value)
		if err != nil {
			b.err = err
			return b
		}
		if err := obj.Set(key, numValue); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ObjectValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// SetNumber 设置数字值（链式调用）
// SetNumber sets a number value (method chaining)
func (b *JSONBuilder) SetNumber(key string, value interface{}) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if obj, ok := b.current.(IObject); ok {
		numValue, err := b.factory.CreateNumber(value)
		if err != nil {
			b.err = err
			return b
		}
		if err := obj.Set(key, numValue); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ObjectValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// SetBool 设置布尔值（链式调用）
// SetBool sets a boolean value (method chaining)
func (b *JSONBuilder) SetBool(key string, value bool) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if obj, ok := b.current.(IObject); ok {
		if err := obj.Set(key, b.factory.CreateBool(value)); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ObjectValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// SetNull 设置null值（链式调用）
// SetNull sets a null value (method chaining)
func (b *JSONBuilder) SetNull(key string) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if obj, ok := b.current.(IObject); ok {
		if err := obj.Set(key, b.factory.CreateNull()); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ObjectValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// SetTime 设置时间值（链式调用）
// SetTime sets a time value (method chaining)
func (b *JSONBuilder) SetTime(key string, value time.Time) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if obj, ok := b.current.(IObject); ok {
		timeStr := value.Format(time.RFC3339)
		if err := obj.Set(key, b.factory.CreateString(timeStr)); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ObjectValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// SetValue 设置任意值（链式调用）
// SetValue sets any value (method chaining)
func (b *JSONBuilder) SetValue(key string, value interface{}) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if obj, ok := b.current.(IObject); ok {
		jsonValue, err := b.factory.CreateFromRaw(value)
		if err != nil {
			b.err = err
			return b
		}
		if err := obj.Set(key, jsonValue); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ObjectValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// AddString 向数组添加字符串值（链式调用）
// AddString adds a string value to array (method chaining)
func (b *JSONBuilder) AddString(value string) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if arr, ok := b.current.(IArray); ok {
		if err := arr.Append(b.factory.CreateString(value)); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ArrayValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// AddInt 向数组添加整数值（链式调用）
// AddInt adds an integer value to array (method chaining)
func (b *JSONBuilder) AddInt(value int) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if arr, ok := b.current.(IArray); ok {
		numValue, err := b.factory.CreateNumber(value)
		if err != nil {
			b.err = err
			return b
		}
		if err := arr.Append(numValue); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ArrayValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// AddBool 向数组添加布尔值（链式调用）
// AddBool adds a boolean value to array (method chaining)
func (b *JSONBuilder) AddBool(value bool) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if arr, ok := b.current.(IArray); ok {
		if err := arr.Append(b.factory.CreateBool(value)); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ArrayValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// AddNull 向数组添加null值（链式调用）
// AddNull adds a null value to array (method chaining)
func (b *JSONBuilder) AddNull() *JSONBuilder {
	if b.err != nil {
		return b
	}

	if arr, ok := b.current.(IArray); ok {
		if err := arr.Append(b.factory.CreateNull()); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ArrayValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// AddValue 向数组添加任意值（链式调用）
// AddValue adds any value to array (method chaining)
func (b *JSONBuilder) AddValue(value interface{}) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if arr, ok := b.current.(IArray); ok {
		jsonValue, err := b.factory.CreateFromRaw(value)
		if err != nil {
			b.err = err
			return b
		}
		if err := arr.Append(jsonValue); err != nil {
			b.err = err
		}
	} else {
		b.err = NewTypeMismatchError(ArrayValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// BeginObject 开始构建嵌套对象（链式调用）
// BeginObject starts building a nested object (method chaining)
func (b *JSONBuilder) BeginObject(key string) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if obj, ok := b.current.(IObject); ok {
		newObj := b.factory.CreateObject()
		if err := obj.Set(key, newObj); err != nil {
			b.err = err
			return b
		}
		b.current = newObj
		b.path = append(b.path, key)
	} else {
		b.err = NewTypeMismatchError(ObjectValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// BeginArray 开始构建嵌套数组（链式调用）
// BeginArray starts building a nested array (method chaining)
func (b *JSONBuilder) BeginArray(key string) *JSONBuilder {
	if b.err != nil {
		return b
	}

	if obj, ok := b.current.(IObject); ok {
		newArr := b.factory.CreateArray()
		if err := obj.Set(key, newArr); err != nil {
			b.err = err
			return b
		}
		b.current = newArr
		b.path = append(b.path, key)
	} else {
		b.err = NewTypeMismatchError(ObjectValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// AddObject 向数组添加对象（链式调用）
// AddObject adds an object to array (method chaining)
func (b *JSONBuilder) AddObject() *JSONBuilder {
	if b.err != nil {
		return b
	}

	if arr, ok := b.current.(IArray); ok {
		newObj := b.factory.CreateObject()
		if err := arr.Append(newObj); err != nil {
			b.err = err
			return b
		}
		b.current = newObj
		b.path = append(b.path, "[]") // 表示数组元素
	} else {
		b.err = NewTypeMismatchError(ArrayValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// AddArray 向数组添加数组（链式调用）
// AddArray adds an array to array (method chaining)
func (b *JSONBuilder) AddArray() *JSONBuilder {
	if b.err != nil {
		return b
	}

	if arr, ok := b.current.(IArray); ok {
		newArr := b.factory.CreateArray()
		if err := arr.Append(newArr); err != nil {
			b.err = err
			return b
		}
		b.current = newArr
		b.path = append(b.path, "[]") // 表示数组元素
	} else {
		b.err = NewTypeMismatchError(ArrayValueType, b.current.Type(), b.getCurrentPath())
	}
	return b
}

// End 结束当前嵌套层级（链式调用）
// End ends the current nesting level (method chaining)
func (b *JSONBuilder) End() *JSONBuilder {
	if b.err != nil {
		return b
	}

	if len(b.path) > 0 {
		// 回到上一级
		b.path = b.path[:len(b.path)-1]
		b.current = b.navigateToPath(b.root, b.path)
	} else {
		// 已经在根级别
		b.current = b.root
	}
	return b
}

// Build 构建最终的JSON对象
// Build builds the final JSON object
func (b *JSONBuilder) Build() (IValue, error) {
	if b.err != nil {
		return nil, b.err
	}
	return b.root, nil
}

// MustBuild 构建JSON对象，如果有错误则panic
// MustBuild builds the JSON object, panics if there's an error
func (b *JSONBuilder) MustBuild() IValue {
	result, err := b.Build()
	if err != nil {
		panic(err)
	}
	return result
}

// Error 获取构建过程中的错误
// Error gets the error occurred during building
func (b *JSONBuilder) Error() error {
	return b.err
}

// Reset 重置构建器状态
// Reset resets the builder state
func (b *JSONBuilder) Reset() *JSONBuilder {
	b.root = b.factory.CreateObject()
	b.current = b.root
	b.path = b.path[:0]
	b.err = nil
	return b
}

// ResetAsArray 重置为数组构建器
// ResetAsArray resets as an array builder
func (b *JSONBuilder) ResetAsArray() *JSONBuilder {
	b.root = b.factory.CreateArray()
	b.current = b.root
	b.path = b.path[:0]
	b.err = nil
	return b
}

// getCurrentPath 获取当前路径字符串
// getCurrentPath gets the current path string
func (b *JSONBuilder) getCurrentPath() string {
	if len(b.path) == 0 {
		return "$"
	}
	result := "$"
	for _, segment := range b.path {
		if segment == "[]" {
			result += "[]"
		} else {
			result += "." + segment
		}
	}
	return result
}

// navigateToPath 导航到指定路径
// navigateToPath navigates to the specified path
func (b *JSONBuilder) navigateToPath(root IValue, path []string) IValue {
	current := root
	for _, segment := range path {
		if segment == "[]" {
			// 数组元素，需要特殊处理
			if arr, ok := current.(IArray); ok {
				length := arr.Length()
				if length > 0 {
					current = arr.Get(length - 1) // 获取最后一个元素
				}
			}
		} else {
			// 对象属性
			if obj, ok := current.(IObject); ok {
				current = obj.Get(segment)
			}
		}
		if current == nil {
			break
		}
	}
	return current
}
