package xyJson

import (
	"fmt"
	"time"
)

// JSONBuilder JSON构建器实现
type JSONBuilder struct {
	root    IValue
	stack   []builderFrame
	current builderFrame
	err     error
}

// builderFrame 构建器帧
type builderFrame struct {
	value    IValue
	type_    ValueType
	key      string
	isArray  bool
	parent   *builderFrame
}

// NewJSONBuilder 创建新的JSON构建器（对象模式）
func NewJSONBuilder() IBuilder {
	return &JSONBuilder{
		current: builderFrame{
			value:   NewObject(),
			type_:   ObjectValueType,
			isArray: false,
		},
	}
}

// NewArrayBuilder 创建新的JSON构建器（数组模式）
func NewArrayBuilder() IBuilder {
	return &JSONBuilder{
		current: builderFrame{
			value:   NewArray(),
			type_:   ArrayValueType,
			isArray: true,
		},
	}
}

// NewBuilderWithCapacity 创建指定容量的构建器
func NewBuilderWithCapacity(capacity int) IBuilder {
	return &JSONBuilder{
		current: builderFrame{
			value:   NewObjectWithCapacity(capacity),
			type_:   ObjectValueType,
			isArray: false,
		},
	}
}

// SetString 设置字符串字段
func (b *JSONBuilder) SetString(key, value string) IBuilder {
	return b.setValue(key, NewStringValue(value))
}

// SetInt 设置整数字段
func (b *JSONBuilder) SetInt(key string, value int) IBuilder {
	numValue, err := NewNumberValue(value)
	if err != nil {
		b.err = err
		return b
	}
	return b.setValue(key, numValue)
}

// SetInt64 设置64位整数字段
func (b *JSONBuilder) SetInt64(key string, value int64) IBuilder {
	numValue, err := NewNumberValue(value)
	if err != nil {
		b.err = err
		return b
	}
	return b.setValue(key, numValue)
}

// SetFloat64 设置浮点数字段
func (b *JSONBuilder) SetFloat64(key string, value float64) IBuilder {
	numValue, err := NewNumberValue(value)
	if err != nil {
		b.err = err
		return b
	}
	return b.setValue(key, numValue)
}

// SetBool 设置布尔字段
func (b *JSONBuilder) SetBool(key string, value bool) IBuilder {
	return b.setValue(key, NewBoolValue(value))
}

// SetNull 设置null字段
func (b *JSONBuilder) SetNull(key string) IBuilder {
	return b.setValue(key, NewNullValue())
}

// SetTime 设置时间字段
func (b *JSONBuilder) SetTime(key string, value time.Time) IBuilder {
	return b.setValue(key, NewStringValue(value.Format(time.RFC3339)))
}

// SetValue 设置任意值字段
func (b *JSONBuilder) SetValue(key string, value IValue) IBuilder {
	return b.setValue(key, value)
}

// AddString 向数组添加字符串
func (b *JSONBuilder) AddString(value string) IBuilder {
	return b.addValue(NewStringValue(value))
}

// AddInt 向数组添加整数
func (b *JSONBuilder) AddInt(value int) IBuilder {
	numValue, err := NewNumberValue(value)
	if err != nil {
		b.err = err
		return b
	}
	return b.addValue(numValue)
}

// AddBool 向数组添加布尔值
func (b *JSONBuilder) AddBool(value bool) IBuilder {
	return b.addValue(NewBoolValue(value))
}

// AddNull 向数组添加null
func (b *JSONBuilder) AddNull() IBuilder {
	return b.addValue(NewNullValue())
}

// AddValue 向数组添加任意值
func (b *JSONBuilder) AddValue(value IValue) IBuilder {
	return b.addValue(value)
}

// BeginObject 开始构建嵌套对象
func (b *JSONBuilder) BeginObject(key string) IBuilder {
	if b.err != nil {
		return b
	}
	
	if b.current.isArray {
		b.err = fmt.Errorf("cannot set object key '%s' in array context", key)
		return b
	}
	
	newObj := NewObject()
	b.pushFrame(builderFrame{
		value:   newObj,
		type_:   ObjectValueType,
		key:     key,
		isArray: false,
		parent:  &b.current,
	})
	
	return b
}

// BeginArray 开始构建嵌套数组
func (b *JSONBuilder) BeginArray(key string) IBuilder {
	if b.err != nil {
		return b
	}
	
	if b.current.isArray {
		b.err = fmt.Errorf("cannot set array key '%s' in array context", key)
		return b
	}
	
	newArr := NewArray()
	b.pushFrame(builderFrame{
		value:   newArr,
		type_:   ArrayValueType,
		key:     key,
		isArray: true,
		parent:  &b.current,
	})
	
	return b
}

// AddObject 向数组添加对象
func (b *JSONBuilder) AddObject() IBuilder {
	if b.err != nil {
		return b
	}
	
	if !b.current.isArray {
		b.err = fmt.Errorf("cannot add object to non-array context")
		return b
	}
	
	newObj := NewObject()
	b.pushFrame(builderFrame{
		value:   newObj,
		type_:   ObjectValueType,
		isArray: false,
		parent:  &b.current,
	})
	
	return b
}

// AddArray 向数组添加数组
func (b *JSONBuilder) AddArray() IBuilder {
	if b.err != nil {
		return b
	}
	
	if !b.current.isArray {
		b.err = fmt.Errorf("cannot add array to non-array context")
		return b
	}
	
	newArr := NewArray()
	b.pushFrame(builderFrame{
		value:   newArr,
		type_:   ArrayValueType,
		isArray: true,
		parent:  &b.current,
	})
	
	return b
}

// End 结束当前嵌套层级
func (b *JSONBuilder) End() IBuilder {
	if b.err != nil {
		return b
	}
	
	if len(b.stack) == 0 {
		b.err = fmt.Errorf("no nested context to end")
		return b
	}
	
	// 将当前值添加到父级
	parent := b.current.parent
	if parent == nil {
		b.err = fmt.Errorf("invalid parent context")
		return b
	}
	
	if parent.isArray {
		// 添加到父数组
		parentArr := parent.value.(IArray)
		parentArr.Append(b.current.value)
	} else {
		// 添加到父对象
		parentObj := parent.value.(IObject)
		parentObj.Set(b.current.key, b.current.value)
	}
	
	// 弹出栈帧
	b.popFrame()
	
	return b
}

// Build 构建最终值
func (b *JSONBuilder) Build() (IValue, error) {
	if b.err != nil {
		return nil, b.err
	}
	
	if len(b.stack) > 0 {
		return nil, fmt.Errorf("unclosed nested contexts: %d", len(b.stack))
	}
	
	return b.current.value, nil
}

// MustBuild 构建最终值，失败时panic
func (b *JSONBuilder) MustBuild() IValue {
	value, err := b.Build()
	if err != nil {
		panic(err)
	}
	return value
}

// Error 获取构建过程中的错误
func (b *JSONBuilder) Error() error {
	return b.err
}

// Reset 重置构建器
func (b *JSONBuilder) Reset() IBuilder {
	b.root = nil
	b.stack = nil
	b.current = builderFrame{
		value:   NewObject(),
		type_:   ObjectValueType,
		isArray: false,
	}
	b.err = nil
	return b
}

// ResetAsArray 重置为数组构建器
func (b *JSONBuilder) ResetAsArray() IBuilder {
	b.root = nil
	b.stack = nil
	b.current = builderFrame{
		value:   NewArray(),
		type_:   ArrayValueType,
		isArray: true,
	}
	b.err = nil
	return b
}

// setValue 设置值（内部方法）
func (b *JSONBuilder) setValue(key string, value IValue) IBuilder {
	if b.err != nil {
		return b
	}
	
	if b.current.isArray {
		b.err = fmt.Errorf("cannot set key '%s' in array context", key)
		return b
	}
	
	obj := b.current.value.(IObject)
	obj.Set(key, value)
	
	return b
}

// addValue 添加值（内部方法）
func (b *JSONBuilder) addValue(value IValue) IBuilder {
	if b.err != nil {
		return b
	}
	
	if !b.current.isArray {
		b.err = fmt.Errorf("cannot add value to non-array context")
		return b
	}
	
	arr := b.current.value.(IArray)
	arr.Append(value)
	
	return b
}

// pushFrame 压入栈帧
func (b *JSONBuilder) pushFrame(frame builderFrame) {
	b.stack = append(b.stack, b.current)
	b.current = frame
}

// popFrame 弹出栈帧
func (b *JSONBuilder) popFrame() {
	if len(b.stack) > 0 {
		b.current = b.stack[len(b.stack)-1]
		b.stack = b.stack[:len(b.stack)-1]
	}
}

// getCurrentPath 获取当前路径（调试用）
func (b *JSONBuilder) getCurrentPath() string {
	path := "$"
	for _, frame := range b.stack {
		if frame.key != "" {
			path += "." + frame.key
		} else if frame.isArray {
			path += "[]"
		}
	}
	if b.current.key != "" {
		path += "." + b.current.key
	} else if b.current.isArray {
		path += "[]"
	}
	return path
}

// navigateToPath 导航到指定路径（高级功能）
func (b *JSONBuilder) navigateToPath(path string) error {
	// 这里可以实现路径导航功能
	// 暂时返回未实现错误
	return fmt.Errorf("path navigation not implemented")
}
