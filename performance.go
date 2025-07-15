package xyJson

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// PerformanceMonitor 性能监控器
// PerformanceMonitor monitors performance metrics
type PerformanceMonitor struct {
	mu                 sync.RWMutex
	parseCount         int64
	serializeCount     int64
	parseTime          int64 // 纳秒
	serializeTime      int64 // 纳秒
	allocCount         int64
	allocBytes         int64
	gcCount            uint32
	maxMemoryUsage     int64
	currentMemoryUsage int64
	errorCount         int64
	lastResetTime      time.Time
	enabled            bool
}

// PerformanceStats 性能统计信息
// PerformanceStats contains performance statistics
type PerformanceStats struct {
	ParseCount         int64         `json:"parse_count"`
	SerializeCount     int64         `json:"serialize_count"`
	AvgParseTime       time.Duration `json:"avg_parse_time"`
	AvgSerializeTime   time.Duration `json:"avg_serialize_time"`
	TotalParseTime     time.Duration `json:"total_parse_time"`
	TotalSerializeTime time.Duration `json:"total_serialize_time"`
	AllocCount         int64         `json:"alloc_count"`
	AllocBytes         int64         `json:"alloc_bytes"`
	GCCount            uint32        `json:"gc_count"`
	MaxMemoryUsage     int64         `json:"max_memory_usage"`
	CurrentMemoryUsage int64         `json:"current_memory_usage"`
	ErrorCount         int64         `json:"error_count"`
	Uptime             time.Duration `json:"uptime"`
	Enabled            bool          `json:"enabled"`
}

// 全局性能监控器实例
// Global performance monitor instance
var (
	globalMonitor     *PerformanceMonitor
	globalMonitorOnce sync.Once
)

// GetGlobalMonitor 获取全局性能监控器
// GetGlobalMonitor gets the global performance monitor
func GetGlobalMonitor() *PerformanceMonitor {
	globalMonitorOnce.Do(func() {
		globalMonitor = NewPerformanceMonitor()
	})
	return globalMonitor
}

// NewPerformanceMonitor 创建新的性能监控器
// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		lastResetTime: time.Now(),
		enabled:       true,
	}
}

// Enable 启用性能监控
// Enable enables performance monitoring
func (pm *PerformanceMonitor) Enable() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.enabled = true
}

// Disable 禁用性能监控
// Disable disables performance monitoring
func (pm *PerformanceMonitor) Disable() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.enabled = false
}

// IsEnabled 检查是否启用性能监控
// IsEnabled checks if performance monitoring is enabled
func (pm *PerformanceMonitor) IsEnabled() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.enabled
}

// RecordParse 记录解析操作
// RecordParse records a parse operation
func (pm *PerformanceMonitor) RecordParse(duration time.Duration, allocBytes int64) {
	if !pm.IsEnabled() {
		return
	}

	atomic.AddInt64(&pm.parseCount, 1)
	atomic.AddInt64(&pm.parseTime, int64(duration))
	if allocBytes > 0 {
		atomic.AddInt64(&pm.allocCount, 1)
		atomic.AddInt64(&pm.allocBytes, allocBytes)
	}

	pm.updateMemoryUsage()
}

// RecordSerialize 记录序列化操作
// RecordSerialize records a serialize operation
func (pm *PerformanceMonitor) RecordSerialize(duration time.Duration, allocBytes int64) {
	if !pm.IsEnabled() {
		return
	}

	atomic.AddInt64(&pm.serializeCount, 1)
	atomic.AddInt64(&pm.serializeTime, int64(duration))
	if allocBytes > 0 {
		atomic.AddInt64(&pm.allocCount, 1)
		atomic.AddInt64(&pm.allocBytes, allocBytes)
	}

	pm.updateMemoryUsage()
}

// RecordError 记录错误
// RecordError records an error
func (pm *PerformanceMonitor) RecordError() {
	if !pm.IsEnabled() {
		return
	}

	atomic.AddInt64(&pm.errorCount, 1)
}

// updateMemoryUsage 更新内存使用情况
// updateMemoryUsage updates memory usage
func (pm *PerformanceMonitor) updateMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	currentUsage := int64(m.Alloc)
	atomic.StoreInt64(&pm.currentMemoryUsage, currentUsage)

	// 更新最大内存使用量
	for {
		maxUsage := atomic.LoadInt64(&pm.maxMemoryUsage)
		if currentUsage <= maxUsage {
			break
		}
		if atomic.CompareAndSwapInt64(&pm.maxMemoryUsage, maxUsage, currentUsage) {
			break
		}
	}

	// 更新GC计数
	atomic.StoreUint32(&pm.gcCount, m.NumGC)
}

// GetStats 获取性能统计信息
// GetStats gets performance statistics
func (pm *PerformanceMonitor) GetStats() PerformanceStats {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	parseCount := atomic.LoadInt64(&pm.parseCount)
	serializeCount := atomic.LoadInt64(&pm.serializeCount)
	totalParseTime := time.Duration(atomic.LoadInt64(&pm.parseTime))
	totalSerializeTime := time.Duration(atomic.LoadInt64(&pm.serializeTime))

	var avgParseTime, avgSerializeTime time.Duration
	if parseCount > 0 {
		avgParseTime = totalParseTime / time.Duration(parseCount)
	}
	if serializeCount > 0 {
		avgSerializeTime = totalSerializeTime / time.Duration(serializeCount)
	}

	return PerformanceStats{
		ParseCount:         parseCount,
		SerializeCount:     serializeCount,
		AvgParseTime:       avgParseTime,
		AvgSerializeTime:   avgSerializeTime,
		TotalParseTime:     totalParseTime,
		TotalSerializeTime: totalSerializeTime,
		AllocCount:         atomic.LoadInt64(&pm.allocCount),
		AllocBytes:         atomic.LoadInt64(&pm.allocBytes),
		GCCount:            atomic.LoadUint32(&pm.gcCount),
		MaxMemoryUsage:     atomic.LoadInt64(&pm.maxMemoryUsage),
		CurrentMemoryUsage: atomic.LoadInt64(&pm.currentMemoryUsage),
		ErrorCount:         atomic.LoadInt64(&pm.errorCount),
		Uptime:             time.Since(pm.lastResetTime),
		Enabled:            pm.enabled,
	}
}

// Reset 重置统计信息
// Reset resets statistics
func (pm *PerformanceMonitor) Reset() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	atomic.StoreInt64(&pm.parseCount, 0)
	atomic.StoreInt64(&pm.serializeCount, 0)
	atomic.StoreInt64(&pm.parseTime, 0)
	atomic.StoreInt64(&pm.serializeTime, 0)
	atomic.StoreInt64(&pm.allocCount, 0)
	atomic.StoreInt64(&pm.allocBytes, 0)
	atomic.StoreUint32(&pm.gcCount, 0)
	atomic.StoreInt64(&pm.maxMemoryUsage, 0)
	atomic.StoreInt64(&pm.currentMemoryUsage, 0)
	atomic.StoreInt64(&pm.errorCount, 0)
	pm.lastResetTime = time.Now()
}

// TimedOperation 计时操作包装器
// TimedOperation is a wrapper for timed operations
type TimedOperation struct {
	monitor   *PerformanceMonitor
	startTime time.Time
	startMem  int64
	opType    string
}

// StartParseTimer 开始解析计时
// StartParseTimer starts a parse timer
func (pm *PerformanceMonitor) StartParseTimer() *TimedOperation {
	if !pm.IsEnabled() {
		return nil
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &TimedOperation{
		monitor:   pm,
		startTime: time.Now(),
		startMem:  int64(m.Alloc),
		opType:    "parse",
	}
}

// StartSerializeTimer 开始序列化计时
// StartSerializeTimer starts a serialize timer
func (pm *PerformanceMonitor) StartSerializeTimer() *TimedOperation {
	if !pm.IsEnabled() {
		return nil
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &TimedOperation{
		monitor:   pm,
		startTime: time.Now(),
		startMem:  int64(m.Alloc),
		opType:    "serialize",
	}
}

// End 结束计时
// End ends the timing
func (to *TimedOperation) End() {
	if to == nil || to.monitor == nil {
		return
	}

	duration := time.Since(to.startTime)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	allocBytes := int64(m.Alloc) - to.startMem
	if allocBytes < 0 {
		allocBytes = 0
	}

	switch to.opType {
	case "parse":
		to.monitor.RecordParse(duration, allocBytes)
	case "serialize":
		to.monitor.RecordSerialize(duration, allocBytes)
	}
}

// EndWithError 结束计时并记录错误
// EndWithError ends the timing and records an error
func (to *TimedOperation) EndWithError() {
	if to == nil || to.monitor == nil {
		return
	}

	to.End()
	to.monitor.RecordError()
}

// MemoryProfiler 内存分析器
// MemoryProfiler analyzes memory usage
type MemoryProfiler struct {
	mu           sync.RWMutex
	snapshots    []MemorySnapshot
	maxSnapshots int
	interval     time.Duration
	stopChan     chan struct{}
	running      bool
}

// MemorySnapshot 内存快照
// MemorySnapshot represents a memory snapshot
type MemorySnapshot struct {
	Timestamp     time.Time `json:"timestamp"`
	Alloc         uint64    `json:"alloc"`
	TotalAlloc    uint64    `json:"total_alloc"`
	Sys           uint64    `json:"sys"`
	NumGC         uint32    `json:"num_gc"`
	GCCPUFraction float64   `json:"gc_cpu_fraction"`
	HeapAlloc     uint64    `json:"heap_alloc"`
	HeapSys       uint64    `json:"heap_sys"`
	HeapInuse     uint64    `json:"heap_inuse"`
	StackInuse    uint64    `json:"stack_inuse"`
}

// NewMemoryProfiler 创建新的内存分析器
// NewMemoryProfiler creates a new memory profiler
func NewMemoryProfiler(maxSnapshots int, interval time.Duration) *MemoryProfiler {
	if maxSnapshots <= 0 {
		maxSnapshots = DefaultMaxSnapshots
	}
	if interval <= 0 {
		interval = DefaultSnapshotInterval
	}

	return &MemoryProfiler{
		snapshots:    make([]MemorySnapshot, 0, maxSnapshots),
		maxSnapshots: maxSnapshots,
		interval:     interval,
		stopChan:     make(chan struct{}),
	}
}

// Start 开始内存分析
// Start begins memory profiling
func (mp *MemoryProfiler) Start() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.running {
		return
	}

	mp.running = true
	go mp.profileLoop()
}

// Stop 停止内存分析
// Stop stops memory profiling
func (mp *MemoryProfiler) Stop() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if !mp.running {
		return
	}

	mp.running = false
	close(mp.stopChan)
	mp.stopChan = make(chan struct{})
}

// profileLoop 分析循环
// profileLoop runs the profiling loop
func (mp *MemoryProfiler) profileLoop() {
	ticker := time.NewTicker(mp.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mp.takeSnapshot()
		case <-mp.stopChan:
			return
		}
	}
}

// takeSnapshot 拍摄内存快照
// takeSnapshot takes a memory snapshot
func (mp *MemoryProfiler) takeSnapshot() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	snapshot := MemorySnapshot{
		Timestamp:     time.Now(),
		Alloc:         m.Alloc,
		TotalAlloc:    m.TotalAlloc,
		Sys:           m.Sys,
		NumGC:         m.NumGC,
		GCCPUFraction: m.GCCPUFraction,
		HeapAlloc:     m.HeapAlloc,
		HeapSys:       m.HeapSys,
		HeapInuse:     m.HeapInuse,
		StackInuse:    m.StackInuse,
	}

	mp.mu.Lock()
	defer mp.mu.Unlock()

	// 添加快照
	mp.snapshots = append(mp.snapshots, snapshot)

	// 保持最大快照数量
	if len(mp.snapshots) > mp.maxSnapshots {
		copy(mp.snapshots, mp.snapshots[1:])
		mp.snapshots = mp.snapshots[:mp.maxSnapshots]
	}
}

// GetSnapshots 获取内存快照
// GetSnapshots gets memory snapshots
func (mp *MemoryProfiler) GetSnapshots() []MemorySnapshot {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	// 返回副本
	result := make([]MemorySnapshot, len(mp.snapshots))
	copy(result, mp.snapshots)
	return result
}

// GetLatestSnapshot 获取最新的内存快照
// GetLatestSnapshot gets the latest memory snapshot
func (mp *MemoryProfiler) GetLatestSnapshot() *MemorySnapshot {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	if len(mp.snapshots) == 0 {
		return nil
	}

	latest := mp.snapshots[len(mp.snapshots)-1]
	return &latest
}

// ClearSnapshots 清除所有快照
// ClearSnapshots clears all snapshots
func (mp *MemoryProfiler) ClearSnapshots() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.snapshots = mp.snapshots[:0]
}

// IsRunning 检查是否正在运行
// IsRunning checks if profiling is running
func (mp *MemoryProfiler) IsRunning() bool {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return mp.running
}

// GetMemoryTrend 获取内存趋势
// GetMemoryTrend gets memory trend
func (mp *MemoryProfiler) GetMemoryTrend() (trend string, growth float64) {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	if len(mp.snapshots) < 2 {
		return "unknown", 0
	}

	first := mp.snapshots[0]
	last := mp.snapshots[len(mp.snapshots)-1]

	if last.Alloc > first.Alloc {
		growth = float64(last.Alloc-first.Alloc) / float64(first.Alloc) * 100
		if growth > 10 {
			trend = "increasing"
		} else {
			trend = "stable"
		}
	} else if last.Alloc < first.Alloc {
		growth = -float64(first.Alloc-last.Alloc) / float64(first.Alloc) * 100
		if growth < -10 {
			trend = "decreasing"
		} else {
			trend = "stable"
		}
	} else {
		trend = "stable"
		growth = 0
	}

	return trend, growth
}

// 全局内存分析器实例
// Global memory profiler instance
var (
	globalProfiler     *MemoryProfiler
	globalProfilerOnce sync.Once
)

// GetGlobalProfiler 获取全局内存分析器
// GetGlobalProfiler gets the global memory profiler
func GetGlobalProfiler() *MemoryProfiler {
	globalProfilerOnce.Do(func() {
		globalProfiler = NewMemoryProfiler(DefaultMaxSnapshots, DefaultSnapshotInterval)
	})
	return globalProfiler
}

// ForceGC 强制垃圾回收
// ForceGC forces garbage collection
func ForceGC() {
	runtime.GC()
	runtime.GC() // 调用两次确保完全清理
}

// GetMemoryStats 获取当前内存统计
// GetMemoryStats gets current memory statistics
func GetMemoryStats() MemorySnapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemorySnapshot{
		Timestamp:     time.Now(),
		Alloc:         m.Alloc,
		TotalAlloc:    m.TotalAlloc,
		Sys:           m.Sys,
		NumGC:         m.NumGC,
		GCCPUFraction: m.GCCPUFraction,
		HeapAlloc:     m.HeapAlloc,
		HeapSys:       m.HeapSys,
		HeapInuse:     m.HeapInuse,
		StackInuse:    m.StackInuse,
	}
}
