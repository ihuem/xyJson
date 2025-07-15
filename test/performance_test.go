package test

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	xyJson "github/ihuem/xyJson"
	"github/ihuem/xyJson/test/testutil"
)

// TestPerformanceMonitorBasic 测试性能监控器基本功能
// TestPerformanceMonitorBasic tests basic performance monitor functionality
func TestPerformanceMonitorBasic(t *testing.T) {
	monitor := xyJson.NewPerformanceMonitor()

	t.Run("initial_state", func(t *testing.T) {
		assert.True(t, monitor.IsEnabled())

		// 等待一小段时间确保Uptime > 0
		time.Sleep(1 * time.Millisecond)

		stats := monitor.GetStats()
		assert.Equal(t, int64(0), stats.ParseCount)
		assert.Equal(t, int64(0), stats.SerializeCount)
		assert.Equal(t, time.Duration(0), stats.AvgParseTime)
		assert.Equal(t, time.Duration(0), stats.AvgSerializeTime)
		assert.Equal(t, int64(0), stats.ErrorCount)
		assert.True(t, stats.Enabled)
		assert.Greater(t, stats.Uptime, time.Duration(0))
	})

	t.Run("enable_disable", func(t *testing.T) {
		monitor.Disable()
		assert.False(t, monitor.IsEnabled())

		stats := monitor.GetStats()
		assert.False(t, stats.Enabled)

		monitor.Enable()
		assert.True(t, monitor.IsEnabled())

		stats = monitor.GetStats()
		assert.True(t, stats.Enabled)
	})
}

// TestPerformanceMonitorRecordParse 测试解析性能记录
// TestPerformanceMonitorRecordParse tests parse performance recording
func TestPerformanceMonitorRecordParse(t *testing.T) {
	monitor := xyJson.NewPerformanceMonitor()

	t.Run("record_single_parse", func(t *testing.T) {
		duration := 100 * time.Microsecond
		allocBytes := int64(1024)

		monitor.RecordParse(duration, allocBytes)

		stats := monitor.GetStats()
		assert.Equal(t, int64(1), stats.ParseCount)
		assert.Equal(t, duration, stats.TotalParseTime)
		assert.Equal(t, duration, stats.AvgParseTime)
		assert.Equal(t, int64(1), stats.AllocCount)
		assert.Equal(t, allocBytes, stats.AllocBytes)
	})

	t.Run("record_multiple_parses", func(t *testing.T) {
		monitor.Reset()

		durations := []time.Duration{
			50 * time.Microsecond,
			100 * time.Microsecond,
			150 * time.Microsecond,
		}
		allocSizes := []int64{512, 1024, 2048}

		for i, duration := range durations {
			monitor.RecordParse(duration, allocSizes[i])
		}

		stats := monitor.GetStats()
		assert.Equal(t, int64(3), stats.ParseCount)
		assert.Equal(t, int64(3), stats.AllocCount)
		assert.Equal(t, int64(512+1024+2048), stats.AllocBytes)

		expectedTotal := 50*time.Microsecond + 100*time.Microsecond + 150*time.Microsecond
		assert.Equal(t, expectedTotal, stats.TotalParseTime)
		assert.Equal(t, expectedTotal/3, stats.AvgParseTime)
	})

	t.Run("record_parse_disabled", func(t *testing.T) {
		monitor.Reset()
		monitor.Disable()

		monitor.RecordParse(100*time.Microsecond, 1024)

		stats := monitor.GetStats()
		assert.Equal(t, int64(0), stats.ParseCount)
		assert.Equal(t, time.Duration(0), stats.TotalParseTime)

		monitor.Enable()
	})
}

// TestPerformanceMonitorRecordSerialize 测试序列化性能记录
// TestPerformanceMonitorRecordSerialize tests serialize performance recording
func TestPerformanceMonitorRecordSerialize(t *testing.T) {
	monitor := xyJson.NewPerformanceMonitor()

	t.Run("record_single_serialize", func(t *testing.T) {
		monitor.Reset()

		duration := 80 * time.Microsecond
		allocBytes := int64(512)

		monitor.RecordSerialize(duration, allocBytes)

		stats := monitor.GetStats()
		assert.Equal(t, int64(1), stats.SerializeCount)
		assert.Equal(t, duration, stats.TotalSerializeTime)
		assert.Equal(t, duration, stats.AvgSerializeTime)
		assert.Equal(t, int64(1), stats.AllocCount)
		assert.Equal(t, allocBytes, stats.AllocBytes)
	})

	t.Run("record_multiple_serializes", func(t *testing.T) {
		monitor.Reset()

		durations := []time.Duration{
			30 * time.Microsecond,
			60 * time.Microsecond,
			90 * time.Microsecond,
		}
		allocSizes := []int64{256, 512, 1024}

		for i, duration := range durations {
			monitor.RecordSerialize(duration, allocSizes[i])
		}

		stats := monitor.GetStats()
		assert.Equal(t, int64(3), stats.SerializeCount)
		assert.Equal(t, int64(3), stats.AllocCount)
		assert.Equal(t, int64(256+512+1024), stats.AllocBytes)

		expectedTotal := 30*time.Microsecond + 60*time.Microsecond + 90*time.Microsecond
		assert.Equal(t, expectedTotal, stats.TotalSerializeTime)
		assert.Equal(t, expectedTotal/3, stats.AvgSerializeTime)
	})

	t.Run("record_serialize_disabled", func(t *testing.T) {
		monitor.Reset()
		monitor.Disable()

		monitor.RecordSerialize(80*time.Microsecond, 512)

		stats := monitor.GetStats()
		assert.Equal(t, int64(0), stats.SerializeCount)
		assert.Equal(t, time.Duration(0), stats.TotalSerializeTime)

		monitor.Enable()
	})
}

// TestPerformanceMonitorRecordError 测试错误记录
// TestPerformanceMonitorRecordError tests error recording
func TestPerformanceMonitorRecordError(t *testing.T) {
	monitor := xyJson.NewPerformanceMonitor()

	t.Run("record_single_error", func(t *testing.T) {
		monitor.Reset()

		monitor.RecordError()

		stats := monitor.GetStats()
		assert.Equal(t, int64(1), stats.ErrorCount)
	})

	t.Run("record_multiple_errors", func(t *testing.T) {
		monitor.Reset()

		for i := 0; i < 5; i++ {
			monitor.RecordError()
		}

		stats := monitor.GetStats()
		assert.Equal(t, int64(5), stats.ErrorCount)
	})

	t.Run("record_error_disabled", func(t *testing.T) {
		monitor.Reset()
		monitor.Disable()

		monitor.RecordError()

		stats := monitor.GetStats()
		assert.Equal(t, int64(0), stats.ErrorCount)

		monitor.Enable()
	})
}

// TestPerformanceMonitorMemoryTracking 测试内存跟踪
// TestPerformanceMonitorMemoryTracking tests memory tracking
func TestPerformanceMonitorMemoryTracking(t *testing.T) {
	monitor := xyJson.NewPerformanceMonitor()
	monitor.Reset()

	t.Run("memory_usage_tracking", func(t *testing.T) {
		// 记录一些操作以触发内存更新
		monitor.RecordParse(100*time.Microsecond, 1024)
		monitor.RecordSerialize(80*time.Microsecond, 512)

		stats := monitor.GetStats()
		assert.GreaterOrEqual(t, stats.CurrentMemoryUsage, int64(0))
		assert.GreaterOrEqual(t, stats.MaxMemoryUsage, stats.CurrentMemoryUsage)
		assert.GreaterOrEqual(t, stats.GCCount, uint32(0))
	})

	t.Run("memory_growth_tracking", func(t *testing.T) {
		initialStats := monitor.GetStats()
		initialMax := initialStats.MaxMemoryUsage

		// 分配一些内存
		data := make([]byte, 1024*1024) // 1MB
		_ = data

		// 记录操作以触发内存检查
		monitor.RecordParse(100*time.Microsecond, int64(len(data)))

		newStats := monitor.GetStats()
		assert.GreaterOrEqual(t, newStats.MaxMemoryUsage, initialMax)
	})
}

// TestPerformanceMonitorReset 测试重置功能
// TestPerformanceMonitorReset tests reset functionality
func TestPerformanceMonitorReset(t *testing.T) {
	monitor := xyJson.NewPerformanceMonitor()

	// 记录一些数据
	monitor.RecordParse(100*time.Microsecond, 1024)
	monitor.RecordSerialize(80*time.Microsecond, 512)
	monitor.RecordError()

	// 验证数据已记录
	stats := monitor.GetStats()
	assert.Greater(t, stats.ParseCount, int64(0))
	assert.Greater(t, stats.SerializeCount, int64(0))
	assert.Greater(t, stats.ErrorCount, int64(0))

	// 重置
	monitor.Reset()

	// 验证数据已清零
	stats = monitor.GetStats()
	assert.Equal(t, int64(0), stats.ParseCount)
	assert.Equal(t, int64(0), stats.SerializeCount)
	assert.Equal(t, int64(0), stats.ErrorCount)
	assert.Equal(t, int64(0), stats.AllocCount)
	assert.Equal(t, int64(0), stats.AllocBytes)
	assert.Equal(t, time.Duration(0), stats.TotalParseTime)
	assert.Equal(t, time.Duration(0), stats.TotalSerializeTime)
	assert.Equal(t, time.Duration(0), stats.AvgParseTime)
	assert.Equal(t, time.Duration(0), stats.AvgSerializeTime)

	// 验证uptime已重置
	assert.Less(t, stats.Uptime, 100*time.Millisecond)
}

// TestPerformanceMonitorTimedOperation 测试计时操作
// TestPerformanceMonitorTimedOperation tests timed operations
func TestPerformanceMonitorTimedOperation(t *testing.T) {
	monitor := xyJson.NewPerformanceMonitor()
	monitor.Reset()

	t.Run("parse_timer", func(t *testing.T) {
		timer := monitor.StartParseTimer()

		// 模拟一些工作
		time.Sleep(10 * time.Millisecond)

		timer.End()

		stats := monitor.GetStats()
		assert.Equal(t, int64(1), stats.ParseCount)
		assert.Greater(t, stats.TotalParseTime, 5*time.Millisecond)
		assert.Less(t, stats.TotalParseTime, 50*time.Millisecond)
	})

	t.Run("serialize_timer", func(t *testing.T) {
		monitor.Reset()

		timer := monitor.StartSerializeTimer()

		// 模拟一些工作
		time.Sleep(5 * time.Millisecond)

		timer.End()

		stats := monitor.GetStats()
		assert.Equal(t, int64(1), stats.SerializeCount)
		assert.Greater(t, stats.TotalSerializeTime, 2*time.Millisecond)
		assert.Less(t, stats.TotalSerializeTime, 25*time.Millisecond)
	})

	t.Run("timer_with_error", func(t *testing.T) {
		monitor.Reset()

		timer := monitor.StartParseTimer()
		time.Sleep(5 * time.Millisecond)
		timer.EndWithError()

		stats := monitor.GetStats()
		assert.Equal(t, int64(1), stats.ParseCount)
		assert.Equal(t, int64(1), stats.ErrorCount)
		assert.Greater(t, stats.TotalParseTime, 2*time.Millisecond)
	})

	t.Run("timer_disabled", func(t *testing.T) {
		monitor.Reset()
		monitor.Disable()

		timer := monitor.StartParseTimer()
		time.Sleep(5 * time.Millisecond)
		timer.End()

		stats := monitor.GetStats()
		assert.Equal(t, int64(0), stats.ParseCount)

		monitor.Enable()
	})
}

// TestPerformanceMonitorConcurrency 测试并发安全性
// TestPerformanceMonitorConcurrency tests concurrency safety
func TestPerformanceMonitorConcurrency(t *testing.T) {
	monitor := xyJson.NewPerformanceMonitor()
	monitor.Reset()

	testutil.RunConcurrently(t, 10, 100, func(goroutineID, iteration int) {
		// 并发记录解析操作
		monitor.RecordParse(time.Duration(iteration)*time.Microsecond, int64(iteration*10))

		// 并发记录序列化操作
		monitor.RecordSerialize(time.Duration(iteration)*time.Microsecond, int64(iteration*5))

		// 并发记录错误
		if iteration%10 == 0 {
			monitor.RecordError()
		}

		// 并发获取统计信息
		stats := monitor.GetStats()
		assert.GreaterOrEqual(t, stats.ParseCount, int64(0))
		assert.GreaterOrEqual(t, stats.SerializeCount, int64(0))
		assert.GreaterOrEqual(t, stats.ErrorCount, int64(0))

		// 并发启用/禁用
		if iteration%20 == 0 {
			monitor.Disable()
			time.Sleep(time.Microsecond)
			monitor.Enable()
		}
	})

	// 验证最终状态
	finalStats := monitor.GetStats()
	assert.GreaterOrEqual(t, finalStats.ParseCount, int64(0))
	assert.GreaterOrEqual(t, finalStats.SerializeCount, int64(0))
	assert.GreaterOrEqual(t, finalStats.ErrorCount, int64(0))
}

// TestGlobalPerformanceMonitor 测试全局性能监控器
// TestGlobalPerformanceMonitor tests global performance monitor
func TestGlobalPerformanceMonitor(t *testing.T) {
	t.Run("singleton_behavior", func(t *testing.T) {
		monitor1 := xyJson.GetGlobalMonitor()
		monitor2 := xyJson.GetGlobalMonitor()

		assert.Same(t, monitor1, monitor2)
		assert.True(t, monitor1.IsEnabled())
	})

	t.Run("global_monitor_functionality", func(t *testing.T) {
		globalMonitor := xyJson.GetGlobalMonitor()
		globalMonitor.Reset()

		globalMonitor.RecordParse(100*time.Microsecond, 1024)

		stats := globalMonitor.GetStats()
		assert.Equal(t, int64(1), stats.ParseCount)
		assert.Equal(t, 100*time.Microsecond, stats.TotalParseTime)
	})
}

// TestMemoryProfiler 测试内存分析器
// TestMemoryProfiler tests memory profiler
func TestMemoryProfiler(t *testing.T) {
	t.Run("basic_functionality", func(t *testing.T) {
		profiler := xyJson.NewMemoryProfiler(10, 50*time.Millisecond)

		assert.False(t, profiler.IsRunning())
		assert.Len(t, profiler.GetSnapshots(), 0)

		// 启动分析器
		profiler.Start()
		assert.True(t, profiler.IsRunning())

		// 等待一些快照
		time.Sleep(150 * time.Millisecond)

		// 停止分析器
		profiler.Stop()
		assert.False(t, profiler.IsRunning())

		// 检查快照
		snapshots := profiler.GetSnapshots()
		assert.Greater(t, len(snapshots), 0)

		// 检查最新快照
		latest := profiler.GetLatestSnapshot()
		assert.NotNil(t, latest)
		assert.Greater(t, latest.Alloc, uint64(0))
	})

	t.Run("memory_trend_analysis", func(t *testing.T) {
		profiler := xyJson.NewMemoryProfiler(5, 10*time.Millisecond)

		profiler.Start()
		time.Sleep(100 * time.Millisecond)
		profiler.Stop()

		trend, growth := profiler.GetMemoryTrend()
		assert.NotEmpty(t, trend)
		assert.GreaterOrEqual(t, growth, 0.0)
	})

	t.Run("clear_snapshots", func(t *testing.T) {
		profiler := xyJson.NewMemoryProfiler(5, 10*time.Millisecond)

		profiler.Start()
		time.Sleep(50 * time.Millisecond)
		profiler.Stop()

		assert.Greater(t, len(profiler.GetSnapshots()), 0)

		profiler.ClearSnapshots()
		assert.Len(t, profiler.GetSnapshots(), 0)
		assert.Nil(t, profiler.GetLatestSnapshot())
	})
}

// TestGlobalMemoryProfiler 测试全局内存分析器
// TestGlobalMemoryProfiler tests global memory profiler
func TestGlobalMemoryProfiler(t *testing.T) {
	t.Run("singleton_behavior", func(t *testing.T) {
		profiler1 := xyJson.GetGlobalProfiler()
		profiler2 := xyJson.GetGlobalProfiler()

		assert.Same(t, profiler1, profiler2)
	})

	t.Run("global_profiler_functionality", func(t *testing.T) {
		globalProfiler := xyJson.GetGlobalProfiler()
		globalProfiler.ClearSnapshots()

		if globalProfiler.IsRunning() {
			globalProfiler.Stop()
		}

		globalProfiler.Start()
		time.Sleep(50 * time.Millisecond)
		globalProfiler.Stop()

		assert.Greater(t, len(globalProfiler.GetSnapshots()), 0)
	})
}

// TestMemoryUtilities 测试内存工具函数
// TestMemoryUtilities tests memory utility functions
func TestMemoryUtilities(t *testing.T) {
	t.Run("force_gc", func(t *testing.T) {
		initialStats := xyJson.GetMemoryStats()

		// 分配一些内存
		data := make([][]byte, 1000)
		for i := range data {
			data[i] = make([]byte, 1024)
		}

		// 强制GC
		xyJson.ForceGC()

		newStats := xyJson.GetMemoryStats()
		assert.GreaterOrEqual(t, newStats.NumGC, initialStats.NumGC)

		// 清理引用以允许GC
		data = nil
		runtime.GC()
	})

	t.Run("get_memory_stats", func(t *testing.T) {
		stats := xyJson.GetMemoryStats()

		assert.Greater(t, stats.Alloc, uint64(0))
		assert.Greater(t, stats.TotalAlloc, uint64(0))
		assert.Greater(t, stats.Sys, uint64(0))
		assert.GreaterOrEqual(t, stats.NumGC, uint32(0))
		assert.GreaterOrEqual(t, stats.GCCPUFraction, 0.0)
		assert.Greater(t, stats.HeapAlloc, uint64(0))
		assert.Greater(t, stats.HeapSys, uint64(0))
		assert.GreaterOrEqual(t, stats.HeapInuse, uint64(0))
		assert.GreaterOrEqual(t, stats.StackInuse, uint64(0))
		assert.False(t, stats.Timestamp.IsZero())
	})
}

// TestPerformanceIntegration 测试性能监控集成
// TestPerformanceIntegration tests performance monitoring integration
func TestPerformanceIntegration(t *testing.T) {
	monitor := xyJson.NewPerformanceMonitor()
	monitor.Reset()

	generator := testutil.NewTestDataGenerator()
	testData := generator.GeneratePerformanceTestData()

	t.Run("parse_with_monitoring", func(t *testing.T) {
		timer := monitor.StartParseTimer()

		_, err := xyJson.ParseString(testData["medium"])
		assert.NoError(t, err)

		timer.End()

		stats := monitor.GetStats()
		assert.Equal(t, int64(1), stats.ParseCount)
		assert.Greater(t, stats.TotalParseTime, time.Duration(0))
	})

	t.Run("serialize_with_monitoring", func(t *testing.T) {
		monitor.Reset()

		obj := xyJson.CreateObject()
		obj.Set("test", "value")

		timer := monitor.StartSerializeTimer()

		_, err := xyJson.SerializeToString(obj)
		assert.NoError(t, err)

		timer.End()

		stats := monitor.GetStats()
		assert.Equal(t, int64(1), stats.SerializeCount)
		assert.Greater(t, stats.TotalSerializeTime, time.Duration(0))
	})

	t.Run("error_handling_with_monitoring", func(t *testing.T) {
		monitor.Reset()

		timer := monitor.StartParseTimer()

		_, err := xyJson.ParseString("invalid json")
		assert.Error(t, err)

		timer.EndWithError()

		stats := monitor.GetStats()
		assert.Equal(t, int64(1), stats.ParseCount)
		assert.Equal(t, int64(1), stats.ErrorCount)
	})
}

// TestPerformanceMemoryUsage 测试性能监控的内存使用
// TestPerformanceMemoryUsage tests memory usage of performance monitoring
func TestPerformanceMemoryUsage(t *testing.T) {
	monitor := xyJson.NewPerformanceMonitor()

	// 测试性能监控本身不会造成内存泄漏
	testutil.AssertNoMemoryLeak(t, func() {
		for i := 0; i < 1000; i++ {
			monitor.RecordParse(time.Duration(i)*time.Microsecond, int64(i))
			monitor.RecordSerialize(time.Duration(i)*time.Microsecond, int64(i))
			if i%100 == 0 {
				monitor.RecordError()
			}
			_ = monitor.GetStats()
		}
		monitor.Reset()
	})
}

// BenchmarkPerformanceMonitor 性能监控器基准测试
// BenchmarkPerformanceMonitor benchmarks performance monitor
func BenchmarkPerformanceMonitor(b *testing.B) {
	monitor := xyJson.NewPerformanceMonitor()

	b.Run("RecordParse", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			monitor.RecordParse(100*time.Microsecond, 1024)
		}
	})

	b.Run("RecordSerialize", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			monitor.RecordSerialize(80*time.Microsecond, 512)
		}
	})

	b.Run("GetStats", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = monitor.GetStats()
		}
	})

	b.Run("TimedOperation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			timer := monitor.StartParseTimer()
			timer.End()
		}
	})

	b.Run("ConcurrentRecording", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				monitor.RecordParse(100*time.Microsecond, 1024)
			}
		})
	})
}
