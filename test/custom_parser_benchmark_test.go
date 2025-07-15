package test

import (
	"encoding/json"
	"testing"
	
	xyJson "github.com/ihuem/xyJson"
)

// TestStruct 用于基准测试的结构体
type TestStruct struct {
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Email   string  `json:"email"`
	Active  bool    `json:"active"`
	Salary  float64 `json:"salary"`
	Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
		Zip    string `json:"zip"`
	} `json:"address"`
	Tags []string `json:"tags"`
}

// 测试数据
var customTestJSONData = []byte(`{
	"name": "John Doe",
	"age": 30,
	"email": "john.doe@example.com",
	"active": true,
	"salary": 75000.50,
	"address": {
		"street": "123 Main St",
		"city": "New York",
		"zip": "10001"
	},
	"tags": ["developer", "golang", "json"]
}`)

var customTestJSONString = string(customTestJSONData)

// BenchmarkOfficialJSON 官方json包基准测试
func BenchmarkOfficialJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result TestStruct
		err := json.Unmarshal(customTestJSONData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkXyJsonFast 快速路径基准测试
func BenchmarkXyJsonFast(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result TestStruct
		err := xyJson.UnmarshalToStructFast(customTestJSONData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkXyJsonCustom 自定义解析器基准测试
func BenchmarkXyJsonCustom(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result TestStruct
		err := xyJson.UnmarshalToStructCustom(customTestJSONData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkXyJsonCustomString 自定义解析器字符串基准测试
func BenchmarkXyJsonCustomString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result TestStruct
		err := xyJson.UnmarshalStringToStructCustom(customTestJSONString, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkXyJsonStandard 标准路径基准测试
func BenchmarkXyJsonStandard(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result TestStruct
		err := xyJson.UnmarshalToStruct(customTestJSONData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// 内存分配基准测试
func BenchmarkOfficialJSONMemory(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result TestStruct
		err := json.Unmarshal(customTestJSONData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkXyJsonFastMemory(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result TestStruct
		err := xyJson.UnmarshalToStructFast(customTestJSONData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkXyJsonCustomMemory(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result TestStruct
		err := xyJson.UnmarshalToStructCustom(customTestJSONData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// 功能测试
func TestCustomParserCorrectness(t *testing.T) {
	// 测试自定义解析器的正确性
	var customResult TestStruct
	err := xyJson.UnmarshalToStructCustom(customTestJSONData, &customResult)
	if err != nil {
		t.Fatalf("Custom parser failed: %v", err)
	}

	// 与官方json包对比
	var officialResult TestStruct
	err = json.Unmarshal(customTestJSONData, &officialResult)
	if err != nil {
		t.Fatalf("Official JSON failed: %v", err)
	}

	// 验证结果一致性
	if customResult.Name != officialResult.Name {
		t.Errorf("Name mismatch: custom=%s, official=%s", customResult.Name, officialResult.Name)
	}
	if customResult.Age != officialResult.Age {
		t.Errorf("Age mismatch: custom=%d, official=%d", customResult.Age, officialResult.Age)
	}
	if customResult.Email != officialResult.Email {
		t.Errorf("Email mismatch: custom=%s, official=%s", customResult.Email, officialResult.Email)
	}
	if customResult.Active != officialResult.Active {
		t.Errorf("Active mismatch: custom=%t, official=%t", customResult.Active, officialResult.Active)
	}
	if customResult.Salary != officialResult.Salary {
		t.Errorf("Salary mismatch: custom=%f, official=%f", customResult.Salary, officialResult.Salary)
	}
	if customResult.Address.Street != officialResult.Address.Street {
		t.Errorf("Address.Street mismatch: custom=%s, official=%s", customResult.Address.Street, officialResult.Address.Street)
	}
	if len(customResult.Tags) != len(officialResult.Tags) {
		t.Errorf("Tags length mismatch: custom=%d, official=%d", len(customResult.Tags), len(officialResult.Tags))
	}
	for i, tag := range customResult.Tags {
		if i < len(officialResult.Tags) && tag != officialResult.Tags[i] {
			t.Errorf("Tags[%d] mismatch: custom=%s, official=%s", i, tag, officialResult.Tags[i])
		}
	}
}

// 测试字符串解析
func TestCustomParserStringCorrectness(t *testing.T) {
	var customResult TestStruct
	err := xyJson.UnmarshalStringToStructCustom(customTestJSONString, &customResult)
	if err != nil {
		t.Fatalf("Custom string parser failed: %v", err)
	}

	// 验证基本字段
	if customResult.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", customResult.Name)
	}
	if customResult.Age != 30 {
		t.Errorf("Expected age 30, got %d", customResult.Age)
	}
	if !customResult.Active {
		t.Error("Expected active to be true")
	}
}
