package test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	xyJson "github.com/ihuem/xyJson"
	"github.com/ihuem/xyJson/test/testutil"
)

// TestIntegrationBasicWorkflow 测试基本工作流程集成
// TestIntegrationBasicWorkflow tests basic workflow integration
func TestIntegrationBasicWorkflow(t *testing.T) {
	t.Run("parse_modify_serialize_workflow", func(t *testing.T) {
		// 1. 解析JSON
		jsonStr := `{
			"users": [
				{"id": 1, "name": "张三", "age": 25, "active": true},
				{"id": 2, "name": "李四", "age": 30, "active": false}
			],
			"metadata": {
				"total": 2,
				"timestamp": "2024-01-01T00:00:00Z"
			}
		}`

		root, err := xyJson.ParseString(jsonStr)
		require.NoError(t, err)
		require.NotNil(t, root)

		// 2. 验证解析结果
		_, ok := root.(xyJson.IObject)
		assert.True(t, ok)
		assert.True(t, xyJson.Exists(root, "$.users"))
		assert.True(t, xyJson.Exists(root, "$.metadata"))
		assert.Equal(t, 2, xyJson.Count(root, "$.users[*]"))

		// 3. 使用JSONPath查询数据
		firstUser, err := xyJson.Get(root, "$.users[0]")
		require.NoError(t, err)
		name, err := xyJson.Get(firstUser, "$.name")
		require.NoError(t, err)
		require.NotNil(t, name)
		assert.Equal(t, "张三", name.String())
		age, err := xyJson.Get(firstUser, "$.age")
		require.NoError(t, err)
		require.NotNil(t, age)
		if scalarAge, ok := age.(xyJson.IScalarValue); ok {
			ageInt, err := scalarAge.Int()
			assert.NoError(t, err)
			assert.Equal(t, 25, ageInt)
		}

		// 4. 修改数据
		err = xyJson.Set(root, "$.users[0].age", xyJson.CreateNumber(26))
		require.NoError(t, err)

		err = xyJson.Set(root, "$.metadata.updated", xyJson.CreateBool(true))
		require.NoError(t, err)

		// 4. 添加新用户
		newUser := xyJson.CreateObject()
		newUser.Set("name", xyJson.CreateString("王五"))
		newUser.Set("age", xyJson.CreateNumber(28))
		newUser.Set("active", xyJson.CreateBool(true))

		users, err := xyJson.Get(root, "$.users")
		require.NoError(t, err)
		if usersArray, ok := users.(xyJson.IArray); ok {
			usersArray.Append(newUser)
		}

		// 5. 验证修改结果
		updatedAge, err := xyJson.Get(root, "$.users[0].age")
		require.NoError(t, err)
		require.NotNil(t, updatedAge)
		if scalarAge, ok := updatedAge.(xyJson.IScalarValue); ok {
			ageInt, err := scalarAge.Int()
			assert.NoError(t, err)
			assert.Equal(t, 26, ageInt)
		}

		assert.True(t, xyJson.Exists(root, "$.metadata.updated"))
		assert.Equal(t, 3, xyJson.Count(root, "$.users[*]"))

		// 6. 序列化回JSON
		resultJSON, err := xyJson.SerializeToString(root)
		require.NoError(t, err)
		assert.NotEmpty(t, resultJSON)

		// 7. 验证序列化结果可以重新解析
		reparsed, err := xyJson.ParseString(resultJSON)
		require.NoError(t, err)

		// 验证数据一致性
		age, err = xyJson.Get(reparsed, "$.users[0].age")
		require.NoError(t, err)
		require.NotNil(t, age)
		if scalarAge, ok := age.(xyJson.IScalarValue); ok {
			ageInt, err := scalarAge.Int()
			assert.NoError(t, err)
			assert.Equal(t, 26, ageInt)
		}

		updated, err := xyJson.Get(reparsed, "$.metadata.updated")
		require.NoError(t, err)
		require.NotNil(t, updated)
		if scalarUpdated, ok := updated.(xyJson.IScalarValue); ok {
			updatedBool, err := scalarUpdated.Bool()
			assert.NoError(t, err)
			assert.True(t, updatedBool)
		}

		name, err = xyJson.Get(reparsed, "$.users[2].name")
		require.NoError(t, err)
		require.NotNil(t, name)
		assert.Equal(t, "王五", name.String())
	})
}

// TestIntegrationPerformanceMonitoring 测试性能监控集成
// TestIntegrationPerformanceMonitoring tests performance monitoring integration
func TestIntegrationPerformanceMonitoring(t *testing.T) {
	// 启用全局性能监控
	xyJson.EnablePerformanceMonitoring()
	monitor := xyJson.GetGlobalMonitor()
	monitor.Reset()

	t.Run("end_to_end_with_monitoring", func(t *testing.T) {
		generator := testutil.NewTestDataGenerator()
		testData := generator.GeneratePerformanceTestData()

		// 解析操作（自动监控）
		root, err := xyJson.ParseString(testData["medium"])
		require.NoError(t, err)

		// JSONPath查询（不直接监控，但会影响整体性能）
		results, err := xyJson.GetAll(root, "$.users[*].name")
		require.NoError(t, err)
		assert.Greater(t, len(results), 0)

		// 修改操作
		err = xyJson.Set(root, "$.processed", true)
		require.NoError(t, err)

		// 序列化操作（自动监控）
		_, err = xyJson.SerializeToString(root)
		require.NoError(t, err)

		// 验证监控数据
		stats := monitor.GetStats()
		assert.Equal(t, int64(1), stats.ParseCount)
		assert.Equal(t, int64(1), stats.SerializeCount)
		assert.Greater(t, stats.TotalParseTime, time.Duration(0))
		assert.Greater(t, stats.TotalSerializeTime, time.Duration(0))
		assert.Equal(t, int64(0), stats.ErrorCount)
	})

	t.Run("error_monitoring", func(t *testing.T) {
		monitor.Reset()

		// 解析错误监控（自动监控）
		_, err := xyJson.ParseString("invalid json")
		assert.Error(t, err)

		// 验证错误统计
		stats := monitor.GetStats()
		assert.Equal(t, int64(1), stats.ParseCount)
		assert.Equal(t, int64(1), stats.ErrorCount)
	})
}

// TestIntegrationObjectPool 测试对象池集成
// TestIntegrationObjectPool tests object pool integration
func TestIntegrationObjectPool(t *testing.T) {
	pool := xyJson.NewObjectPool()

	t.Run("pool_with_complex_operations", func(t *testing.T) {
		initialStats := pool.GetStats()

		// 使用池创建复杂对象结构
		root := pool.GetObject()
		defer pool.PutObject(root)

		// 创建用户数组
		users := pool.GetArray()
		defer pool.PutArray(users)

		for i := 0; i < 5; i++ {
			user := pool.GetObject()
			user.Set("id", xyJson.CreateNumber(float64(i+1)))
			user.Set("name", xyJson.CreateString(fmt.Sprintf("用户%d", i+1)))
			user.Set("active", xyJson.CreateBool(i%2 == 0))

			// 创建用户的标签数组
			tags := pool.GetArray()
			tags.Append(xyJson.CreateString("tag1"))
			tags.Append(xyJson.CreateString("tag2"))
			user.Set("tags", tags)

			users.Append(user)

			// 注意：这里不立即放回池中，因为对象还在使用
		}

		root.Set("users", users)
		root.Set("total", xyJson.CreateNumber(float64(users.Length())))

		// 序列化测试
		jsonStr, err := xyJson.SerializeToString(root)
		require.NoError(t, err)
		assert.Contains(t, jsonStr, "用户1")

		// 验证池统计
		stats := pool.GetStats()
		assert.Greater(t, stats.TotalAllocated, initialStats.TotalAllocated)
		assert.Greater(t, stats.CurrentInUse, int64(0))

		// 清理：将所有用户对象放回池中
		for i := 0; i < users.Length(); i++ {
			user := users.Get(i)
			if userObj, ok := user.(xyJson.IObject); ok {
				// 获取并放回标签数组
				if tags := userObj.Get("tags"); tags != nil && tags.Type() == xyJson.ArrayValueType {
					if tagsArray, ok := tags.(xyJson.IArray); ok {
						pool.PutArray(tagsArray)
					}
				}
				pool.PutObject(userObj)
			}
		}
	})

	t.Run("pool_memory_efficiency", func(t *testing.T) {
		// 测试对象池的内存效率
		var beforeMem, afterMem runtime.MemStats

		runtime.GC()
		runtime.ReadMemStats(&beforeMem)

		// 大量对象创建和回收
		for i := 0; i < 1000; i++ {
			obj := pool.GetObject()
			obj.Set("id", xyJson.CreateNumber(float64(i)))
			obj.Set("data", fmt.Sprintf("data_%d", i))

			arr := pool.GetArray()
			for j := 0; j < 10; j++ {
				arr.Append(xyJson.CreateNumber(float64(j)))
			}
			obj.Set("numbers", arr)

			pool.PutObject(obj)
			pool.PutArray(arr)
		}

		runtime.GC()
		runtime.ReadMemStats(&afterMem)

		// 验证内存使用合理
		memoryIncrease := afterMem.Alloc - beforeMem.Alloc
		assert.Less(t, memoryIncrease, uint64(1024*1024)) // 应该小于1MB

		// 验证池统计
		stats := pool.GetStats()
		assert.Greater(t, stats.TotalReused, int64(0))
		assert.Equal(t, int64(0), stats.CurrentInUse) // 所有对象都应该已归还
	})
}

// TestIntegrationLargeDataProcessing 测试大数据量处理
// TestIntegrationLargeDataProcessing tests large data processing
func TestIntegrationLargeDataProcessing(t *testing.T) {
	t.Run("large_json_processing", func(t *testing.T) {
		generator := testutil.NewTestDataGenerator()
		performanceData := generator.GeneratePerformanceTestData()

		// 解析大型JSON
		start := time.Now()
		root, err := xyJson.ParseString(performanceData["large"])
		require.NoError(t, err)
		parseTime := time.Since(start)

		// 验证解析时间合理（应该在合理范围内）
		assert.Less(t, parseTime, 5*time.Second)

		// 执行复杂查询
		start = time.Now()
		allNames, err := xyJson.GetAll(root, "$.data[*].name")
		require.NoError(t, err)
		queryTime := time.Since(start)

		assert.Greater(t, len(allNames), 0)
		assert.Less(t, queryTime, 1*time.Second)

		// 批量修改数据
		start = time.Now()
		for i := 0; i < 100 && i < len(allNames); i++ {
			path := fmt.Sprintf("$.data[%d].processed", i)
			err = xyJson.Set(root, path, xyJson.CreateBool(true))
			assert.NoError(t, err)
		}
		modifyTime := time.Since(start)

		assert.Less(t, modifyTime, 2*time.Second)

		// 序列化大型对象
		start = time.Now()
		_, err = xyJson.SerializeToString(root)
		require.NoError(t, err)
		serializeTime := time.Since(start)

		assert.Less(t, serializeTime, 5*time.Second)
	})

	t.Run("memory_usage_monitoring", func(t *testing.T) {
		var beforeMem, afterMem runtime.MemStats

		runtime.GC()
		runtime.ReadMemStats(&beforeMem)

		// 处理多个大型JSON文档
		generator := testutil.NewTestDataGenerator()
		performanceData := generator.GeneratePerformanceTestData()

		for i := 0; i < 10; i++ {
			root, err := xyJson.ParseString(performanceData["medium"])
			require.NoError(t, err)

			// 执行一些操作
			_, err = xyJson.GetAll(root, "$.data[*]")
			require.NoError(t, err)

			err = xyJson.Set(root, "$.batch", xyJson.CreateNumber(float64(i)))
			require.NoError(t, err)

			_, err = xyJson.SerializeToString(root)
			require.NoError(t, err)

			// 强制GC以释放内存
			if i%3 == 0 {
				runtime.GC()
			}
		}

		runtime.GC()
		runtime.ReadMemStats(&afterMem)

		// 验证内存增长合理
		// 使用有符号整数避免uint64下溢
		var memoryIncrease int64
		if afterMem.Alloc >= beforeMem.Alloc {
			memoryIncrease = int64(afterMem.Alloc - beforeMem.Alloc)
		} else {
			// 如果内存减少了，说明GC工作良好，设置为0
			memoryIncrease = 0
		}
		assert.Less(t, memoryIncrease, int64(50*1024*1024)) // 应该小于50MB
	})
}

// TestIntegrationConcurrentOperations 测试并发操作集成
// TestIntegrationConcurrentOperations tests concurrent operations integration
func TestIntegrationConcurrentOperations(t *testing.T) {
	t.Run("concurrent_parse_and_query", func(t *testing.T) {
		generator := testutil.NewTestDataGenerator()
		jsonStr := generator.GenerateJSONPathTestData()

		// 解析一次，多个goroutine并发查询
		root, err := xyJson.ParseString(jsonStr)
		require.NoError(t, err)

		testutil.RunConcurrently(t, 10, 50, func(goroutineID, iteration int) {
			// 并发读取操作
			result, err := xyJson.Get(root, "$.store.book[0].title")
			assert.NoError(t, err)
			assert.NotNil(t, result)

			// 并发查询操作
			results, err := xyJson.GetAll(root, "$.users[*].name")
			assert.NoError(t, err)
			assert.Greater(t, len(results), 0)

			// 并发存在性检查
			exists := xyJson.Exists(root, "$.store.bicycle")
			assert.True(t, exists)

			// 并发计数操作
			count := xyJson.Count(root, "$.store.book[*]")
			assert.Greater(t, count, 0)
		})
	})

	t.Run("concurrent_parse_different_data", func(t *testing.T) {
		generator := testutil.NewTestDataGenerator()
		testData := generator.GeneratePerformanceTestData()

		testutil.RunConcurrently(t, 5, 20, func(goroutineID, iteration int) {
			// 每个goroutine解析不同的数据
			dataKey := []string{"small", "medium"}[iteration%2]

			root, err := xyJson.ParseString(testData[dataKey])
			assert.NoError(t, err)
			assert.NotNil(t, root)

			// 执行一些操作
			err = xyJson.Set(root, "$.goroutine_id", goroutineID)
			assert.NoError(t, err)

			err = xyJson.Set(root, "$.iteration", iteration)
			assert.NoError(t, err)

			// 序列化
			_, err = xyJson.SerializeToString(root)
			assert.NoError(t, err)
		})
	})

	t.Run("concurrent_object_pool_usage", func(t *testing.T) {
		pool := xyJson.NewObjectPool()

		testutil.RunConcurrently(t, 8, 25, func(goroutineID, iteration int) {
			// 并发使用对象池
			obj := pool.GetObject()
			obj.Set("goroutine", goroutineID)
			obj.Set("iteration", iteration)
			obj.Set("timestamp", time.Now().Unix())

			arr := pool.GetArray()
			for i := 0; i < 5; i++ {
				arr.Append(i + iteration)
			}
			obj.Set("data", arr)

			// 序列化
			_, err := xyJson.SerializeToString(obj)
			assert.NoError(t, err)

			// 归还到池中
			pool.PutObject(obj)
			pool.PutArray(arr)
		})

		// 验证池状态
		stats := pool.GetStats()
		assert.Equal(t, int64(0), stats.CurrentInUse)
		assert.Greater(t, stats.TotalReused, int64(0))
	})
}

// TestIntegrationLongRunningStability 测试长时间运行稳定性
// TestIntegrationLongRunningStability tests long-running stability
func TestIntegrationLongRunningStability(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过长时间运行测试")
	}

	t.Run("continuous_operations", func(t *testing.T) {
		monitor := xyJson.NewPerformanceMonitor()
		pool := xyJson.NewObjectPool()
		generator := testutil.NewTestDataGenerator()
		testData := generator.GeneratePerformanceTestData()

		// 运行5分钟的连续操作
		startTime := time.Now()
		duration := 30 * time.Second // 缩短测试时间
		operationCount := 0
		errorCount := 0

		for time.Since(startTime) < duration {
			// 解析操作
			timer := monitor.StartParseTimer()
			root, err := xyJson.ParseString(testData["medium"])
			if err != nil {
				errorCount++
				timer.EndWithError()
				continue
			}
			timer.End()

			// JSONPath查询
			_, err = xyJson.GetAll(root, "$.data[*].name")
			if err != nil {
				errorCount++
				continue
			}

			// 修改操作
			err = xyJson.Set(root, "$.processed_at", time.Now().Unix())
			if err != nil {
				errorCount++
				continue
			}

			// 序列化操作
			serializeTimer := monitor.StartSerializeTimer()
			_, err = xyJson.SerializeToString(root)
			if err != nil {
				errorCount++
				serializeTimer.EndWithError()
				continue
			}
			serializeTimer.End()

			// 使用对象池创建新对象
			obj := pool.GetObject()
			obj.Set("operation", operationCount)
			obj.Set("timestamp", time.Now().Unix())
			pool.PutObject(obj)

			operationCount++

			// 定期强制GC
			if operationCount%100 == 0 {
				runtime.GC()
			}
		}

		// 验证稳定性
		assert.Greater(t, operationCount, 100)        // 至少完成100次操作
		assert.Less(t, errorCount, operationCount/10) // 错误率应该小于10%

		// 验证监控数据
		stats := monitor.GetStats()
		assert.Greater(t, stats.ParseCount, int64(0))
		assert.Greater(t, stats.SerializeCount, int64(0))
		assert.Less(t, stats.ErrorCount, int64(operationCount/10))

		// 验证池状态
		poolStats := pool.GetStats()
		assert.Equal(t, int64(0), poolStats.CurrentInUse)
		assert.Greater(t, poolStats.TotalReused, int64(0))
	})
}

// TestIntegrationMemoryLeakDetection 测试内存泄漏检测
// TestIntegrationMemoryLeakDetection tests memory leak detection
func TestIntegrationMemoryLeakDetection(t *testing.T) {
	t.Run("no_memory_leak_in_operations", func(t *testing.T) {
		testutil.AssertNoMemoryLeak(t, func() {
			generator := testutil.NewTestDataGenerator()
			testData := generator.GeneratePerformanceTestData()
			pool := xyJson.NewObjectPool()

			for i := 0; i < 1000; i++ {
				// 解析
				root, err := xyJson.ParseString(testData["small"])
				assert.NoError(t, err)

				// 查询
				_, err = xyJson.Get(root, "$.name")
				assert.NoError(t, err)

				// 修改
				err = xyJson.Set(root, "$.iteration", xyJson.CreateNumber(float64(i)))
				assert.NoError(t, err)

				// 序列化
				_, err = xyJson.SerializeToString(root)
				assert.NoError(t, err)

				// 使用对象池
				obj := pool.GetObject()
				obj.Set("test", xyJson.CreateNumber(float64(i)))
				pool.PutObject(obj)

				// 定期清理
				if i%100 == 0 {

					runtime.GC()
				}
			}
		})
	})
}

// TestIntegrationErrorRecovery 测试错误恢复
// TestIntegrationErrorRecovery tests error recovery
func TestIntegrationErrorRecovery(t *testing.T) {
	t.Run("graceful_error_handling", func(t *testing.T) {
		monitor := xyJson.NewPerformanceMonitor()
		monitor.Reset()

		// 测试各种错误情况下的恢复能力
		errorCases := []struct {
			name string
			data string
			path string
		}{
			{"invalid_json", "invalid json", "$.test"},
			{"empty_json", "", "$.test"},
			{"malformed_json", `{"key": }`, "$.key"},
			{"valid_json_invalid_path", `{"key": "value"}`, "invalid.path"},
			{"valid_json_nonexistent_path", `{"key": "value"}`, "$.nonexistent"},
		}

		successCount := 0

		for _, tc := range errorCases {
			t.Run(tc.name, func(t *testing.T) {
				// 尝试解析
				timer := monitor.StartParseTimer()
				root, err := xyJson.ParseString(tc.data)
				if err != nil {
					timer.EndWithError()
					// 解析失败是预期的，继续测试其他功能
					return
				}
				timer.End()

				// 尝试查询
				_, err = xyJson.Get(root, tc.path)
				if err != nil {
					// 查询失败也是预期的
					return
				}

				successCount++
			})
		}

		// 验证错误处理后系统仍然正常工作
		validJSON := `{"test": "value", "number": 42}`
		root, err := xyJson.ParseString(validJSON)
		require.NoError(t, err)

		result, err := xyJson.Get(root, "$.test")
		require.NoError(t, err)
		assert.Equal(t, "value", result.String())

		// 验证监控数据记录了错误
		stats := monitor.GetStats()
		assert.Greater(t, stats.ErrorCount, int64(0))
	})
}

// TestIntegrationDataConsistency 测试数据一致性
// TestIntegrationDataConsistency tests data consistency
func TestIntegrationDataConsistency(t *testing.T) {
	t.Run("round_trip_consistency", func(t *testing.T) {
		// 测试数据在解析-修改-序列化-重新解析过程中的一致性
		originalJSON := `{
			"string": "测试字符串",
			"number": 42.5,
			"boolean": true,
			"null": null,
			"array": [1, 2, 3, "four", true],
			"object": {
				"nested": "value",
				"count": 10
			}
		}`

		// 第一次解析
		root1, err := xyJson.ParseString(originalJSON)
		require.NoError(t, err)

		// 修改数据
		err = xyJson.Set(root1, "$.modified", true)
		require.NoError(t, err)

		err = xyJson.Set(root1, "$.object.new_field", "new_value")
		require.NoError(t, err)

		// 序列化
		modifiedJSON, err := xyJson.SerializeToString(root1)
		require.NoError(t, err)

		// 重新解析
		root2, err := xyJson.ParseString(modifiedJSON)
		require.NoError(t, err)

		// 验证原始数据保持一致
		strVal, _ := xyJson.Get(root2, "$.string")
		assert.Equal(t, "测试字符串", strVal.String())

		numVal, _ := xyJson.Get(root2, "$.number")
		if scalarNum, ok := numVal.(xyJson.IScalarValue); ok {
			numFloat, err := scalarNum.Float64()
			assert.NoError(t, err)
			assert.Equal(t, 42.5, numFloat)
		}

		boolVal, _ := xyJson.Get(root2, "$.boolean")
		if scalarBool, ok := boolVal.(xyJson.IScalarValue); ok {
			boolResult, err := scalarBool.Bool()
			assert.NoError(t, err)
			assert.True(t, boolResult)
		}

		nullVal, _ := xyJson.Get(root2, "$.null")
		assert.True(t, nullVal.IsNull())

		arrVal, _ := xyJson.Get(root2, "$.array")
		if arrInterface, ok := arrVal.(xyJson.IArray); ok {
			assert.Equal(t, 5, arrInterface.Length())
		}

		nestedVal, _ := xyJson.Get(root2, "$.object.nested")
		assert.Equal(t, "value", nestedVal.String())
		countVal, _ := xyJson.Get(root2, "$.object.count")
		if scalarCount, ok := countVal.(xyJson.IScalarValue); ok {
			countInt, err := scalarCount.Int()
			assert.NoError(t, err)
			assert.Equal(t, 10, countInt)
		}

		// 验证修改的数据
		modifiedVal, _ := xyJson.Get(root2, "$.modified")
		if scalarModified, ok := modifiedVal.(xyJson.IScalarValue); ok {
			modifiedBool, err := scalarModified.Bool()
			assert.NoError(t, err)
			assert.True(t, modifiedBool)
		}

		newFieldVal, _ := xyJson.Get(root2, "$.object.new_field")
		assert.Equal(t, "new_value", newFieldVal.String())
	})

	t.Run("concurrent_modification_safety", func(t *testing.T) {
		// 测试并发修改时的数据安全性
		originalJSON := `{"counter": 0, "data": []}`
		root, err := xyJson.ParseString(originalJSON)
		require.NoError(t, err)

		var wg sync.WaitGroup
		var mu sync.Mutex
		errorCount := 0

		// 启动多个goroutine并发修改
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				for j := 0; j < 10; j++ {
					// 使用互斥锁保护并发修改
					mu.Lock()

					// 读取当前计数器
					counter, err := xyJson.Get(root, "$.counter")
					if err != nil {
						errorCount++
						mu.Unlock()
						continue
					}

					// 增加计数器
					var newValue int
					if scalarCounter, ok := counter.(xyJson.IScalarValue); ok {
						currentInt, intErr := scalarCounter.Int()
						if intErr != nil {
							errorCount++
							mu.Unlock()
							continue
						}
						newValue = currentInt + 1
					}
					err = xyJson.Set(root, "$.counter", xyJson.CreateNumber(float64(newValue)))
					if err != nil {
						errorCount++
					}

					// 添加数据到数组
					data, err := xyJson.Get(root, "$.data")
					if err == nil && data.Type() == xyJson.ArrayValueType {
						if dataArray, ok := data.(xyJson.IArray); ok {
							dataArray.Append(xyJson.CreateString(fmt.Sprintf("item_%d_%d", id, j)))
						}
					}

					mu.Unlock()
					time.Sleep(time.Microsecond) // 短暂休眠
				}
			}(i)
		}

		wg.Wait()

		// 验证最终状态
		finalCounter, err := xyJson.Get(root, "$.counter")
		require.NoError(t, err)
		if scalarFinalCounter, ok := finalCounter.(xyJson.IScalarValue); ok {
			finalCountInt, err := scalarFinalCounter.Int()
			assert.NoError(t, err)
			assert.Equal(t, 100, finalCountInt) // 10个goroutine * 10次操作
		}

		finalData, err := xyJson.Get(root, "$.data")
		require.NoError(t, err)
		if finalDataArray, ok := finalData.(xyJson.IArray); ok {
			assert.Equal(t, 100, finalDataArray.Length())
		}

		assert.Equal(t, 0, errorCount) // 不应该有错误
	})
}
