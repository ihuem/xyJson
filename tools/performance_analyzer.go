package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	xyJson "github/ihuem/xyJson" // 导入xyJson包
)

// PerformanceAnalyzer 性能分析器
type PerformanceAnalyzer struct {
	config    *AnalyzerConfig
	results   *AnalysisResults
	testData  []TestCase
	startTime time.Time
	memBefore runtime.MemStats
	memAfter  runtime.MemStats
}

// AnalyzerConfig 分析器配置
type AnalyzerConfig struct {
	Iterations   int           // 迭代次数
	WarmupRounds int           // 预热轮数
	CPUProfile   bool          // 是否启用CPU性能分析
	MemProfile   bool          // 是否启用内存性能分析
	GCProfile    bool          // 是否启用GC分析
	OutputDir    string        // 输出目录
	TestDataSize int           // 测试数据大小
	Concurrency  int           // 并发数
	Timeout      time.Duration // 超时时间
	Verbose      bool          // 详细输出
}

// TestCase 测试用例
type TestCase struct {
	Name        string
	Description string
	JSONData    []byte
	Size        int
	Complexity  string // simple, medium, complex
}

// AnalysisResults 分析结果
type AnalysisResults struct {
	ParseResults       []BenchmarkResult
	SerializeResults   []BenchmarkResult
	JSONPathResults    []BenchmarkResult
	MemoryResults      []MemoryResult
	ConcurrencyResults []ConcurrencyResult
	Recommendations    []string
	OverallScore       float64
	Timestamp          time.Time
}

// BenchmarkResult 基准测试结果
type BenchmarkResult struct {
	TestCase     string
	Operations   int64
	Duration     time.Duration
	OpsPerSecond float64
	AvgLatency   time.Duration
	MinLatency   time.Duration
	MaxLatency   time.Duration
	MemoryUsed   int64
	Allocations  int64
}

// MemoryResult 内存分析结果
type MemoryResult struct {
	TestCase       string
	HeapAlloc      uint64
	HeapSys        uint64
	HeapInuse      uint64
	StackInuse     uint64
	GCCycles       uint32
	GCPauseTotal   time.Duration
	PoolHitRate    float64
	PoolEfficiency float64
}

// ConcurrencyResult 并发测试结果
type ConcurrencyResult struct {
	Goroutines     int
	Operations     int64
	Duration       time.Duration
	Throughput     float64
	ErrorRate      float64
	ContentionRate float64
}

func main() {
	config := &AnalyzerConfig{
		Iterations:   10000,
		WarmupRounds: 1000,
		CPUProfile:   true,
		MemProfile:   true,
		GCProfile:    true,
		OutputDir:    "./performance_reports",
		TestDataSize: 1000,
		Concurrency:  10,
		Timeout:      time.Minute * 5,
		Verbose:      true,
	}

	analyzer := NewPerformanceAnalyzer(config)
	if err := analyzer.Run(); err != nil {
		log.Fatalf("性能分析失败: %v", err)
	}

	analyzer.GenerateReport()
	analyzer.PrintSummary()
}

// NewPerformanceAnalyzer 创建性能分析器
func NewPerformanceAnalyzer(config *AnalyzerConfig) *PerformanceAnalyzer {
	return &PerformanceAnalyzer{
		config: config,
		results: &AnalysisResults{
			Timestamp: time.Now(),
		},
		testData: generateTestData(config.TestDataSize),
	}
}

// Run 运行性能分析
func (pa *PerformanceAnalyzer) Run() error {
	// 创建输出目录
	if err := os.MkdirAll(pa.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 配置xyJson
	pa.setupXyJson()

	// 启动性能分析
	if pa.config.CPUProfile {
		if err := pa.startCPUProfile(); err != nil {
			return err
		}
		defer pa.stopCPUProfile()
	}

	log.Println("开始性能分析...")
	pa.startTime = time.Now()
	runtime.ReadMemStats(&pa.memBefore)

	// 预热
	pa.warmup()

	// 运行基准测试
	pa.runParseTests()
	pa.runSerializeTests()
	pa.runJSONPathTests()
	pa.runMemoryTests()
	pa.runConcurrencyTests()

	// 内存分析
	if pa.config.MemProfile {
		pa.generateMemProfile()
	}

	// 生成建议
	pa.generateRecommendations()

	runtime.ReadMemStats(&pa.memAfter)
	log.Printf("性能分析完成，耗时: %v", time.Since(pa.startTime))

	return nil
}

// setupXyJson 配置xyJson
func (pa *PerformanceAnalyzer) setupXyJson() {
	// 使用高性能配置
	config := xyJson.HighPerformanceConfig()
	xyJson.SetGlobalConfig(config)

	// 启用性能监控
	xyJson.EnablePerformanceMonitoring()

	// 重置统计
	xyJson.ResetPerformanceStats()

	log.Println("xyJson配置完成")
}

// warmup 预热
func (pa *PerformanceAnalyzer) warmup() {
	if pa.config.Verbose {
		log.Printf("预热 %d 轮...", pa.config.WarmupRounds)
	}

	for i := 0; i < pa.config.WarmupRounds; i++ {
		for _, testCase := range pa.testData {
			// 解析
			value, _ := xyJson.Parse(testCase.JSONData)
			// 序列化
			xyJson.Serialize(value)
			// JSONPath查询
			xyJson.Get(value, "$.name")
		}
	}

	// 强制GC
	runtime.GC()
	time.Sleep(time.Millisecond * 100)
}

// runParseTests 运行解析测试
func (pa *PerformanceAnalyzer) runParseTests() {
	if pa.config.Verbose {
		log.Println("运行解析性能测试...")
	}

	for _, testCase := range pa.testData {
		result := pa.benchmarkParse(testCase)
		pa.results.ParseResults = append(pa.results.ParseResults, result)
	}
}

// benchmarkParse 基准测试解析
func (pa *PerformanceAnalyzer) benchmarkParse(testCase TestCase) BenchmarkResult {
	var totalDuration time.Duration
	var minLatency = time.Hour
	var maxLatency time.Duration
	var memUsed int64
	var allocations int64

	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	start := time.Now()
	for i := 0; i < pa.config.Iterations; i++ {
		opStart := time.Now()
		_, err := xyJson.Parse(testCase.JSONData)
		opDuration := time.Since(opStart)

		if err != nil {
			continue
		}

		totalDuration += opDuration
		if opDuration < minLatency {
			minLatency = opDuration
		}
		if opDuration > maxLatency {
			maxLatency = opDuration
		}
	}
	duration := time.Since(start)

	runtime.ReadMemStats(&memAfter)
	memUsed = int64(memAfter.HeapAlloc - memBefore.HeapAlloc)
	allocations = int64(memAfter.Mallocs - memBefore.Mallocs)

	opsPerSecond := float64(pa.config.Iterations) / duration.Seconds()
	avgLatency := totalDuration / time.Duration(pa.config.Iterations)

	return BenchmarkResult{
		TestCase:     testCase.Name,
		Operations:   int64(pa.config.Iterations),
		Duration:     duration,
		OpsPerSecond: opsPerSecond,
		AvgLatency:   avgLatency,
		MinLatency:   minLatency,
		MaxLatency:   maxLatency,
		MemoryUsed:   memUsed,
		Allocations:  allocations,
	}
}

// runSerializeTests 运行序列化测试
func (pa *PerformanceAnalyzer) runSerializeTests() {
	if pa.config.Verbose {
		log.Println("运行序列化性能测试...")
	}

	for _, testCase := range pa.testData {
		value, _ := xyJson.Parse(testCase.JSONData)
		result := pa.benchmarkSerialize(testCase, value)
		pa.results.SerializeResults = append(pa.results.SerializeResults, result)
	}
}

// benchmarkSerialize 基准测试序列化
func (pa *PerformanceAnalyzer) benchmarkSerialize(testCase TestCase, value xyJson.IValue) BenchmarkResult {
	var totalDuration time.Duration
	var minLatency = time.Hour
	var maxLatency time.Duration
	var memUsed int64
	var allocations int64

	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	start := time.Now()
	for i := 0; i < pa.config.Iterations; i++ {
		opStart := time.Now()
		_, err := xyJson.Serialize(value)
		opDuration := time.Since(opStart)

		if err != nil {
			continue
		}

		totalDuration += opDuration
		if opDuration < minLatency {
			minLatency = opDuration
		}
		if opDuration > maxLatency {
			maxLatency = opDuration
		}
	}
	duration := time.Since(start)

	runtime.ReadMemStats(&memAfter)
	memUsed = int64(memAfter.HeapAlloc - memBefore.HeapAlloc)
	allocations = int64(memAfter.Mallocs - memBefore.Mallocs)

	opsPerSecond := float64(pa.config.Iterations) / duration.Seconds()
	avgLatency := totalDuration / time.Duration(pa.config.Iterations)

	return BenchmarkResult{
		TestCase:     testCase.Name + "_serialize",
		Operations:   int64(pa.config.Iterations),
		Duration:     duration,
		OpsPerSecond: opsPerSecond,
		AvgLatency:   avgLatency,
		MinLatency:   minLatency,
		MaxLatency:   maxLatency,
		MemoryUsed:   memUsed,
		Allocations:  allocations,
	}
}

// runJSONPathTests 运行JSONPath测试
func (pa *PerformanceAnalyzer) runJSONPathTests() {
	if pa.config.Verbose {
		log.Println("运行JSONPath性能测试...")
	}

	jsonPathQueries := []string{
		"$.name",
		"$.users[*].name",
		"$..email",
		"$.users[?(@.age > 25)]",
		"$.users[0:3]",
	}

	for _, testCase := range pa.testData {
		value, _ := xyJson.Parse(testCase.JSONData)
		for _, query := range jsonPathQueries {
			result := pa.benchmarkJSONPath(testCase, value, query)
			pa.results.JSONPathResults = append(pa.results.JSONPathResults, result)
		}
	}
}

// benchmarkJSONPath 基准测试JSONPath
func (pa *PerformanceAnalyzer) benchmarkJSONPath(testCase TestCase, value xyJson.IValue, query string) BenchmarkResult {
	var totalDuration time.Duration
	var minLatency = time.Hour
	var maxLatency time.Duration

	start := time.Now()
	for i := 0; i < pa.config.Iterations/10; i++ { // JSONPath测试减少迭代次数
		opStart := time.Now()
		_, err := xyJson.Get(value, query)
		opDuration := time.Since(opStart)

		if err != nil {
			continue
		}

		totalDuration += opDuration
		if opDuration < minLatency {
			minLatency = opDuration
		}
		if opDuration > maxLatency {
			maxLatency = opDuration
		}
	}
	duration := time.Since(start)

	iterations := pa.config.Iterations / 10
	opsPerSecond := float64(iterations) / duration.Seconds()
	avgLatency := totalDuration / time.Duration(iterations)

	return BenchmarkResult{
		TestCase:     fmt.Sprintf("%s_jsonpath_%s", testCase.Name, query),
		Operations:   int64(iterations),
		Duration:     duration,
		OpsPerSecond: opsPerSecond,
		AvgLatency:   avgLatency,
		MinLatency:   minLatency,
		MaxLatency:   maxLatency,
	}
}

// runMemoryTests 运行内存测试
func (pa *PerformanceAnalyzer) runMemoryTests() {
	if pa.config.Verbose {
		log.Println("运行内存性能测试...")
	}

	for _, testCase := range pa.testData {
		result := pa.analyzeMemoryUsage(testCase)
		pa.results.MemoryResults = append(pa.results.MemoryResults, result)
	}
}

// analyzeMemoryUsage 分析内存使用
func (pa *PerformanceAnalyzer) analyzeMemoryUsage(testCase TestCase) MemoryResult {
	var memBefore, memAfter runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memBefore)

	// 执行操作
	for i := 0; i < 1000; i++ {
		value, _ := xyJson.Parse(testCase.JSONData)
		xyJson.Serialize(value)
	}

	runtime.GC()
	runtime.ReadMemStats(&memAfter)

	// 获取对象池统计
	poolStats := xyJson.GetDefaultPool().GetStats()

	return MemoryResult{
		TestCase:       testCase.Name,
		HeapAlloc:      memAfter.HeapAlloc - memBefore.HeapAlloc,
		HeapSys:        memAfter.HeapSys - memBefore.HeapSys,
		HeapInuse:      memAfter.HeapInuse - memBefore.HeapInuse,
		StackInuse:     memAfter.StackInuse - memBefore.StackInuse,
		GCCycles:       memAfter.NumGC - memBefore.NumGC,
		GCPauseTotal:   time.Duration(memAfter.PauseTotalNs - memBefore.PauseTotalNs),
		PoolHitRate:    poolStats.PoolHitRate,
		PoolEfficiency: float64(poolStats.TotalReused) / float64(poolStats.TotalAllocated) * 100,
	}
}

// runConcurrencyTests 运行并发测试
func (pa *PerformanceAnalyzer) runConcurrencyTests() {
	if pa.config.Verbose {
		log.Println("运行并发性能测试...")
	}

	goroutineCounts := []int{1, 2, 4, 8, 16, 32}
	for _, count := range goroutineCounts {
		if count > pa.config.Concurrency {
			break
		}
		result := pa.benchmarkConcurrency(count)
		pa.results.ConcurrencyResults = append(pa.results.ConcurrencyResults, result)
	}
}

// benchmarkConcurrency 基准测试并发
func (pa *PerformanceAnalyzer) benchmarkConcurrency(goroutines int) ConcurrencyResult {
	operationsPerGoroutine := pa.config.Iterations / goroutines
	totalOperations := int64(operationsPerGoroutine * goroutines)

	done := make(chan bool, goroutines)
	errorCount := int64(0)

	start := time.Now()
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			for j := 0; j < operationsPerGoroutine; j++ {
				testCase := pa.testData[j%len(pa.testData)]
				value, err := xyJson.Parse(testCase.JSONData)
				if err != nil {
					errorCount++
					continue
				}
				_, err = xyJson.Serialize(value)
				if err != nil {
					errorCount++
				}
			}
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < goroutines; i++ {
		<-done
	}
	duration := time.Since(start)

	throughput := float64(totalOperations) / duration.Seconds()
	errorRate := float64(errorCount) / float64(totalOperations) * 100

	return ConcurrencyResult{
		Goroutines:     goroutines,
		Operations:     totalOperations,
		Duration:       duration,
		Throughput:     throughput,
		ErrorRate:      errorRate,
		ContentionRate: 0, // 简化实现
	}
}

// generateRecommendations 生成优化建议
func (pa *PerformanceAnalyzer) generateRecommendations() {
	recommendations := []string{}

	// 分析解析性能
	if len(pa.results.ParseResults) > 0 {
		avgOps := 0.0
		for _, result := range pa.results.ParseResults {
			avgOps += result.OpsPerSecond
		}
		avgOps /= float64(len(pa.results.ParseResults))

		if avgOps < 10000 {
			recommendations = append(recommendations, "解析性能较低，建议启用对象池优化")
		}
	}

	// 分析内存使用
	if len(pa.results.MemoryResults) > 0 {
		avgPoolHitRate := 0.0
		for _, result := range pa.results.MemoryResults {
			avgPoolHitRate += result.PoolHitRate
		}
		avgPoolHitRate /= float64(len(pa.results.MemoryResults))

		if avgPoolHitRate < 0.8 {
			recommendations = append(recommendations, "对象池命中率较低，建议调整池大小配置")
		}
	}

	// 分析并发性能
	if len(pa.results.ConcurrencyResults) > 0 {
		for _, result := range pa.results.ConcurrencyResults {
			if result.ErrorRate > 1.0 {
				recommendations = append(recommendations, "并发错误率较高，建议检查线程安全性")
				break
			}
		}
	}

	// 通用建议
	recommendations = append(recommendations, "建议在生产环境中启用性能监控")
	recommendations = append(recommendations, "定期分析性能指标并调整配置")
	recommendations = append(recommendations, "对于大型JSON数据，考虑使用流式处理")

	pa.results.Recommendations = recommendations

	// 计算总体评分
	pa.calculateOverallScore()
}

// calculateOverallScore 计算总体评分
func (pa *PerformanceAnalyzer) calculateOverallScore() {
	score := 100.0

	// 解析性能评分
	if len(pa.results.ParseResults) > 0 {
		avgOps := 0.0
		for _, result := range pa.results.ParseResults {
			avgOps += result.OpsPerSecond
		}
		avgOps /= float64(len(pa.results.ParseResults))

		if avgOps < 5000 {
			score -= 20
		} else if avgOps < 10000 {
			score -= 10
		}
	}

	// 内存效率评分
	if len(pa.results.MemoryResults) > 0 {
		avgPoolHitRate := 0.0
		for _, result := range pa.results.MemoryResults {
			avgPoolHitRate += result.PoolHitRate
		}
		avgPoolHitRate /= float64(len(pa.results.MemoryResults))

		if avgPoolHitRate < 0.6 {
			score -= 15
		} else if avgPoolHitRate < 0.8 {
			score -= 8
		}
	}

	// 并发性能评分
	if len(pa.results.ConcurrencyResults) > 0 {
		avgErrorRate := 0.0
		for _, result := range pa.results.ConcurrencyResults {
			avgErrorRate += result.ErrorRate
		}
		avgErrorRate /= float64(len(pa.results.ConcurrencyResults))

		if avgErrorRate > 2.0 {
			score -= 20
		} else if avgErrorRate > 1.0 {
			score -= 10
		}
	}

	if score < 0 {
		score = 0
	}

	pa.results.OverallScore = score
}

// startCPUProfile 启动CPU性能分析
func (pa *PerformanceAnalyzer) startCPUProfile() error {
	filename := fmt.Sprintf("%s/cpu_profile_%d.prof", pa.config.OutputDir, time.Now().Unix())
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	return pprof.StartCPUProfile(f)
}

// stopCPUProfile 停止CPU性能分析
func (pa *PerformanceAnalyzer) stopCPUProfile() {
	pprof.StopCPUProfile()
}

// generateMemProfile 生成内存性能分析
func (pa *PerformanceAnalyzer) generateMemProfile() {
	filename := fmt.Sprintf("%s/mem_profile_%d.prof", pa.config.OutputDir, time.Now().Unix())
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("创建内存性能分析文件失败: %v", err)
		return
	}
	defer f.Close()

	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Printf("写入内存性能分析失败: %v", err)
	}
}

// GenerateReport 生成详细报告
func (pa *PerformanceAnalyzer) GenerateReport() {
	filename := fmt.Sprintf("%s/performance_report_%d.json", pa.config.OutputDir, time.Now().Unix())

	// 使用xyJson生成报告
	report := xyJson.NewBuilder().
		SetTime("timestamp", pa.results.Timestamp).
		SetFloat64("overall_score", pa.results.OverallScore).
		BeginObject("configuration").
		SetInt("iterations", pa.config.Iterations).
		SetInt("warmup_rounds", pa.config.WarmupRounds).
		SetInt("test_data_size", pa.config.TestDataSize).
		SetInt("concurrency", pa.config.Concurrency).
		End().
		BeginArray("parse_results")

	for _, result := range pa.results.ParseResults {
		report.AddObject().
			SetString("test_case", result.TestCase).
			SetInt("operations", int(result.Operations)).
			SetString("duration", result.Duration.String()).
			SetFloat64("ops_per_second", result.OpsPerSecond).
			SetString("avg_latency", result.AvgLatency.String()).
			SetInt("memory_used", int(result.MemoryUsed)).
			End()
	}

	report.End(). // 结束parse_results数组
			BeginArray("recommendations")

	for _, rec := range pa.results.Recommendations {
		report.AddString(rec)
	}

	report.End() // 结束recommendations数组

	reportValue := report.MustBuild()
	reportJSON, _ := xyJson.Pretty(reportValue)

	if err := os.WriteFile(filename, []byte(reportJSON), 0644); err != nil {
		log.Printf("生成报告失败: %v", err)
	} else {
		log.Printf("性能报告已生成: %s", filename)
	}
}

// PrintSummary 打印摘要
func (pa *PerformanceAnalyzer) PrintSummary() {
	fmt.Println("\n=== xyJson 性能分析报告 ===")
	fmt.Printf("分析时间: %v\n", pa.results.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("总体评分: %.1f/100\n", pa.results.OverallScore)
	fmt.Printf("测试配置: %d次迭代, %d个测试用例\n", pa.config.Iterations, len(pa.testData))

	if len(pa.results.ParseResults) > 0 {
		fmt.Println("\n--- 解析性能 ---")
		for _, result := range pa.results.ParseResults {
			fmt.Printf("%s: %.0f ops/sec, 平均延迟: %v\n",
				result.TestCase, result.OpsPerSecond, result.AvgLatency)
		}
	}

	if len(pa.results.MemoryResults) > 0 {
		fmt.Println("\n--- 内存使用 ---")
		for _, result := range pa.results.MemoryResults {
			fmt.Printf("%s: 堆分配: %d bytes, 池命中率: %.1f%%\n",
				result.TestCase, result.HeapAlloc, result.PoolHitRate*100)
		}
	}

	if len(pa.results.ConcurrencyResults) > 0 {
		fmt.Println("\n--- 并发性能 ---")
		for _, result := range pa.results.ConcurrencyResults {
			fmt.Printf("%d goroutines: %.0f ops/sec, 错误率: %.2f%%\n",
				result.Goroutines, result.Throughput, result.ErrorRate)
		}
	}

	if len(pa.results.Recommendations) > 0 {
		fmt.Println("\n--- 优化建议 ---")
		for i, rec := range pa.results.Recommendations {
			fmt.Printf("%d. %s\n", i+1, rec)
		}
	}

	fmt.Printf("\n详细报告已保存到: %s\n", pa.config.OutputDir)
}

// generateTestData 生成测试数据
func generateTestData(size int) []TestCase {
	testCases := []TestCase{
		{
			Name:        "simple_object",
			Description: "简单JSON对象",
			Complexity:  "simple",
			JSONData:    []byte(`{"name":"测试","age":25,"active":true}`),
		},
		{
			Name:        "nested_object",
			Description: "嵌套JSON对象",
			Complexity:  "medium",
			JSONData:    []byte(`{"user":{"name":"张三","profile":{"email":"test@example.com","settings":{"theme":"dark","notifications":true}}}}`),
		},
		{
			Name:        "large_array",
			Description: "大型数组",
			Complexity:  "complex",
		},
	}

	// 生成大型数组测试数据
	largeArrayBuilder := xyJson.NewArrayBuilder()
	for i := 0; i < size; i++ {
		user := xyJson.NewBuilder().
			SetInt("id", i).
			SetString("name", fmt.Sprintf("用户%d", i)).
			SetString("email", fmt.Sprintf("user%d@example.com", i)).
			SetInt("age", 20+i%50).
			SetBool("active", i%2 == 0).
			MustBuild()
		largeArrayBuilder.AddValue(user)
	}
	largeArray := largeArrayBuilder.MustBuild()
	largeArrayJSON, _ := xyJson.Serialize(largeArray)
	testCases[2].JSONData = largeArrayJSON
	testCases[2].Size = len(largeArrayJSON)

	// 设置大小信息
	for i := range testCases {
		if testCases[i].Size == 0 {
			testCases[i].Size = len(testCases[i].JSONData)
		}
	}

	return testCases
}
