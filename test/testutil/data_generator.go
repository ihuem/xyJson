package testutil

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"
)

// TestDataGenerator 测试数据生成器
// TestDataGenerator generates test data for various testing scenarios
type TestDataGenerator struct {
	rand *rand.Rand
}

// NewTestDataGenerator 创建测试数据生成器
// NewTestDataGenerator creates a new test data generator
func NewTestDataGenerator() *TestDataGenerator {
	return &TestDataGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateSimpleObject 生成简单对象
// GenerateSimpleObject generates a simple JSON object for testing
func (g *TestDataGenerator) GenerateSimpleObject() map[string]interface{} {
	return map[string]interface{}{
		"name":     "测试用户",
		"age":      g.rand.Intn(100),
		"active":   g.rand.Float32() > 0.5,
		"score":    g.rand.Float64() * 100,
		"tags":     []string{"tag1", "tag2", "tag3"},
		"metadata": nil,
	}
}

// GenerateComplexObject 生成复杂对象
// GenerateComplexObject generates a complex nested JSON object
func (g *TestDataGenerator) GenerateComplexObject(depth int) interface{} {
	if depth <= 0 {
		return g.GenerateSimpleObject()
	}

	obj := map[string]interface{}{
		"id":        g.rand.Int63(),
		"timestamp": time.Now().Unix(),
		"data":      g.GenerateSimpleObject(),
		"children":  make([]interface{}, g.rand.Intn(5)+1),
	}

	// 递归生成子对象
	children := obj["children"].([]interface{})
	for i := range children {
		children[i] = g.GenerateComplexObject(depth - 1)
	}

	return obj
}

// GenerateLargeArray 生成大数组
// GenerateLargeArray generates a large array for performance testing
func (g *TestDataGenerator) GenerateLargeArray(size int) []interface{} {
	arr := make([]interface{}, size)
	for i := 0; i < size; i++ {
		arr[i] = g.GenerateSimpleObject()
	}
	return arr
}

// GenerateJSONString 生成JSON字符串
// GenerateJSONString generates a JSON string from the given data
func (g *TestDataGenerator) GenerateJSONString(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// GenerateInvalidJSON 生成无效的JSON字符串用于错误测试
// GenerateInvalidJSON generates invalid JSON strings for error testing
func (g *TestDataGenerator) GenerateInvalidJSON() []string {
	return []string{
		`{"name":}`,                    // 缺少值
		`{"name":"test",}`,             // 多余逗号
		`{"name":"test"`,               // 缺少右括号
		`"name":"test"}`,               // 缺少左括号
		`{"name":"test", "age":}`,      // 缺少值
		`{"name":"test" "age":25}`,     // 缺少逗号
		`{"name":"test", "age":25,}`,   // 多余逗号
		`[1,2,3,]`,                     // 数组多余逗号
		`[1,2,3`,                       // 数组缺少右括号
		`1,2,3]`,                       // 数组缺少左括号
		`{"name":"test", "name":"duplicate"}`, // 重复键
	}
}

// GenerateUnicodeJSON 生成包含Unicode字符的JSON
// GenerateUnicodeJSON generates JSON with Unicode characters
func (g *TestDataGenerator) GenerateUnicodeJSON() map[string]interface{} {
	return map[string]interface{}{
		"中文名称":   "张三",
		"emoji":   "😀🎉🚀",
		"special": "\u0048\u0065\u006C\u006C\u006F", // "Hello" in Unicode
		"mixed":   "Hello 世界 🌍",
		"escape":  "Line1\nLine2\tTabbed",
	}
}

// GenerateEdgeCases 生成边界情况测试数据
// GenerateEdgeCases generates edge case test data
func (g *TestDataGenerator) GenerateEdgeCases() []interface{} {
	return []interface{}{
		nil,
		"",
		0,
		-1,
		1.7976931348623157e+308, // 最大float64
		4.9406564584124654e-324, // 最小正float64
		[]interface{}{},
		map[string]interface{}{},
		strings.Repeat("a", 10000), // 长字符串
		map[string]interface{}{
			"deeply": map[string]interface{}{
				"nested": map[string]interface{}{
					"object": map[string]interface{}{
						"value": "deep",
					},
				},
			},
		},
	}
}

// GeneratePerformanceTestData 生成性能测试数据
// GeneratePerformanceTestData generates data for performance benchmarks
func (g *TestDataGenerator) GeneratePerformanceTestData() map[string]string {
	data := make(map[string]string)
	
	// 小型JSON (< 1KB)
	smallObj := g.GenerateSimpleObject()
	smallJSON, _ := g.GenerateJSONString(smallObj)
	data["small"] = smallJSON
	
	// 中型JSON (10-100KB)
	mediumObj := map[string]interface{}{
		"users": g.GenerateLargeArray(100),
		"meta":  g.GenerateComplexObject(3),
	}
	mediumJSON, _ := g.GenerateJSONString(mediumObj)
	data["medium"] = mediumJSON
	
	// 大型JSON (> 1MB)
	largeObj := map[string]interface{}{
		"data":     g.GenerateLargeArray(1000),
		"metadata": g.GenerateComplexObject(5),
		"config":   g.GenerateComplexObject(4),
	}
	largeJSON, _ := g.GenerateJSONString(largeObj)
	data["large"] = largeJSON
	
	return data
}

// GenerateJSONPathTestData 生成JSONPath测试数据
// GenerateJSONPathTestData generates test data for JSONPath queries
func (g *TestDataGenerator) GenerateJSONPathTestData() string {
	data := map[string]interface{}{
		"store": map[string]interface{}{
			"book": []interface{}{
				map[string]interface{}{
					"category": "reference",
					"author":   "Nigel Rees",
					"title":    "Sayings of the Century",
					"price":    8.95,
				},
				map[string]interface{}{
					"category": "fiction",
					"author":   "Evelyn Waugh",
					"title":    "Sword of Honour",
					"price":    12.99,
				},
				map[string]interface{}{
					"category": "fiction",
					"author":   "Herman Melville",
					"title":    "Moby Dick",
					"isbn":     "0-553-21311-3",
					"price":    8.99,
				},
			},
			"bicycle": map[string]interface{}{
				"color": "red",
				"price": 19.95,
			},
		},
		"users": []interface{}{
			map[string]interface{}{
				"name": "张三",
				"age":  25,
				"city": "北京",
			},
			map[string]interface{}{
				"name": "李四",
				"age":  30,
				"city": "上海",
			},
			map[string]interface{}{
				"name": "王五",
				"age":  35,
				"city": "广州",
			},
		},
	}
	
	jsonStr, _ := g.GenerateJSONString(data)
	return jsonStr
}