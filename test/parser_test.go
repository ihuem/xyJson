package test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	xyJson "github.com/ihuem/xyJson"
	"github.com/ihuem/xyJson/test/testutil"
)

// TestParseString 测试字符串解析功能
// TestParseString tests string parsing functionality
func TestParseString(t *testing.T) {
	tests := []testutil.TableTest{
		{
			Name:     "simple_object",
			Input:    `{"name":"test"}`,
			Expected: map[string]interface{}{"name": "test"},
			Error:    false,
		},
		{
			Name:     "simple_array",
			Input:    `[1,2,3]`,
			Expected: []interface{}{1.0, 2.0, 3.0},
			Error:    false,
		},
		{
			Name:     "boolean_true",
			Input:    `true`,
			Expected: true,
			Error:    false,
		},
		{
			Name:     "boolean_false",
			Input:    `false`,
			Expected: false,
			Error:    false,
		},
		{
			Name:     "null_value",
			Input:    `null`,
			Expected: nil,
			Error:    false,
		},
		{
			Name:     "number_integer",
			Input:    `42`,
			Expected: 42.0,
			Error:    false,
		},
		{
			Name:     "number_float",
			Input:    `3.14`,
			Expected: 3.14,
			Error:    false,
		},
		{
			Name:     "empty_string",
			Input:    `""`,
			Expected: "",
			Error:    false,
		},
		{
			Name:     "invalid_json_missing_quote",
			Input:    `{"name":test}`,
			Expected: nil,
			Error:    true,
		},
		{
			Name:     "invalid_json_trailing_comma",
			Input:    `{"name":"test",}`,
			Expected: nil,
			Error:    true,
		},
		{
			Name:     "empty_input",
			Input:    ``,
			Expected: nil,
			Error:    true,
		},
	}

	testutil.RunTableTests(t, tests, func(input interface{}) (interface{}, error) {
		jsonStr := input.(string)
		value, err := xyJson.ParseString(jsonStr)
		if err != nil {
			return nil, err
		}
		return convertIValueToInterface(value), nil
	})
}

// TestParseComplexObject 测试复杂对象解析
// TestParseComplexObject tests complex object parsing
func TestParseComplexObject(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	complexJSON := generator.GenerateJSONPathTestData()
	value, err := xyJson.ParseString(complexJSON)
	require.NoError(t, err)
	require.NotNil(t, value)
	// 验证解析结果的结构
	assert.Equal(t, xyJson.ObjectValueType, value.Type())

	// 测试嵌套访问
	storeValue, err := xyJson.Get(value, "$.store")
	assert.NoError(t, err)
	assert.Equal(t, xyJson.ObjectValueType, storeValue.Type())

	booksValue, err := xyJson.Get(value, "$.store.book")
	assert.NoError(t, err)
	assert.Equal(t, xyJson.ArrayValueType, booksValue.Type())
	// 检查book数组的长度
	if booksArray, ok := booksValue.(xyJson.IArray); ok {
		assert.Equal(t, 3, booksArray.Length()) // 实际应该是3个书籍
	}

	usersValue, err := xyJson.Get(value, "$.users")
	assert.NoError(t, err)
	assert.Equal(t, xyJson.ArrayValueType, usersValue.Type())
}

// TestParseUnicodeStrings 测试Unicode字符串解析
// TestParseUnicodeStrings tests Unicode string parsing
func TestParseUnicodeStrings(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	unicodeData := generator.GenerateUnicodeJSON()

	jsonStr, err := generator.GenerateJSONString(unicodeData)
	require.NoError(t, err)

	value, err := xyJson.ParseString(jsonStr)
	require.NoError(t, err)
	require.NotNil(t, value)

	// 验证中文字符
	chineseName, err := xyJson.Get(value, "$.中文名称")
	assert.NoError(t, err)
	assert.Equal(t, "张三", chineseName.String())

	// 验证emoji字符
	emojiValue, err := xyJson.Get(value, "$.emoji")
	assert.NoError(t, err)
	assert.Equal(t, "😀🎉🚀", emojiValue.String())

	// 验证转义字符
	escapeValue, err := xyJson.Get(value, "$.escape")
	assert.NoError(t, err)
	assert.Contains(t, escapeValue.String(), "\n")
	assert.Contains(t, escapeValue.String(), "\t")
}

// TestParseInvalidJSON 测试无效JSON处理
// TestParseInvalidJSON tests invalid JSON handling
func TestParseInvalidJSON(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	invalidJSONs := generator.GenerateInvalidJSON()

	for _, invalidJSON := range invalidJSONs {
		t.Run("invalid_"+invalidJSON, func(t *testing.T) {
			value, err := xyJson.ParseString(invalidJSON)
			assert.Error(t, err, "Should return error for invalid JSON: %s", invalidJSON)
			assert.Nil(t, value, "Should return nil value for invalid JSON")

			// 验证错误类型
			assert.IsType(t, &xyJson.JSONError{}, err)
		})
	}
}

// TestParseEdgeCases 测试边界情况
// TestParseEdgeCases tests edge cases
func TestParseEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		hasError bool
	}{
		{
			name:     "very_large_number",
			input:    `1.7976931348623157e+308`,
			expected: 1.7976931348623157e+308,
			hasError: false,
		},
		{
			name:     "very_small_number",
			input:    `4.9406564584124654e-324`,
			expected: 4.9406564584124654e-324,
			hasError: false,
		},
		{
			name:     "negative_zero",
			input:    `-0`,
			expected: 0.0,
			hasError: false,
		},
		{
			name:     "empty_object",
			input:    `{}`,
			expected: map[string]interface{}{},
			hasError: false,
		},
		{
			name:     "empty_array",
			input:    `[]`,
			expected: []interface{}{},
			hasError: false,
		},
		{
			name:     "whitespace_only",
			input:    `   `,
			expected: nil,
			hasError: true,
		},
		{
			name:     "long_string",
			input:    `"` + strings.Repeat("a", 10000) + `"`,
			expected: strings.Repeat("a", 10000),
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := xyJson.ParseString(tt.input)

			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, value)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, value)
				actual := convertIValueToInterface(value)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

// TestParseMaxDepth 测试最大深度限制
// TestParseMaxDepth tests maximum depth limitation
func TestParseMaxDepth(t *testing.T) {
	parser := xyJson.NewParser()
	parser.SetMaxDepth(5)

	// 创建深度为6的嵌套对象
	deepJSON := `{"a":{"b":{"c":{"d":{"e":{"f":"too_deep"}}}}}}`

	value, err := parser.ParseString(deepJSON)
	assert.Error(t, err)
	assert.Nil(t, value)
	assert.Contains(t, err.Error(), "maximum depth")

	// 测试深度为5的对象应该成功
	validJSON := `{"a":{"b":{"c":{"d":{"e":"ok"}}}}}`
	value, err = parser.ParseString(validJSON)
	assert.NoError(t, err)
	assert.NotNil(t, value)
}

// TestParseWithCustomFactory 测试自定义工厂
// TestParseWithCustomFactory tests parsing with custom factory
func TestParseWithCustomFactory(t *testing.T) {
	factory := xyJson.NewValueFactory()
	parser := xyJson.NewParserWithFactory(factory)

	jsonStr := `{"name":"test","age":25}`
	value, err := parser.ParseString(jsonStr)

	assert.NoError(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, xyJson.ObjectValueType, value.Type())
}

// TestParsePerformanceMonitoring 测试性能监控集成
// TestParsePerformanceMonitoring tests performance monitoring integration
func TestParsePerformanceMonitoring(t *testing.T) {
	// 启用性能监控
	xyJson.EnablePerformanceMonitoring()
	monitor := xyJson.GetGlobalMonitor()
	monitor.Reset()

	// 创建更复杂的JSON字符串
	complexJSON := `{
		"users": [
			{"id": 1, "name": "Alice", "email": "alice@example.com", "age": 25, "active": true},
			{"id": 2, "name": "Bob", "email": "bob@example.com", "age": 30, "active": false},
			{"id": 3, "name": "Charlie", "email": "charlie@example.com", "age": 35, "active": true}
		],
		"metadata": {
			"total": 3,
			"page": 1,
			"limit": 10,
			"timestamp": "2024-01-01T00:00:00Z"
		},
		"config": {
			"debug": true,
			"version": "1.0.0",
			"features": ["auth", "logging", "monitoring"]
		}
	}`

	// 执行多次解析操作以产生可测量的时间
	for i := 0; i < 50; i++ {
		value, err := xyJson.ParseString(complexJSON)
		assert.NoError(t, err)
		assert.NotNil(t, value)
	}

	// 验证性能统计
	stats := monitor.GetStats()
	assert.Greater(t, stats.ParseCount, int64(0))
	assert.Greater(t, stats.TotalParseTime, time.Duration(0))
	assert.Equal(t, int64(50), stats.ParseCount)
}

// TestParseConcurrency 测试并发解析
// TestParseConcurrency tests concurrent parsing
func TestParseConcurrency(t *testing.T) {
	// 使用固定的JSON字符串避免并发问题
	jsonStr := `{
		"store": {
			"book": [
				{
					"category": "reference",
					"author": "Nigel Rees",
					"title": "Sayings of the Century",
					"price": 8.95
				},
				{
					"category": "fiction",
					"author": "Evelyn Waugh",
					"title": "Sword of Honour",
					"price": 12.99
				}
			],
			"bicycle": {
				"color": "red",
				"price": 19.95
			}
		},
		"users": [
			{
				"name": "张三",
				"age": 25,
				"city": "北京"
			},
			{
				"name": "李四",
				"age": 30,
				"city": "上海"
			}
		]
	}`

	testutil.RunConcurrently(t, 10, 100, func(goroutineID, iteration int) {
		value, err := xyJson.ParseString(jsonStr)
		if err != nil {
			t.Errorf("Goroutine %d, iteration %d: Parse error: %v", goroutineID, iteration, err)
			return
		}
		if value == nil {
			t.Errorf("Goroutine %d, iteration %d: Parsed value is nil", goroutineID, iteration)
			return
		}
		if value.Type() != xyJson.ObjectValueType {
			t.Errorf("Goroutine %d, iteration %d: Expected ObjectValueType, got %v", goroutineID, iteration, value.Type())
			return
		}
	})
}

// TestParseMemoryUsage 测试内存使用
// TestParseMemoryUsage tests memory usage
func TestParseMemoryUsage(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	performanceData := generator.GeneratePerformanceTestData()

	// 测试小型JSON内存使用
	testutil.AssertNoMemoryLeak(t, func() {
		for i := 0; i < 1000; i++ {
			value, err := xyJson.ParseString(performanceData["small"])
			assert.NoError(t, err)
			assert.NotNil(t, value)
		}
	})
}

// TestParseErrorRecovery 测试错误恢复
// TestParseErrorRecovery tests error recovery
func TestParseErrorRecovery(t *testing.T) {
	// 解析无效JSON后，解析器应该能够正常处理有效JSON
	invalidJSON := `{"name":}`
	validJSON := `{"name":"test"}`

	// 第一次解析无效JSON
	value1, err1 := xyJson.ParseString(invalidJSON)
	assert.Error(t, err1)
	assert.Nil(t, value1)

	// 第二次解析有效JSON应该成功
	value2, err2 := xyJson.ParseString(validJSON)
	assert.NoError(t, err2)
	assert.NotNil(t, value2)
	assert.Equal(t, xyJson.ObjectValueType, value2.Type())
}

// convertIValueToInterface 将IValue转换为interface{}用于测试比较
// convertIValueToInterface converts IValue to interface{} for test comparison
func convertIValueToInterface(value xyJson.IValue) interface{} {
	if value == nil {
		return nil
	}

	switch value.Type() {
	case xyJson.NullValueType:
		return nil
	case xyJson.BoolValueType:
		if scalarValue, ok := value.(xyJson.IScalarValue); ok {
			boolVal, _ := scalarValue.Bool()
			return boolVal
		}
		return nil
	case xyJson.NumberValueType:
		if scalarValue, ok := value.(xyJson.IScalarValue); ok {
			floatVal, _ := scalarValue.Float64()
			return floatVal
		}
		return nil
	case xyJson.StringValueType:
		return value.String()
	case xyJson.ArrayValueType:
		if arrayValue, ok := value.(xyJson.IArray); ok {
			arr := make([]interface{}, arrayValue.Length())
			for i := 0; i < arrayValue.Length(); i++ {
				item := arrayValue.Get(i)
				arr[i] = convertIValueToInterface(item)
			}
			return arr
		}
		return nil
	case xyJson.ObjectValueType:
		if objectValue, ok := value.(xyJson.IObject); ok {
			obj := make(map[string]interface{})
			keys := objectValue.Keys()
			for _, key := range keys {
				val := objectValue.Get(key)
				obj[key] = convertIValueToInterface(val)
			}
			return obj
		}
		return nil
	default:
		return nil
	}
}

// TestParseStringEscaping 测试字符串转义处理
// TestParseStringEscaping tests string escaping handling
func TestParseStringEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic_escape",
			input:    `"Hello\nWorld"`,
			expected: "Hello\nWorld",
		},
		{
			name:     "tab_escape",
			input:    `"Hello\tWorld"`,
			expected: "Hello\tWorld",
		},
		{
			name:     "quote_escape",
			input:    `"Say \"Hello\""`,
			expected: `Say "Hello"`,
		},
		{
			name:     "backslash_escape",
			input:    `"Path\\to\\file"`,
			expected: `Path\to\file`,
		},
		{
			name:     "unicode_escape",
			input:    `"\u0048\u0065\u006C\u006C\u006F"`,
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := xyJson.ParseString(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, value.String())
		})
	}
}

// TestParseNumberFormats 测试数字格式解析
// TestParseNumberFormats tests number format parsing
func TestParseNumberFormats(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "integer",
			input:    `42`,
			expected: 42.0,
		},
		{
			name:     "negative_integer",
			input:    `-42`,
			expected: -42.0,
		},
		{
			name:     "float",
			input:    `3.14`,
			expected: 3.14,
		},
		{
			name:     "negative_float",
			input:    `-3.14`,
			expected: -3.14,
		},
		{
			name:     "scientific_notation",
			input:    `1.23e10`,
			expected: 1.23e10,
		},
		{
			name:     "scientific_negative_exp",
			input:    `1.23e-10`,
			expected: 1.23e-10,
		},
		{
			name:     "zero",
			input:    `0`,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := xyJson.ParseString(tt.input)
			assert.NoError(t, err)
			if scalarValue, ok := value.(xyJson.IScalarValue); ok {
				floatVal, err := scalarValue.Float64()
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, floatVal)
			}
		})
	}
}
