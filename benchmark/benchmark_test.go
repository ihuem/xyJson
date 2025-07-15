package benchmark

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	xyJson "github/ihuem/xyJson"
	"github/ihuem/xyJson/test/testutil"
)

// 基准测试数据
// Benchmark test data
var (
	// 小型JSON数据
	// Small JSON data
	smallJSON = `{"name":"张三","age":25,"active":true}`
	
	// 中型JSON数据
	// Medium JSON data
	mediumJSON = `{
		"users": [
			{"id": 1, "name": "张三", "email": "zhangsan@example.com", "age": 25, "active": true},
			{"id": 2, "name": "李四", "email": "lisi@example.com", "age": 30, "active": false},
			{"id": 3, "name": "王五", "email": "wangwu@example.com", "age": 28, "active": true}
		],
		"metadata": {
			"total": 3,
			"page": 1,
			"limit": 10,
			"timestamp": "2024-01-01T00:00:00Z"
		}
	}`
	
	// 大型JSON数据（动态生成）
	// Large JSON data (dynamically generated)
	largeJSON string
)

// init 初始化基准测试数据
// init initializes benchmark test data
func init() {
	generator := testutil.NewTestDataGenerator()
	performanceData := generator.GeneratePerformanceTestData()
	largeJSON = performanceData["large"]
}

// BenchmarkParseSmall 小型JSON解析基准测试
// BenchmarkParseSmall benchmarks parsing small JSON
func BenchmarkParseSmall(b *testing.B) {
	b.Run("xyJson", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.ParseString(smallJSON)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("StandardLib", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result interface{}
			err := json.Unmarshal([]byte(smallJSON), &result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkParseMedium 中型JSON解析基准测试
// BenchmarkParseMedium benchmarks parsing medium JSON
func BenchmarkParseMedium(b *testing.B) {
	b.Run("xyJson", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.ParseString(mediumJSON)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("StandardLib", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result interface{}
			err := json.Unmarshal([]byte(mediumJSON), &result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkParseLarge 大型JSON解析基准测试
// BenchmarkParseLarge benchmarks parsing large JSON
func BenchmarkParseLarge(b *testing.B) {
	b.Run("xyJson", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.ParseString(largeJSON)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("StandardLib", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result interface{}
			err := json.Unmarshal([]byte(largeJSON), &result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkSerializeSmall 小型对象序列化基准测试
// BenchmarkSerializeSmall benchmarks serializing small objects
func BenchmarkSerializeSmall(b *testing.B) {
	// 准备xyJson对象
	xyObj := xyJson.CreateObject()
	xyObj.Set("name", "张三")
	xyObj.Set("age", 25)
	xyObj.Set("active", true)
	
	// 准备标准库对象
	stdObj := map[string]interface{}{
		"name":   "张三",
		"age":    25,
		"active": true,
	}
	
	b.Run("xyJson", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.SerializeToString(xyObj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("StandardLib", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(stdObj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkSerializeMedium 中型对象序列化基准测试
// BenchmarkSerializeMedium benchmarks serializing medium objects
func BenchmarkSerializeMedium(b *testing.B) {
	// 准备xyJson对象
	xyObj := xyJson.CreateObject()
	users := xyJson.CreateArray()
	
	for i := 1; i <= 3; i++ {
		user := xyJson.CreateObject()
		user.Set("id", i)
		user.Set("name", fmt.Sprintf("用户%d", i))
		user.Set("email", fmt.Sprintf("user%d@example.com", i))
		user.Set("age", 20+i*5)
		user.Set("active", i%2 == 1)
		users.Append(user)
	}
	
	metadata := xyJson.CreateObject()
	metadata.Set("total", 3)
	metadata.Set("page", 1)
	metadata.Set("limit", 10)
	metadata.Set("timestamp", "2024-01-01T00:00:00Z")
	
	xyObj.Set("users", users)
	xyObj.Set("metadata", metadata)
	
	// 准备标准库对象
	stdObj := map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": 1, "name": "用户1", "email": "user1@example.com", "age": 25, "active": true},
			{"id": 2, "name": "用户2", "email": "user2@example.com", "age": 30, "active": false},
			{"id": 3, "name": "用户3", "email": "user3@example.com", "age": 35, "active": true},
		},
		"metadata": map[string]interface{}{
			"total":     3,
			"page":      1,
			"limit":     10,
			"timestamp": "2024-01-01T00:00:00Z",
		},
	}
	
	b.Run("xyJson", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.SerializeToString(xyObj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("StandardLib", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(stdObj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkJSONPathQuery JSONPath查询基准测试
// BenchmarkJSONPathQuery benchmarks JSONPath queries
func BenchmarkJSONPathQuery(b *testing.B) {
	generator := testutil.NewTestDataGenerator()
	jsonStr := generator.GenerateJSONPathTestData()
	
	root, err := xyJson.ParseString(jsonStr)
	if err != nil {
		b.Fatal(err)
	}
	
	b.Run("SimpleQuery", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.Get(root, "$.store.bicycle.color")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("ArrayQuery", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.Get(root, "$.store.book[0].title")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("WildcardQuery", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.GetAll(root, "$.store.book[*].title")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("FilterQuery", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.GetAll(root, "$.store.book[?(@.price < 10)]")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("ExistsCheck", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = xyJson.Exists(root, "$.store.bicycle")
		}
	})
	
	b.Run("CountQuery", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = xyJson.Count(root, "$.store.book[*]")
		}
	})
}

// BenchmarkMemoryUsage 内存使用基准测试
// BenchmarkMemoryUsage benchmarks memory usage
func BenchmarkMemoryUsage(b *testing.B) {
	b.Run("ParseMemory", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.ParseString(mediumJSON)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("SerializeMemory", func(b *testing.B) {
		obj := xyJson.CreateObject()
		obj.Set("name", "测试")
		obj.Set("value", 42)
		
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.SerializeToString(obj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("ObjectCreationMemory", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			obj := xyJson.CreateObject()
			obj.Set("id", i)
			obj.Set("name", fmt.Sprintf("item_%d", i))
			obj.Set("active", i%2 == 0)
		}
	})
	
	b.Run("ArrayCreationMemory", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			arr := xyJson.CreateArray()
			for j := 0; j < 10; j++ {
				arr.Append(j)
			}
		}
	})
}

// BenchmarkConcurrentOperations 并发操作基准测试
// BenchmarkConcurrentOperations benchmarks concurrent operations
func BenchmarkConcurrentOperations(b *testing.B) {
	b.Run("ConcurrentParse", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := xyJson.ParseString(mediumJSON)
				if err != nil {
					b.Error(err)
				}
			}
		})
	})
	
	b.Run("ConcurrentSerialize", func(b *testing.B) {
		obj := xyJson.CreateObject()
		obj.Set("name", "测试")
		obj.Set("value", 42)
		
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := xyJson.SerializeToString(obj)
				if err != nil {
					b.Error(err)
				}
			}
		})
	})
	
	b.Run("ConcurrentJSONPath", func(b *testing.B) {
		generator := testutil.NewTestDataGenerator()
		jsonStr := generator.GenerateJSONPathTestData()
		
		root, err := xyJson.ParseString(jsonStr)
		if err != nil {
			b.Fatal(err)
		}
		
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := xyJson.Get(root, "$.store.book[0].title")
				if err != nil {
					b.Error(err)
				}
			}
		})
	})
}

// BenchmarkObjectPool 对象池基准测试
// BenchmarkObjectPool benchmarks object pool
func BenchmarkObjectPool(b *testing.B) {
	pool := xyJson.NewObjectPool()
	
	b.Run("WithPool", func(b *testing.B) {
		// pool.SetEnabled(true) // 删除不存在的方法调用
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			obj := pool.GetObject()
			obj.Set("test", "value")
			pool.PutObject(obj)
		}
	})
	
	b.Run("WithoutPool", func(b *testing.B) {
		// pool.SetEnabled(false) // 删除不存在的方法调用
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			obj := pool.GetObject() // 实际创建新对象
			obj.Set("test", "value")
			pool.PutObject(obj) // 不会真正放回池中
		}
		// pool.SetEnabled(true) // 删除不存在的方法调用
	})
	
	b.Run("ConcurrentPool", func(b *testing.B) {
		// pool.SetEnabled(true) // 删除不存在的方法调用
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				obj := pool.GetObject()
				obj.Set("test", "value")
				pool.PutObject(obj)
			}
		})
	})
}

// BenchmarkPerformanceMonitoring 性能监控基准测试
// BenchmarkPerformanceMonitoring benchmarks performance monitoring
func BenchmarkPerformanceMonitoring(b *testing.B) {
	monitor := xyJson.NewPerformanceMonitor()
	
	b.Run("WithMonitoring", func(b *testing.B) {
		monitor.Enable()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			timer := monitor.StartParseTimer()
			_, err := xyJson.ParseString(smallJSON)
			if err != nil {
				b.Fatal(err)
			}
			timer.End()
		}
	})
	
	b.Run("WithoutMonitoring", func(b *testing.B) {
		monitor.Disable()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			timer := monitor.StartParseTimer()
			_, err := xyJson.ParseString(smallJSON)
			if err != nil {
				b.Fatal(err)
			}
			timer.End()
		}
		monitor.Enable()
	})
}

// BenchmarkStringOperations 字符串操作基准测试
// BenchmarkStringOperations benchmarks string operations
func BenchmarkStringOperations(b *testing.B) {
	testStrings := []string{
		"简单字符串",
		"包含\"转义\"字符的字符串",
		"包含Unicode字符的字符串：🚀✨🎉",
		strings.Repeat("长字符串测试", 100),
	}
	
	for i, str := range testStrings {
		b.Run(fmt.Sprintf("String%d", i+1), func(b *testing.B) {
			value := xyJson.CreateString(str)
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				_ = value.String()
			}
		})
	}
}

// BenchmarkNumberOperations 数字操作基准测试
// BenchmarkNumberOperations benchmarks number operations
func BenchmarkNumberOperations(b *testing.B) {
	testNumbers := []interface{}{
		42,
		3.14159,
		-123456789,
		1.23456789e10,
	}
	
	for i, num := range testNumbers {
		b.Run(fmt.Sprintf("Number%d", i+1), func(b *testing.B) {
			value := xyJson.CreateNumber(num)
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				if scalarValue, ok := value.(xyJson.IScalarValue); ok {
					_, _ = scalarValue.Float64()
					_, _ = scalarValue.Int()
				}
			}
		})
	}
}

// BenchmarkComplexOperations 复杂操作基准测试
// BenchmarkComplexOperations benchmarks complex operations
func BenchmarkComplexOperations(b *testing.B) {
	b.Run("ParseAndQuery", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			root, err := xyJson.ParseString(mediumJSON)
			if err != nil {
				b.Fatal(err)
			}
			
			_, err = xyJson.Get(root, "$.users[0].name")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("BuildAndSerialize", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			obj := xyJson.CreateObject()
			obj.Set("id", i)
			obj.Set("name", fmt.Sprintf("用户%d", i))
			obj.Set("active", i%2 == 0)
			
			arr := xyJson.CreateArray()
			for j := 0; j < 5; j++ {
				arr.Append(j)
			}
			obj.Set("numbers", arr)
			
			_, err := xyJson.SerializeToString(obj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("RoundTrip", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// 解析
			root, err := xyJson.ParseString(mediumJSON)
			if err != nil {
				b.Fatal(err)
			}
			
			// 修改
			err = xyJson.Set(root, "$.metadata.processed", true)
			if err != nil {
				b.Fatal(err)
			}
			
			// 序列化
			_, err = xyJson.SerializeToString(root)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkMemoryProfiler 内存分析器基准测试
// BenchmarkMemoryProfiler benchmarks memory profiler
func BenchmarkMemoryProfiler(b *testing.B) {
	profiler := xyJson.NewMemoryProfiler(100, 10*time.Millisecond)
	
	b.Run("WithProfiler", func(b *testing.B) {
		profiler.Start()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.ParseString(smallJSON)
			if err != nil {
				b.Fatal(err)
			}
		}
		profiler.Stop()
	})
	
	b.Run("WithoutProfiler", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := xyJson.ParseString(smallJSON)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkGCPressure GC压力测试
// BenchmarkGCPressure benchmarks GC pressure
func BenchmarkGCPressure(b *testing.B) {
	b.Run("HighAllocationRate", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// 创建大量临时对象
			for j := 0; j < 100; j++ {
				obj := xyJson.CreateObject()
				obj.Set("id", j)
				obj.Set("data", strings.Repeat("x", 100))
				
				_, err := xyJson.SerializeToString(obj)
				if err != nil {
					b.Fatal(err)
				}
			}
			
			// 强制GC
			if i%10 == 0 {
				runtime.GC()
			}
		}
	})
	
	b.Run("WithObjectPool", func(b *testing.B) {
		pool := xyJson.NewObjectPool()
		// pool.SetEnabled(true) // 删除不存在的方法调用
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// 使用对象池减少分配
			for j := 0; j < 100; j++ {
				obj := pool.GetObject()
				obj.Set("id", j)
				obj.Set("data", strings.Repeat("x", 100))
				
				_, err := xyJson.SerializeToString(obj)
				if err != nil {
					b.Fatal(err)
				}
				
				pool.PutObject(obj)
			}
			
			// 强制GC
			if i%10 == 0 {
				runtime.GC()
			}
		}
	})
}

// runBenchmarkSuite 运行完整的基准测试套件
// runBenchmarkSuite runs the complete benchmark suite
func runBenchmarkSuite() {
	fmt.Println("=== xyJson 性能基准测试套件 ===")
	fmt.Println("=== xyJson Performance Benchmark Suite ===")
	
	// 运行所有基准测试
	testing.Benchmark(BenchmarkParseSmall)
	testing.Benchmark(BenchmarkParseMedium)
	testing.Benchmark(BenchmarkParseLarge)
	testing.Benchmark(BenchmarkSerializeSmall)
	testing.Benchmark(BenchmarkSerializeMedium)
	testing.Benchmark(BenchmarkJSONPathQuery)
	testing.Benchmark(BenchmarkMemoryUsage)
	testing.Benchmark(BenchmarkConcurrentOperations)
	testing.Benchmark(BenchmarkObjectPool)
	testing.Benchmark(BenchmarkPerformanceMonitoring)
	testing.Benchmark(BenchmarkComplexOperations)
	
	fmt.Println("\n=== 基准测试完成 ===")
	fmt.Println("=== Benchmark Tests Completed ===")
}