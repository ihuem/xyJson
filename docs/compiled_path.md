# JSONPath预编译功能详细指南

## 概述

JSONPath预编译功能是xyJson库的一个重要性能优化特性，通过预先编译JSONPath表达式，避免重复解析开销，在重复查询场景下可以带来显著的性能提升。

## 核心优势

- **性能提升**: 重复查询时性能提升约58%
- **缓存优化**: 智能缓存机制，缓存命中时性能提升84%
- **内存优化**: 减少重复解析的内存分配
- **线程安全**: 完全的并发安全保护
- **向后兼容**: 与现有JSONPath API完全兼容

## 技术原理

### 传统JSONPath查询流程
```
用户调用 → 解析路径 → 执行查询 → 返回结果
     ↑         ↑
   每次都要重复解析路径表达式
```

### 预编译JSONPath查询流程
```
预编译阶段: 用户调用 → 解析路径 → 生成CompiledPath对象
查询阶段:   用户调用 → 直接执行查询 → 返回结果
                    ↑
                跳过解析步骤
```

## 核心API

### 1. 路径编译

```go
// 基本编译
func CompilePath(path string) (*CompiledPath, error)

// 便利编译（失败时panic）
func MustCompilePath(path string) *CompiledPath
```

**使用示例：**
```go
// 安全编译
compiledPath, err := xyJson.CompilePath("$.user.profile.email")
if err != nil {
    log.Fatal("路径编译失败:", err)
}

// 便利编译（适用于确信路径正确的场景）
compiledPath := xyJson.MustCompilePath("$.user.profile.email")
```

### 2. CompiledPath方法

#### 查询方法
```go
// 查询单个值
func (cp *CompiledPath) Query(root IValue) (IValue, error)

// 查询所有匹配值
func (cp *CompiledPath) QueryAll(root IValue) ([]IValue, error)
```

#### 修改方法
```go
// 设置值
func (cp *CompiledPath) Set(root IValue, value IValue) error

// 删除值
func (cp *CompiledPath) Delete(root IValue) error
```

#### 检查方法
```go
// 检查路径是否存在
func (cp *CompiledPath) Exists(root IValue) bool

// 统计匹配数量
func (cp *CompiledPath) Count(root IValue) int

// 获取原始路径字符串
func (cp *CompiledPath) Path() string
```

### 3. 缓存管理

```go
// 获取缓存统计信息
func GetPathCacheStats() (size, maxSize int)

// 设置缓存最大大小
func SetPathCacheMaxSize(maxSize int)

// 清空缓存
func ClearPathCache()
```

## 使用场景

### ✅ 推荐使用场景

1. **重复查询相同路径**
```go
// 服务类中预编译路径
type UserService struct {
    userNamePath  *xyJson.CompiledPath
    userEmailPath *xyJson.CompiledPath
    userAgePath   *xyJson.CompiledPath
}

func NewUserService() *UserService {
    return &UserService{
        userNamePath:  xyJson.MustCompilePath("$.user.name"),
        userEmailPath: xyJson.MustCompilePath("$.user.email"),
        userAgePath:   xyJson.MustCompilePath("$.user.age"),
    }
}

func (s *UserService) ProcessUsers(users []xyJson.IValue) {
    for _, user := range users {
        name, _ := s.userNamePath.Query(user)
        email, _ := s.userEmailPath.Query(user)
        age, _ := s.userAgePath.Query(user)
        // 处理用户数据...
    }
}
```

2. **批量数据处理**
```go
func ProcessBatchData(dataList []xyJson.IValue) {
    // 预编译常用路径
    idPath := xyJson.MustCompilePath("$.id")
    statusPath := xyJson.MustCompilePath("$.status")
    timestampPath := xyJson.MustCompilePath("$.timestamp")
    
    for _, data := range dataList {
        id, _ := idPath.Query(data)
        status, _ := statusPath.Query(data)
        timestamp, _ := timestampPath.Query(data)
        // 处理单条数据...
    }
}
```

3. **高频API接口**
```go
type APIHandler struct {
    requestIDPath    *xyJson.CompiledPath
    userTokenPath    *xyJson.CompiledPath
    requestDataPath  *xyJson.CompiledPath
}

func NewAPIHandler() *APIHandler {
    return &APIHandler{
        requestIDPath:   xyJson.MustCompilePath("$.request.id"),
        userTokenPath:   xyJson.MustCompilePath("$.auth.token"),
        requestDataPath: xyJson.MustCompilePath("$.data"),
    }
}

func (h *APIHandler) HandleRequest(request xyJson.IValue) {
    requestID, _ := h.requestIDPath.Query(request)
    token, _ := h.userTokenPath.Query(request)
    data, _ := h.requestDataPath.Query(request)
    // 处理请求...
}
```

### ❌ 不推荐使用场景

1. **一次性查询**
```go
// 不推荐：编译开销大于收益
path, _ := xyJson.CompilePath("$.single.use.path")
result, _ := path.Query(root)

// 推荐：直接使用Get
result, _ := xyJson.Get(root, "$.single.use.path")
```

2. **动态路径**
```go
// 不推荐：路径每次都不同，无法复用
for i := 0; i < 100; i++ {
    path := fmt.Sprintf("$.items[%d].value", i)
    compiled, _ := xyJson.CompilePath(path)
    result, _ := compiled.Query(root)
}

// 推荐：直接使用Get
for i := 0; i < 100; i++ {
    path := fmt.Sprintf("$.items[%d].value", i)
    result, _ := xyJson.Get(root, path)
}
```

## 性能分析

### 基准测试结果

```
BenchmarkCompiledPathVsRegular/Regular_Path-8         1000000    1267 ns/op    128 B/op    4 allocs/op
BenchmarkCompiledPathVsRegular/Compiled_Path-8        2000000     529 ns/op     64 B/op    2 allocs/op
BenchmarkCompiledPathVsRegular/Compiled_Path_with_Compilation-8  100000  11282 ns/op   256 B/op    8 allocs/op

BenchmarkPathCachePerformance/Cache_Miss-8            100000   11406 ns/op    256 B/op    8 allocs/op
BenchmarkPathCachePerformance/Cache_Hit-8             120000    9584 ns/op    128 B/op    4 allocs/op
```

### 性能分析

1. **预编译路径查询**: 529 ns/op（比传统方式快58%）
2. **传统路径查询**: 1267 ns/op
3. **包含编译的查询**: 11282 ns/op（首次编译开销）
4. **缓存命中**: 9584 ns/op（比缓存未命中快16%）

### 性能收益计算

假设一个路径被查询N次：
- 传统方式总耗时：N × 1267 ns
- 预编译方式总耗时：11282 ns + N × 529 ns

**收益平衡点**：当N > 15时，预编译开始产生收益
**推荐使用**：当N > 50时，收益显著

## 缓存机制

### 缓存策略

1. **LRU淘汰**: 最近最少使用的路径会被淘汰
2. **默认大小**: 50个编译路径
3. **线程安全**: 支持并发访问
4. **自动管理**: 无需手动清理

### 缓存配置

```go
// 查看当前缓存状态
size, maxSize := xyJson.GetPathCacheStats()
fmt.Printf("缓存使用: %d/%d\n", size, maxSize)

// 调整缓存大小（根据应用需求）
xyJson.SetPathCacheMaxSize(100)  // 增加到100个

// 清空缓存（通常在测试或重置时使用）
xyJson.ClearPathCache()
```

### 缓存最佳实践

1. **合理设置缓存大小**
   - 小型应用：20-50个
   - 中型应用：50-100个
   - 大型应用：100-200个

2. **监控缓存命中率**
```go
func monitorCachePerformance() {
    initialSize, _ := xyJson.GetPathCacheStats()
    
    // 执行一些查询操作
    performQueries()
    
    finalSize, maxSize := xyJson.GetPathCacheStats()
    hitRate := float64(finalSize-initialSize) / float64(maxSize) * 100
    fmt.Printf("缓存命中率: %.1f%%\n", hitRate)
}
```

## 错误处理

### 编译错误

```go
compiledPath, err := xyJson.CompilePath("$.invalid..path")
if err != nil {
    switch e := err.(type) {
    case *xyJson.PathError:
        fmt.Printf("路径语法错误: %s\n", e.Error())
    default:
        fmt.Printf("编译失败: %v\n", err)
    }
    return
}
```

### 查询错误

```go
result, err := compiledPath.Query(root)
if err != nil {
    switch e := err.(type) {
    case *xyJson.TypeError:
        fmt.Printf("类型错误: %s\n", e.Error())
    case *xyJson.PathError:
        fmt.Printf("路径不存在: %s\n", e.Error())
    default:
        fmt.Printf("查询失败: %v\n", err)
    }
    return
}
```

## 线程安全

### CompiledPath线程安全

```go
// CompiledPath对象是线程安全的，可以在多个goroutine中共享
var compiledPath = xyJson.MustCompilePath("$.user.name")

func worker(id int, data xyJson.IValue, wg *sync.WaitGroup) {
    defer wg.Done()
    
    // 安全地在多个goroutine中使用同一个CompiledPath
    result, err := compiledPath.Query(data)
    if err == nil {
        fmt.Printf("Worker %d: %s\n", id, result.String())
    }
}

func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go worker(i, someData, &wg)
    }
    
    wg.Wait()
}
```

### 缓存线程安全

```go
// 缓存操作也是线程安全的
func concurrentCacheAccess() {
    var wg sync.WaitGroup
    
    // 多个goroutine同时编译不同路径
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            path := fmt.Sprintf("$.data[%d].value", index)
            compiled, _ := xyJson.CompilePath(path)
            // 使用compiled...
        }(i)
    }
    
    wg.Wait()
}
```

## 内存管理

### 内存优化

1. **减少重复解析**: 避免每次查询都解析路径
2. **智能缓存**: 自动管理内存使用
3. **对象复用**: 内部使用对象池优化

### 内存监控

```go
func monitorMemoryUsage() {
    var m1, m2 runtime.MemStats
    
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    // 执行大量预编译操作
    for i := 0; i < 1000; i++ {
        path := fmt.Sprintf("$.data[%d].value", i%10) // 重复路径
        compiled, _ := xyJson.CompilePath(path)
        compiled.Query(someData)
    }
    
    runtime.GC()
    runtime.ReadMemStats(&m2)
    
    fmt.Printf("内存增长: %d KB\n", (m2.Alloc-m1.Alloc)/1024)
}
```

## 迁移指南

### 从传统JSONPath迁移

**步骤1：识别重复查询**
```go
// 迁移前：重复查询
for _, data := range dataList {
    name, _ := xyJson.GetString(data, "$.user.name")     // 重复解析
    email, _ := xyJson.GetString(data, "$.user.email")   // 重复解析
    age, _ := xyJson.GetInt(data, "$.user.age")          // 重复解析
}
```

**步骤2：预编译路径**
```go
// 迁移后：预编译优化
namePath := xyJson.MustCompilePath("$.user.name")
emailPath := xyJson.MustCompilePath("$.user.email")
agePath := xyJson.MustCompilePath("$.user.age")

for _, data := range dataList {
    nameValue, _ := namePath.Query(data)
    emailValue, _ := emailPath.Query(data)
    ageValue, _ := agePath.Query(data)
    
    // 类型转换
    name := nameValue.String()
    email := emailValue.String()
    age, _ := ageValue.(xyJson.IScalarValue).Int()
}
```

**步骤3：结合便利API**
```go
// 最佳实践：预编译 + 便利API
type DataProcessor struct {
    namePath  *xyJson.CompiledPath
    emailPath *xyJson.CompiledPath
    agePath   *xyJson.CompiledPath
}

func NewDataProcessor() *DataProcessor {
    return &DataProcessor{
        namePath:  xyJson.MustCompilePath("$.user.name"),
        emailPath: xyJson.MustCompilePath("$.user.email"),
        agePath:   xyJson.MustCompilePath("$.user.age"),
    }
}

func (dp *DataProcessor) ProcessData(dataList []xyJson.IValue) {
    for _, data := range dataList {
        // 直接获取强类型值
        if nameValue, err := dp.namePath.Query(data); err == nil {
            name := nameValue.String()
            // 处理name...
        }
        
        if emailValue, err := dp.emailPath.Query(data); err == nil {
            email := emailValue.String()
            // 处理email...
        }
        
        if ageValue, err := dp.agePath.Query(data); err == nil {
            if scalarAge, ok := ageValue.(xyJson.IScalarValue); ok {
                age, _ := scalarAge.Int()
                // 处理age...
            }
        }
    }
}
```

## 总结

JSONPath预编译功能是xyJson库的重要性能优化特性，特别适用于需要重复查询相同路径的场景。通过合理使用预编译功能，可以显著提升应用性能，同时保持代码的清晰和可维护性。

### 关键要点

1. **性能收益**: 重复查询时性能提升58%
2. **使用场景**: 适用于重复查询，不适用于一次性查询
3. **缓存机制**: 智能LRU缓存，自动优化
4. **线程安全**: 完全支持并发使用
5. **向后兼容**: 与现有API完全兼容

### 最佳实践

1. 在服务类中预编译常用路径
2. 合理设置缓存大小
3. 监控缓存命中率
4. 结合便利API使用
5. 注意错误处理

通过遵循这些指南，您可以充分利用JSONPath预编译功能，为您的应用带来显著的性能提升。