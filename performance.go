package xyJson

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// PerformanceMonitor 性能监控器
type PerformanceMonitor struct {
	config   PerformanceConfig
	stats    *performanceStats
	enabled  int64 // 使用原子操作
	startTime time.Time
	mu       sync.RWMutex
}

// performanceStats 性能统计信息
type performanceStats struct {
	// 解析统计
	parseCount     uint64
	parseTime      uint64 // 纳秒
	parseErrors    uint64
	parseBytes     uint64
	
	// 序列化统计
	serializeCount uint64
	serializeTime  uint64 // 纳秒
	serializeErrors uint64
	serializeBytes uint64
	
	// JSONPath统计
	jsonPathCount  uint64
	jsonPathTime   uint64 // 纳秒
	jsonPathErrors uint64
	
	// 内存统计
	memoryAllocated uint64
	memoryFreed     uint64
	gcCount         uint64
	
	// 对象池统计
	poolHits   uint64
	poolMisses uint64
}

// PerformanceStats 性能统计信息（导出版本）
type PerformanceStats struct {
	// 解析统计
	ParseCount      uint64        `json:"parse_count"`
	AvgParseTime    time.Duration `json:"avg_parse_time"`
	ParseErrors     uint64        `json:"parse_errors"`
	ParseThroughput float64       `json:"parse_throughput"` // MB/s
	
	// 序列化统计
	SerializeCount      uint64        `json:"serialize_count"`
	AvgSerializeTime    time.Duration `json:"avg_serialize_time"`
	SerializeErrors     uint64        `json:"serialize_errors"`
	SerializeThroughput float64       `json:"serialize_throughput"` // MB/s
	
	// JSONPath统计
	JSONPathCount    uint64        `json:"jsonpath_count"`
	AvgJSONPathTime  time.Duration `json:"avg_jsonpath_time"`
	JSONPathErrors   uint64        `json:"jsonpath_errors"`
	
	// 内存统计
	MemoryAllocated uint64  `json:"memory_allocated"`
	MemoryFreed     uint64  `json:"memory_freed"`
	GCCount         uint64  `json:"gc_count"`
	MemoryUsage     float64 `json:"memory_usage"` // MB
	
	// 对象池统计
	PoolHitRate float64 `json:"pool_hit_rate"`
	
	// 运行时统计
	Uptime       time.Duration `json:"uptime"`
	Goroutines   int           `json:"goroutines"`
	CPUUsage     float64       `json:"cpu_usage"`
}

// 全局性能监控器
var globalMonitor *PerformanceMonitor
var monitorOnce sync.Once

// GetPerformanceMonitor 获取全局性能监控器
func GetPerformanceMonitor() *PerformanceMonitor {
	monitorOnce.Do(func() {
		globalMonitor = NewPerformanceMonitor()
	})
	return globalMonitor
}

// NewPerformanceMonitor 创建新的性能监控器
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		config:    GetGlobalConfig().Performance,
		stats:     &performanceStats{},
		startTime: time.Now(),
	}
}

// Enable 启用性能监控
func (pm *PerformanceMonitor) Enable() {
	atomic.StoreInt64(&pm.enabled, 1)
}

// Disable 禁用性能监控
func (pm *PerformanceMonitor) Disable() {
	atomic.StoreInt64(&pm.enabled, 0)
}

// IsEnabled 检查是否启用
func (pm *PerformanceMonitor) IsEnabled() bool {
	return atomic.LoadInt64(&pm.enabled) == 1
}

// RecordParse 记录解析操作
func (pm *PerformanceMonitor) RecordParse(duration time.Duration, bytes int, err error) {
	if !pm.IsEnabled() {
		return
	}
	
	atomic.AddUint64(&pm.stats.parseCount, 1)
	atomic.AddUint64(&pm.stats.parseTime, uint64(duration.Nanoseconds()))
	atomic.AddUint64(&pm.stats.parseBytes, uint64(bytes))
	
	if err != nil {
		atomic.AddUint64(&pm.stats.parseErrors, 1)
	}
}

// RecordSerialize 记录序列化操作
func (pm *PerformanceMonitor) RecordSerialize(duration time.Duration, bytes int, err error) {
	if !pm.IsEnabled() {
		return
	}
	
	atomic.AddUint64(&pm.stats.serializeCount, 1)
	atomic.AddUint64(&pm.stats.serializeTime, uint64(duration.Nanoseconds()))
	atomic.AddUint64(&pm.stats.serializeBytes, uint64(bytes))
	
	if err != nil {
		atomic.AddUint64(&pm.stats.serializeErrors, 1)
	}
}

// RecordJSONPath 记录JSONPath操作
func (pm *PerformanceMonitor) RecordJSONPath(duration time.Duration, err error) {
	if !pm.IsEnabled() {
		return
	}
	
	atomic.AddUint64(&pm.stats.jsonPathCount, 1)
	atomic.AddUint64(&pm.stats.jsonPathTime, uint64(duration.Nanoseconds()))
	
	if err != nil {
		atomic.AddUint64(&pm.stats.jsonPathErrors, 1)
	}
}

// RecordMemory 记录内存操作
func (pm *PerformanceMonitor) RecordMemory(allocated, freed uint64) {
	if !pm.IsEnabled() {
		return
	}
	
	atomic.AddUint64(&pm.stats.memoryAllocated, allocated)
	atomic.AddUint64(&pm.stats.memoryFreed, freed)
}

// RecordPoolHit 记录对象池命中
func (pm *PerformanceMonitor) RecordPoolHit() {
	if !pm.IsEnabled() {
		return
	}
	
	atomic.AddUint64(&pm.stats.poolHits, 1)
}

// RecordPoolMiss 记录对象池未命中
func (pm *PerformanceMonitor) RecordPoolMiss() {
	if !pm.IsEnabled() {
		return
	}
	
	atomic.AddUint64(&pm.stats.poolMisses, 1)
}

// GetStats 获取性能统计信息
func (pm *PerformanceMonitor) GetStats() PerformanceStats {
	stats := PerformanceStats{
		ParseCount:      atomic.LoadUint64(&pm.stats.parseCount),
		ParseErrors:     atomic.LoadUint64(&pm.stats.parseErrors),
		SerializeCount:  atomic.LoadUint64(&pm.stats.serializeCount),
		SerializeErrors: atomic.LoadUint64(&pm.stats.serializeErrors),
		JSONPathCount:   atomic.LoadUint64(&pm.stats.jsonPathCount),
		JSONPathErrors:  atomic.LoadUint64(&pm.stats.jsonPathErrors),
		MemoryAllocated: atomic.LoadUint64(&pm.stats.memoryAllocated),
		MemoryFreed:     atomic.LoadUint64(&pm.stats.memoryFreed),
		Uptime:          time.Since(pm.startTime),
		Goroutines:      runtime.NumGoroutine(),
	}
	
	// 计算平均时间
	if stats.ParseCount > 0 {
		totalParseTime := atomic.LoadUint64(&pm.stats.parseTime)
		stats.AvgParseTime = time.Duration(totalParseTime / stats.ParseCount)
		
		// 计算吞吐量 (MB/s)
		totalParseBytes := atomic.LoadUint64(&pm.stats.parseBytes)
		if totalParseTime > 0 {
			stats.ParseThroughput = float64(totalParseBytes) / (float64(totalParseTime) / 1e9) / 1024 / 1024
		}
	}
	
	if stats.SerializeCount > 0 {
		totalSerializeTime := atomic.LoadUint64(&pm.stats.serializeTime)
		stats.AvgSerializeTime = time.Duration(totalSerializeTime / stats.SerializeCount)
		
		// 计算吞吐量 (MB/s)
		totalSerializeBytes := atomic.LoadUint64(&pm.stats.serializeBytes)
		if totalSerializeTime > 0 {
			stats.SerializeThroughput = float64(totalSerializeBytes) / (float64(totalSerializeTime) / 1e9) / 1024 / 1024
		}
	}
	
	if stats.JSONPathCount > 0 {
		totalJSONPathTime := atomic.LoadUint64(&pm.stats.jsonPathTime)
		stats.AvgJSONPathTime = time.Duration(totalJSONPathTime / stats.JSONPathCount)
	}
	
	// 计算对象池命中率
	poolHits := atomic.LoadUint64(&pm.stats.poolHits)
	poolMisses := atomic.LoadUint64(&pm.stats.poolMisses)
	if poolHits+poolMisses > 0 {
		stats.PoolHitRate = float64(poolHits) / float64(poolHits+poolMisses)
	}
	
	// 获取内存统计
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	stats.MemoryUsage = float64(memStats.Alloc) / 1024 / 1024 // MB
	stats.GCCount = uint64(memStats.NumGC)
	
	return stats
}

// Reset 重置统计信息
func (pm *PerformanceMonitor) Reset() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.stats = &performanceStats{}
	pm.startTime = time.Now()
}

// 包级别的便捷函数

// EnablePerformanceMonitoring 启用性能监控
func EnablePerformanceMonitoring() {
	GetPerformanceMonitor().Enable()
}

// DisablePerformanceMonitoring 禁用性能监控
func DisablePerformanceMonitoring() {
	GetPerformanceMonitor().Disable()
}

// GetPerformanceStats 获取性能统计信息
func GetPerformanceStats() PerformanceStats {
	return GetPerformanceMonitor().GetStats()
}

// ResetPerformanceStats 重置性能统计信息
func ResetPerformanceStats() {
	GetPerformanceMonitor().Reset()
}

// GetMemoryStats 获取内存统计信息
func GetMemoryStats() runtime.MemStats {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats
}
