package test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	. "github.com/ihuem/xyJson"
)

// 测试用的结构体定义
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	Zip    string `json:"zip"`
}

type User struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Active  bool     `json:"active"`
	Balance float64  `json:"balance"`
	Address Address  `json:"address"`
	Tags    []string `json:"tags"`
	Scores  []int    `json:"scores"`
}

type UserWithPointers struct {
	ID      *int     `json:"id"`
	Name    *string  `json:"name"`
	Active  *bool    `json:"active"`
	Balance *float64 `json:"balance"`
	Address *Address `json:"address"`
}

type UserWithTags struct {
	ID       int    `json:"user_id"`
	Name     string `json:"full_name"`
	Email    string `json:"-"`                  // 忽略字段
	Password string `json:"password,omitempty"` // 自定义字段名和空值忽略
	Active   bool   `json:"is_active"`
}

type UserWithTime struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserWithMap struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	Metadata map[string]string `json:"metadata"`
	Settings map[string]int    `json:"settings"`
}

// TestSerializeToStruct_BasicTypes 测试基础类型序列化
func TestSerializeToStruct_BasicTypes(t *testing.T) {
	jsonData := `{"name":"Alice","age":25}`
	value, err := ParseString(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	var person Person
	err = SerializeToStruct(value, &person)
	if err != nil {
		t.Fatalf("Failed to serialize to struct: %v", err)
	}

	if person.Name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", person.Name)
	}
	if person.Age != 25 {
		t.Errorf("Expected age 25, got %d", person.Age)
	}
}

// TestSerializeToStruct_ComplexStruct 测试复杂结构体序列化
func TestSerializeToStruct_ComplexStruct(t *testing.T) {
	jsonData := `{
		"id": 123,
		"name": "John Doe",
		"email": "john@example.com",
		"active": true,
		"balance": 1234.56,
		"address": {
			"street": "123 Main St",
			"city": "New York",
			"zip": "10001"
		},
		"tags": ["admin", "user"],
		"scores": [85, 92, 78]
	}`

	value, err := ParseString(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	var user User
	err = SerializeToStruct(value, &user)
	if err != nil {
		t.Fatalf("Failed to serialize to struct: %v", err)
	}

	// 验证基础字段
	if user.ID != 123 {
		t.Errorf("Expected ID 123, got %d", user.ID)
	}
	if user.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", user.Email)
	}
	if !user.Active {
		t.Error("Expected active to be true")
	}
	if user.Balance != 1234.56 {
		t.Errorf("Expected balance 1234.56, got %f", user.Balance)
	}

	// 验证嵌套结构体
	if user.Address.Street != "123 Main St" {
		t.Errorf("Expected street '123 Main St', got '%s'", user.Address.Street)
	}
	if user.Address.City != "New York" {
		t.Errorf("Expected city 'New York', got '%s'", user.Address.City)
	}
	if user.Address.Zip != "10001" {
		t.Errorf("Expected zip '10001', got '%s'", user.Address.Zip)
	}

	// 验证切片
	expectedTags := []string{"admin", "user"}
	if !reflect.DeepEqual(user.Tags, expectedTags) {
		t.Errorf("Expected tags %v, got %v", expectedTags, user.Tags)
	}

	expectedScores := []int{85, 92, 78}
	if !reflect.DeepEqual(user.Scores, expectedScores) {
		t.Errorf("Expected scores %v, got %v", expectedScores, user.Scores)
	}
}

// TestSerializeToStruct_WithPointers 测试指针类型序列化
func TestSerializeToStruct_WithPointers(t *testing.T) {
	jsonData := `{
		"id": 456,
		"name": "Jane Smith",
		"active": false,
		"balance": 789.12,
		"address": {
			"street": "456 Oak Ave",
			"city": "Boston",
			"zip": "02101"
		}
	}`

	value, err := ParseString(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	var user UserWithPointers
	err = SerializeToStruct(value, &user)
	if err != nil {
		t.Fatalf("Failed to serialize to struct: %v", err)
	}

	// 验证指针字段
	if user.ID == nil || *user.ID != 456 {
		t.Errorf("Expected ID 456, got %v", user.ID)
	}
	if user.Name == nil || *user.Name != "Jane Smith" {
		t.Errorf("Expected name 'Jane Smith', got %v", user.Name)
	}
	if user.Active == nil || *user.Active != false {
		t.Errorf("Expected active false, got %v", user.Active)
	}
	if user.Balance == nil || *user.Balance != 789.12 {
		t.Errorf("Expected balance 789.12, got %v", user.Balance)
	}
	if user.Address == nil {
		t.Error("Expected address to be non-nil")
	} else {
		if user.Address.Street != "456 Oak Ave" {
			t.Errorf("Expected street '456 Oak Ave', got '%s'", user.Address.Street)
		}
	}
}

// TestSerializeToStruct_JSONTags 测试JSON标签处理
func TestSerializeToStruct_JSONTags(t *testing.T) {
	jsonData := `{
		"user_id": 789,
		"full_name": "Bob Wilson",
		"email": "bob@example.com",
		"password": "secret123",
		"is_active": true
	}`

	value, err := ParseString(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	var user UserWithTags
	err = SerializeToStruct(value, &user)
	if err != nil {
		t.Fatalf("Failed to serialize to struct: %v", err)
	}

	// 验证标签映射
	if user.ID != 789 {
		t.Errorf("Expected ID 789, got %d", user.ID)
	}
	if user.Name != "Bob Wilson" {
		t.Errorf("Expected name 'Bob Wilson', got '%s'", user.Name)
	}
	if user.Active != true {
		t.Errorf("Expected active true, got %t", user.Active)
	}
	// Email字段应该被忽略，保持零值
	if user.Email != "" {
		t.Errorf("Expected email to be empty (ignored), got '%s'", user.Email)
	}
	// Password字段应该被设置
	if user.Password != "secret123" {
		t.Errorf("Expected password 'secret123', got '%s'", user.Password)
	}
}

// TestSerializeToStruct_TimeFields 测试时间字段序列化
func TestSerializeToStruct_TimeFields(t *testing.T) {
	jsonData := `{
		"id": 999,
		"name": "Time User",
		"created_at": "2023-01-15T10:30:00Z",
		"updated_at": "2023-12-25T15:45:30Z"
	}`

	value, err := ParseString(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	var user UserWithTime
	err = SerializeToStruct(value, &user)
	if err != nil {
		t.Fatalf("Failed to serialize to struct: %v", err)
	}

	// 验证时间字段
	expectedCreated, _ := time.Parse(time.RFC3339, "2023-01-15T10:30:00Z")
	expectedUpdated, _ := time.Parse(time.RFC3339, "2023-12-25T15:45:30Z")

	if !user.CreatedAt.Equal(expectedCreated) {
		t.Errorf("Expected created_at %v, got %v", expectedCreated, user.CreatedAt)
	}
	if !user.UpdatedAt.Equal(expectedUpdated) {
		t.Errorf("Expected updated_at %v, got %v", expectedUpdated, user.UpdatedAt)
	}
}

// TestSerializeToStruct_MapFields 测试Map字段序列化
func TestSerializeToStruct_MapFields(t *testing.T) {
	jsonData := `{
		"id": 111,
		"name": "Map User",
		"metadata": {
			"department": "engineering",
			"role": "developer",
			"location": "remote"
		},
		"settings": {
			"theme": 1,
			"notifications": 0,
			"timeout": 300
		}
	}`

	value, err := ParseString(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	var user UserWithMap
	err = SerializeToStruct(value, &user)
	if err != nil {
		t.Fatalf("Failed to serialize to struct: %v", err)
	}

	// 验证Map字段
	expectedMetadata := map[string]string{
		"department": "engineering",
		"role":       "developer",
		"location":   "remote",
	}
	if !reflect.DeepEqual(user.Metadata, expectedMetadata) {
		t.Errorf("Expected metadata %v, got %v", expectedMetadata, user.Metadata)
	}

	expectedSettings := map[string]int{
		"theme":         1,
		"notifications": 0,
		"timeout":       300,
	}
	if !reflect.DeepEqual(user.Settings, expectedSettings) {
		t.Errorf("Expected settings %v, got %v", expectedSettings, user.Settings)
	}
}

// TestSerializeToStruct_ErrorCases 测试错误情况
func TestSerializeToStruct_ErrorCases(t *testing.T) {
	value, _ := ParseString(`{"name":"test"}`)

	// 测试nil target
	err := SerializeToStruct(value, nil)
	if err == nil {
		t.Error("Expected error for nil target")
	}

	// 测试非指针target
	var person Person
	err = SerializeToStruct(value, person)
	if err == nil {
		t.Error("Expected error for non-pointer target")
	}

	// 测试指向非结构体的指针
	var str string
	err = SerializeToStruct(value, &str)
	if err == nil {
		t.Error("Expected error for pointer to non-struct")
	}

	// 测试nil value
	err = SerializeToStruct(nil, &person)
	if err == nil {
		t.Error("Expected error for nil value")
	}
}

// TestMustSerializeToStruct 测试Must版本函数
func TestMustSerializeToStruct(t *testing.T) {
	jsonData := `{"name":"Test User","age":30}`
	value, _ := ParseString(jsonData)

	var person Person
	// 这应该不会panic
	MustSerializeToStruct(value, &person)

	if person.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", person.Name)
	}
	if person.Age != 30 {
		t.Errorf("Expected age 30, got %d", person.Age)
	}
}

// TestUnmarshalToStruct 测试直接解析到结构体
func TestUnmarshalToStruct(t *testing.T) {
	jsonData := []byte(`{"name":"Direct User","age":35}`)

	var person Person
	err := UnmarshalToStruct(jsonData, &person)
	if err != nil {
		t.Fatalf("Failed to unmarshal to struct: %v", err)
	}

	if person.Name != "Direct User" {
		t.Errorf("Expected name 'Direct User', got '%s'", person.Name)
	}
	if person.Age != 35 {
		t.Errorf("Expected age 35, got %d", person.Age)
	}
}

// TestUnmarshalStringToStruct 测试从字符串解析到结构体
func TestUnmarshalStringToStruct(t *testing.T) {
	jsonData := `{"name":"String User","age":40}`

	var person Person
	err := UnmarshalStringToStruct(jsonData, &person)
	if err != nil {
		t.Fatalf("Failed to unmarshal string to struct: %v", err)
	}

	if person.Name != "String User" {
		t.Errorf("Expected name 'String User', got '%s'", person.Name)
	}
	if person.Age != 40 {
		t.Errorf("Expected age 40, got %d", person.Age)
	}
}

// BenchmarkSerializeToStruct 性能基准测试
func BenchmarkSerializeToStruct(b *testing.B) {
	jsonData := `{
		"id": 123,
		"name": "Benchmark User",
		"email": "bench@example.com",
		"active": true,
		"balance": 1000.50,
		"address": {
			"street": "123 Bench St",
			"city": "Test City",
			"zip": "12345"
		},
		"tags": ["test", "benchmark"],
		"scores": [90, 85, 95]
	}`

	value, _ := ParseString(jsonData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user User
		_ = SerializeToStruct(value, &user)
	}
}

// BenchmarkUnmarshalToStruct 解析+序列化性能基准测试
func BenchmarkUnmarshalToStruct(b *testing.B) {
	jsonData := []byte(`{
		"id": 123,
		"name": "Benchmark User",
		"email": "bench@example.com",
		"active": true,
		"balance": 1000.50,
		"address": {
			"street": "123 Bench St",
			"city": "Test City",
			"zip": "12345"
		},
		"tags": ["test", "benchmark"],
		"scores": [90, 85, 95]
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user User
		_ = UnmarshalToStruct(jsonData, &user)
	}
}

// TestCompareWithStandardJSON 与官方json包的对比测试
func TestCompareWithStandardJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "基础类型",
			jsonData: `{"name":"Alice","age":25}`,
		},
		{
			name: "复杂结构体",
			jsonData: `{
				"id": 123,
				"name": "John Doe",
				"email": "john@example.com",
				"active": true,
				"balance": 1234.56,
				"address": {
					"street": "123 Main St",
					"city": "New York",
					"zip": "10001"
				},
				"tags": ["admin", "user"],
				"scores": [85, 92, 78]
			}`,
		},
		{
			name: "带时间字段",
			jsonData: `{
				"id": 999,
				"name": "Time User",
				"created_at": "2023-01-15T10:30:00Z",
				"updated_at": "2023-12-25T15:45:30Z"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 根据测试用例选择合适的结构体类型
			switch tt.name {
			case "基础类型":
				var xyJsonResult, stdJsonResult Person

				// 使用xyJson解析
				err := UnmarshalStringToStruct(tt.jsonData, &xyJsonResult)
				if err != nil {
					t.Fatalf("xyJson解析失败: %v", err)
				}

				// 使用官方json包解析
				err = json.Unmarshal([]byte(tt.jsonData), &stdJsonResult)
				if err != nil {
					t.Fatalf("官方json包解析失败: %v", err)
				}

				// 比较结果
				if !reflect.DeepEqual(xyJsonResult, stdJsonResult) {
					t.Errorf("解析结果不一致:\nxyJson: %+v\n官方json: %+v", xyJsonResult, stdJsonResult)
				}

			case "复杂结构体":
				var xyJsonResult, stdJsonResult User

				// 使用xyJson解析
				err := UnmarshalStringToStruct(tt.jsonData, &xyJsonResult)
				if err != nil {
					t.Fatalf("xyJson解析失败: %v", err)
				}

				// 使用官方json包解析
				err = json.Unmarshal([]byte(tt.jsonData), &stdJsonResult)
				if err != nil {
					t.Fatalf("官方json包解析失败: %v", err)
				}

				// 比较结果
				if !reflect.DeepEqual(xyJsonResult, stdJsonResult) {
					t.Errorf("解析结果不一致:\nxyJson: %+v\n官方json: %+v", xyJsonResult, stdJsonResult)
				}

			case "带时间字段":
				var xyJsonResult, stdJsonResult UserWithTime

				// 使用xyJson解析
				err := UnmarshalStringToStruct(tt.jsonData, &xyJsonResult)
				if err != nil {
					t.Fatalf("xyJson解析失败: %v", err)
				}

				// 使用官方json包解析
				err = json.Unmarshal([]byte(tt.jsonData), &stdJsonResult)
				if err != nil {
					t.Fatalf("官方json包解析失败: %v", err)
				}

				// 比较结果
				if !reflect.DeepEqual(xyJsonResult, stdJsonResult) {
					t.Errorf("解析结果不一致:\nxyJson: %+v\n官方json: %+v", xyJsonResult, stdJsonResult)
				}
			}
		})
	}
}

// BenchmarkCompareWithStandardJSON 与官方json包的性能对比基准测试
func BenchmarkCompareWithStandardJSON(b *testing.B) {
	jsonData := []byte(`{
		"id": 123,
		"name": "Benchmark User",
		"email": "bench@example.com",
		"active": true,
		"balance": 1000.50,
		"address": {
			"street": "123 Bench St",
			"city": "Test City",
			"zip": "12345"
		},
		"tags": ["test", "benchmark"],
		"scores": [90, 85, 95]
	}`)

	b.Run("xyJson", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var user User
			_ = UnmarshalToStruct(jsonData, &user)
		}
	})

	b.Run("官方json包", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var user User
			_ = json.Unmarshal(jsonData, &user)
		}
	})
}

// BenchmarkCompareSerializationOnly 仅序列化部分的性能对比（不包括解析）
func BenchmarkCompareSerializationOnly(b *testing.B) {
	jsonData := `{
		"id": 123,
		"name": "Benchmark User",
		"email": "bench@example.com",
		"active": true,
		"balance": 1000.50,
		"address": {
			"street": "123 Bench St",
			"city": "Test City",
			"zip": "12345"
		},
		"tags": ["test", "benchmark"],
		"scores": [90, 85, 95]
	}`

	// 预解析数据
	value, _ := ParseString(jsonData)

	b.Run("xyJson序列化", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var user User
			_ = SerializeToStruct(value, &user)
		}
	})

	// 注意：官方json包没有直接从已解析的值序列化到struct的功能
	// 这里我们比较的是完整的解析+反序列化过程
	b.Run("官方json包完整过程", func(b *testing.B) {
		jsonBytes := []byte(jsonData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var user User
			_ = json.Unmarshal(jsonBytes, &user)
		}
	})
}

// BenchmarkCompareFastPath 对比快速路径的性能
func BenchmarkCompareFastPath(b *testing.B) {
	testData := `{"name":"Alice","age":25,"email":"alice@example.com","active":true}`
	type Person struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Email  string `json:"email"`
		Active bool   `json:"active"`
	}

	data := []byte(testData)

	b.Run("xyJson原始方法", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var person Person
			err := UnmarshalToStruct(data, &person)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("xyJson快速路径", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var person Person
			err := UnmarshalToStructFast(data, &person)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("官方json包", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var person Person
			err := json.Unmarshal(data, &person)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkUnmarshalToStructFast 测试UnmarshalToStructFast的性能
func BenchmarkUnmarshalToStructFast(b *testing.B) {
	jsonData := []byte(`{
		"id": 123,
		"name": "Fast User",
		"email": "fast@example.com",
		"active": true,
		"balance": 1000.50,
		"address": {
			"street": "123 Fast St",
			"city": "Speed City",
			"zip": "12345"
		},
		"tags": ["fast", "benchmark"],
		"scores": [95, 88, 92]
	}`)

	b.Run("xyJson快速路径", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var user User
			_ = UnmarshalToStructFast(jsonData, &user)
		}
	})

	b.Run("xyJson原始方法", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var user User
			_ = UnmarshalToStruct(jsonData, &user)
		}
	})

	b.Run("官方json包", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var user User
			_ = json.Unmarshal(jsonData, &user)
		}
	})
}
