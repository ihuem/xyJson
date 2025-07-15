package xyJson

import (
	"sync"
	"sync/atomic"
	"time"
)

// ObjectPool 对象池实现
type ObjectPool struct {
	config       ObjectPoolConfig
	objectPool   sync.Pool
	arrayPool    sync.Pool
	stats        *poolStats
	cleanupTimer *time.Timer
	mu           sync.RWMutex
}

// poolStats 池统计信息
type poolStats struct {
	totalAllocated uint64
	totalReused    uint64
	currentInUse   uint64
	objectHits     uint64
	objectMisses   uint64
	arrayHits      uint64
	arrayMisses    uint64
}

// NewObjectPool 创建新的对象池
func NewObjectPool() IObjectPool {
	return NewObjectPoolWithConfig(GetGlobalConfig().ObjectPool)
}

// NewObjectPoolWithConfig 使用指定配置创建对象池
func NewObjectPoolWithConfig(config ObjectPoolConfig) IObjectPool {
	pool := &ObjectPool{
		config: config,
		stats:  &poolStats{},
	}
	
	// 初始化对象池
	pool.objectPool = sync.Pool{
		New: func() interface{} {
			atomic.AddUint64(&pool.stats.totalAllocated, 1)
			atomic.AddUint64(&pool.stats.objectMisses, 1)
			return NewObject()
		},
	}
	
	// 初始化数组池
	pool.arrayPool = sync.Pool{
		New: func() interface{} {
			atomic.AddUint64(&pool.stats.totalAllocated, 1)
			atomic.AddUint64(&pool.stats.arrayMisses, 1)
			return NewArray()
		},
	}
	
	// 预填充池
	pool.preFillPool()
	
	// 启动清理定时器
	if config.CleanupInterval > 0 {
		pool.startCleanupTimer()
	}
	
	return pool
}

// GetObject 从池中获取对象
func (p *ObjectPool) GetObject() IObject {
	atomic.AddUint64(&p.stats.currentInUse, 1)
	atomic.AddUint64(&p.stats.objectHits, 1)
	
	obj := p.objectPool.Get().(IObject)
	
	// 重置对象状态
	if resettable, ok := obj.(*objectValue); ok {
		resettable.reset()
	}
	
	return obj
}

// PutObject 将对象放回池中
func (p *ObjectPool) PutObject(obj IObject) {
	if obj == nil {
		return
	}
	
	atomic.AddUint64(&p.stats.currentInUse, ^uint64(0)) // 减1
	atomic.AddUint64(&p.stats.totalReused, 1)
	
	// 清理对象
	obj.Clear()
	
	p.objectPool.Put(obj)
}

// GetArray 从池中获取数组
func (p *ObjectPool) GetArray() IArray {
	atomic.AddUint64(&p.stats.currentInUse, 1)
	atomic.AddUint64(&p.stats.arrayHits, 1)
	
	arr := p.arrayPool.Get().(IArray)
	
	// 重置数组状态
	if resettable, ok := arr.(*arrayValue); ok {
		resettable.reset()
	}
	
	return arr
}

// PutArray 将数组放回池中
func (p *ObjectPool) PutArray(arr IArray) {
	if arr == nil {
		return
	}
	
	atomic.AddUint64(&p.stats.currentInUse, ^uint64(0)) // 减1
	atomic.AddUint64(&p.stats.totalReused, 1)
	
	// 清理数组
	arr.Clear()
	
	p.arrayPool.Put(arr)
}

// GetStats 获取池统计信息
func (p *ObjectPool) GetStats() PoolStats {
	totalAllocated := atomic.LoadUint64(&p.stats.totalAllocated)
	totalReused := atomic.LoadUint64(&p.stats.totalReused)
	currentInUse := atomic.LoadUint64(&p.stats.currentInUse)
	
	var hitRate float64
	if totalAllocated > 0 {
		hitRate = float64(totalReused) / float64(totalAllocated)
	}
	
	return PoolStats{
		TotalAllocated: totalAllocated,
		TotalReused:    totalReused,
		CurrentInUse:   currentInUse,
		PoolHitRate:    hitRate,
	}
}

// Clear 清空池
func (p *ObjectPool) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// 重新创建池
	p.objectPool = sync.Pool{
		New: func() interface{} {
			atomic.AddUint64(&p.stats.totalAllocated, 1)
			return NewObject()
		},
	}
	
	p.arrayPool = sync.Pool{
		New: func() interface{} {
			atomic.AddUint64(&p.stats.totalAllocated, 1)
			return NewArray()
		},
	}
	
	// 重置统计
	atomic.StoreUint64(&p.stats.totalAllocated, 0)
	atomic.StoreUint64(&p.stats.totalReused, 0)
	atomic.StoreUint64(&p.stats.currentInUse, 0)
	atomic.StoreUint64(&p.stats.objectHits, 0)
	atomic.StoreUint64(&p.stats.objectMisses, 0)
	atomic.StoreUint64(&p.stats.arrayHits, 0)
	atomic.StoreUint64(&p.stats.arrayMisses, 0)
}

// preFillPool 预填充池
func (p *ObjectPool) preFillPool() {
	if p.config.InitialSize <= 0 {
		return
	}
	
	// 预创建对象
	for i := 0; i < p.config.InitialSize; i++ {
		obj := NewObject()
		p.objectPool.Put(obj)
	}
	
	// 预创建数组
	for i := 0; i < p.config.InitialSize; i++ {
		arr := NewArray()
		p.arrayPool.Put(arr)
	}
}

// startCleanupTimer 启动清理定时器
func (p *ObjectPool) startCleanupTimer() {
	p.cleanupTimer = time.AfterFunc(p.config.CleanupInterval, func() {
		p.cleanup()
		// 重新启动定时器
		p.startCleanupTimer()
	})
}

// cleanup 清理过期对象
func (p *ObjectPool) cleanup() {
	if !p.config.AutoResize {
		return
	}
	
	// 简单的清理策略：如果使用率低，减少池大小
	stats := p.GetStats()
	if stats.CurrentInUse < stats.TotalAllocated/4 {
		// 使用率低于25%，清理一些对象
		p.mu.Lock()
		// 这里可以实现更复杂的清理逻辑
		p.mu.Unlock()
	}
}

// 全局对象池实例
var defaultPool IObjectPool
var poolOnce sync.Once

// GetDefaultPool 获取默认对象池
func GetDefaultPool() IObjectPool {
	poolOnce.Do(func() {
		defaultPool = NewObjectPool()
	})
	return defaultPool
}

// SetDefaultPool 设置默认对象池
func SetDefaultPool(pool IObjectPool) {
	defaultPool = pool
}
