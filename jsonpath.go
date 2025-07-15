package xyJson

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// JSONPathQuery JSONPath查询实现
type JSONPathQuery struct {
	config JSONPathConfig
	cache  map[string]*compiledPath
}

// compiledPath 编译后的路径
type compiledPath struct {
	original string
	steps    []pathStep
}

// pathStep 路径步骤
type pathStep struct {
	type_     stepType
	key       string
	index     int
	start     int
	end       int
	filter    string
	recursive bool
}

// stepType 步骤类型
type stepType int

const (
	stepRoot stepType = iota
	stepKey
	stepIndex
	stepSlice
	stepWildcard
	stepFilter
	stepRecursive
)

// NewJSONPathQuery 创建新的JSONPath查询器
func NewJSONPathQuery() IPathQuery {
	return &JSONPathQuery{
		config: GetGlobalConfig().JSONPath,
		cache:  make(map[string]*compiledPath),
	}
}

// NewJSONPathQueryWithConfig 使用指定配置创建JSONPath查询器
func NewJSONPathQueryWithConfig(config JSONPathConfig) IPathQuery {
	return &JSONPathQuery{
		config: config,
		cache:  make(map[string]*compiledPath),
	}
}

// SelectOne 选择单个值
func (q *JSONPathQuery) SelectOne(root IValue, path string) (IValue, error) {
	results, err := q.SelectAll(root, path)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, NewPathError(path, "no matching values found")
	}
	return results[0], nil
}

// SelectAll 选择所有匹配的值
func (q *JSONPathQuery) SelectAll(root IValue, path string) ([]IValue, error) {
	compiled, err := q.compilePath(path)
	if err != nil {
		return nil, err
	}
	
	results := []IValue{}
	q.executeSteps(root, compiled.steps, 0, &results)
	return results, nil
}

// Set 设置值
func (q *JSONPathQuery) Set(root IValue, path string, value IValue) error {
	compiled, err := q.compilePath(path)
	if err != nil {
		return err
	}
	
	return q.setValue(root, compiled.steps, 0, value)
}

// Delete 删除值
func (q *JSONPathQuery) Delete(root IValue, path string) error {
	compiled, err := q.compilePath(path)
	if err != nil {
		return err
	}
	
	return q.deleteValue(root, compiled.steps, 0)
}

// Exists 检查路径是否存在
func (q *JSONPathQuery) Exists(root IValue, path string) bool {
	results, err := q.SelectAll(root, path)
	return err == nil && len(results) > 0
}

// Count 统计匹配数量
func (q *JSONPathQuery) Count(root IValue, path string) int {
	results, err := q.SelectAll(root, path)
	if err != nil {
		return 0
	}
	return len(results)
}

// compilePath 编译路径
func (q *JSONPathQuery) compilePath(path string) (*compiledPath, error) {
	if compiled, exists := q.cache[path]; exists {
		return compiled, nil
	}
	
	steps, err := q.parsePath(path)
	if err != nil {
		return nil, err
	}
	
	compiled := &compiledPath{
		original: path,
		steps:    steps,
	}
	
	// 缓存编译结果
	if len(q.cache) < q.config.CacheSize {
		q.cache[path] = compiled
	}
	
	return compiled, nil
}

// parsePath 解析路径
func (q *JSONPathQuery) parsePath(path string) ([]pathStep, error) {
	if !strings.HasPrefix(path, "$") {
		return nil, NewPathError(path, "path must start with '$'")
	}
	
	steps := []pathStep{{type_: stepRoot}}
	remaining := path[1:] // 跳过 '$'
	
	for len(remaining) > 0 {
		if strings.HasPrefix(remaining, "..") {
			// 递归下降
			steps = append(steps, pathStep{type_: stepRecursive})
			remaining = remaining[2:]
		} else if strings.HasPrefix(remaining, ".") {
			// 点号访问
			remaining = remaining[1:]
			step, consumed, err := q.parseStep(remaining)
			if err != nil {
				return nil, err
			}
			steps = append(steps, step)
			remaining = remaining[consumed:]
		} else if strings.HasPrefix(remaining, "[") {
			// 括号访问
			step, consumed, err := q.parseBracketStep(remaining)
			if err != nil {
				return nil, err
			}
			steps = append(steps, step)
			remaining = remaining[consumed:]
		} else {
			return nil, NewPathError(path, fmt.Sprintf("unexpected character at position %d", len(path)-len(remaining)))
		}
	}
	
	return steps, nil
}

// parseStep 解析步骤
func (q *JSONPathQuery) parseStep(remaining string) (pathStep, int, error) {
	if len(remaining) == 0 {
		return pathStep{}, 0, NewPathError("", "unexpected end of path")
	}
	
	if remaining[0] == '*' {
		return pathStep{type_: stepWildcard}, 1, nil
	}
	
	// 解析键名
	i := 0
	for i < len(remaining) && remaining[i] != '.' && remaining[i] != '[' {
		i++
	}
	
	if i == 0 {
		return pathStep{}, 0, NewPathError("", "empty key name")
	}
	
	return pathStep{type_: stepKey, key: remaining[:i]}, i, nil
}

// parseBracketStep 解析括号步骤
func (q *JSONPathQuery) parseBracketStep(remaining string) (pathStep, int, error) {
	if !strings.HasPrefix(remaining, "[") {
		return pathStep{}, 0, NewPathError("", "expected '['")
	}
	
	closeIndex := strings.Index(remaining, "]")
	if closeIndex == -1 {
		return pathStep{}, 0, NewPathError("", "missing closing ']'")
	}
	
	content := remaining[1:closeIndex]
	
	// 通配符
	if content == "*" {
		return pathStep{type_: stepWildcard}, closeIndex + 1, nil
	}
	
	// 过滤器
	if strings.HasPrefix(content, "?(") && strings.HasSuffix(content, ")") {
		filter := content[2 : len(content)-1]
		return pathStep{type_: stepFilter, filter: filter}, closeIndex + 1, nil
	}
	
	// 切片
	if strings.Contains(content, ":") {
		parts := strings.Split(content, ":")
		if len(parts) != 2 {
			return pathStep{}, 0, NewPathError("", "invalid slice syntax")
		}
		
		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return pathStep{}, 0, NewPathError("", "invalid slice start")
		}
		
		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return pathStep{}, 0, NewPathError("", "invalid slice end")
		}
		
		return pathStep{type_: stepSlice, start: start, end: end}, closeIndex + 1, nil
	}
	
	// 索引或键名
	if index, err := strconv.Atoi(content); err == nil {
		return pathStep{type_: stepIndex, index: index}, closeIndex + 1, nil
	}
	
	// 字符串键名（去掉引号）
	if (strings.HasPrefix(content, "'") && strings.HasSuffix(content, "'")) ||
		(strings.HasPrefix(content, "\"") && strings.HasSuffix(content, "\"")) {
		key := content[1 : len(content)-1]
		return pathStep{type_: stepKey, key: key}, closeIndex + 1, nil
	}
	
	return pathStep{type_: stepKey, key: content}, closeIndex + 1, nil
}

// executeSteps 执行步骤
func (q *JSONPathQuery) executeSteps(current IValue, steps []pathStep, stepIndex int, results *[]IValue) {
	if stepIndex >= len(steps) {
		*results = append(*results, current)
		return
	}
	
	step := steps[stepIndex]
	
	switch step.type_ {
	case stepRoot:
		q.executeSteps(current, steps, stepIndex+1, results)
	
	case stepKey:
		if current.IsObject() {
			obj := current.(IObject)
			if value, exists := obj.Get(step.key); exists {
				q.executeSteps(value, steps, stepIndex+1, results)
			}
		}
	
	case stepIndex:
		if current.IsArray() {
			arr := current.(IArray)
			index := step.index
			if index < 0 {
				index = arr.Length() + index // 负索引
			}
			if value, exists := arr.Get(index); exists {
				q.executeSteps(value, steps, stepIndex+1, results)
			}
		}
	
	case stepSlice:
		if current.IsArray() {
			arr := current.(IArray)
			start, end := step.start, step.end
			length := arr.Length()
			
			if start < 0 {
				start = length + start
			}
			if end < 0 {
				end = length + end
			}
			if start < 0 {
				start = 0
			}
			if end > length {
				end = length
			}
			
			for i := start; i < end; i++ {
				if value, exists := arr.Get(i); exists {
					q.executeSteps(value, steps, stepIndex+1, results)
				}
			}
		}
	
	case stepWildcard:
		if current.IsObject() {
			obj := current.(IObject)
			obj.Range(func(key string, value IValue) bool {
				q.executeSteps(value, steps, stepIndex+1, results)
				return true
			})
		} else if current.IsArray() {
			arr := current.(IArray)
			arr.Range(func(index int, value IValue) bool {
				q.executeSteps(value, steps, stepIndex+1, results)
				return true
			})
		}
	
	case stepFilter:
		if current.IsArray() {
			arr := current.(IArray)
			arr.Range(func(index int, value IValue) bool {
				if q.evaluateFilter(value, step.filter) {
					q.executeSteps(value, steps, stepIndex+1, results)
				}
				return true
			})
		}
	
	case stepRecursive:
		// 递归搜索
		q.recursiveSearch(current, steps, stepIndex+1, results)
	}
}

// recursiveSearch 递归搜索
func (q *JSONPathQuery) recursiveSearch(current IValue, steps []pathStep, stepIndex int, results *[]IValue) {
	// 在当前节点执行剩余步骤
	q.executeSteps(current, steps, stepIndex, results)
	
	// 递归搜索子节点
	if current.IsObject() {
		obj := current.(IObject)
		obj.Range(func(key string, value IValue) bool {
			q.recursiveSearch(value, steps, stepIndex, results)
			return true
		})
	} else if current.IsArray() {
		arr := current.(IArray)
		arr.Range(func(index int, value IValue) bool {
			q.recursiveSearch(value, steps, stepIndex, results)
			return true
		})
	}
}

// evaluateFilter 评估过滤器
func (q *JSONPathQuery) evaluateFilter(value IValue, filter string) bool {
	// 简单的过滤器实现
	// 支持基本的比较操作
	
	// 匹配 @.key == 'value' 格式
	eqRegex := regexp.MustCompile(`@\.([a-zA-Z_][a-zA-Z0-9_]*) == ['"]([^'"]*)['"]`)
	if matches := eqRegex.FindStringSubmatch(filter); len(matches) == 3 {
		key := matches[1]
		expectedValue := matches[2]
		
		if value.IsObject() {
			obj := value.(IObject)
			if fieldValue, exists := obj.Get(key); exists && fieldValue.IsString() {
				return fieldValue.String() == expectedValue
			}
		}
		return false
	}
	
	// 匹配 @.key > number 格式
	gtRegex := regexp.MustCompile(`@\.([a-zA-Z_][a-zA-Z0-9_]*) > ([0-9.]+)`)
	if matches := gtRegex.FindStringSubmatch(filter); len(matches) == 3 {
		key := matches[1]
		threshold, _ := strconv.ParseFloat(matches[2], 64)
		
		if value.IsObject() {
			obj := value.(IObject)
			if fieldValue, exists := obj.Get(key); exists && fieldValue.IsNumber() {
				if num, err := fieldValue.Float64(); err == nil {
					return num > threshold
				}
			}
		}
		return false
	}
	
	// 匹配 @.key < number 格式
	ltRegex := regexp.MustCompile(`@\.([a-zA-Z_][a-zA-Z0-9_]*) < ([0-9.]+)`)
	if matches := ltRegex.FindStringSubmatch(filter); len(matches) == 3 {
		key := matches[1]
		threshold, _ := strconv.ParseFloat(matches[2], 64)
		
		if value.IsObject() {
			obj := value.(IObject)
			if fieldValue, exists := obj.Get(key); exists && fieldValue.IsNumber() {
				if num, err := fieldValue.Float64(); err == nil {
					return num < threshold
				}
			}
		}
		return false
	}
	
	return false
}

// setValue 设置值
func (q *JSONPathQuery) setValue(root IValue, steps []pathStep, stepIndex int, value IValue) error {
	// 简化实现：只支持简单路径的设置
	if stepIndex >= len(steps)-1 {
		return NewPathError("", "cannot set value at root")
	}
	
	// 这里需要更复杂的实现来支持完整的路径设置
	return NewPathError("", "setValue not fully implemented")
}

// deleteValue 删除值
func (q *JSONPathQuery) deleteValue(root IValue, steps []pathStep, stepIndex int) error {
	// 简化实现：只支持简单路径的删除
	if stepIndex >= len(steps)-1 {
		return NewPathError("", "cannot delete root")
	}
	
	// 这里需要更复杂的实现来支持完整的路径删除
	return NewPathError("", "deleteValue not fully implemented")
}
