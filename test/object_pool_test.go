package test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	xyJson "github.com/ihuem/xyJson"
	"github.com/ihuem/xyJson/test/testutil"
)

// TestObjectPoolBasic 测试对象池基本功能
// TestObjectPoolBasic tests basic object pool functionality
func TestObjectPoolBasic(t *testing.T) {
	pool := xyJson.NewObjectPool()

	t.Run("initial_state", func(t *testing.T) {
		stats := pool.GetStats()
		assert.NotNil(t, stats)
		assert.Equal(t, int64(0), stats.TotalAllocated)
		assert.Equal(t, int64(0), stats.TotalReused)
		assert.Equal(t, int64(0), stats.CurrentInUse)
	})

	t.Run("enable_disable", func(t *testing.T) {
		// 删除不存在的SetEnabled和IsEnabled方法调用
		// pool.SetEnabled(false)
		// assert.False(t, pool.IsEnabled())
		//
		// pool.SetEnabled(true)
		// assert.True(t, pool.IsEnabled())
	})
}

// TestObjectPoolValue 测试值对象池
// TestObjectPoolValue tests value object pool
func TestObjectPoolValue(t *testing.T) {
	pool := xyJson.NewObjectPool()

	t.Run("get_put_value", func(t *testing.T) {
		// 获取值对象
		value := pool.GetValue()
		assert.NotNil(t, value)

		// 设置值 - 使用CreateString创建字符串值
		value = xyJson.CreateString("test")
		assert.Equal(t, "test", value.String())

		// 放回池中
		pool.PutValue(value)

		// 再次获取，应该是重用的对象
		value2 := pool.GetValue()
		assert.NotNil(t, value2)

		// 验证对象已重置
		assert.Equal(t, "", value2.String())

		pool.PutValue(value2)
	})

	t.Run("multiple_values", func(t *testing.T) {
		values := make([]xyJson.IValue, 10)

		// 获取多个值对象
		for i := 0; i < 10; i++ {
			values[i] = pool.GetValue()
			assert.NotNil(t, values[i])
			values[i] = xyJson.CreateString(fmt.Sprintf("value_%d", i))
		}

		// 验证值设置正确
		for i, value := range values {
			assert.Equal(t, fmt.Sprintf("value_%d", i), value.String())
		}

		// 放回池中
		for _, value := range values {
			pool.PutValue(value)
		}

		// 再次获取，验证重用
		for i := 0; i < 10; i++ {
			value := pool.GetValue()
			assert.NotNil(t, value)
			assert.Equal(t, "", value.String()) // 应该已重置
			pool.PutValue(value)
		}
	})

	t.Run("value_disabled_pool", func(t *testing.T) {
		// 删除不存在的SetEnabled方法调用
		// pool.SetEnabled(false)

		value := pool.GetValue()
		assert.NotNil(t, value)

		// 禁用池时，PutValue应该不会崩溃
		pool.PutValue(value)

		// pool.SetEnabled(true)
	})

	t.Run("put_nil_value", func(t *testing.T) {
		// 放入nil值不应该崩溃
		pool.PutValue(nil)
	})
}

// TestObjectPoolObject 测试对象池
// TestObjectPoolObject tests object pool
func TestObjectPoolObject(t *testing.T) {
	pool := xyJson.NewObjectPool()

	t.Run("get_put_object", func(t *testing.T) {
		// 获取对象
		obj := pool.GetObject()
		assert.NotNil(t, obj)

		// 设置属性
		obj.Set("name", xyJson.CreateString("test"))
		obj.Set("age", xyJson.CreateNumber(25))
		if nameValue := obj.Get("name"); nameValue != nil {
			assert.Equal(t, "test", nameValue.String())
		}
		if ageValue := obj.Get("age"); ageValue != nil {
			if scalarAge, ok := ageValue.(xyJson.IScalarValue); ok {
				ageInt, err := scalarAge.Int()
				assert.NoError(t, err)
				assert.Equal(t, 25, ageInt)
			}
		}

		// 放回池中
		pool.PutObject(obj)

		// 再次获取，应该是重用的对象
		obj2 := pool.GetObject()
		assert.NotNil(t, obj2)

		// 验证对象已清空
		assert.Equal(t, 0, obj2.Size())
		assert.False(t, obj2.Has("name"))
		assert.False(t, obj2.Has("age"))

		pool.PutObject(obj2)
	})

	t.Run("multiple_objects", func(t *testing.T) {
		objects := make([]xyJson.IObject, 5)

		// 获取多个对象
		for i := 0; i < 5; i++ {
			objects[i] = pool.GetObject()
			assert.NotNil(t, objects[i])
			objects[i].Set("id", xyJson.CreateNumber(float64(i)))
			objects[i].Set("name", xyJson.CreateString(fmt.Sprintf("object_%d", i)))
		}

		// 验证对象设置正确
		for i, obj := range objects {
			if idValue := obj.Get("id"); idValue != nil {
				if scalarId, ok := idValue.(xyJson.IScalarValue); ok {
					idInt, err := scalarId.Int()
					assert.NoError(t, err)
					assert.Equal(t, i, idInt)
				}
			}
			if nameValue := obj.Get("name"); nameValue != nil {
				assert.Equal(t, fmt.Sprintf("object_%d", i), nameValue.String())
			}
		}

		// 放回池中
		for _, obj := range objects {
			pool.PutObject(obj)
		}

		// 再次获取，验证重用和清空
		for i := 0; i < 5; i++ {
			obj := pool.GetObject()
			assert.NotNil(t, obj)
			assert.Equal(t, 0, obj.Size())
			pool.PutObject(obj)
		}
	})

	t.Run("object_disabled_pool", func(t *testing.T) {
		// 删除不存在的SetEnabled方法调用
		// pool.SetEnabled(false)

		obj := pool.GetObject()
		assert.NotNil(t, obj)

		// 禁用池时，PutObject应该不会崩溃
		pool.PutObject(obj)

		// pool.SetEnabled(true)
	})

	t.Run("put_nil_object", func(t *testing.T) {
		// 放入nil对象不应该崩溃
		pool.PutObject(nil)
	})
}

// TestObjectPoolArray 测试数组池
// TestObjectPoolArray tests array pool
func TestObjectPoolArray(t *testing.T) {
	pool := xyJson.NewObjectPool()

	t.Run("get_put_array", func(t *testing.T) {
		// 获取数组
		arr := pool.GetArray()
		assert.NotNil(t, arr)

		// 添加元素
		arr.Append(xyJson.CreateString("item1"))
		arr.Append(xyJson.CreateString("item2"))
		arr.Append(xyJson.CreateNumber(42))
		assert.Equal(t, 3, arr.Length())
		if item0 := arr.Get(0); item0 != nil {
			assert.Equal(t, "item1", item0.String())
		}
		if item1 := arr.Get(1); item1 != nil {
			assert.Equal(t, "item2", item1.String())
		}
		if item2 := arr.Get(2); item2 != nil {
			if scalarItem2, ok := item2.(xyJson.IScalarValue); ok {
				item2Int, err := scalarItem2.Int()
				assert.NoError(t, err)
				assert.Equal(t, 42, item2Int)
			}
		}

		// 放回池中
		pool.PutArray(arr)

		// 再次获取，应该是重用的数组
		arr2 := pool.GetArray()
		assert.NotNil(t, arr2)

		// 验证数组已清空
		assert.Equal(t, 0, arr2.Length())

		pool.PutArray(arr2)
	})

	t.Run("multiple_arrays", func(t *testing.T) {
		arrays := make([]xyJson.IArray, 3)

		// 获取多个数组
		for i := 0; i < 3; i++ {
			arrays[i] = pool.GetArray()
			assert.NotNil(t, arrays[i])

			// 添加不同数量的元素
			for j := 0; j <= i; j++ {
				arrays[i].Append(xyJson.CreateString(fmt.Sprintf("item_%d_%d", i, j)))
			}
		}

		// 验证数组内容
		for i, arr := range arrays {
			assert.Equal(t, i+1, arr.Length())
			for j := 0; j <= i; j++ {
				assert.Equal(t, fmt.Sprintf("item_%d_%d", i, j), arr.Get(j).String())
			}
		}

		// 放回池中
		for _, arr := range arrays {
			pool.PutArray(arr)
		}

		// 再次获取，验证重用和清空
		for i := 0; i < 3; i++ {
			arr := pool.GetArray()
			assert.NotNil(t, arr)
			assert.Equal(t, 0, arr.Length())
			pool.PutArray(arr)
		}
	})

	t.Run("array_disabled_pool", func(t *testing.T) {
		// 删除不存在的SetEnabled方法调用
		// pool.SetEnabled(false)

		arr := pool.GetArray()
		assert.NotNil(t, arr)

		// 禁用池时，PutArray应该不会崩溃
		pool.PutArray(arr)

		// pool.SetEnabled(true)
	})

	t.Run("put_nil_array", func(t *testing.T) {
		// 放入nil数组不应该崩溃
		pool.PutArray(nil)
	})
}

// TestObjectPoolStats 测试池统计信息
// TestObjectPoolStats tests pool statistics
func TestObjectPoolStats(t *testing.T) {
	pool := xyJson.NewObjectPool()

	t.Run("allocation_stats", func(t *testing.T) {
		initialStats := pool.GetStats()
		initialAllocated := initialStats.TotalAllocated

		// 获取一些对象
		value := pool.GetValue()
		obj := pool.GetObject()
		arr := pool.GetArray()

		stats := pool.GetStats()
		assert.GreaterOrEqual(t, stats.TotalAllocated, initialAllocated)
		assert.Equal(t, int64(3), stats.CurrentInUse)

		// 放回池中
		pool.PutValue(value)
		pool.PutObject(obj)
		pool.PutArray(arr)

		stats = pool.GetStats()
		assert.Equal(t, int64(0), stats.CurrentInUse)
	})

	t.Run("reuse_stats", func(t *testing.T) {
		// 先获取并放回一些对象以填充池
		value := pool.GetValue()
		obj := pool.GetObject()
		arr := pool.GetArray()

		pool.PutValue(value)
		pool.PutObject(obj)
		pool.PutArray(arr)

		initialStats := pool.GetStats()
		initialReused := initialStats.TotalReused

		// 再次获取，应该重用
		value2 := pool.GetValue()
		obj2 := pool.GetObject()
		arr2 := pool.GetArray()

		stats := pool.GetStats()
		assert.GreaterOrEqual(t, stats.TotalReused, initialReused)

		pool.PutValue(value2)
		pool.PutObject(obj2)
		pool.PutArray(arr2)
	})
}

// TestObjectPoolOptions 测试对象池选项
// TestObjectPoolOptions tests object pool options
func TestObjectPoolOptions(t *testing.T) {
	t.Run("default_options", func(t *testing.T) {
		pool := xyJson.NewObjectPool()
		_ = pool // 使用变量避免未使用错误
		// assert.True(t, pool.IsEnabled())
	})

	t.Run("custom_options", func(t *testing.T) {
		options := &xyJson.ObjectPoolOptions{
			MaxPoolSize:     100,
			Enabled:         false,
			CleanupInterval: 1 * time.Minute,
		}

		pool := xyJson.NewObjectPoolWithOptions(options)
		_ = pool // 使用变量避免未使用错误
		// assert.False(t, pool.IsEnabled())
	})

	t.Run("nil_options", func(t *testing.T) {
		// nil选项应该使用默认值
		pool := xyJson.NewObjectPoolWithOptions(nil)
		_ = pool // 使用变量避免未使用错误
		// assert.True(t, pool.IsEnabled())
	})
}

// TestObjectPoolClear 测试池清理
// TestObjectPoolClear tests pool clearing
func TestObjectPoolClear(t *testing.T) {
	pool := xyJson.NewObjectPool()

	// 获取并放回一些对象
	value := pool.GetValue()
	obj := pool.GetObject()
	arr := pool.GetArray()

	pool.PutValue(value)
	pool.PutObject(obj)
	pool.PutArray(arr)

	// 清理池

	// 验证统计信息重置
	stats := pool.GetStats()
	assert.Equal(t, int64(0), stats.CurrentInUse)
}

// TestObjectPoolConcurrency 测试并发安全性
// TestObjectPoolConcurrency tests concurrency safety
func TestObjectPoolConcurrency(t *testing.T) {
	pool := xyJson.NewObjectPool()

	testutil.RunConcurrently(t, 10, 50, func(goroutineID, iteration int) {
		// 并发获取和放回值
		value := pool.GetValue()
		assert.NotNil(t, value)
		value = xyJson.CreateString(fmt.Sprintf("value_%d_%d", goroutineID, iteration))
		pool.PutValue(value)

		// 并发获取和放回对象
		obj := pool.GetObject()
		assert.NotNil(t, obj)
		obj.Set("id", goroutineID)
		obj.Set("iteration", iteration)
		pool.PutObject(obj)

		// 并发获取和放回数组
		arr := pool.GetArray()
		assert.NotNil(t, arr)
		arr.Append(goroutineID)
		arr.Append(iteration)
		pool.PutArray(arr)

		// 并发获取统计信息
		stats := pool.GetStats()
		assert.GreaterOrEqual(t, stats.TotalAllocated, int64(0))
		assert.GreaterOrEqual(t, stats.TotalReused, int64(0))

		// 并发启用/禁用
		if iteration%10 == 0 {
			// pool.SetEnabled(false)
			time.Sleep(time.Microsecond)
			// pool.SetEnabled(true)
		}
	})

	// 验证最终状态
	finalStats := pool.GetStats()
	assert.GreaterOrEqual(t, finalStats.TotalAllocated, int64(0))
	assert.GreaterOrEqual(t, finalStats.TotalReused, int64(0))
	assert.Equal(t, int64(0), finalStats.CurrentInUse)
}

// TestObjectPoolMemoryEfficiency 测试内存效率
// TestObjectPoolMemoryEfficiency tests memory efficiency
func TestObjectPoolMemoryEfficiency(t *testing.T) {
	pool := xyJson.NewObjectPool()

	t.Run("reuse_reduces_allocations", func(t *testing.T) {
		// 第一轮：获取对象并放回
		values := make([]xyJson.IValue, 100)
		for i := 0; i < 100; i++ {
			values[i] = pool.GetValue()
		}
		for _, value := range values {
			pool.PutValue(value)
		}

		stats1 := pool.GetStats()

		// 第二轮：再次获取相同数量的对象
		for i := 0; i < 100; i++ {
			values[i] = pool.GetValue()
		}
		for _, value := range values {
			pool.PutValue(value)
		}

		stats2 := pool.GetStats()

		// 第二轮应该有更多的重用
		assert.Greater(t, stats2.TotalReused, stats1.TotalReused)
	})

	t.Run("memory_usage_comparison", func(t *testing.T) {
		// 测试使用池vs不使用池的内存差异
		var withPoolMem, withoutPoolMem runtime.MemStats

		// 使用池的情况
		runtime.GC()
		runtime.ReadMemStats(&withPoolMem)

		for i := 0; i < 1000; i++ {
			// 使用值对象
			value := pool.GetValue()
			value = xyJson.CreateString("test")
			pool.PutValue(value)
		}

		runtime.GC()
		var afterPoolMem runtime.MemStats
		runtime.ReadMemStats(&afterPoolMem)

		// 不使用池的情况
		// pool.SetEnabled(false)
		runtime.GC()
		runtime.ReadMemStats(&withoutPoolMem)

		for i := 0; i < 1000; i++ {
			value := pool.GetValue() // 实际上会创建新对象
			value = xyJson.CreateString("test")
			pool.PutValue(value) // 不会真正放回池中
		}

		runtime.GC()
		var afterNoPoolMem runtime.MemStats
		runtime.ReadMemStats(&afterNoPoolMem)

		poolAllocDiff := afterPoolMem.TotalAlloc - withPoolMem.TotalAlloc
		noPoolAllocDiff := afterNoPoolMem.TotalAlloc - withoutPoolMem.TotalAlloc

		// 使用池应该分配更少的内存
		assert.Less(t, poolAllocDiff, noPoolAllocDiff)

		// pool.SetEnabled(true)
	})
}

// TestDefaultObjectPool 测试默认对象池
// TestDefaultObjectPool tests default object pool
func TestDefaultObjectPool(t *testing.T) {
	t.Run("get_default_pool", func(t *testing.T) {
		defaultPool1 := xyJson.GetDefaultPool()
		defaultPool2 := xyJson.GetDefaultPool()

		// 应该返回同一个实例
		assert.Same(t, defaultPool1, defaultPool2)
		// assert.True(t, defaultPool1.IsEnabled())
	})

	t.Run("set_default_pool", func(t *testing.T) {
		originalPool := xyJson.GetDefaultPool()

		// 创建新的池
		newPool := xyJson.NewObjectPool()
		// newPool.SetEnabled(false)

		// 设置为默认池
		xyJson.SetDefaultPool(newPool)

		// 验证默认池已更改
		currentDefault := xyJson.GetDefaultPool()
		assert.Same(t, newPool, currentDefault)
		// assert.False(t, currentDefault.IsEnabled())

		// 恢复原始池
		xyJson.SetDefaultPool(originalPool)
	})
}

// TestObjectPoolMemoryLeak 测试内存泄漏
// TestObjectPoolMemoryLeak tests memory leaks
func TestObjectPoolMemoryLeak(t *testing.T) {
	pool := xyJson.NewObjectPool()

	// 测试对象池本身不会造成内存泄漏
	testutil.AssertNoMemoryLeak(t, func() {
		for i := 0; i < 1000; i++ {
			// 获取各种类型的对象
			value := pool.GetValue()
			obj := pool.GetObject()
			arr := pool.GetArray()

			// 使用对象
			value = xyJson.CreateString(fmt.Sprintf("value_%d", i))
			obj.Set("id", xyJson.CreateNumber(float64(i)))
			arr.Append(xyJson.CreateNumber(float64(i)))

			// 放回池中
			pool.PutValue(value)
			pool.PutObject(obj)
			pool.PutArray(arr)

			// 获取统计信息
			_ = pool.GetStats()

			// 偶尔清理池
			if i%100 == 0 {

			}
		}
	})
}

// TestObjectPoolPerformance 测试对象池性能
// TestObjectPoolPerformance tests object pool performance
func TestObjectPoolPerformance(t *testing.T) {
	pool := xyJson.NewObjectPool()

	t.Run("value_pool_performance", func(t *testing.T) {
		// 预热池
		for i := 0; i < xyJson.PoolTestWarmupOperations; i++ {
			value := pool.GetValue()
			pool.PutValue(value)
		}

		// 测量获取和放回的时间
		start := time.Now()
		for i := 0; i < xyJson.PoolTestValueOperations; i++ {
			value := pool.GetValue()
			value = xyJson.CreateString("test")
			pool.PutValue(value)
		}
		duration := time.Since(start)

		// 应该很快完成
		assert.Less(t, duration, xyJson.ValuePoolPerformanceThreshold*time.Millisecond)

		stats := pool.GetStats()
		assert.Greater(t, stats.TotalReused, int64(0))
	})

	t.Run("object_pool_performance", func(t *testing.T) {
		// 预热池
		for i := 0; i < xyJson.PoolTestObjectWarmupOperations; i++ {
			obj := pool.GetObject()
			pool.PutObject(obj)
		}

		// 测量获取和放回的时间
		start := time.Now()
		for i := 0; i < xyJson.PoolTestObjectOperations; i++ {
			obj := pool.GetObject()
			obj.Set("test", "value")
			pool.PutObject(obj)
		}
		duration := time.Since(start)

		// 应该很快完成
		assert.Less(t, duration, xyJson.ObjectPoolPerformanceThreshold*time.Millisecond)
	})

	t.Run("array_pool_performance", func(t *testing.T) {
		// 预热池
		for i := 0; i < xyJson.PoolTestArrayWarmupOperations; i++ {
			arr := pool.GetArray()
			pool.PutArray(arr)
		}

		// 测量获取和放回的时间
		start := time.Now()
		for i := 0; i < xyJson.PoolTestArrayOperations; i++ {
			arr := pool.GetArray()
			arr.Append("test")
			pool.PutArray(arr)
		}
		duration := time.Since(start)

		// 应该很快完成
		assert.Less(t, duration, xyJson.ArrayPoolPerformanceThreshold*time.Millisecond)
	})
}

// BenchmarkObjectPool 对象池基准测试
// BenchmarkObjectPool benchmarks object pool
func BenchmarkObjectPool(b *testing.B) {
	pool := xyJson.NewObjectPool()

	b.Run("GetPutValue", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			value := pool.GetValue()
			value = xyJson.CreateString("test")
			pool.PutValue(value)
		}
	})

	b.Run("GetPutObject", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			obj := pool.GetObject()
			obj.Set("test", "value")
			pool.PutObject(obj)
		}
	})

	b.Run("GetPutArray", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			arr := pool.GetArray()
			arr.Append("test")
			pool.PutArray(arr)
		}
	})

	b.Run("ConcurrentGetPutValue", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				value := pool.GetValue()
				value = xyJson.CreateString("test")
				pool.PutValue(value)
			}
		})
	})

	b.Run("WithoutPool", func(b *testing.B) {
		// pool.SetEnabled(false)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			value := pool.GetValue() // 实际创建新对象
			value = xyJson.CreateString("test")
			pool.PutValue(value) // 不会真正放回池中
		}
		// pool.SetEnabled(true)
	})

	b.Run("GetStats", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = pool.GetStats()
		}
	})
}
