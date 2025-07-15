package test

import (
	xyJson "github/ihuem/xyJson"
	"testing"
)

// TestConvenienceGetMethods 测试便利的Get方法
func TestConvenienceGetMethods(t *testing.T) {
	data := `{
		"user": {
			"name": "Alice",
			"age": 30,
			"height": 165.5,
			"active": true,
			"profile": {
				"email": "alice@example.com"
			},
			"hobbies": ["reading", "swimming"]
		},
		"product": {
			"price": 29.99,
			"count": 100
		}
	}`

	root, err := xyJson.ParseString(data)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// 测试GetString
	t.Run("GetString", func(t *testing.T) {
		name, err := xyJson.GetString(root, "$.user.name")
		if err != nil {
			t.Errorf("GetString failed: %v", err)
		}
		if name != "Alice" {
			t.Errorf("Expected 'Alice', got '%s'", name)
		}

		email, err := xyJson.GetString(root, "$.user.profile.email")
		if err != nil {
			t.Errorf("GetString failed for nested path: %v", err)
		}
		if email != "alice@example.com" {
			t.Errorf("Expected 'alice@example.com', got '%s'", email)
		}
	})

	// 测试GetInt
	t.Run("GetInt", func(t *testing.T) {
		age, err := xyJson.GetInt(root, "$.user.age")
		if err != nil {
			t.Errorf("GetInt failed: %v", err)
		}
		if age != 30 {
			t.Errorf("Expected 30, got %d", age)
		}

		count, err := xyJson.GetInt(root, "$.product.count")
		if err != nil {
			t.Errorf("GetInt failed: %v", err)
		}
		if count != 100 {
			t.Errorf("Expected 100, got %d", count)
		}
	})

	// 测试GetFloat64
	t.Run("GetFloat64", func(t *testing.T) {
		height, err := xyJson.GetFloat64(root, "$.user.height")
		if err != nil {
			t.Errorf("GetFloat64 failed: %v", err)
		}
		if height != 165.5 {
			t.Errorf("Expected 165.5, got %f", height)
		}

		price, err := xyJson.GetFloat64(root, "$.product.price")
		if err != nil {
			t.Errorf("GetFloat64 failed: %v", err)
		}
		if price != 29.99 {
			t.Errorf("Expected 29.99, got %f", price)
		}
	})

	// 测试GetBool
	t.Run("GetBool", func(t *testing.T) {
		active, err := xyJson.GetBool(root, "$.user.active")
		if err != nil {
			t.Errorf("GetBool failed: %v", err)
		}
		if !active {
			t.Errorf("Expected true, got %t", active)
		}
	})

	// 测试GetObject
	t.Run("GetObject", func(t *testing.T) {
		profile, err := xyJson.GetObject(root, "$.user.profile")
		if err != nil {
			t.Errorf("GetObject failed: %v", err)
		}
		if profile == nil {
			t.Error("Expected object, got nil")
		}
		if profile.Size() != 1 {
			t.Errorf("Expected object size 1, got %d", profile.Size())
		}
	})

	// 测试GetArray
	t.Run("GetArray", func(t *testing.T) {
		hobbies, err := xyJson.GetArray(root, "$.user.hobbies")
		if err != nil {
			t.Errorf("GetArray failed: %v", err)
		}
		if hobbies == nil {
			t.Error("Expected array, got nil")
		}
		if hobbies.Length() != 2 {
			t.Errorf("Expected array length 2, got %d", hobbies.Length())
		}
	})
}

// TestMustGetMethods 测试Must版本的Get方法
func TestMustGetMethods(t *testing.T) {
	data := `{
		"user": {
			"name": "Bob",
			"age": 25,
			"salary": 50000.50,
			"married": false
		}
	}`

	root, err := xyJson.ParseString(data)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// 测试MustGetString
	t.Run("MustGetString", func(t *testing.T) {
		name := xyJson.MustGetString(root, "$.user.name")
		if name != "Bob" {
			t.Errorf("Expected 'Bob', got '%s'", name)
		}
	})

	// 测试MustGetInt
	t.Run("MustGetInt", func(t *testing.T) {
		age := xyJson.MustGetInt(root, "$.user.age")
		if age != 25 {
			t.Errorf("Expected 25, got %d", age)
		}
	})

	// 测试MustGetFloat64
	t.Run("MustGetFloat64", func(t *testing.T) {
		salary := xyJson.MustGetFloat64(root, "$.user.salary")
		if salary != 50000.50 {
			t.Errorf("Expected 50000.50, got %f", salary)
		}
	})

	// 测试MustGetBool
	t.Run("MustGetBool", func(t *testing.T) {
		married := xyJson.MustGetBool(root, "$.user.married")
		if married {
			t.Errorf("Expected false, got %t", married)
		}
	})
}

// TestConvenienceAPIErrorHandling 测试便利API的错误处理
func TestConvenienceAPIErrorHandling(t *testing.T) {
	data := `{"user": {"name": "Charlie"}}`
	root, err := xyJson.ParseString(data)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// 测试路径不存在的情况
	t.Run("PathNotFound", func(t *testing.T) {
		_, err := xyJson.GetString(root, "$.user.nonexistent")
		if err == nil {
			t.Error("Expected error for non-existent path")
		}
	})

	// 测试类型转换错误
	t.Run("TypeConversionError", func(t *testing.T) {
		_, err := xyJson.GetInt(root, "$.user.name") // name是字符串，不能转换为int
		if err == nil {
			t.Error("Expected error for type conversion")
		}
	})
}

// TestConvenienceAPIComparison 对比新旧API的使用方式
func TestConvenienceAPIComparison(t *testing.T) {
	data := `{"product": {"price": 99.99}}`
	root, err := xyJson.ParseString(data)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// 旧的方式：需要类型断言
	t.Run("OldWay", func(t *testing.T) {
		priceValue, err := xyJson.Get(root, "$.product.price")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		// 需要类型断言
		scalarValue, ok := priceValue.(xyJson.IScalarValue)
		if !ok {
			t.Fatal("Failed to cast to IScalarValue")
		}

		price, err := scalarValue.Float64()
		if err != nil {
			t.Fatalf("Float64 conversion failed: %v", err)
		}

		if price != 99.99 {
			t.Errorf("Expected 99.99, got %f", price)
		}
	})

	// 新的方式：直接获取类型化的值
	t.Run("NewWay", func(t *testing.T) {
		price, err := xyJson.GetFloat64(root, "$.product.price")
		if err != nil {
			t.Fatalf("GetFloat64 failed: %v", err)
		}

		if price != 99.99 {
			t.Errorf("Expected 99.99, got %f", price)
		}
	})

	// 最简洁的方式：Must版本（适用于确信路径存在的场景）
	t.Run("MustWay", func(t *testing.T) {
		price := xyJson.MustGetFloat64(root, "$.product.price")
		if price != 99.99 {
			t.Errorf("Expected 99.99, got %f", price)
		}
	})
}
