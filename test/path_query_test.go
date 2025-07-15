package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	xyJson "github/ihuem/xyJson"
	"github/ihuem/xyJson/test/testutil"
)

// TestJSONPathBasicQuery 测试基本JSONPath查询
// TestJSONPathBasicQuery tests basic JSONPath queries
func TestJSONPathBasicQuery(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	jsonStr := generator.GenerateJSONPathTestData()

	root, err := xyJson.ParseString(jsonStr)
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		expected interface{}
		hasError bool
	}{
		{
			name:     "root_access",
			path:     "$",
			expected: root,
			hasError: false,
		},
		{
			name:     "simple_property",
			path:     "$.store",
			expected: "object",
			hasError: false,
		},
		{
			name:     "nested_property",
			path:     "$.store.bicycle.color",
			expected: "red",
			hasError: false,
		},
		{
			name:     "array_index",
			path:     "$.store.book[0].title",
			expected: "Sayings of the Century",
			hasError: false,
		},
		{
			name:     "array_last_index",
			path:     "$.store.book[2].author",
			expected: "Herman Melville",
			hasError: false,
		},
		{
			name:     "user_name",
			path:     "$.users[0].name",
			expected: "张三",
			hasError: false,
		},
		{
			name:     "user_city",
			path:     "$.users[1].city",
			expected: "上海",
			hasError: false,
		},
		{
			name:     "nonexistent_property",
			path:     "$.nonexistent",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid_array_index",
			path:     "$.users[999]",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid_path_format",
			path:     "invalid.path",
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := xyJson.Get(root, tt.path)

			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				if tt.expected == root {
					assert.Equal(t, root, result)
				} else if tt.expected == "object" {
					assert.Equal(t, xyJson.ObjectValueType, result.Type())
				} else {
					assert.Equal(t, tt.expected, result.String())
				}
			}
		})
	}
}

// TestJSONPathWildcardQuery 测试通配符查询
// TestJSONPathWildcardQuery tests wildcard queries
func TestJSONPathWildcardQuery(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	jsonStr := generator.GenerateJSONPathTestData()

	root, err := xyJson.ParseString(jsonStr)
	require.NoError(t, err)

	t.Run("array_wildcard", func(t *testing.T) {
		results, err := xyJson.GetAll(root, "$.store.book[*].title")
		assert.NoError(t, err)
		assert.Len(t, results, 3)

		titles := make([]string, len(results))
		for i, result := range results {
			titles[i] = result.String()
		}

		assert.Contains(t, titles, "Sayings of the Century")
		assert.Contains(t, titles, "Sword of Honour")
		assert.Contains(t, titles, "Moby Dick")
	})

	t.Run("object_wildcard", func(t *testing.T) {
		results, err := xyJson.GetAll(root, "$.users[*].name")
		assert.NoError(t, err)
		assert.Len(t, results, 3)

		names := make([]string, len(results))
		for i, result := range results {
			names[i] = result.String()
		}

		assert.Contains(t, names, "张三")
		assert.Contains(t, names, "李四")
		assert.Contains(t, names, "王五")
	})
}

// TestJSONPathFilterQuery 测试过滤器查询
// TestJSONPathFilterQuery tests filter queries
func TestJSONPathFilterQuery(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	jsonStr := generator.GenerateJSONPathTestData()

	root, err := xyJson.ParseString(jsonStr)
	require.NoError(t, err)

	t.Run("price_filter", func(t *testing.T) {
		// 查找价格小于10的书籍
		results, err := xyJson.GetAll(root, "$.store.book[?(@.price < 10)]")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)

		// 验证所有结果的价格都小于10
		for _, result := range results {
			price, err := xyJson.Get(result, "$.price")
			assert.NoError(t, err)
			if scalarPrice, ok := price.(xyJson.IScalarValue); ok {
				priceVal, err := scalarPrice.Float64()
				assert.NoError(t, err)
				assert.Less(t, priceVal, 10.0)
			}
		}
	})

	t.Run("category_filter", func(t *testing.T) {
		// 查找类别为fiction的书籍
		results, err := xyJson.GetAll(root, "$.store.book[?(@.category == 'fiction')]")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)

		// 验证所有结果的类别都是fiction
		for _, result := range results {
			category, err := xyJson.Get(result, "$.category")
			assert.NoError(t, err)
			assert.Equal(t, "fiction", category.String())
		}
	})

	t.Run("age_filter", func(t *testing.T) {
		// 查找年龄大于25的用户
		results, err := xyJson.GetAll(root, "$.users[?(@.age > 25)]")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)

		// 验证所有结果的年龄都大于25
		for _, result := range results {
			age, err := xyJson.Get(result, "$.age")
			assert.NoError(t, err)
			if scalarAge, ok := age.(xyJson.IScalarValue); ok {
				ageVal, err := scalarAge.Float64()
				assert.NoError(t, err)
				assert.Greater(t, ageVal, 25.0)
			}
		}
	})
}

// TestJSONPathExists 测试路径存在性检查
// TestJSONPathExists tests path existence checking
func TestJSONPathExists(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	jsonStr := generator.GenerateJSONPathTestData()

	root, err := xyJson.ParseString(jsonStr)
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "existing_root",
			path:     "$",
			expected: true,
		},
		{
			name:     "existing_property",
			path:     "$.store",
			expected: true,
		},
		{
			name:     "existing_nested_property",
			path:     "$.store.bicycle.color",
			expected: true,
		},
		{
			name:     "existing_array_element",
			path:     "$.users[0]",
			expected: true,
		},
		{
			name:     "nonexistent_property",
			path:     "$.nonexistent",
			expected: false,
		},
		{
			name:     "nonexistent_nested_property",
			path:     "$.store.nonexistent",
			expected: false,
		},
		{
			name:     "nonexistent_array_element",
			path:     "$.users[999]",
			expected: false,
		},
		{
			name:     "invalid_path",
			path:     "invalid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := xyJson.Exists(root, tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestJSONPathCount 测试路径计数
// TestJSONPathCount tests path counting
func TestJSONPathCount(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	jsonStr := generator.GenerateJSONPathTestData()

	root, err := xyJson.ParseString(jsonStr)
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		expected int
	}{
		{
			name:     "root_count",
			path:     "$",
			expected: 1,
		},
		{
			name:     "books_count",
			path:     "$.store.book[*]",
			expected: 3,
		},
		{
			name:     "users_count",
			path:     "$.users[*]",
			expected: 3,
		},
		{
			name:     "book_titles_count",
			path:     "$.store.book[*].title",
			expected: 3,
		},
		{
			name:     "nonexistent_count",
			path:     "$.nonexistent[*]",
			expected: 0,
		},
		{
			name:     "invalid_path_count",
			path:     "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := xyJson.Count(root, tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestJSONPathSet 测试路径设置值
// TestJSONPathSet tests setting values by path
func TestJSONPathSet(t *testing.T) {
	// 创建测试对象
	obj := xyJson.CreateObject()
	obj.Set("name", "原始名称")
	obj.Set("age", 25)

	address := xyJson.CreateObject()
	address.Set("city", "北京")
	obj.Set("address", address)

	skills := xyJson.CreateArray()
	skills.Append("Go")
	skills.Append("JSON")
	obj.Set("skills", skills)

	t.Run("set_simple_property", func(t *testing.T) {
		newName := xyJson.CreateString("新名称")
		err := xyJson.Set(obj, "$.name", newName)
		assert.NoError(t, err)

		// 验证设置成功
		result, err := xyJson.Get(obj, "$.name")
		assert.NoError(t, err)
		assert.Equal(t, "新名称", result.String())
	})

	t.Run("set_nested_property", func(t *testing.T) {
		newCity := xyJson.CreateString("上海")
		err := xyJson.Set(obj, "$.address.city", newCity)
		assert.NoError(t, err)

		// 验证设置成功
		result, err := xyJson.Get(obj, "$.address.city")
		assert.NoError(t, err)
		assert.Equal(t, "上海", result.String())
	})

	t.Run("set_array_element", func(t *testing.T) {
		newSkill := xyJson.CreateString("Python")
		err := xyJson.Set(obj, "$.skills[0]", newSkill)
		assert.NoError(t, err)

		// 验证设置成功
		result, err := xyJson.Get(obj, "$.skills[0]")
		assert.NoError(t, err)
		assert.Equal(t, "Python", result.String())
	})

	t.Run("set_new_property", func(t *testing.T) {
		newValue := xyJson.CreateString("新属性值")
		err := xyJson.Set(obj, "$.newProperty", newValue)
		assert.NoError(t, err)

		// 验证设置成功
		result, err := xyJson.Get(obj, "$.newProperty")
		assert.NoError(t, err)
		assert.Equal(t, "新属性值", result.String())
	})

	t.Run("set_invalid_path", func(t *testing.T) {
		newValue := xyJson.CreateString("值")
		err := xyJson.Set(obj, "invalid.path", newValue)
		assert.Error(t, err)
	})

	t.Run("set_root_error", func(t *testing.T) {
		newValue := xyJson.CreateString("值")
		err := xyJson.Set(obj, "$", newValue)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "root")
	})
}

// TestJSONPathDelete 测试路径删除值
// TestJSONPathDelete tests deleting values by path
func TestJSONPathDelete(t *testing.T) {
	// 创建测试对象
	obj := xyJson.CreateObject()
	obj.Set("name", "测试")
	obj.Set("age", 25)
	obj.Set("temp", "临时值")

	address := xyJson.CreateObject()
	address.Set("city", "北京")
	address.Set("zipcode", "100000")
	obj.Set("address", address)

	skills := xyJson.CreateArray()
	skills.Append("Go")
	skills.Append("JSON")
	skills.Append("Test")
	obj.Set("skills", skills)

	t.Run("delete_simple_property", func(t *testing.T) {
		// 确认属性存在
		assert.True(t, xyJson.Exists(obj, "$.temp"))

		// 删除属性
		err := xyJson.Delete(obj, "$.temp")
		assert.NoError(t, err)

		// 验证删除成功
		assert.False(t, xyJson.Exists(obj, "$.temp"))
	})

	t.Run("delete_nested_property", func(t *testing.T) {
		// 确认属性存在
		assert.True(t, xyJson.Exists(obj, "$.address.zipcode"))

		// 删除属性
		err := xyJson.Delete(obj, "$.address.zipcode")
		assert.NoError(t, err)

		// 验证删除成功
		assert.False(t, xyJson.Exists(obj, "$.address.zipcode"))
		// 确认父对象仍然存在
		assert.True(t, xyJson.Exists(obj, "$.address"))
		assert.True(t, xyJson.Exists(obj, "$.address.city"))
	})

	t.Run("delete_array_element", func(t *testing.T) {
		// 确认数组元素存在
		assert.True(t, xyJson.Exists(obj, "$.skills[1]"))
		originalCount := xyJson.Count(obj, "$.skills[*]")

		// 删除数组元素
		err := xyJson.Delete(obj, "$.skills[1]")
		assert.NoError(t, err)

		// 验证删除成功
		newCount := xyJson.Count(obj, "$.skills[*]")
		assert.Equal(t, originalCount-1, newCount)
	})

	t.Run("delete_nonexistent_property", func(t *testing.T) {
		err := xyJson.Delete(obj, "$.nonexistent")
		assert.Error(t, err)
	})

	t.Run("delete_invalid_path", func(t *testing.T) {
		err := xyJson.Delete(obj, "invalid.path")
		assert.Error(t, err)
	})

	t.Run("delete_root_error", func(t *testing.T) {
		err := xyJson.Delete(obj, "$")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "root")
	})
}

// TestJSONPathRecursiveDescent 测试递归下降查询
// TestJSONPathRecursiveDescent tests recursive descent queries
func TestJSONPathRecursiveDescent(t *testing.T) {
	// 创建深度嵌套的测试数据
	obj := xyJson.CreateObject()
	obj.Set("name", "root")

	level1 := xyJson.CreateObject()
	level1.Set("name", "level1")
	obj.Set("child", level1)

	level2 := xyJson.CreateObject()
	level2.Set("name", "level2")
	level1.Set("child", level2)

	level3 := xyJson.CreateObject()
	level3.Set("name", "level3")
	level2.Set("child", level3)

	t.Run("recursive_descent_all_names", func(t *testing.T) {
		results, err := xyJson.GetAll(obj, "$..name")
		assert.NoError(t, err)
		assert.Len(t, results, 4)

		names := make([]string, len(results))
		for i, result := range results {
			names[i] = result.String()
		}

		assert.Contains(t, names, "root")
		assert.Contains(t, names, "level1")
		assert.Contains(t, names, "level2")
		assert.Contains(t, names, "level3")
	})

	t.Run("recursive_descent_specific_property", func(t *testing.T) {
		results, err := xyJson.GetAll(obj, "$..child")
		assert.NoError(t, err)
		assert.Len(t, results, 3) // level1, level2, level3
	})
}

// TestJSONPathComplexQueries 测试复杂查询
// TestJSONPathComplexQueries tests complex queries
func TestJSONPathComplexQueries(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	jsonStr := generator.GenerateJSONPathTestData()

	root, err := xyJson.ParseString(jsonStr)
	require.NoError(t, err)

	t.Run("multiple_conditions", func(t *testing.T) {
		// 查找价格大于8且小于15的书籍
		results, err := xyJson.GetAll(root, "$.store.book[?(@.price > 8 && @.price < 15)]")
		assert.NoError(t, err)

		for _, result := range results {
			price, err := xyJson.Get(result, "$.price")
			assert.NoError(t, err)
			if scalarPrice, ok := price.(xyJson.IScalarValue); ok {
				priceVal, err := scalarPrice.Float64()
				assert.NoError(t, err)
				assert.Greater(t, priceVal, 8.0)
				assert.Less(t, priceVal, 15.0)
			}
		}
	})

	t.Run("nested_array_access", func(t *testing.T) {
		// 获取所有书籍的所有属性
		results, err := xyJson.GetAll(root, "$.store.book[*].*")
		assert.NoError(t, err)
		assert.Greater(t, len(results), 0)
	})

	t.Run("combined_wildcard_and_filter", func(t *testing.T) {
		// 查找所有年龄大于等于30的用户的城市
		results, err := xyJson.GetAll(root, "$.users[?(@.age >= 30)].city")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
	})
}

// TestJSONPathErrorHandling 测试错误处理
// TestJSONPathErrorHandling tests error handling
func TestJSONPathErrorHandling(t *testing.T) {
	obj := xyJson.CreateObject()
	obj.Set("test", "value")

	t.Run("nil_root", func(t *testing.T) {
		result, err := xyJson.Get(nil, "$.test")
		assert.Error(t, err)
		assert.Nil(t, result)

		exists := xyJson.Exists(nil, "$.test")
		assert.False(t, exists)

		count := xyJson.Count(nil, "$.test")
		assert.Equal(t, 0, count)
	})

	t.Run("empty_path", func(t *testing.T) {
		result, err := xyJson.Get(obj, "")
		assert.NoError(t, err)
		assert.Equal(t, obj, result)

		exists := xyJson.Exists(obj, "")
		assert.True(t, exists)
	})

	t.Run("malformed_path", func(t *testing.T) {
		invalidPaths := []string{
			"invalid",
			"$.[",
			"$.[abc",
			"$.test[",
			"$.test[abc]",
		}

		for _, path := range invalidPaths {
			result, err := xyJson.Get(obj, path)
			assert.Error(t, err, "Path should be invalid: %s", path)
			assert.Nil(t, result)
		}
	})
}

// TestJSONPathPerformance 测试JSONPath性能
// TestJSONPathPerformance tests JSONPath performance
func TestJSONPathPerformance(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	performanceData := generator.GeneratePerformanceTestData()

	// 解析大型JSON
	root, err := xyJson.ParseString(performanceData["large"])
	require.NoError(t, err)

	// 测试简单查询性能
	testutil.MeasureTime(func() {
		for i := 0; i < 1000; i++ {
			_, err := xyJson.Get(root, "$.data[0]")
			assert.NoError(t, err)
		}
	})

	// 测试复杂查询性能
	testutil.MeasureTime(func() {
		for i := 0; i < 100; i++ {
			_, err := xyJson.GetAll(root, "$.data[*].name")
			assert.NoError(t, err)
		}
	})
}

// TestJSONPathConcurrency 测试并发访问
// TestJSONPathConcurrency tests concurrent access
func TestJSONPathConcurrency(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	jsonStr := generator.GenerateJSONPathTestData()

	root, err := xyJson.ParseString(jsonStr)
	require.NoError(t, err)

	testutil.RunConcurrently(t, 10, 100, func(goroutineID, iteration int) {
		// 并发读取
		result, err := xyJson.Get(root, "$.store.book[0].title")
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// 并发存在性检查
		exists := xyJson.Exists(root, "$.users[0].name")
		assert.True(t, exists)

		// 并发计数
		count := xyJson.Count(root, "$.store.book[*]")
		assert.Equal(t, 3, count)
	})
}

// TestJSONPathMemoryUsage 测试内存使用
// TestJSONPathMemoryUsage tests memory usage
func TestJSONPathMemoryUsage(t *testing.T) {
	generator := testutil.NewTestDataGenerator()
	jsonStr := generator.GenerateJSONPathTestData()

	root, err := xyJson.ParseString(jsonStr)
	require.NoError(t, err)

	// 测试查询操作的内存使用
	testutil.AssertNoMemoryLeak(t, func() {
		for i := 0; i < 1000; i++ {
			_, err := xyJson.Get(root, "$.store.book[0].title")
			assert.NoError(t, err)

			_, err = xyJson.GetAll(root, "$.users[*].name")
			assert.NoError(t, err)

			_ = xyJson.Exists(root, "$.store.bicycle")
			_ = xyJson.Count(root, "$.store.book[*]")
		}
	})
}

// TestJSONPathCustomFactory 测试自定义工厂
// TestJSONPathCustomFactory tests custom factory
func TestJSONPathCustomFactory(t *testing.T) {
	factory := xyJson.NewValueFactory()
	pathQuery := xyJson.NewPathQueryWithFactory(factory)

	obj := xyJson.CreateObject()
	obj.Set("test", "value")

	result, err := pathQuery.SelectOne(obj, "$.test")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "value", result.String())

	// 测试nil工厂
	pathQueryNil := xyJson.NewPathQueryWithFactory(nil)
	result, err = pathQueryNil.SelectOne(obj, "$.test")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
