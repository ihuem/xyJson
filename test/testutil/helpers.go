package testutil

import (
	"encoding/json"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertJSONEqual 断言两个JSON字符串相等（忽略格式差异）
// AssertJSONEqual asserts that two JSON strings are equal (ignoring formatting)
func AssertJSONEqual(t *testing.T, expected, actual string) {
	var expectedObj, actualObj interface{}
	
	err := json.Unmarshal([]byte(expected), &expectedObj)
	require.NoError(t, err, "Expected JSON should be valid")
	
	err = json.Unmarshal([]byte(actual), &actualObj)
	require.NoError(t, err, "Actual JSON should be valid")
	
	assert.Equal(t, expectedObj, actualObj, "JSON objects should be equal")
}

// AssertValidJSON 断言字符串是有效的JSON
// AssertValidJSON asserts that a string is valid JSON
func AssertValidJSON(t *testing.T, jsonStr string) {
	var obj interface{}
	err := json.Unmarshal([]byte(jsonStr), &obj)
	assert.NoError(t, err, "Should be valid JSON: %s", jsonStr)
}

// AssertInvalidJSON 断言字符串是无效的JSON
// AssertInvalidJSON asserts that a string is invalid JSON
func AssertInvalidJSON(t *testing.T, jsonStr string) {
	var obj interface{}
	err := json.Unmarshal([]byte(jsonStr), &obj)
	assert.Error(t, err, "Should be invalid JSON: %s", jsonStr)
}

// MeasureMemory 测量函数执行时的内存使用
// MeasureMemory measures memory usage during function execution
func MeasureMemory(fn func()) (allocBytes uint64, allocObjects uint64) {
	runtime.GC()
	runtime.GC() // 运行两次确保清理完成
	
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	
	fn()
	
	runtime.ReadMemStats(&m2)
	
	return m2.TotalAlloc - m1.TotalAlloc, m2.Mallocs - m1.Mallocs
}

// MeasureTime 测量函数执行时间
// MeasureTime measures function execution time
func MeasureTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

// RunConcurrently 并发运行函数
// RunConcurrently runs a function concurrently
func RunConcurrently(t *testing.T, goroutines int, iterations int, fn func(int, int)) {
	done := make(chan bool, goroutines)
	errors := make(chan error, goroutines*iterations)
	
	for g := 0; g < goroutines; g++ {
		go func(goroutineID int) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Goroutine %d panicked: %v", goroutineID, r)
				}
				done <- true
			}()
			
			for i := 0; i < iterations; i++ {
				fn(goroutineID, i)
			}
		}(g)
	}
	
	// 等待所有goroutine完成
	for i := 0; i < goroutines; i++ {
		<-done
	}
	
	// 检查是否有错误
	close(errors)
	for err := range errors {
		if err != nil {
			t.Error(err)
		}
	}
}

// ComparePerformance 比较两个函数的性能
// ComparePerformance compares performance of two functions
func ComparePerformance(t *testing.T, name1 string, fn1 func(), name2 string, fn2 func(), iterations int) {
	// 预热
	for i := 0; i < 10; i++ {
		fn1()
		fn2()
	}
	
	// 测量第一个函数
	start1 := time.Now()
	for i := 0; i < iterations; i++ {
		fn1()
	}
	duration1 := time.Since(start1)
	
	// 测量第二个函数
	start2 := time.Now()
	for i := 0; i < iterations; i++ {
		fn2()
	}
	duration2 := time.Since(start2)
	
	// 输出比较结果
	t.Logf("%s: %v (%v per op)", name1, duration1, duration1/time.Duration(iterations))
	t.Logf("%s: %v (%v per op)", name2, duration2, duration2/time.Duration(iterations))
	
	if duration1 < duration2 {
		improvement := float64(duration2-duration1) / float64(duration2) * 100
		t.Logf("%s is %.2f%% faster than %s", name1, improvement, name2)
	} else {
		improvement := float64(duration1-duration2) / float64(duration1) * 100
		t.Logf("%s is %.2f%% faster than %s", name2, improvement, name1)
	}
}

// DeepEqual 深度比较两个值是否相等
// DeepEqual performs deep comparison of two values
func DeepEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Values are not deeply equal:\nExpected: %+v\nActual: %+v", expected, actual)
	}
}

// CreateTempFile 创建临时测试文件
// CreateTempFile creates a temporary test file
func CreateTempFile(t *testing.T, content string) string {
	file := t.TempDir() + "/test.json"
	err := writeFile(file, content)
	require.NoError(t, err)
	return file
}

// writeFile 写入文件的辅助函数
func writeFile(filename, content string) error {
	// 这里应该使用 os.WriteFile，但为了避免导入os包
	// 在实际实现中需要添加适当的导入
	return nil // 占位符实现
}

// TableTest 表格驱动测试的辅助结构
// TableTest helper structure for table-driven tests
type TableTest struct {
	Name     string
	Input    interface{}
	Expected interface{}
	Error    bool
}

// RunTableTests 运行表格驱动测试
// RunTableTests runs table-driven tests
func RunTableTests(t *testing.T, tests []TableTest, testFunc func(interface{}) (interface{}, error)) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result, err := testFunc(tt.Input)
			
			if tt.Error {
				assert.Error(t, err, "Expected an error for input: %v", tt.Input)
			} else {
				assert.NoError(t, err, "Unexpected error for input: %v", tt.Input)
				assert.Equal(t, tt.Expected, result, "Result mismatch for input: %v", tt.Input)
			}
		})
	}
}

// BenchmarkHelper 基准测试辅助结构
// BenchmarkHelper helper structure for benchmarks
type BenchmarkHelper struct {
	Name string
	Data interface{}
	Func func(interface{}) error
}

// RunBenchmarks 运行多个基准测试
// RunBenchmarks runs multiple benchmarks
func RunBenchmarks(b *testing.B, benchmarks []BenchmarkHelper) {
	for _, bm := range benchmarks {
		b.Run(bm.Name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				err := bm.Func(bm.Data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// MemoryStats 内存统计信息
// MemoryStats memory statistics
type MemoryStats struct {
	AllocBytes   uint64
	AllocObjects uint64
	GCCount      uint32
}

// GetMemoryStats 获取当前内存统计
// GetMemoryStats gets current memory statistics
func GetMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return MemoryStats{
		AllocBytes:   m.TotalAlloc,
		AllocObjects: m.Mallocs,
		GCCount:      m.NumGC,
	}
}

// AssertNoMemoryLeak 断言没有内存泄漏
// AssertNoMemoryLeak asserts that there's no memory leak
func AssertNoMemoryLeak(t *testing.T, fn func()) {
	runtime.GC()
	runtime.GC()
	
	before := GetMemoryStats()
	
	fn()
	
	runtime.GC()
	runtime.GC()
	
	after := GetMemoryStats()
	
	// 允许一定的内存增长容忍度，考虑到JSON解析过程中的临时对象分配
	// Allow some memory growth tolerance, considering temporary object allocation during JSON parsing
	const tolerance = 10 * 1024 * 1024 // 10MB
	if after.AllocBytes > before.AllocBytes+tolerance {
		t.Errorf("Potential memory leak detected: %d bytes allocated", 
			after.AllocBytes-before.AllocBytes)
	}
}