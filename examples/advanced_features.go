// Package examples demonstrates advanced features of xyJson library
// 本包演示了 xyJson 库的高级功能
package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ihuem/xyJson"
)

func main2() {
	// 性能监控示例 / Performance monitoring example
	performanceMonitoringExample()

	// 对象池优化示例 / Object pool optimization example
	objectPoolExample()

	// 并发安全示例 / Concurrency safety example
	concurrencyExample()

	// 内存分析示例 / Memory profiling example
	memoryProfilingExample()

	// 自定义序列化选项 / Custom serialization options
	customSerializationExample()

	// 复杂JSONPath查询 / Complex JSONPath queries
	complexJsonPathExample()
}

// performanceMonitoringExample 演示性能监控功能
// performanceMonitoringExample demonstrates performance monitoring features
func performanceMonitoringExample() {
	fmt.Println("=== Performance Monitoring Example ===")

	// 启用性能监控
	// Enable performance monitoring
	xyJson.EnablePerformanceMonitoring()

	// 执行一些JSON操作
	// Perform some JSON operations
	for i := 0; i < 100; i++ {
		jsonStr := fmt.Sprintf(`{"id":%d,"name":"user%d","active":true}`, i, i)
		value, err := xyJson.ParseString(jsonStr)
		if err != nil {
			continue
		}

		// 序列化回字符串
		_, err = xyJson.SerializeToString(value)
		if err != nil {
			continue
		}
	}

	// 获取性能统计
	// Get performance statistics
	stats := xyJson.GetPerformanceStats()
	fmt.Printf("Parse operations: %d\n", stats.ParseCount)
	fmt.Printf("Serialize operations: %d\n", stats.SerializeCount)
	fmt.Printf("Total parse time: %v\n", stats.TotalParseTime)
	fmt.Printf("Total serialize time: %v\n", stats.TotalSerializeTime)
	fmt.Printf("Average parse time: %v\n", stats.AvgParseTime)
	fmt.Printf("Average serialize time: %v\n", stats.AvgSerializeTime)
	fmt.Printf("Peak memory usage: %d bytes\n", stats.MaxMemoryUsage)

	// 重置统计
	// Reset statistics
	xyJson.ResetPerformanceStats()
	fmt.Println("Performance statistics reset")

	fmt.Println()
}

// objectPoolExample 演示对象池优化功能
// objectPoolExample demonstrates object pool optimization features
func objectPoolExample() {
	fmt.Println("=== Object Pool Example ===")

	// 创建自定义对象池
	// Create custom object pool
	pool := xyJson.NewObjectPool()
	factory := xyJson.NewValueFactoryWithPool(pool)

	// 使用对象池创建多个对象
	// Create multiple objects using object pool
	for i := 0; i < 10; i++ {
		obj := factory.CreateObject()
		obj.Set("id", i)
		obj.Set("name", fmt.Sprintf("item%d", i))

		// 使用完后归还到池中（通常由库自动处理）
		// Return to pool after use (usually handled automatically by library)
		pool.PutObject(obj)
	}

	// 获取池统计信息
	// Get pool statistics
	stats := pool.GetStats()
	fmt.Printf("Total allocated: %d\n", stats.TotalAllocated)
	fmt.Printf("Total reused: %d\n", stats.TotalReused)
	fmt.Printf("Currently in use: %d\n", stats.CurrentInUse)
	fmt.Printf("Pool hit rate: %.2f%%\n", stats.PoolHitRate*100)

	fmt.Println()
}

// concurrencyExample 演示并发安全功能
// concurrencyExample demonstrates concurrency safety features
func concurrencyExample() {
	fmt.Println("=== Concurrency Example ===")

	// 创建共享的JSON对象
	// Create shared JSON object
	sharedData := `{"counter":0,"items":[]}`
	root, err := xyJson.ParseString(sharedData)
	if err != nil {
		log.Fatal("Parse error:", err)
	}

	var wg sync.WaitGroup
	var mu sync.RWMutex

	// 启动多个goroutine进行并发读写
	// Start multiple goroutines for concurrent read/write
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 10; j++ {
				// 读操作
				// Read operation
				mu.RLock()
				counter, _ := xyJson.Get(root, "$.counter")
				fmt.Printf("Goroutine %d read counter: %s\n", id, counter.String())
				mu.RUnlock()

				// 写操作（通过重新解析实现，实际应用中可能需要更复杂的同步机制）
				// Write operation (implemented via re-parsing, real applications may need more complex synchronization)
				mu.Lock()
				newData := fmt.Sprintf(`{"counter":%d,"items":["item%d_%d"]}`, j, id, j)
				newRoot, parseErr := xyJson.ParseString(newData)
				if parseErr == nil {
					root = newRoot
				}
				mu.Unlock()

				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Concurrency test completed")

	fmt.Println()
}

// memoryProfilingExample 演示内存分析功能
// memoryProfilingExample demonstrates memory profiling features
func memoryProfilingExample() {
	fmt.Println("=== Memory Profiling Example ===")

	// 开始内存分析
	// Start memory profiling
	xyJson.StartMemoryProfiling()

	// 执行一些内存密集型操作
	// Perform memory-intensive operations
	for i := 0; i < 50; i++ {
		// 创建大型JSON结构
		// Create large JSON structure
		obj := xyJson.CreateObject()
		for j := 0; j < 100; j++ {
			arr := xyJson.CreateArray()
			for k := 0; k < 10; k++ {
				arr.Append(fmt.Sprintf("item_%d_%d_%d", i, j, k))
			}
			obj.Set(fmt.Sprintf("field_%d", j), arr)
		}

		// 序列化大对象
		// Serialize large object
		_, err := xyJson.SerializeToString(obj)
		if err != nil {
			fmt.Printf("Serialization error: %v\n", err)
		}

		if i%10 == 0 {
			fmt.Printf("Processed %d iterations\n", i+1)
		}
	}

	// 停止内存分析
	// Stop memory profiling
	xyJson.StopMemoryProfiling()

	// 获取内存快照
	// Get memory snapshots
	snapshots := xyJson.GetMemorySnapshots()
	fmt.Printf("Total memory snapshots: %d\n", len(snapshots))

	if len(snapshots) > 0 {
		latest := xyJson.GetLatestMemorySnapshot()
		if latest != nil {
			fmt.Printf("Latest memory usage: %d bytes\n", latest.TotalAlloc)
			fmt.Printf("当前堆分配: %d bytes\n", latest.HeapAlloc)
			fmt.Printf("GC cycles: %d\n", latest.NumGC)
		}

		// 获取内存趋势
		// Get memory trend
		trend, growth := xyJson.GetMemoryTrend()
		fmt.Printf("Memory trend: %s (%.2f%% growth)\n", trend, growth*100)
	}

	fmt.Println()
}

// customSerializationExample 演示自定义序列化选项
// customSerializationExample demonstrates custom serialization options
func customSerializationExample() {
	fmt.Println("=== Custom Serialization Example ===")

	// 创建测试对象
	// Create test object
	obj := xyJson.CreateObject()
	obj.Set("zebra", "animal")
	obj.Set("apple", "fruit")
	obj.Set("banana", "fruit")
	obj.Set("cat", "animal")
	obj.Set("html_content", "<script>alert('test')</script>")

	// 默认序列化
	// Default serialization
	defaultJson, _ := xyJson.SerializeToString(obj)
	fmt.Println("Default serialization:")
	fmt.Println(defaultJson)

	// 创建自定义序列化器
	// Create custom serializer
	customSerializer := xyJson.NewSerializerWithOptions(&xyJson.SerializeOptions{
		Indent:     "  ",  // 2-space indentation
		Compact:    false, // Pretty format
		EscapeHTML: true,  // Escape HTML characters
		SortKeys:   true,  // Sort object keys
		MaxDepth:   10,    // Maximum nesting depth
	})

	// 使用自定义选项序列化
	// Serialize with custom options
	customJson, _ := customSerializer.SerializeToString(obj)
	fmt.Println("\nCustom serialization (sorted keys, HTML escaped, pretty):")
	fmt.Println(customJson)

	// 紧凑格式
	// Compact format
	compactSerializer := xyJson.CompactSerializer()
	compactJson, _ := compactSerializer.SerializeToString(obj)
	fmt.Println("\nCompact serialization:")
	fmt.Println(compactJson)

	fmt.Println()
}

// complexJsonPathExample 演示复杂的JSONPath查询
// complexJsonPathExample demonstrates complex JSONPath queries
func complexJsonPathExample() {
	fmt.Println("=== Complex JSONPath Example ===")

	// 创建复杂的JSON数据
	// Create complex JSON data
	jsonStr := `{
		"company": {
			"name": "TechCorp",
			"departments": [
				{
					"name": "Engineering",
					"employees": [
						{"name": "Alice", "salary": 80000, "skills": ["Go", "Python", "Docker"]},
						{"name": "Bob", "salary": 75000, "skills": ["JavaScript", "React", "Node.js"]},
						{"name": "Charlie", "salary": 90000, "skills": ["Go", "Kubernetes", "AWS"]}
					]
				},
				{
					"name": "Marketing",
					"employees": [
						{"name": "David", "salary": 60000, "skills": ["SEO", "Content", "Analytics"]},
						{"name": "Eve", "salary": 65000, "skills": ["Design", "Branding", "Social Media"]}
					]
				}
			]
		}
	}`

	root, err := xyJson.ParseString(jsonStr)
	if err != nil {
		log.Fatal("Parse error:", err)
	}

	// 复杂查询示例
	// Complex query examples

	// 1. 获取所有员工姓名
	// 1. Get all employee names
	names, _ := xyJson.GetAll(root, "$..employees[*].name")
	fmt.Println("All employee names:")
	for _, name := range names {
		fmt.Printf("  - %s\n", name.String())
	}

	// 2. 获取工程部门的所有员工
	// 2. Get all employees in Engineering department
	engEmployees, _ := xyJson.GetAll(root, "$.company.departments[?(@.name=='Engineering')].employees[*].name")
	fmt.Println("\nEngineering employees:")
	for _, emp := range engEmployees {
		fmt.Printf("  - %s\n", emp.String())
	}

	// 3. 获取所有薪资
	// 3. Get all salaries
	salaries, _ := xyJson.GetAll(root, "$..employees[*].salary")
	fmt.Println("\nAll salaries:")
	totalSalary := 0.0
	for _, salary := range salaries {
		if val, err := salary.(xyJson.IScalarValue).Float64(); err == nil {
			totalSalary += val
			fmt.Printf("  - $%.0f\n", val)
		}
	}
	fmt.Printf("Total salary budget: $%.0f\n", totalSalary)

	// 4. 获取掌握Go技能的员工
	// 4. Get employees with Go skills (this would require more complex filtering in a real implementation)
	fmt.Println("\nEmployees with Go skills:")
	allEmployees, _ := xyJson.GetAll(root, "$..employees[*]")
	for _, emp := range allEmployees {
		if empObj, ok := emp.(xyJson.IObject); ok {
			skills, _ := xyJson.GetAll(empObj, "$.skills[*]")
			for _, skill := range skills {
				if skill.String() == "\"Go\"" {
					name := empObj.Get("name")
					fmt.Printf("  - %s\n", name.String())
					break
				}
			}
		}
	}

	// 5. 统计各部门员工数量
	// 5. Count employees by department
	fmt.Println("\nEmployee count by department:")
	departments, _ := xyJson.GetAll(root, "$.company.departments[*]")
	for _, dept := range departments {
		if deptObj, ok := dept.(xyJson.IObject); ok {
			deptName := deptObj.Get("name").String()
			empCount := xyJson.Count(deptObj, "$.employees[*]")
			fmt.Printf("  %s: %d employees\n", deptName, empCount)
		}
	}

	fmt.Println()
}
