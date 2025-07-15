package xyJson

import (
	"sync"
	"sync/atomic"
	"time"
)

// objectPool 对象池实现
// objectPool implements the IObjectPool interface
type objectPool struct {
	valuePool  sync.Pool
	objectPool sync.Pool
	arrayPool  sync.Pool

	// 统计信息
	stats struct {
		totalAllocated int64
		totalReused    int64
		currentInUse   int64
	}

	// 配置选项
	maxPoolSize int
	enabled     bool
}

// NewObjectPool 创建新的对象池
// NewObjectPool creates a new object pool
func NewObjectPool() IObjectPool {
	return NewObjectPoolWithOptions(DefaultObjectPoolOptions())
}

// ObjectPoolOptions 对象池配置选项
// ObjectPoolOptions represents object pool configuration options
type ObjectPoolOptions struct {
	// MaxPoolSize 最大池大小（0表示无限制）
	// MaxPoolSize is the maximum pool size (0 means unlimited)
	MaxPoolSize int

	// Enabled 是否启用对象池
	// Enabled indicates whether the object pool is enabled
	Enabled bool

	// CleanupInterval 清理间隔
	// CleanupInterval is the cleanup interval
	CleanupInterval time.Duration
}

// DefaultObjectPoolOptions 返回默认对象池选项
// DefaultObjectPoolOptions returns default object pool options
func DefaultObjectPoolOptions() *ObjectPoolOptions {
	return &ObjectPoolOptions{
		MaxPoolSize:     1000, // 默认最大1000个对象
		Enabled:         true,
		CleanupInterval: 5 * time.Minute,
	}
}

// NewObjectPoolWithOptions 使用指定选项创建对象池
// NewObjectPoolWithOptions creates an object pool with specified options
func NewObjectPoolWithOptions(options *ObjectPoolOptions) IObjectPool {
	if options == nil {
		options = DefaultObjectPoolOptions()
	}

	pool := &objectPool{
		maxPoolSize: options.MaxPoolSize,
		enabled:     options.Enabled,
	}

	// 初始化sync.Pool
	pool.valuePool.New = func() interface{} {
		atomic.AddInt64(&pool.stats.totalAllocated, 1)
		return &scalarValue{}
	}

	pool.objectPool.New = func() interface{} {
		atomic.AddInt64(&pool.stats.totalAllocated, 1)
		return NewObject()
	}

	pool.arrayPool.New = func() interface{} {
		atomic.AddInt64(&pool.stats.totalAllocated, 1)
		return NewArray()
	}

	// 启动清理协程
	if options.CleanupInterval > 0 {
		go pool.cleanupRoutine(options.CleanupInterval)
	}

	return pool
}

// GetValue 从池中获取值对象
// GetValue gets a value object from the pool
func (p *objectPool) GetValue() IValue {
	if !p.enabled {
		return &scalarValue{}
	}

	atomic.AddInt64(&p.stats.currentInUse, 1)

	if value := p.valuePool.Get(); value != nil {
		atomic.AddInt64(&p.stats.totalReused, 1)
		if sv, ok := value.(*scalarValue); ok {
			sv.reset()
			return sv
		}
	}

	atomic.AddInt64(&p.stats.totalAllocated, 1)
	return &scalarValue{}
}

// PutValue 将值对象放回池中
// PutValue puts a value object back to the pool
func (p *objectPool) PutValue(value IValue) {
	if !p.enabled || value == nil {
		return
	}

	atomic.AddInt64(&p.stats.currentInUse, -1)

	// 只回收标量值
	if sv, ok := value.(*scalarValue); ok {
		sv.reset()
		p.valuePool.Put(sv)
	}
}

// GetObject 从池中获取对象
// GetObject gets an object from the pool
func (p *objectPool) GetObject() IObject {
	if !p.enabled {
		return NewObject()
	}

	atomic.AddInt64(&p.stats.currentInUse, 1)

	if obj := p.objectPool.Get(); obj != nil {
		atomic.AddInt64(&p.stats.totalReused, 1)
		if ov, ok := obj.(*objectValue); ok {
			ov.reset()
			return ov
		}
		if iobj, ok := obj.(IObject); ok {
			iobj.Clear()
			return iobj
		}
	}

	atomic.AddInt64(&p.stats.totalAllocated, 1)
	return NewObject()
}

// PutObject 将对象放回池中
// PutObject puts an object back to the pool
func (p *objectPool) PutObject(obj IObject) {
	if !p.enabled || obj == nil {
		return
	}

	atomic.AddInt64(&p.stats.currentInUse, -1)

	// 清空对象并放回池中
	obj.Clear()
	p.objectPool.Put(obj)
}

// GetArray 从池中获取数组
// GetArray gets an array from the pool
func (p *objectPool) GetArray() IArray {
	if !p.enabled {
		return NewArray()
	}

	atomic.AddInt64(&p.stats.currentInUse, 1)

	if arr := p.arrayPool.Get(); arr != nil {
		atomic.AddInt64(&p.stats.totalReused, 1)
		if av, ok := arr.(*arrayValue); ok {
			av.reset()
			return av
		}
		if iarr, ok := arr.(IArray); ok {
			iarr.Clear()
			return iarr
		}
	}

	atomic.AddInt64(&p.stats.totalAllocated, 1)
	return NewArray()
}

// PutArray 将数组放回池中
// PutArray puts an array back to the pool
func (p *objectPool) PutArray(arr IArray) {
	if !p.enabled || arr == nil {
		return
	}

	atomic.AddInt64(&p.stats.currentInUse, -1)

	// 清空数组并放回池中
	arr.Clear()
	p.arrayPool.Put(arr)
}

// GetStats 获取池统计信息
// GetStats gets pool statistics
func (p *objectPool) GetStats() *PoolStats {
	totalAllocated := atomic.LoadInt64(&p.stats.totalAllocated)
	totalReused := atomic.LoadInt64(&p.stats.totalReused)
	currentInUse := atomic.LoadInt64(&p.stats.currentInUse)

	var hitRate float64
	if totalAllocated > 0 {
		hitRate = float64(totalReused) / float64(totalAllocated) * 100
	}

	return &PoolStats{
		TotalAllocated: totalAllocated,
		TotalReused:    totalReused,
		CurrentInUse:   currentInUse,
		PoolHitRate:    hitRate,
	}
}

// reset 重置标量值状态
// reset resets the scalar value state
func (sv *scalarValue) reset() {
	sv.valueType = NullValueType
	sv.rawData = nil
}

// cleanupRoutine 定期清理协程
// cleanupRoutine is the periodic cleanup routine
func (p *objectPool) cleanupRoutine(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		// 这里可以添加清理逻辑，比如释放长时间未使用的对象
		// 由于sync.Pool会自动进行GC清理，这里主要用于统计信息重置
		p.resetStatsIfNeeded()
	}
}

// resetStatsIfNeeded 在需要时重置统计信息
// resetStatsIfNeeded resets statistics when needed
func (p *objectPool) resetStatsIfNeeded() {
	// 如果分配的对象数量过大，重置统计信息以避免溢出
	const maxStats = 1000000000 // 10亿

	if atomic.LoadInt64(&p.stats.totalAllocated) > maxStats {
		atomic.StoreInt64(&p.stats.totalAllocated, 0)
		atomic.StoreInt64(&p.stats.totalReused, 0)
		// 不重置currentInUse，因为它表示当前状态
	}
}

// SetEnabled 设置对象池是否启用
// SetEnabled sets whether the object pool is enabled
func (p *objectPool) SetEnabled(enabled bool) {
	p.enabled = enabled
}

// IsEnabled 检查对象池是否启用
// IsEnabled checks if the object pool is enabled
func (p *objectPool) IsEnabled() bool {
	return p.enabled
}

// Clear 清空对象池
// Clear clears the object pool
func (p *objectPool) Clear() {
	// 创建新的sync.Pool实例来清空池
	p.valuePool = sync.Pool{
		New: func() interface{} {
			atomic.AddInt64(&p.stats.totalAllocated, 1)
			return &scalarValue{}
		},
	}

	p.objectPool = sync.Pool{
		New: func() interface{} {
			atomic.AddInt64(&p.stats.totalAllocated, 1)
			return NewObject()
		},
	}

	p.arrayPool = sync.Pool{
		New: func() interface{} {
			atomic.AddInt64(&p.stats.totalAllocated, 1)
			return NewArray()
		},
	}

	// 重置统计信息
	atomic.StoreInt64(&p.stats.totalAllocated, 0)
	atomic.StoreInt64(&p.stats.totalReused, 0)
	atomic.StoreInt64(&p.stats.currentInUse, 0)
}

// 全局默认对象池实例
// Global default object pool instance
var defaultPool IObjectPool = NewObjectPool()

// GetDefaultPool 获取默认对象池
// GetDefaultPool gets the default object pool
func GetDefaultPool() IObjectPool {
	return defaultPool
}

// SetDefaultPool 设置默认对象池
// SetDefaultPool sets the default object pool
func SetDefaultPool(pool IObjectPool) {
	if pool != nil {
		defaultPool = pool
	}
}
