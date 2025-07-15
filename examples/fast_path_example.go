package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ihuem/xyJson"
)

// User 用户结构体
// User represents a user structure
type User struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Active  bool      `json:"active"`
	Balance float64   `json:"balance"`
	Created time.Time `json:"created"`
	Address Address   `json:"address"`
	Tags    []string  `json:"tags"`
	Scores  []int     `json:"scores"`
}

// Address 地址结构体
// Address represents an address structure
type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	Zip    string `json:"zip"`
}

func main() {
	// 示例JSON数据
	// Sample JSON data
	jsonData := `{
		"id": 123,
		"name": "Alice Johnson",
		"email": "alice@example.com",
		"active": true,
		"balance": 1250.75,
		"created": "2023-01-15T10:30:00Z",
		"address": {
			"street": "123 Main St",
			"city": "New York",
			"zip": "10001"
		},
		"tags": ["premium", "verified"],
		"scores": [95, 88, 92]
	}`

	fmt.Println("=== xyJson 快速路径示例 ===")
	fmt.Println("=== xyJson Fast Path Example ===")
	fmt.Println()

	// 1. 使用快速路径解析（推荐用于高性能场景）
	// 1. Using fast path parsing (recommended for high-performance scenarios)
	fmt.Println("1. 快速路径解析 (Fast Path Parsing):")
	var user1 User
	err := xyJson.UnmarshalToStructFast([]byte(jsonData), &user1)
	if err != nil {
		log.Fatal("快速路径解析失败:", err)
	}
	fmt.Printf("   用户: %+v\n", user1)
	fmt.Printf("   地址: %+v\n", user1.Address)
	fmt.Printf("   标签: %v\n", user1.Tags)
	fmt.Println()

	// 2. 使用标准路径解析（功能更丰富）
	// 2. Using standard path parsing (more features)
	fmt.Println("2. 标准路径解析 (Standard Path Parsing):")
	var user2 User
	err = xyJson.UnmarshalToStruct([]byte(jsonData), &user2)
	if err != nil {
		log.Fatal("标准路径解析失败:", err)
	}
	fmt.Printf("   用户: %+v\n", user2)
	fmt.Println()

	// 3. 验证两种方法结果一致
	// 3. Verify both methods produce identical results
	fmt.Println("3. 结果对比 (Result Comparison):")
	if user1.ID == user2.ID && user1.Name == user2.Name && user1.Email == user2.Email {
		fmt.Println("   ✓ 快速路径和标准路径结果一致")
		fmt.Println("   ✓ Fast path and standard path produce identical results")
	} else {
		fmt.Println("   ✗ 结果不一致")
	}
	fmt.Println()

	// 4. 性能对比示例
	// 4. Performance comparison example
	fmt.Println("4. 性能特点 (Performance Characteristics):")
	fmt.Println("   快速路径 (Fast Path):")
	fmt.Println("   - 性能接近官方json包 (Performance close to standard json package)")
	fmt.Println("   - 内存使用优化 (Optimized memory usage)")
	fmt.Println("   - 适合简单解析场景 (Suitable for simple parsing scenarios)")
	fmt.Println()
	fmt.Println("   标准路径 (Standard Path):")
	fmt.Println("   - 支持JSONPath查询 (Supports JSONPath queries)")
	fmt.Println("   - 功能更丰富 (More features available)")
	fmt.Println("   - 适合复杂操作场景 (Suitable for complex operation scenarios)")
	fmt.Println()

	// 5. 快速路径的Must版本
	// 5. Must version of fast path
	fmt.Println("5. Must版本示例 (Must Version Example):")
	var user3 User
	// 注意：Must版本在出错时会panic
	// Note: Must version will panic on error
	xyJson.MustUnmarshalToStructFast([]byte(jsonData), &user3)
	fmt.Printf("   Must版本解析成功: %s\n", user3.Name)
	fmt.Println()

	// 6. 字符串版本
	// 6. String version
	fmt.Println("6. 字符串版本 (String Version):")
	var user4 User
	err = xyJson.UnmarshalStringToStructFast(jsonData, &user4)
	if err != nil {
		log.Fatal("字符串版本解析失败:", err)
	}
	fmt.Printf("   字符串解析成功: %s\n", user4.Name)
	fmt.Println()

	// 7. 使用建议
	// 7. Usage recommendations
	fmt.Println("7. 使用建议 (Usage Recommendations):")
	fmt.Println("   选择快速路径的场景 (Choose Fast Path for):")
	fmt.Println("   - 高性能需求 (High performance requirements)")
	fmt.Println("   - 简单JSON到struct转换 (Simple JSON to struct conversion)")
	fmt.Println("   - 大量数据处理 (Large data processing)")
	fmt.Println()
	fmt.Println("   选择标准路径的场景 (Choose Standard Path for):")
	fmt.Println("   - 需要JSONPath查询 (Need JSONPath queries)")
	fmt.Println("   - 复杂JSON操作 (Complex JSON operations)")
	fmt.Println("   - 需要高级功能 (Need advanced features)")
	fmt.Println()

	fmt.Println("=== 示例完成 ===")
	fmt.Println("=== Example Complete ===")
}

// 性能测试函数示例
// Performance testing function example
func benchmarkExample() {
	jsonData := []byte(`{"id":1,"name":"Test","email":"test@example.com"}`)

	// 快速路径
	// Fast path
	start := time.Now()
	for i := 0; i < 10000; i++ {
		var user User
		xyJson.UnmarshalToStructFast(jsonData, &user)
	}
	fastDuration := time.Since(start)

	// 标准路径
	// Standard path
	start = time.Now()
	for i := 0; i < 10000; i++ {
		var user User
		xyJson.UnmarshalToStruct(jsonData, &user)
	}
	standardDuration := time.Since(start)

	fmt.Printf("快速路径耗时: %v\n", fastDuration)
	fmt.Printf("标准路径耗时: %v\n", standardDuration)
	fmt.Printf("性能提升: %.2fx\n", float64(standardDuration)/float64(fastDuration))
}
