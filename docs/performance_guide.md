# xyJson 性能优化指南 / Performance Optimization Guide

本文档提供了使用 xyJson 库时的性能优化建议和最佳实践。

This document provides performance optimization recommendations and best practices when using the xyJson library.

## 目录 / Table of Contents

1. [基本性能原则](#基本性能原则--basic-performance-principles)
2. [内存优化](#内存优化--memory-optimization)
3. [解析优化](#解析优化--parsing-optimization)
4. [序列化优化](#序列化优化--serialization-optimization)
5. [JSONPath查询优化](#jsonpath查询优化--jsonpath-query-optimization)
6. [并发优化](#并发优化--concurrency-optimization)
7. [监控和分析](#监控和分析--monitoring-and-analysis)
8. [常见性能陷阱](#常见性能陷阱--common-performance-pitfalls)

## 基本性能原则 / Basic Performance Principles

### 1. 重用对象 / Reuse Objects

```go
// ✅ 好的做法：重用解析器和序列化器
// ✅ Good practice: Reuse parsers and serializers
parser := xyJson.NewParser()
serializer := xyJson.NewSerializer()

for _, data := range jsonDataList {
    value, err := parser.Parse(data)
    if err != nil {
        continue
    }
    result, _ := serializer.Serialize(value)
    // 处理结果...
}

// ❌ 避免：每次都创建新实例
// ❌ Avoid: Creating new instances every time
for _, data := range jsonDataList {
    value, err := xyJson.Parse(data) // 内部创建新的解析器
    // ...
}
```

### 2. 启用对象池 / Enable Object Pool

```go
// 创建带对象池的工厂
// Create factory with object pool
pool := xyJson.NewObjectPoolWithOptions(&xyJson.ObjectPoolOptions{
    MaxValuePoolSize:  1000,
    MaxObjectPoolSize: 500,
    MaxArrayPoolSize:  500,
    EnablePooling:     true,
})

factory := xyJson.NewValueFactoryWithPool(pool)
parser := xyJson.NewParserWithFactory(factory)

// 定期检查池统计
// Regularly check pool statistics
stats := pool.GetStats()
fmt.Printf("Pool hit rate: %.2f%%\n", stats.PoolHitRate*100)
```

## 内存优化 / Memory Optimization

### 1. 控制最大深度 / Control Maximum Depth

```go
// 设置合理的最大深度以防止栈溢出
// Set reasonable maximum depth to prevent stack overflow
parser := xyJson.NewParser()
parser.SetMaxDepth(50) // 根据实际需求调整

serializer := xyJson.NewSerializerWithOptions(&xyJson.SerializeOptions{
    MaxDepth: 50,
})
```

### 2. 及时清理大对象 / Clean Up Large Objects Promptly

```go
// 处理大型JSON后及时清理
// Clean up promptly after processing large JSON
func processLargeJson(data []byte) {
    value, err := xyJson.Parse(data)
    if err != nil {
        return
    }
    
    // 处理数据...
    processData(value)
    
    // 显式清理（如果是对象或数组）
    // Explicit cleanup (if object or array)
    if obj, ok := value.(xyJson.IObject); ok {
        obj.Clear()
    } else if arr, ok := value.(xyJson.IArray); ok {
        arr.Clear()
    }
    
    // 建议手动触发GC（仅在处理非常大的数据时）
    // Suggest manual GC trigger (only when processing very large data)
    runtime.GC()
}
```

### 3. 使用流式处理 / Use Streaming Processing

```go
// 对于大型数组，考虑分批处理
// For large arrays, consider batch processing
func processBigArray(arr xyJson.IArray) {
    batchSize := 100
    length := arr.Length()
    
    for i := 0; i < length; i += batchSize {
        end := i + batchSize
        if end > length {
            end = length
        }
        
        // 处理批次
        // Process batch
        for j := i; j < end; j++ {
            item := arr.Get(j)
            processItem(item)
        }
        
        // 可选：在批次间暂停以允许GC
        // Optional: pause between batches to allow GC
        if i%1000 == 0 {
            runtime.GC()
            time.Sleep(time.Millisecond)
        }
    }
}
```

## 解析优化 / Parsing Optimization

### 1. 预分配缓冲区 / Pre-allocate Buffers

```go
// 如果知道大概的JSON大小，可以预分配
// Pre-allocate if you know approximate JSON size
func parseWithPreallocation(data []byte) {
    // 估算需要的容量
    estimatedSize := len(data) / 4 // 经验值
    
    obj := xyJson.CreateObjectWithCapacity(estimatedSize)
    // 或者
    arr := xyJson.CreateArrayWithCapacity(estimatedSize)
}
```

### 2. 避免重复解析 / Avoid Repeated Parsing

```go
// ✅ 缓存解析结果
// ✅ Cache parsing results
type JsonCache struct {
    cache map[string]xyJson.IValue
    mutex sync.RWMutex
}

func (c *JsonCache) Parse(jsonStr string) (xyJson.IValue, error) {
    c.mutex.RLock()
    if cached, exists := c.cache[jsonStr]; exists {
        c.mutex.RUnlock()
        return cached.Clone(), nil // 返回克隆以保证安全
    }
    c.mutex.RUnlock()
    
    value, err := xyJson.ParseString(jsonStr)
    if err != nil {
        return nil, err
    }
    
    c.mutex.Lock()
    c.cache[jsonStr] = value.Clone()
    c.mutex.Unlock()
    
    return value, nil
}
```

### 3. 使用字节切片而非字符串 / Use Byte Slices Instead of Strings

```go
// ✅ 更高效：直接使用字节切片
// ✅ More efficient: Use byte slices directly
value, err := xyJson.Parse(jsonBytes)

// ❌ 较低效：字符串需要额外的内存分配
// ❌ Less efficient: Strings require additional memory allocation
value, err := xyJson.ParseString(string(jsonBytes))
```

## 序列化优化 / Serialization Optimization

### 1. 选择合适的序列化选项 / Choose Appropriate Serialization Options

```go
// 对于网络传输，使用紧凑格式
// For network transmission, use compact format
compactSerializer := xyJson.CompactSerializer()

// 对于调试或日志，使用美化格式
// For debugging or logging, use pretty format
prettySerializer := xyJson.PrettySerializer("  ")

// 自定义选项以平衡性能和需求
// Custom options to balance performance and requirements
customSerializer := xyJson.NewSerializerWithOptions(&xyJson.SerializeOptions{
    Compact:    true,  // 紧凑格式更快
    EscapeHTML: false, // 如果不需要HTML转义，禁用以提高性能
    SortKeys:   false, // 如果不需要键排序，禁用以提高性能
})
```

### 2. 重用序列化器 / Reuse Serializers

```go
// 创建全局序列化器实例
// Create global serializer instances
var (
    compactSerializer = xyJson.CompactSerializer()
    prettySerializer  = xyJson.PrettySerializer("  ")
)

func serializeForAPI(value xyJson.IValue) ([]byte, error) {
    return compactSerializer.Serialize(value)
}

func serializeForLogging(value xyJson.IValue) (string, error) {
    return prettySerializer.SerializeToString(value)
}
```

## JSONPath查询优化 / JSONPath Query Optimization

### 1. 缓存编译的路径 / Cache Compiled Paths

```go
// 如果频繁使用相同的JSONPath，考虑缓存
// If frequently using the same JSONPath, consider caching
type PathQueryCache struct {
    queries map[string]xyJson.IPathQuery
    mutex   sync.RWMutex
}

func (c *PathQueryCache) Query(root xyJson.IValue, path string) (xyJson.IValue, error) {
    c.mutex.RLock()
    query, exists := c.queries[path]
    c.mutex.RUnlock()
    
    if !exists {
        query = xyJson.NewPathQuery()
        c.mutex.Lock()
        c.queries[path] = query
        c.mutex.Unlock()
    }
    
    return query.SelectOne(root, path)
}
```

### 2. 优化路径表达式 / Optimize Path Expressions

```go
// ✅ 更具体的路径更快
// ✅ More specific paths are faster
value, _ := xyJson.Get(root, "$.users[0].name")

// ❌ 避免过于宽泛的搜索
// ❌ Avoid overly broad searches
values, _ := xyJson.GetAll(root, "$..name") // 可能很慢

// ✅ 如果知道结构，使用精确路径
// ✅ If you know the structure, use precise paths
value, _ := xyJson.Get(root, "$.data.users[0].profile.name")
```

### 3. 批量查询优化 / Batch Query Optimization

```go
// 如果需要多个值，考虑一次性获取父对象
// If you need multiple values, consider getting the parent object once
func getMultipleFields(root xyJson.IValue) {
    // ❌ 多次查询
    // ❌ Multiple queries
    name, _ := xyJson.Get(root, "$.user.name")
    age, _ := xyJson.Get(root, "$.user.age")
    email, _ := xyJson.Get(root, "$.user.email")
    
    // ✅ 一次获取父对象
    // ✅ Get parent object once
    user, _ := xyJson.Get(root, "$.user")
    if userObj, ok := user.(xyJson.IObject); ok {
        name := userObj.Get("name")
        age := userObj.Get("age")
        email := userObj.Get("email")
    }
}
```

## 并发优化 / Concurrency Optimization

### 1. 读写分离 / Separate Read and Write Operations

```go
// 使用读写锁保护共享JSON数据
// Use read-write locks to protect shared JSON data
type SafeJsonData struct {
    data  xyJson.IValue
    mutex sync.RWMutex
}

func (s *SafeJsonData) Read(path string) (xyJson.IValue, error) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    return xyJson.Get(s.data, path)
}

func (s *SafeJsonData) Write(path string, value xyJson.IValue) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    return xyJson.Set(s.data, path, value)
}
```

### 2. 工作池模式 / Worker Pool Pattern

```go
// 使用工作池处理大量JSON数据
// Use worker pool to process large amounts of JSON data
func processJsonConcurrently(jsonDataList [][]byte, numWorkers int) {
    jobs := make(chan []byte, len(jsonDataList))
    results := make(chan xyJson.IValue, len(jsonDataList))
    
    // 启动工作者
    // Start workers
    for i := 0; i < numWorkers; i++ {
        go func() {
            parser := xyJson.NewParser() // 每个工作者有自己的解析器
            for data := range jobs {
                if value, err := parser.Parse(data); err == nil {
                    results <- value
                }
            }
        }()
    }
    
    // 发送任务
    // Send jobs
    for _, data := range jsonDataList {
        jobs <- data
    }
    close(jobs)
    
    // 收集结果
    // Collect results
    for i := 0; i < len(jsonDataList); i++ {
        result := <-results
        // 处理结果...
    }
}
```

## 监控和分析 / Monitoring and Analysis

### 1. 启用性能监控 / Enable Performance Monitoring

```go
// 在应用启动时启用监控
// Enable monitoring at application startup
func init() {
    xyJson.EnablePerformanceMonitoring()
}

// 定期检查性能统计
// Regularly check performance statistics
func logPerformanceStats() {
    stats := xyJson.GetPerformanceStats()
    log.Printf("JSON Performance Stats:")
    log.Printf("  Parse operations: %d", stats.ParseCount)
    log.Printf("  Average parse time: %v", stats.AverageParseTime())
    log.Printf("  Serialize operations: %d", stats.SerializeCount)
    log.Printf("  Average serialize time: %v", stats.AverageSerializeTime())
    log.Printf("  Peak memory usage: %d bytes", stats.PeakMemoryUsage)
    log.Printf("  Error count: %d", stats.ErrorCount)
}
```

### 2. 内存分析 / Memory Profiling

```go
// 在处理大量数据前启动内存分析
// Start memory profiling before processing large amounts of data
func analyzeMemoryUsage() {
    xyJson.StartMemoryProfiling()
    
    // 执行JSON操作...
    processLargeJsonData()
    
    xyJson.StopMemoryProfiling()
    
    // 分析内存趋势
    // Analyze memory trends
    trend, growth := xyJson.GetMemoryTrend()
    if growth > 0.1 { // 如果内存增长超过10%
        log.Printf("Warning: Memory growth detected: %s (%.2f%%)", trend, growth*100)
    }
}
```

### 3. 基准测试 / Benchmarking

```go
// 编写基准测试来测量性能
// Write benchmarks to measure performance
func BenchmarkJsonParsing(b *testing.B) {
    data := []byte(`{"name":"test","value":123,"active":true}`)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := xyJson.Parse(data)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkJsonSerialization(b *testing.B) {
    obj := xyJson.CreateObject()
    obj.Set("name", "test")
    obj.Set("value", 123)
    obj.Set("active", true)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := xyJson.Serialize(obj)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## 常见性能陷阱 / Common Performance Pitfalls

### 1. 避免频繁的类型转换 / Avoid Frequent Type Conversions

```go
// ❌ 避免：频繁的类型检查和转换
// ❌ Avoid: Frequent type checking and conversion
func processValue(value xyJson.IValue) {
    if value.Type() == xyJson.StringType {
        str, _ := xyJson.ToString(value)
        // 处理字符串...
    } else if value.Type() == xyJson.NumberType {
        num, _ := xyJson.ToFloat64(value)
        // 处理数字...
    }
    // 重复的类型检查...
}

// ✅ 更好：一次性类型检查
// ✅ Better: One-time type checking
func processValueEfficient(value xyJson.IValue) {
    switch value.Type() {
    case xyJson.StringType:
        if scalar, ok := value.(xyJson.IScalarValue); ok {
            str, _ := scalar.String()
            // 处理字符串...
        }
    case xyJson.NumberType:
        if scalar, ok := value.(xyJson.IScalarValue); ok {
            num, _ := scalar.Float64()
            // 处理数字...
        }
    }
}
```

### 2. 避免不必要的深拷贝 / Avoid Unnecessary Deep Copies

```go
// ❌ 避免：不必要的克隆
// ❌ Avoid: Unnecessary cloning
func processReadOnly(value xyJson.IValue) {
    cloned := value.Clone() // 如果只是读取，不需要克隆
    // 只读操作...
}

// ✅ 更好：只在需要修改时克隆
// ✅ Better: Clone only when modification is needed
func processWithModification(value xyJson.IValue) {
    // 先尝试只读操作
    readOnlyResult := readValue(value)
    
    // 只有在需要修改时才克隆
    if needsModification {
        cloned := value.Clone()
        modifyValue(cloned)
    }
}
```

### 3. 避免字符串拼接 / Avoid String Concatenation

```go
// ❌ 避免：在循环中进行字符串拼接
// ❌ Avoid: String concatenation in loops
func buildJsonString(items []string) string {
    result := "["
    for i, item := range items {
        if i > 0 {
            result += ","
        }
        result += `"` + item + `"`
    }
    result += "]"
    return result
}

// ✅ 更好：使用StringBuilder或直接构建JSON对象
// ✅ Better: Use StringBuilder or build JSON object directly
func buildJsonEfficient(items []string) (string, error) {
    arr := xyJson.CreateArray()
    for _, item := range items {
        arr.Append(item)
    }
    return xyJson.SerializeToString(arr)
}
```

## 性能测试建议 / Performance Testing Recommendations

### 1. 建立基准 / Establish Baselines

```bash
# 运行基准测试
# Run benchmarks
go test -bench=. -benchmem ./benchmark/

# 生成性能分析文件
# Generate profiling files
go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof

# 分析性能瓶颈
# Analyze performance bottlenecks
go tool pprof cpu.prof
go tool pprof mem.prof
```

### 2. 持续监控 / Continuous Monitoring

```go
// 在生产环境中定期收集性能指标
// Regularly collect performance metrics in production
func collectMetrics() {
    stats := xyJson.GetPerformanceStats()
    
    // 发送到监控系统
    // Send to monitoring system
    metrics.Gauge("json.parse.count").Set(float64(stats.ParseCount))
    metrics.Gauge("json.parse.avg_time").Set(float64(stats.AverageParseTime().Nanoseconds()))
    metrics.Gauge("json.serialize.count").Set(float64(stats.SerializeCount))
    metrics.Gauge("json.memory.peak").Set(float64(stats.PeakMemoryUsage))
}
```

通过遵循这些性能优化指南，您可以最大化 xyJson 库的性能，确保应用程序在处理JSON数据时保持高效和稳定。

By following these performance optimization guidelines, you can maximize the performance of the xyJson library and ensure your application remains efficient and stable when processing JSON data.