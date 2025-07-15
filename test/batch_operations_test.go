package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	xyJson "github.com/ihuem/xyJson"
)

// TestGetBatch 测试批量获取功能
// TestGetBatch tests batch get functionality
func TestGetBatch(t *testing.T) {
	t.Run("正常情况 - Normal case", func(t *testing.T) {
		jsonStr := `{
			"user": {
				"name": "Alice",
				"age": 25,
				"active": true
			},
			"settings": {
				"theme": "dark",
				"notifications": false
			},
			"data": [1, 2, 3, 4, 5]
		}`
		root, err := xyJson.ParseString(jsonStr)
		require.NoError(t, err)

		paths := []string{
			"$.user.name",
			"$.user.age",
			"$.user.active",
			"$.settings.theme",
			"$.data[2]",
		}

		results := xyJson.GetBatch(root, paths)

		// 验证结果数量
		assert.Equal(t, 5, len(results))

		// 验证每个结果
		assert.Equal(t, "$.user.name", results[0].Path)
		assert.NoError(t, results[0].Error)
		assert.Equal(t, "Alice", results[0].Value.AsString())

		assert.Equal(t, "$.user.age", results[1].Path)
		assert.NoError(t, results[1].Error)
		assert.Equal(t, 25, results[1].Value.AsInt())

		assert.Equal(t, "$.user.active", results[2].Path)
		assert.NoError(t, results[2].Error)
		assert.Equal(t, true, results[2].Value.AsBool())

		assert.Equal(t, "$.settings.theme", results[3].Path)
		assert.NoError(t, results[3].Error)
		assert.Equal(t, "dark", results[3].Value.AsString())

		assert.Equal(t, "$.data[2]", results[4].Path)
		assert.NoError(t, results[4].Error)
		assert.Equal(t, 3, results[4].Value.AsInt())
	})

	t.Run("包含错误路径 - With invalid paths", func(t *testing.T) {
		jsonStr := `{"a": 1, "b": {"c": 2}}`
		root, err := xyJson.ParseString(jsonStr)
		require.NoError(t, err)

		paths := []string{
			"$.a",           // 有效路径
			"$.nonexistent", // 无效路径
			"$.b.c",         // 有效路径
			"$.invalid[0]",  // 无效路径
		}

		results := xyJson.GetBatch(root, paths)

		assert.Equal(t, 4, len(results))

		// 第一个路径应该成功
		assert.Equal(t, "$.a", results[0].Path)
		assert.NoError(t, results[0].Error)
		assert.Equal(t, 1, results[0].Value.AsInt())

		// 第二个路径应该失败
		assert.Equal(t, "$.nonexistent", results[1].Path)
		assert.Error(t, results[1].Error)
		assert.Nil(t, results[1].Value)

		// 第三个路径应该成功
		assert.Equal(t, "$.b.c", results[2].Path)
		assert.NoError(t, results[2].Error)
		assert.Equal(t, 2, results[2].Value.AsInt())

		// 第四个路径应该失败
		assert.Equal(t, "$.invalid[0]", results[3].Path)
		assert.Error(t, results[3].Error)
		assert.Nil(t, results[3].Value)
	})

	t.Run("空路径数组 - Empty paths array", func(t *testing.T) {
		root := xyJson.CreateObject()
		paths := []string{}

		results := xyJson.GetBatch(root, paths)

		assert.Equal(t, 0, len(results))
	})

	t.Run("nil根值 - Nil root value", func(t *testing.T) {
		paths := []string{"$.test"}

		results := xyJson.GetBatch(nil, paths)

		assert.Equal(t, 1, len(results))
		assert.Equal(t, "$.test", results[0].Path)
		assert.Error(t, results[0].Error)
		assert.Nil(t, results[0].Value)
	})
}

// TestSetBatch 测试批量设置功能
// TestSetBatch tests batch set functionality
func TestSetBatch(t *testing.T) {
	t.Run("正常情况 - Normal case", func(t *testing.T) {
		root := xyJson.CreateObject()

		operations := []xyJson.BatchSetOperation{
			{Path: "$.user.name", Value: "Bob"},
			{Path: "$.user.age", Value: 30},
			{Path: "$.user.active", Value: true},
			{Path: "$.settings.theme", Value: "light"},
			{Path: "$.data[0]", Value: 100},
		}

		results := xyJson.SetBatch(root, operations)

		// 验证结果数量
		assert.Equal(t, 5, len(results))

		// 验证所有操作都成功
		for i, result := range results {
			assert.Equal(t, operations[i].Path, result.Path)
			assert.NoError(t, result.Error, "Operation %d should succeed", i)
		}

		// 验证值是否正确设置
		name, err := xyJson.GetString(root, "$.user.name")
		assert.NoError(t, err)
		assert.Equal(t, "Bob", name)

		age, err := xyJson.GetInt(root, "$.user.age")
		assert.NoError(t, err)
		assert.Equal(t, 30, age)

		active, err := xyJson.GetBool(root, "$.user.active")
		assert.NoError(t, err)
		assert.Equal(t, true, active)

		theme, err := xyJson.GetString(root, "$.settings.theme")
		assert.NoError(t, err)
		assert.Equal(t, "light", theme)

		value, err := xyJson.GetInt(root, "$.data[0]")
		assert.NoError(t, err)
		assert.Equal(t, 100, value)
	})

	t.Run("包含无效操作 - With invalid operations", func(t *testing.T) {
		root := xyJson.CreateObject()

		operations := []xyJson.BatchSetOperation{
			{Path: "$.valid", Value: "test"},     // 有效操作
			{Path: "$.invalid[abc]", Value: 123}, // 无效路径
			{Path: "$.another", Value: true},     // 有效操作
		}

		results := xyJson.SetBatch(root, operations)

		assert.Equal(t, 3, len(results))

		// 第一个操作应该成功
		assert.Equal(t, "$.valid", results[0].Path)
		assert.NoError(t, results[0].Error)

		// 第二个操作应该失败
		assert.Equal(t, "$.invalid[abc]", results[1].Path)
		assert.Error(t, results[1].Error)

		// 第三个操作应该成功
		assert.Equal(t, "$.another", results[2].Path)
		assert.NoError(t, results[2].Error)

		// 验证成功的操作确实设置了值
		value, err := xyJson.GetString(root, "$.valid")
		assert.NoError(t, err)
		assert.Equal(t, "test", value)

		another, err := xyJson.GetBool(root, "$.another")
		assert.NoError(t, err)
		assert.Equal(t, true, another)
	})

	t.Run("空操作数组 - Empty operations array", func(t *testing.T) {
		root := xyJson.CreateObject()
		operations := []xyJson.BatchSetOperation{}

		results := xyJson.SetBatch(root, operations)

		assert.Equal(t, 0, len(results))
	})

	t.Run("nil根值 - Nil root value", func(t *testing.T) {
		operations := []xyJson.BatchSetOperation{
			{Path: "$.test", Value: "value"},
		}

		results := xyJson.SetBatch(nil, operations)

		assert.Equal(t, 1, len(results))
		assert.Equal(t, "$.test", results[0].Path)
		assert.Error(t, results[0].Error)
	})

	t.Run("覆盖现有值 - Overwrite existing values", func(t *testing.T) {
		jsonStr := `{"user": {"name": "Alice", "age": 25}}`
		root, err := xyJson.ParseString(jsonStr)
		require.NoError(t, err)

		operations := []xyJson.BatchSetOperation{
			{Path: "$.user.name", Value: "Bob"},
			{Path: "$.user.age", Value: 30},
			{Path: "$.user.email", Value: "bob@example.com"},
		}

		results := xyJson.SetBatch(root, operations)

		// 验证所有操作都成功
		for _, result := range results {
			assert.NoError(t, result.Error)
		}

		// 验证值被正确覆盖和添加
		name, err := xyJson.GetString(root, "$.user.name")
		assert.NoError(t, err)
		assert.Equal(t, "Bob", name)

		age, err := xyJson.GetInt(root, "$.user.age")
		assert.NoError(t, err)
		assert.Equal(t, 30, age)

		email, err := xyJson.GetString(root, "$.user.email")
		assert.NoError(t, err)
		assert.Equal(t, "bob@example.com", email)
	})
}

// TestBatchOperationsIntegration 测试批量操作的集成场景
// TestBatchOperationsIntegration tests integration scenarios for batch operations
func TestBatchOperationsIntegration(t *testing.T) {
	t.Run("配置文件批量读取 - Configuration batch reading", func(t *testing.T) {
		configJson := `{
			"database": {
				"host": "localhost",
				"port": 5432,
				"name": "myapp",
				"ssl": true
			},
			"redis": {
				"url": "redis://localhost:6379",
				"timeout": 30
			},
			"logging": {
				"level": "info",
				"file": "/var/log/app.log"
			}
		}`

		root, err := xyJson.ParseString(configJson)
		require.NoError(t, err)

		// 批量读取配置项
		configPaths := []string{
			"$.database.host",
			"$.database.port",
			"$.database.name",
			"$.database.ssl",
			"$.redis.url",
			"$.redis.timeout",
			"$.logging.level",
			"$.logging.file",
		}

		results := xyJson.GetBatch(root, configPaths)

		// 验证所有配置项都能正确读取
		for _, result := range results {
			assert.NoError(t, result.Error, "Failed to read config: %s", result.Path)
			assert.NotNil(t, result.Value, "Config value should not be nil: %s", result.Path)
		}

		// 验证具体的配置值
		assert.Equal(t, "localhost", results[0].Value.AsString())
		assert.Equal(t, 5432, results[1].Value.AsInt())
		assert.Equal(t, "myapp", results[2].Value.AsString())
		assert.Equal(t, true, results[3].Value.AsBool())
		assert.Equal(t, "redis://localhost:6379", results[4].Value.AsString())
		assert.Equal(t, 30, results[5].Value.AsInt())
		assert.Equal(t, "info", results[6].Value.AsString())
		assert.Equal(t, "/var/log/app.log", results[7].Value.AsString())
	})

	t.Run("用户数据批量更新 - User data batch update", func(t *testing.T) {
		userJson := `{
			"users": [
				{"id": 1, "name": "Alice", "active": true},
				{"id": 2, "name": "Bob", "active": false},
				{"id": 3, "name": "Charlie", "active": true}
			]
		}`

		root, err := xyJson.ParseString(userJson)
		require.NoError(t, err)

		// 批量更新用户状态
		updateOperations := []xyJson.BatchSetOperation{
			{Path: "$.users[0].active", Value: false},
			{Path: "$.users[1].active", Value: true},
			{Path: "$.users[2].name", Value: "Charles"},
			{Path: "$.users[0].lastLogin", Value: "2024-01-15"},
			{Path: "$.users[1].email", Value: "bob@example.com"},
		}

		results := xyJson.SetBatch(root, updateOperations)

		// 验证所有更新操作都成功
		for _, result := range results {
			assert.NoError(t, result.Error, "Failed to update: %s", result.Path)
		}

		// 验证更新后的值
		user0Active, err := xyJson.GetBool(root, "$.users[0].active")
		assert.NoError(t, err)
		assert.Equal(t, false, user0Active)

		user1Active, err := xyJson.GetBool(root, "$.users[1].active")
		assert.NoError(t, err)
		assert.Equal(t, true, user1Active)

		user2Name, err := xyJson.GetString(root, "$.users[2].name")
		assert.NoError(t, err)
		assert.Equal(t, "Charles", user2Name)

		lastLogin, err := xyJson.GetString(root, "$.users[0].lastLogin")
		assert.NoError(t, err)
		assert.Equal(t, "2024-01-15", lastLogin)

		email, err := xyJson.GetString(root, "$.users[1].email")
		assert.NoError(t, err)
		assert.Equal(t, "bob@example.com", email)
	})
}

// BenchmarkGetBatch 批量获取性能基准测试
// BenchmarkGetBatch benchmarks batch get performance
func BenchmarkGetBatch(b *testing.B) {
	// 创建测试数据
	jsonStr := `{
		"users": [
			{"id": 1, "name": "Alice", "age": 25, "active": true},
			{"id": 2, "name": "Bob", "age": 30, "active": false},
			{"id": 3, "name": "Charlie", "age": 35, "active": true},
			{"id": 4, "name": "David", "age": 28, "active": true},
			{"id": 5, "name": "Eve", "age": 32, "active": false}
		],
		"settings": {
			"theme": "dark",
			"language": "en",
			"notifications": true,
			"autoSave": false
		}
	}`

	root, err := xyJson.ParseString(jsonStr)
	if err != nil {
		b.Fatal(err)
	}

	paths := []string{
		"$.users[0].name",
		"$.users[1].age",
		"$.users[2].active",
		"$.users[3].id",
		"$.users[4].name",
		"$.settings.theme",
		"$.settings.language",
		"$.settings.notifications",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = xyJson.GetBatch(root, paths)
	}
}

// BenchmarkSetBatch 批量设置性能基准测试
// BenchmarkSetBatch benchmarks batch set performance
func BenchmarkSetBatch(b *testing.B) {
	operations := []xyJson.BatchSetOperation{
		{Path: "$.user1.name", Value: "Alice"},
		{Path: "$.user1.age", Value: 25},
		{Path: "$.user2.name", Value: "Bob"},
		{Path: "$.user2.age", Value: 30},
		{Path: "$.settings.theme", Value: "dark"},
		{Path: "$.settings.language", Value: "en"},
		{Path: "$.config.debug", Value: true},
		{Path: "$.config.timeout", Value: 5000},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root := xyJson.CreateObject()
		_ = xyJson.SetBatch(root, operations)
	}
}

// BenchmarkGetBatchVsIndividual 比较批量获取与单独获取的性能
// BenchmarkGetBatchVsIndividual compares batch get vs individual get performance
func BenchmarkGetBatchVsIndividual(b *testing.B) {
	jsonStr := `{
		"data": {
			"field1": "value1",
			"field2": 42,
			"field3": true,
			"field4": [1, 2, 3],
			"field5": {"nested": "data"}
		}
	}`

	root, err := xyJson.ParseString(jsonStr)
	if err != nil {
		b.Fatal(err)
	}

	paths := []string{
		"$.data.field1",
		"$.data.field2",
		"$.data.field3",
		"$.data.field4[0]",
		"$.data.field5.nested",
	}

	b.Run("BatchGet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xyJson.GetBatch(root, paths)
		}
	})

	b.Run("IndividualGet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, path := range paths {
				_, _ = xyJson.Get(root, path)
			}
		}
	})
}
