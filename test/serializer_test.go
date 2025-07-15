package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	xyJson "github.com/ihuem/xyJson"
	"github.com/ihuem/xyJson/test/testutil"
)

// TestSerializeToString 测试基本序列化功能
// TestSerializeToString tests basic serialization functionality
func TestSerializeToString(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() xyJson.IValue
		expected string
		hasError bool
	}{
		{
			name: "simple_object",
			setup: func() xyJson.IValue {
				obj := xyJson.CreateObject()
				obj.Set("name", "测试")
				obj.Set("age", 25)
				return obj
			},
			expected: `{"name":"测试","age":25}`,
			hasError: false,
		},
		{
			name: "simple_array",
			setup: func() xyJson.IValue {
				arr := xyJson.CreateArray()
				arr.Append(1)
				arr.Append(2)
				arr.Append(3)
				return arr
			},
			expected: `[1,2,3]`,
			hasError: false,
		},
		{
			name: "boolean_true",
			setup: func() xyJson.IValue {
				return xyJson.CreateBool(true)
			},
			expected: `true`,
			hasError: false,
		},
		{
			name: "boolean_false",
			setup: func() xyJson.IValue {
				return xyJson.CreateBool(false)
			},
			expected: `false`,
			hasError: false,
		},
		{
			name: "null_value",
			setup: func() xyJson.IValue {
				return xyJson.CreateNull()
			},
			expected: `null`,
			hasError: false,
		},
		{
			name: "string_value",
			setup: func() xyJson.IValue {
				return xyJson.CreateString("Hello World")
			},
			expected: `"Hello World"`,
			hasError: false,
		},
		{
			name: "number_integer",
			setup: func() xyJson.IValue {
				return xyJson.CreateNumber(42)
			},
			expected: `42`,
			hasError: false,
		},
		{
			name: "number_float",
			setup: func() xyJson.IValue {
				return xyJson.CreateNumber(3.14)
			},
			expected: `3.14`,
			hasError: false,
		},
		{
			name: "empty_object",
			setup: func() xyJson.IValue {
				return xyJson.CreateObject()
			},
			expected: `{}`,
			hasError: false,
		},
		{
			name: "empty_array",
			setup: func() xyJson.IValue {
				return xyJson.CreateArray()
			},
			expected: `[]`,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := tt.setup()
			result, err := xyJson.SerializeToString(value)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				testutil.AssertJSONEqual(t, tt.expected, result)
			}
		})
	}
}

// TestSerializeComplexObject 测试复杂对象序列化
// TestSerializeComplexObject tests complex object serialization
func TestSerializeComplexObject(t *testing.T) {
	// 创建复杂嵌套对象
	obj := xyJson.CreateObject()
	obj.Set("name", "张三")
	obj.Set("age", 25)
	obj.Set("active", true)

	// 添加嵌套对象
	address := xyJson.CreateObject()
	address.Set("city", "北京")
	address.Set("zipcode", "100000")
	obj.Set("address", address)

	// 添加数组
	skills := xyJson.CreateArray()
	skills.Append("Go")
	skills.Append("JSON")
	skills.Append("测试")
	obj.Set("skills", skills)

	// 使用不转义HTML的序列化器
	serializer := xyJson.MinimalSerializer()
	result, err := serializer.SerializeToString(obj)
	require.NoError(t, err)

	// 验证结果是有效的JSON
	testutil.AssertValidJSON(t, result)

	// 验证包含预期内容
	assert.Contains(t, result, "张三")
	assert.Contains(t, result, "北京")
	assert.Contains(t, result, "Go")
	assert.Contains(t, result, "测试")

	// 反序列化验证
	parsed, err := xyJson.ParseString(result)
	require.NoError(t, err)

	name, err := xyJson.Get(parsed, "$.name")
	assert.NoError(t, err)
	assert.Equal(t, "张三", name.String())

	city, err := xyJson.Get(parsed, "$.address.city")
	assert.NoError(t, err)
	assert.Equal(t, "北京", city.String())
}

// TestSerializeFormatOptions 测试格式化选项
// TestSerializeFormatOptions tests formatting options
func TestSerializeFormatOptions(t *testing.T) {
	obj := xyJson.CreateObject()
	obj.Set("name", "test")
	obj.Set("age", 25)

	t.Run("compact_format", func(t *testing.T) {
		result, err := xyJson.Compact(obj)
		assert.NoError(t, err)
		assert.NotContains(t, result, "\n")
		assert.NotContains(t, result, "  ")
	})

	t.Run("pretty_format", func(t *testing.T) {
		serializer := xyJson.PrettySerializer("  ")
		result, err := serializer.SerializeToString(obj)
		assert.NoError(t, err)
		assert.Contains(t, result, "\n")
		assert.Contains(t, result, "  ")
	})

	t.Run("html_safe_format", func(t *testing.T) {
		htmlObj := xyJson.CreateObject()
		htmlObj.Set("html", "<script>alert('test')</script>")

		serializer := xyJson.HTMLSafeSerializer()
		result, err := serializer.SerializeToString(htmlObj)
		assert.NoError(t, err)
		assert.NotContains(t, result, "<script>")
		assert.Contains(t, result, "\\u003c")
	})
}

// TestSerializeStringEscaping 测试字符串转义
// TestSerializeStringEscaping tests string escaping
func TestSerializeStringEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic_string",
			input:    "Hello World",
			expected: `"Hello World"`,
		},
		{
			name:     "newline_escape",
			input:    "Hello\nWorld",
			expected: `"Hello\nWorld"`,
		},
		{
			name:     "tab_escape",
			input:    "Hello\tWorld",
			expected: `"Hello\tWorld"`,
		},
		{
			name:     "quote_escape",
			input:    `Say "Hello"`,
			expected: `"Say \"Hello\""`,
		},
		{
			name:     "backslash_escape",
			input:    `Path\to\file`,
			expected: `"Path\\to\\file"`,
		},
		{
			name:     "unicode_characters",
			input:    "Hello 世界 🌍",
			expected: `"Hello 世界 🌍"`,
		},
		{
			name:     "control_characters",
			input:    "\x00\x01\x1f",
			expected: `"\u0000\u0001\u001f"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := xyJson.CreateString(tt.input)
			// 对于unicode字符测试，使用不转义HTML的序列化器
			if tt.name == "unicode_characters" {
				result, err := xyJson.MinimalSerializer().SerializeToString(value)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			} else {
				result, err := xyJson.SerializeToString(value)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestSerializeNumberFormats 测试数字格式序列化
// TestSerializeNumberFormats tests number format serialization
func TestSerializeNumberFormats(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "integer",
			input:    42,
			expected: "42",
		},
		{
			name:     "negative_integer",
			input:    -42,
			expected: "-42",
		},
		{
			name:     "zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "float",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "negative_float",
			input:    -3.14,
			expected: "-3.14",
		},
		{
			name:     "scientific_notation",
			input:    1.23e10,
			expected: "12300000000",
		},
		{
			name:     "very_large_number",
			input:    1.7976931348623157e+308,
			expected: "1.7976931348623157e+308",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := xyJson.CreateNumber(tt.input)
			result, err := xyJson.SerializeToString(value)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSerializeCircularReference 测试循环引用检测
// TestSerializeCircularReference tests circular reference detection
func TestSerializeCircularReference(t *testing.T) {
	// 创建循环引用
	obj1 := xyJson.CreateObject()
	obj2 := xyJson.CreateObject()

	obj1.Set("name", "obj1")
	obj1.Set("ref", obj2)

	obj2.Set("name", "obj2")
	obj2.Set("ref", obj1) // 循环引用

	// 序列化应该检测到循环引用并返回错误
	result, err := xyJson.SerializeToString(obj1)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "circular reference")
}

// TestSerializeMaxDepth 测试最大深度限制
// TestSerializeMaxDepth tests maximum depth limitation
func TestSerializeMaxDepth(t *testing.T) {
	// 创建深度嵌套对象
	obj := xyJson.CreateObject()
	current := obj

	// 创建深度为10的嵌套结构
	for i := 0; i < 100; i++ {
		next := xyJson.CreateObject()
		next.Set("level", xyJson.CreateNumber(float64(i+1)))
		current.Set("next", next)
		current = next
	}

	// 使用默认深度应该成功
	result, err := xyJson.SerializeToString(obj)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// 使用较小的最大深度应该失败
	options := &xyJson.SerializeOptions{
		MaxDepth: 5,
		Compact:  true,
	}
	serializer := xyJson.NewSerializerWithOptions(options)
	result, err = serializer.SerializeToString(obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum")
}

// TestSerializeNilValue 测试nil值处理
// TestSerializeNilValue tests nil value handling
func TestSerializeNilValue(t *testing.T) {
	result, err := xyJson.SerializeToString(nil)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "nil value")
}

// TestSerializePerformanceMonitoring 测试性能监控集成
// TestSerializePerformanceMonitoring tests performance monitoring integration
func TestSerializePerformanceMonitoring(t *testing.T) {
	// 启用性能监控
	xyJson.EnablePerformanceMonitoring()
	monitor := xyJson.GetGlobalMonitor()
	monitor.Reset()

	// 创建测试对象
	obj := xyJson.CreateObject()
	obj.Set("name", xyJson.CreateString("test"))
	obj.Set("age", xyJson.CreateNumber(25))
	obj.Set("active", xyJson.CreateBool(true))

	// 执行序列化操作
	result, err := xyJson.SerializeToString(obj)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// 验证性能统计
	stats := monitor.GetStats()
	assert.Greater(t, stats.SerializeCount, int64(0))
	assert.Greater(t, stats.TotalSerializeTime, int64(0))
}

// TestSerializeConcurrency 测试并发序列化
// TestSerializeConcurrency tests concurrent serialization
func TestSerializeConcurrency(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	complexObj := generator.GenerateComplexObject(3)

	// 先解析为IValue
	jsonStr, err := generator.GenerateJSONString(complexObj)
	require.NoError(t, err)

	value, err := xyJson.ParseString(jsonStr)
	require.NoError(t, err)

	testutil.RunConcurrently(t, 10, 100, func(goroutineID, iteration int) {
		result, err := xyJson.SerializeToString(value)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		testutil.AssertValidJSON(t, result)
	})
}

// TestSerializeMemoryUsage 测试内存使用
// TestSerializeMemoryUsage tests memory usage
func TestSerializeMemoryUsage(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	performanceData := generator.GeneratePerformanceTestData()

	// 解析测试数据
	value, err := xyJson.ParseString(performanceData["medium"])
	require.NoError(t, err)

	// 测试序列化内存使用
	testutil.AssertNoMemoryLeak(t, func() {
		for i := 0; i < 1000; i++ {
			result, err := xyJson.SerializeToString(value)
			assert.NoError(t, err)
			assert.NotEmpty(t, result)
		}
	})
}

// TestSerializeRoundTrip 测试往返转换
// TestSerializeRoundTrip tests round-trip conversion
func TestSerializeRoundTrip(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	testCases := []interface{}{
		generator.GenerateSimpleObject(),
		generator.GenerateComplexObject(3),
		generator.GenerateUnicodeJSON(),
		generator.GenerateLargeArray(100),
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			// 原始数据 -> JSON字符串
			originalJSON, err := generator.GenerateJSONString(testCase)
			require.NoError(t, err)

			// JSON字符串 -> IValue
			value, err := xyJson.ParseString(originalJSON)
			require.NoError(t, err)

			// IValue -> JSON字符串
			serializedJSON, err := xyJson.SerializeToString(value)
			require.NoError(t, err)

			// 验证往返转换的一致性
			testutil.AssertJSONEqual(t, originalJSON, serializedJSON)
		})
	}
}

// TestSerializeCustomOptions 测试自定义序列化选项
// TestSerializeCustomOptions tests custom serialization options
func TestSerializeCustomOptions(t *testing.T) {
	t.Run("sort_keys", func(t *testing.T) {
		obj := xyJson.CreateObject()
		obj.Set("zebra", xyJson.CreateString("last"))
		obj.Set("alpha", xyJson.CreateString("first"))
		obj.Set("beta", xyJson.CreateString("middle"))

		options := &xyJson.SerializeOptions{
			SortKeys: true,
			Compact:  true,
			MaxDepth: 1000,
		}
		serializer := xyJson.NewSerializerWithOptions(options)
		result, err := serializer.SerializeToString(obj)
		assert.NoError(t, err)

		// 验证键是否按字母顺序排列
		alphaIndex := strings.Index(result, "alpha")
		betaIndex := strings.Index(result, "beta")
		zebraIndex := strings.Index(result, "zebra")

		assert.True(t, alphaIndex < betaIndex)
		assert.True(t, betaIndex < zebraIndex)
	})

	t.Run("no_sort_keys", func(t *testing.T) {
		obj := xyJson.CreateObject()
		obj.Set("zebra", xyJson.CreateString("last"))
		obj.Set("alpha", xyJson.CreateString("first"))
		obj.Set("beta", xyJson.CreateString("middle"))

		options := &xyJson.SerializeOptions{
			SortKeys: false,
			Compact:  true,
			MaxDepth: 1000,
		}
		serializer := xyJson.NewSerializerWithOptions(options)
		result, err := serializer.SerializeToString(obj)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})
}

// TestSerializeErrorHandling 测试错误处理
// TestSerializeErrorHandling tests error handling
func TestSerializeErrorHandling(t *testing.T) {
	t.Run("nil_serializer_options", func(t *testing.T) {
		serializer := xyJson.NewSerializerWithOptions(nil)
		obj := xyJson.CreateObject()
		obj.Set("test", xyJson.CreateString("value"))

		result, err := serializer.SerializeToString(obj)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("get_set_options", func(t *testing.T) {
		serializer := xyJson.NewSerializer()

		// 获取默认选项
		options := serializer.GetOptions()
		assert.NotNil(t, options)

		// 设置新选项
		newOptions := &xyJson.SerializeOptions{
			Indent:   "  ",
			SortKeys: true,
		}
		serializer.SetOptions(newOptions)

		updatedOptions := serializer.GetOptions()
		assert.Equal(t, "  ", updatedOptions.Indent)
		assert.True(t, updatedOptions.SortKeys)
	})
}
