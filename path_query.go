package xyJson

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// pathQuery JSONPath查询实现
// pathQuery implements JSONPath query functionality
type pathQuery struct {
	factory IValueFactory
}

// pathSegment 路径段
// pathSegment represents a path segment
type pathSegment struct {
	Type      SegmentType
	Key       string
	Index     int
	Filter    *pathFilter
	Wildcard  bool
	Recursive bool
}

// pathFilter 路径过滤器
// pathFilter represents a path filter
type pathFilter struct {
	Expression string
	Operator   string
	Value      interface{}
	Compiled   *regexp.Regexp
}

// NewPathQuery 创建新的JSONPath查询器
// NewPathQuery creates a new JSONPath query
func NewPathQuery() IPathQuery {
	return &pathQuery{
		factory: NewValueFactory(),
	}
}

// NewPathQueryWithFactory 使用指定工厂创建JSONPath查询器
// NewPathQueryWithFactory creates a JSONPath query with specified factory
func NewPathQueryWithFactory(factory IValueFactory) IPathQuery {
	if factory == nil {
		factory = NewValueFactory()
	}
	return &pathQuery{
		factory: factory,
	}
}

// SelectOne 根据路径选择单个值
// SelectOne selects a single value by path
func (pq *pathQuery) SelectOne(root IValue, path string) (IValue, error) {
	if root == nil {
		return nil, NewPathNotFoundError(path)
	}

	if path == "" || path == "$" {
		return root, nil
	}

	segments, err := pq.parsePath(path)
	if err != nil {
		return nil, err
	}

	results := pq.executeQuery(root, segments, false)
	if len(results) == 0 {
		return nil, NewPathNotFoundError(path)
	}

	return results[0], nil
}

// SelectAll 根据路径选择所有匹配的值
// SelectAll selects all matching values by path
func (pq *pathQuery) SelectAll(root IValue, path string) ([]IValue, error) {
	if root == nil {
		return nil, NewPathNotFoundError(path)
	}

	if path == "" || path == "$" {
		return []IValue{root}, nil
	}

	segments, err := pq.parsePath(path)
	if err != nil {
		return nil, err
	}

	return pq.executeQuery(root, segments, true), nil
}

// Set 根据路径设置值
// Set sets a value by path
func (pq *pathQuery) Set(root IValue, path string, value IValue) error {
	if root == nil {
		return NewPathNotFoundError(path)
	}

	if path == "" || path == "$" {
		return NewInvalidJSONError("cannot set root value", nil)
	}

	segments, err := pq.parsePath(path)
	if err != nil {
		return err
	}

	return pq.setValueAtPath(root, segments, value)
}

// Delete 根据路径删除值
// Delete deletes a value by path
func (pq *pathQuery) Delete(root IValue, path string) error {
	if root == nil {
		return NewPathNotFoundError(path)
	}

	if path == "" || path == "$" {
		return NewInvalidJSONError("cannot delete root value", nil)
	}

	segments, err := pq.parsePath(path)
	if err != nil {
		return err
	}

	return pq.deleteValueAtPath(root, segments)
}

// Exists 检查路径是否存在
// Exists checks if a path exists
func (pq *pathQuery) Exists(root IValue, path string) bool {
	if root == nil {
		return false
	}

	if path == "" || path == "$" {
		return true
	}

	segments, err := pq.parsePath(path)
	if err != nil {
		return false
	}

	results := pq.executeQuery(root, segments, false)
	return len(results) > 0
}

// Count 统计匹配路径的数量
// Count counts the number of matching paths
func (pq *pathQuery) Count(root IValue, path string) int {
	if root == nil {
		return 0
	}

	if path == "" || path == "$" {
		return 1
	}

	segments, err := pq.parsePath(path)
	if err != nil {
		return 0
	}

	results := pq.executeQuery(root, segments, true)
	return len(results)
}

// parsePath 解析JSONPath路径
// parsePath parses a JSONPath string
func (pq *pathQuery) parsePath(path string) ([]*pathSegment, error) {
	if !strings.HasPrefix(path, "$") {
		return nil, NewInvalidJSONError("path must start with '$'", nil)
	}

	path = path[1:] // 移除 '$'
	if path == "" {
		return []*pathSegment{}, nil
	}

	var segments []*pathSegment
	i := 0

	for i < len(path) {
		if path[i] == '.' {
			i++ // 跳过 '.'
			if i >= len(path) {
				break
			}

			// 处理递归下降 '..'
			if i < len(path) && path[i] == '.' {
				i++ // 跳过第二个 '.'
				segments = append(segments, &pathSegment{
					Type:      PropertySegmentType,
					Recursive: true,
				})
				continue
			}

			// 处理通配符 '*'
			if i < len(path) && path[i] == '*' {
				i++
				segments = append(segments, &pathSegment{
					Type:     PropertySegmentType,
					Wildcard: true,
				})
				continue
			}

			// 解析属性名
			start := i
			for i < len(path) && path[i] != '.' && path[i] != '[' {
				i++
			}
			if i > start {
				key := path[start:i]
				segments = append(segments, &pathSegment{
					Type: PropertySegmentType,
					Key:  key,
				})
			}
		} else if path[i] == '[' {
			// 解析数组索引或过滤器
			segment, newPos, err := pq.parseBracketExpression(path, i)
			if err != nil {
				return nil, err
			}
			segments = append(segments, segment)
			i = newPos
		} else {
			// 直接属性名（无前导点）
			start := i
			for i < len(path) && path[i] != '.' && path[i] != '[' {
				i++
			}
			if i > start {
				key := path[start:i]
				segments = append(segments, &pathSegment{
					Type: PropertySegmentType,
					Key:  key,
				})
			}
		}
	}

	return segments, nil
}

// parseBracketExpression 解析方括号表达式
// parseBracketExpression parses bracket expressions
func (pq *pathQuery) parseBracketExpression(path string, start int) (*pathSegment, int, error) {
	if path[start] != '[' {
		return nil, start, NewInvalidJSONError("expected '['", nil)
	}

	end := strings.Index(path[start:], "]")
	if end == -1 {
		return nil, start, NewInvalidJSONError("unclosed bracket", nil)
	}
	end += start

	expr := path[start+1 : end]
	segment := &pathSegment{}

	// 空表达式或通配符
	if expr == "" || expr == "*" {
		segment.Type = IndexSegmentType
		segment.Wildcard = true
		return segment, end + 1, nil
	}

	// 数字索引
	if index, err := strconv.Atoi(expr); err == nil {
		segment.Type = IndexSegmentType
		segment.Index = index
		return segment, end + 1, nil
	}

	// 字符串键（带引号）
	if (strings.HasPrefix(expr, "'") && strings.HasSuffix(expr, "'")) ||
		(strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"")) {
		key := expr[1 : len(expr)-1]
		segment.Type = PropertySegmentType
		segment.Key = key
		return segment, end + 1, nil
	}

	// 过滤器表达式
	if strings.HasPrefix(expr, "?") {
		filter, err := pq.parseFilter(expr[1:])
		if err != nil {
			return nil, start, err
		}
		segment.Type = FilterSegmentType
		segment.Filter = filter
		return segment, end + 1, nil
	}

	// 切片表达式（简化实现）
	if strings.Contains(expr, ":") {
		// 暂时不支持切片，返回通配符
		segment.Type = IndexSegmentType
		segment.Wildcard = true
		return segment, end + 1, nil
	}

	return nil, start, NewInvalidJSONError("invalid bracket expression: "+expr, nil)
}

// parseFilter 解析过滤器表达式
// parseFilter parses filter expressions
func (pq *pathQuery) parseFilter(expr string) (*pathFilter, error) {
	// 简化的过滤器解析
	// 支持基本的比较操作：==, !=, <, >, <=, >=
	operators := []string{"==", "!=", "<=", ">=", "<", ">"}

	for _, op := range operators {
		if idx := strings.Index(expr, op); idx != -1 {
			left := strings.TrimSpace(expr[:idx])
			right := strings.TrimSpace(expr[idx+len(op):])

			// 解析右侧值
			var value interface{}
			if strings.HasPrefix(right, "'") && strings.HasSuffix(right, "'") {
				value = right[1 : len(right)-1] // 字符串
			} else if strings.HasPrefix(right, "\"") && strings.HasSuffix(right, "\"") {
				value = right[1 : len(right)-1] // 字符串
			} else if num, err := strconv.ParseFloat(right, 64); err == nil {
				value = num // 数字
			} else if right == "true" {
				value = true
			} else if right == "false" {
				value = false
			} else if right == "null" {
				value = nil
			} else {
				value = right // 默认为字符串
			}

			return &pathFilter{
				Expression: left,
				Operator:   op,
				Value:      value,
			}, nil
		}
	}

	return nil, NewInvalidJSONError("invalid filter expression: "+expr, nil)
}

// executeQuery 执行查询
// executeQuery executes the query
func (pq *pathQuery) executeQuery(root IValue, segments []*pathSegment, selectAll bool) []IValue {
	if len(segments) == 0 {
		return []IValue{root}
	}

	current := []IValue{root}

	for _, segment := range segments {
		var next []IValue

		for _, value := range current {
			if value == nil {
				continue
			}

			switch segment.Type {
			case PropertySegmentType:
				next = append(next, pq.selectProperty(value, segment, selectAll)...)
			case IndexSegmentType:
				next = append(next, pq.selectIndex(value, segment, selectAll)...)
			case FilterSegmentType:
				next = append(next, pq.selectFilter(value, segment, selectAll)...)
			}

			// 递归下降
			if segment.Recursive {
				next = append(next, pq.selectRecursive(value, segment, selectAll)...)
			}
		}

		current = next
		if !selectAll && len(current) > 0 {
			break
		}
	}

	return current
}

// selectProperty 选择属性
// selectProperty selects properties
func (pq *pathQuery) selectProperty(value IValue, segment *pathSegment, selectAll bool) []IValue {
	if segment.Wildcard {
		// 通配符：选择所有属性
		switch v := value.(type) {
		case IObject:
			var results []IValue
			for _, key := range v.Keys() {
				if val := v.Get(key); val != nil {
					results = append(results, val)
					if !selectAll {
						break
					}
				}
			}
			return results
		case IArray:
			var results []IValue
			for i := 0; i < v.Length(); i++ {
				if val := v.Get(i); val != nil {
					results = append(results, val)
					if !selectAll {
						break
					}
				}
			}
			return results
		}
	} else {
		// 具体属性名
		if obj, ok := value.(IObject); ok {
			if val := obj.Get(segment.Key); val != nil {
				return []IValue{val}
			}
		}
	}

	return nil
}

// selectIndex 选择索引
// selectIndex selects by index
func (pq *pathQuery) selectIndex(value IValue, segment *pathSegment, selectAll bool) []IValue {
	arr, ok := value.(IArray)
	if !ok {
		return nil
	}

	if segment.Wildcard {
		// 通配符：选择所有元素
		var results []IValue
		for i := 0; i < arr.Length(); i++ {
			if val := arr.Get(i); val != nil {
				results = append(results, val)
				if !selectAll {
					break
				}
			}
		}
		return results
	} else {
		// 具体索引
		index := segment.Index
		if index < 0 {
			index = arr.Length() + index // 负索引
		}
		if index >= 0 && index < arr.Length() {
			if val := arr.Get(index); val != nil {
				return []IValue{val}
			}
		}
	}

	return nil
}

// selectFilter 选择过滤器匹配的值
// selectFilter selects values matching the filter
func (pq *pathQuery) selectFilter(value IValue, segment *pathSegment, selectAll bool) []IValue {
	arr, ok := value.(IArray)
	if !ok {
		return nil
	}

	var results []IValue
	for i := 0; i < arr.Length(); i++ {
		elem := arr.Get(i)
		if elem != nil && pq.evaluateFilter(elem, segment.Filter) {
			results = append(results, elem)
			if !selectAll {
				break
			}
		}
	}

	return results
}

// selectRecursive 递归选择
// selectRecursive recursively selects values
func (pq *pathQuery) selectRecursive(value IValue, segment *pathSegment, selectAll bool) []IValue {
	var results []IValue

	switch v := value.(type) {
	case IObject:
		for _, key := range v.Keys() {
			if val := v.Get(key); val != nil {
				results = append(results, pq.selectRecursive(val, segment, selectAll)...)
				if !selectAll && len(results) > 0 {
					break
				}
			}
		}
	case IArray:
		for i := 0; i < v.Length(); i++ {
			if val := v.Get(i); val != nil {
				results = append(results, pq.selectRecursive(val, segment, selectAll)...)
				if !selectAll && len(results) > 0 {
					break
				}
			}
		}
	}

	return results
}

// evaluateFilter 评估过滤器
// evaluateFilter evaluates a filter
func (pq *pathQuery) evaluateFilter(value IValue, filter *pathFilter) bool {
	if filter == nil {
		return true
	}

	// 获取要比较的值
	var compareValue interface{}
	if filter.Expression == "@" {
		// 当前值
		compareValue = value.Raw()
	} else {
		// 属性值
		if obj, ok := value.(IObject); ok {
			if val := obj.Get(filter.Expression); val != nil {
				compareValue = val.Raw()
			}
		}
	}

	// 执行比较
	return pq.compareValues(compareValue, filter.Operator, filter.Value)
}

// compareValues 比较值
// compareValues compares two values
func (pq *pathQuery) compareValues(left interface{}, operator string, right interface{}) bool {
	switch operator {
	case "==":
		return pq.valuesEqual(left, right)
	case "!=":
		return !pq.valuesEqual(left, right)
	case "<", "<=", ">", ">=":
		return pq.compareNumeric(left, operator, right)
	default:
		return false
	}
}

// valuesEqual 判断值是否相等
// valuesEqual checks if values are equal
func (pq *pathQuery) valuesEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}

	// 类型转换和比较
	switch l := left.(type) {
	case string:
		if r, ok := right.(string); ok {
			return l == r
		}
	case float64:
		if r, ok := right.(float64); ok {
			return l == r
		}
		if r, ok := right.(int64); ok {
			return l == float64(r)
		}
	case int64:
		if r, ok := right.(int64); ok {
			return l == r
		}
		if r, ok := right.(float64); ok {
			return float64(l) == r
		}
	case bool:
		if r, ok := right.(bool); ok {
			return l == r
		}
	}

	return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
}

// compareNumeric 数值比较
// compareNumeric performs numeric comparison
func (pq *pathQuery) compareNumeric(left interface{}, operator string, right interface{}) bool {
	// 转换为float64进行比较
	leftNum, leftOk := pq.toFloat64(left)
	rightNum, rightOk := pq.toFloat64(right)

	if !leftOk || !rightOk {
		return false
	}

	switch operator {
	case "<":
		return leftNum < rightNum
	case "<=":
		return leftNum <= rightNum
	case ">":
		return leftNum > rightNum
	case ">=":
		return leftNum >= rightNum
	default:
		return false
	}
}

// toFloat64 转换为float64
// toFloat64 converts to float64
func (pq *pathQuery) toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int64:
		return float64(v), true
	case int:
		return float64(v), true
	case string:
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num, true
		}
	}
	return 0, false
}

// setValueAtPath 在指定路径设置值
// setValueAtPath sets value at the specified path
func (pq *pathQuery) setValueAtPath(root IValue, segments []*pathSegment, value IValue) error {
	if len(segments) == 0 {
		return NewInvalidJSONError("cannot set root value", nil)
	}

	current := root
	for i, segment := range segments[:len(segments)-1] {
		next, err := pq.navigateSegment(current, segment)
		if err != nil {
			return err
		}
		if next == nil {
			// 创建中间路径
			next, err = pq.createIntermediatePath(current, segment, segments[i+1])
			if err != nil {
				return err
			}
		}
		current = next
	}

	// 设置最终值
	lastSegment := segments[len(segments)-1]
	return pq.setFinalValue(current, lastSegment, value)
}

// deleteValueAtPath 删除指定路径的值
// deleteValueAtPath deletes value at the specified path
func (pq *pathQuery) deleteValueAtPath(root IValue, segments []*pathSegment) error {
	if len(segments) == 0 {
		return NewInvalidJSONError("cannot delete root value", nil)
	}

	current := root
	for _, segment := range segments[:len(segments)-1] {
		next, err := pq.navigateSegment(current, segment)
		if err != nil || next == nil {
			return NewPathNotFoundError("path not found")
		}
		current = next
	}

	// 删除最终值
	lastSegment := segments[len(segments)-1]
	return pq.deleteFinalValue(current, lastSegment)
}

// navigateSegment 导航到下一个段
// navigateSegment navigates to the next segment
func (pq *pathQuery) navigateSegment(value IValue, segment *pathSegment) (IValue, error) {
	switch segment.Type {
	case PropertySegmentType:
		if obj, ok := value.(IObject); ok {
			return obj.Get(segment.Key), nil
		}
		return nil, NewTypeMismatchError(ObjectValueType, value.Type(), "")
	case IndexSegmentType:
		if arr, ok := value.(IArray); ok {
			index := segment.Index
			if index < 0 {
				index = arr.Length() + index
			}
			if index >= 0 && index < arr.Length() {
				return arr.Get(index), nil
			}
			return nil, NewIndexOutOfRangeError(index, arr.Length(), "array index out of range")
		}
		return nil, NewTypeMismatchError(ArrayValueType, value.Type(), "")
	default:
		return nil, NewInvalidJSONError("unsupported segment type", nil)
	}
}

// createIntermediatePath 创建中间路径
// createIntermediatePath creates intermediate path
func (pq *pathQuery) createIntermediatePath(parent IValue, current *pathSegment, next *pathSegment) (IValue, error) {
	var newValue IValue

	// 根据下一个段的类型决定创建对象还是数组
	if next.Type == IndexSegmentType {
		newValue = pq.factory.CreateArray()
	} else {
		newValue = pq.factory.CreateObject()
	}

	// 设置到父级
	switch current.Type {
	case PropertySegmentType:
		if obj, ok := parent.(IObject); ok {
			return newValue, obj.Set(current.Key, newValue)
		}
		return nil, NewTypeMismatchError(ObjectValueType, parent.Type(), "")
	case IndexSegmentType:
		if arr, ok := parent.(IArray); ok {
			// 扩展数组到所需大小
			for arr.Length() <= current.Index {
				if err := arr.Append(pq.factory.CreateNull()); err != nil {
					return nil, err
				}
			}
			return newValue, arr.Set(current.Index, newValue)
		}
		return nil, NewTypeMismatchError(ArrayValueType, parent.Type(), "")
	default:
		return nil, NewInvalidJSONError("unsupported segment type", nil)
	}
}

// setFinalValue 设置最终值
// setFinalValue sets the final value
func (pq *pathQuery) setFinalValue(parent IValue, segment *pathSegment, value IValue) error {
	switch segment.Type {
	case PropertySegmentType:
		if obj, ok := parent.(IObject); ok {
			return obj.Set(segment.Key, value)
		}
		return NewTypeMismatchError(ObjectValueType, parent.Type(), "")
	case IndexSegmentType:
		if arr, ok := parent.(IArray); ok {
			index := segment.Index
			if index < 0 {
				index = arr.Length() + index
			}
			// 扩展数组到所需大小
			for arr.Length() <= index {
				if err := arr.Append(pq.factory.CreateNull()); err != nil {
					return err
				}
			}
			return arr.Set(index, value)
		}
		return NewTypeMismatchError(ArrayValueType, parent.Type(), "")
	default:
		return NewInvalidJSONError("unsupported segment type", nil)
	}
}

// deleteFinalValue 删除最终值
// deleteFinalValue deletes the final value
func (pq *pathQuery) deleteFinalValue(parent IValue, segment *pathSegment) error {
	switch segment.Type {
	case PropertySegmentType:
		if obj, ok := parent.(IObject); ok {
			obj.Delete(segment.Key)
			return nil
		}
		return NewTypeMismatchError(ObjectValueType, parent.Type(), "")
	case IndexSegmentType:
		if arr, ok := parent.(IArray); ok {
			index := segment.Index
			if index < 0 {
				index = arr.Length() + index
			}
			if index >= 0 && index < arr.Length() {
				arr.Delete(index)
				return nil
			}
			return NewIndexOutOfRangeError(index, arr.Length(), "array index out of range")
		}
		return NewTypeMismatchError(ArrayValueType, parent.Type(), "")
	default:
		return NewInvalidJSONError("unsupported segment type", nil)
	}
}
