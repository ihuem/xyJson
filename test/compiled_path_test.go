package test

import (
	"testing"

	xyJson "github.com/ihuem/xyJson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCompilePath 测试路径预编译功能
func TestCompilePath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "simple path",
			path:        "$.name",
			expectError: false,
		},
		{
			name:        "nested path",
			path:        "$.user.profile.name",
			expectError: false,
		},
		{
			name:        "array index",
			path:        "$.items[0]",
			expectError: false,
		},
		{
			name:        "wildcard",
			path:        "$.items[*].name",
			expectError: false,
		},
		{
			name:        "root path",
			path:        "$",
			expectError: false,
		},
		{
			name:        "empty path",
			path:        "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled, err := xyJson.CompilePath(tt.path)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, compiled)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, compiled)
				assert.Equal(t, tt.path, compiled.Path())
			}
		})
	}
}

// TestCompiledPathQuery 测试预编译路径查询功能
func TestCompiledPathQuery(t *testing.T) {
	// 准备测试数据
	jsonData := `{
		"name": "John",
		"age": 30,
		"active": true,
		"user": {
			"profile": {
				"name": "John Doe",
				"email": "john@example.com"
			}
		},
		"items": [
			{"name": "item1", "value": 100},
			{"name": "item2", "value": 200},
			{"name": "item3", "value": 300}
		]
	}`

	root, err := xyJson.ParseString(jsonData)
	require.NoError(t, err)

	tests := []struct {
		name        string
		path        string
		expected    interface{}
		expectError bool
	}{
		{
			name:        "simple field",
			path:        "$.name",
			expected:    "John",
			expectError: false,
		},
		{
			name:        "nested field",
			path:        "$.user.profile.name",
			expected:    "John Doe",
			expectError: false,
		},
		{
			name:        "array element",
			path:        "$.items[0].name",
			expected:    "item1",
			expectError: false,
		},
		{
			name:        "non-existent path",
			path:        "$.nonexistent",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "root path",
			path:        "$",
			expected:    root,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled, err := xyJson.CompilePath(tt.path)
			require.NoError(t, err)

			result, err := compiled.Query(root)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.path == "$" {
					assert.Equal(t, tt.expected, result)
				} else {
					assert.Equal(t, tt.expected, result.AsString())
				}
			}
		})
	}
}

// TestCompiledPathQueryAll 测试预编译路径查询所有匹配值
func TestCompiledPathQueryAll(t *testing.T) {
	jsonData := `{
		"items": [
			{"name": "item1", "value": 100},
			{"name": "item2", "value": 200},
			{"name": "item3", "value": 300}
		]
	}`

	root, err := xyJson.ParseString(jsonData)
	require.NoError(t, err)

	// 测试通配符查询
	compiled, err := xyJson.CompilePath("$.items[*].name")
	require.NoError(t, err)

	results, err := compiled.QueryAll(root)
	assert.NoError(t, err)
	assert.Len(t, results, 3)
	assert.Equal(t, "item1", results[0].AsString())
	assert.Equal(t, "item2", results[1].AsString())
	assert.Equal(t, "item3", results[2].AsString())
}

// TestCompiledPathSet 测试预编译路径设置值
func TestCompiledPathSet(t *testing.T) {
	jsonData := `{"name": "John", "age": 30}`

	root, err := xyJson.ParseString(jsonData)
	require.NoError(t, err)

	// 测试设置现有字段
	compiled, err := xyJson.CompilePath("$.name")
	require.NoError(t, err)

	newValue := xyJson.CreateString("Jane")
	err = compiled.Set(root, newValue)
	assert.NoError(t, err)

	// 验证值已更新
	result, err := compiled.Query(root)
	assert.NoError(t, err)
	assert.Equal(t, "Jane", result.AsString())
}

// TestCompiledPathDelete 测试预编译路径删除值
func TestCompiledPathDelete(t *testing.T) {
	jsonData := `{"name": "John", "age": 30, "active": true}`

	root, err := xyJson.ParseString(jsonData)
	require.NoError(t, err)

	// 测试删除字段
	compiled, err := xyJson.CompilePath("$.age")
	require.NoError(t, err)

	err = compiled.Delete(root)
	assert.NoError(t, err)

	// 验证字段已删除
	assert.False(t, compiled.Exists(root))
}

// TestCompiledPathExists 测试预编译路径存在性检查
func TestCompiledPathExists(t *testing.T) {
	jsonData := `{"name": "John", "age": 30}`

	root, err := xyJson.ParseString(jsonData)
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing field", "$.name", true},
		{"non-existing field", "$.nonexistent", false},
		{"root path", "$", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled, err := xyJson.CompilePath(tt.path)
			require.NoError(t, err)

			result := compiled.Exists(root)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCompiledPathCount 测试预编译路径计数功能
func TestCompiledPathCount(t *testing.T) {
	jsonData := `{
		"items": [
			{"name": "item1"},
			{"name": "item2"},
			{"name": "item3"}
		]
	}`

	root, err := xyJson.ParseString(jsonData)
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		expected int
	}{
		{"array elements", "$.items[*]", 3},
		{"single field", "$.items", 1},
		{"non-existing", "$.nonexistent", 0},
		{"root", "$", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled, err := xyJson.CompilePath(tt.path)
			require.NoError(t, err)

			result := compiled.Count(root)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPathCache 测试路径缓存功能
func TestPathCache(t *testing.T) {
	// 清空缓存
	xyJson.ClearPathCache()

	// 检查初始状态
	size, maxSize := xyJson.GetPathCacheStats()
	assert.Equal(t, 0, size)
	assert.Equal(t, xyJson.DefaultPathCacheSize, maxSize)

	// 编译一些路径
	paths := []string{"$.name", "$.age", "$.items[0]"}
	for _, path := range paths {
		_, err := xyJson.CompilePath(path)
		assert.NoError(t, err)
	}

	// 检查缓存大小
	size, _ = xyJson.GetPathCacheStats()
	assert.Equal(t, len(paths), size)

	// 测试缓存命中
	compiled1, err := xyJson.CompilePath("$.name")
	assert.NoError(t, err)
	compiled2, err := xyJson.CompilePath("$.name")
	assert.NoError(t, err)
	// 应该是同一个实例（从缓存获取）
	assert.Equal(t, compiled1, compiled2)

	// 测试设置缓存大小
	xyJson.SetPathCacheMaxSize(2)
	size, maxSize = xyJson.GetPathCacheStats()
	assert.Equal(t, 2, maxSize)

	// 清空缓存
	xyJson.ClearPathCache()
	size, _ = xyJson.GetPathCacheStats()
	assert.Equal(t, 0, size)
}

// TestCompiledPathWithNilRoot 测试预编译路径处理nil根值
func TestCompiledPathWithNilRoot(t *testing.T) {
	compiled, err := xyJson.CompilePath("$.name")
	require.NoError(t, err)

	// 测试Query
	result, err := compiled.Query(nil)
	assert.Error(t, err)
	assert.Nil(t, result)

	// 测试QueryAll
	results, err := compiled.QueryAll(nil)
	assert.Error(t, err)
	assert.Nil(t, results)

	// 测试Set
	err = compiled.Set(nil, xyJson.CreateString("test"))
	assert.Error(t, err)

	// 测试Delete
	err = compiled.Delete(nil)
	assert.Error(t, err)

	// 测试Exists
	assert.False(t, compiled.Exists(nil))

	// 测试Count
	assert.Equal(t, 0, compiled.Count(nil))
}

// BenchmarkCompiledPathVsRegular 比较预编译路径和常规路径的性能
func BenchmarkCompiledPathVsRegular(b *testing.B) {
	jsonData := `{
		"user": {
			"profile": {
				"name": "John Doe",
				"email": "john@example.com"
			}
		},
		"items": [
			{"name": "item1", "value": 100},
			{"name": "item2", "value": 200}
		]
	}`

	root, _ := xyJson.ParseString(jsonData)
	path := "$.user.profile.name"

	b.Run("Regular Path", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = xyJson.Get(root, path)
		}
	})

	b.Run("Compiled Path", func(b *testing.B) {
		compiled, _ := xyJson.CompilePath(path)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = compiled.Query(root)
		}
	})

	b.Run("Compiled Path with Compilation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			compiled, _ := xyJson.CompilePath(path)
			_, _ = compiled.Query(root)
		}
	})
}

// BenchmarkPathCachePerformance 测试路径缓存性能
func BenchmarkPathCachePerformance(b *testing.B) {
	xyJson.ClearPathCache()
	paths := []string{
		"$.name",
		"$.user.profile.name",
		"$.items[0].name",
		"$.items[*].value",
		"$.user.profile.email",
	}

	b.Run("Cache Miss", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			xyJson.ClearPathCache()
			_, _ = xyJson.CompilePath(paths[i%len(paths)])
		}
	})

	b.Run("Cache Hit", func(b *testing.B) {
		// 预热缓存
		for _, path := range paths {
			_, _ = xyJson.CompilePath(path)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = xyJson.CompilePath(paths[i%len(paths)])
		}
	})
}
