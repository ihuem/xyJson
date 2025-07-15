package test

import (
	"testing"

	xyJson "github.com/ihuem/xyJson"
)

// TestGetWithDefaultMethods 测试GetXXXWithDefault系列方法
func TestGetWithDefaultMethods(t *testing.T) {
	// 准备测试数据
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
		}
	}`

	root, err := xyJson.ParseString(data)
	if err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	// 测试GetStringWithDefault - 存在的值
	name := xyJson.GetStringWithDefault(root, "$.user.name", "Unknown")
	if name != "Alice" {
		t.Errorf("期望name为'Alice'，实际为'%s'", name)
	}

	// 测试GetStringWithDefault - 不存在的值
	city := xyJson.GetStringWithDefault(root, "$.user.city", "Unknown")
	if city != "Unknown" {
		t.Errorf("期望city为'Unknown'，实际为'%s'", city)
	}

	// 测试GetIntWithDefault - 存在的值
	age := xyJson.GetIntWithDefault(root, "$.user.age", 0)
	if age != 30 {
		t.Errorf("期望age为30，实际为%d", age)
	}

	// 测试GetIntWithDefault - 不存在的值
	score := xyJson.GetIntWithDefault(root, "$.user.score", 100)
	if score != 100 {
		t.Errorf("期望score为100，实际为%d", score)
	}

	// 测试GetFloat64WithDefault - 存在的值
	height := xyJson.GetFloat64WithDefault(root, "$.user.height", 0.0)
	if height != 165.5 {
		t.Errorf("期望height为165.5，实际为%f", height)
	}

	// 测试GetFloat64WithDefault - 不存在的值
	weight := xyJson.GetFloat64WithDefault(root, "$.user.weight", 60.0)
	if weight != 60.0 {
		t.Errorf("期望weight为60.0，实际为%f", weight)
	}

	// 测试GetBoolWithDefault - 存在的值
	active := xyJson.GetBoolWithDefault(root, "$.user.active", false)
	if !active {
		t.Error("期望active为true")
	}

	// 测试GetBoolWithDefault - 不存在的值
	verified := xyJson.GetBoolWithDefault(root, "$.user.verified", false)
	if verified {
		t.Error("期望verified为false")
	}

	// 测试GetObjectWithDefault - 存在的值
	profile := xyJson.GetObjectWithDefault(root, "$.user.profile", nil)
	if profile == nil {
		t.Error("profile对象不应该为nil")
	} else if profile.Size() != 1 {
		t.Errorf("期望profile对象有1个属性，实际有%d个", profile.Size())
	}

	// 测试GetObjectWithDefault - 不存在的值
	settings := xyJson.GetObjectWithDefault(root, "$.user.settings", nil)
	if settings == nil {
		t.Error("settings不应该为nil，应该返回空对象")
	} else if settings.Size() != 0 {
		t.Errorf("settings应该是空对象，实际有%d个属性", settings.Size())
	}

	// 测试GetArrayWithDefault - 存在的值
	hobbies := xyJson.GetArrayWithDefault(root, "$.user.hobbies", nil)
	if hobbies == nil {
		t.Error("hobbies数组不应该为nil")
	} else if hobbies.Length() != 2 {
		t.Errorf("期望hobbies数组有2个元素，实际有%d个", hobbies.Length())
	}

	// 测试GetArrayWithDefault - 不存在的值
	skills := xyJson.GetArrayWithDefault(root, "$.user.skills", nil)
	if skills == nil {
		t.Error("skills不应该为nil，应该返回空数组")
	} else if skills.Length() != 0 {
		t.Errorf("skills应该是空数组，实际有%d个元素", skills.Length())
	}
}

// TestGetWithDefaultTypeConversion 测试类型转换失败时的默认值返回
func TestGetWithDefaultTypeConversion(t *testing.T) {
	data := `{"user": {"name": "Alice", "age": "thirty"}}`
	root, err := xyJson.ParseString(data)
	if err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	// 测试类型转换失败时返回默认值
	age := xyJson.GetIntWithDefault(root, "$.user.age", 25)
	if age != 25 {
		t.Errorf("类型转换失败时应该返回默认值25，实际为%d", age)
	}

	// 测试将字符串转换为布尔值失败时的默认值
	active := xyJson.GetBoolWithDefault(root, "$.user.name", false)
	// 注意：根据之前的分析，字符串"Alice"会被转换为true，所以这里应该是true
	if !active {
		t.Error("字符串'Alice'应该被转换为true")
	}

	// 测试不存在路径的默认值
	height := xyJson.GetFloat64WithDefault(root, "$.user.height", 170.0)
	if height != 170.0 {
		t.Errorf("不存在路径应该返回默认值170.0，实际为%f", height)
	}
}

// TestGetWithDefaultComplexTypes 测试复杂类型的默认值
func TestGetWithDefaultComplexTypes(t *testing.T) {
	data := `{"user": {"name": "Alice"}}`
	root, err := xyJson.ParseString(data)
	if err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	// 创建默认对象
	defaultObj := xyJson.CreateObject()
	defaultObj.Set("default", xyJson.CreateString("value"))

	// 测试对象默认值
	profile := xyJson.GetObjectWithDefault(root, "$.user.profile", defaultObj)
	if profile == nil {
		t.Error("应该返回默认对象")
	} else if profile.Size() != 1 {
		t.Errorf("默认对象应该有1个属性，实际有%d个", profile.Size())
	}

	// 创建默认数组
	defaultArr := xyJson.CreateArray()
	defaultArr.Append(xyJson.CreateString("default1"))
	defaultArr.Append(xyJson.CreateString("default2"))

	// 测试数组默认值
	hobbies := xyJson.GetArrayWithDefault(root, "$.user.hobbies", defaultArr)
	if hobbies == nil {
		t.Error("应该返回默认数组")
	} else if hobbies.Length() != 2 {
		t.Errorf("默认数组应该有2个元素，实际有%d个", hobbies.Length())
	}
}

// TestGetWithDefaultUsageComparison 对比不同方法的使用
func TestGetWithDefaultUsageComparison(t *testing.T) {
	data := `{"config": {"timeout": 30, "debug": true}}`
	root, err := xyJson.ParseString(data)
	if err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	// 方法1：使用Get方法（需要错误处理）
	timeout1, err1 := xyJson.GetInt(root, "$.config.timeout")
	if err1 != nil {
		timeout1 = 60 // 默认值
	}

	// 方法2：使用TryGet方法（需要判断布尔值）
	timeout2 := 60 // 默认值
	if t, ok := xyJson.TryGetInt(root, "$.config.timeout"); ok {
		timeout2 = t
	}

	// 方法3：使用GetWithDefault方法（最简洁）
	timeout3 := xyJson.GetIntWithDefault(root, "$.config.timeout", 60)

	// 验证结果一致
	if timeout1 != timeout2 || timeout2 != timeout3 {
		t.Errorf("三种方法结果不一致: %d, %d, %d", timeout1, timeout2, timeout3)
	}

	// 测试不存在的配置项
	maxRetries1, err := xyJson.GetInt(root, "$.config.maxRetries")
	if err != nil {
		maxRetries1 = 3 // 默认值
	}

	maxRetries2 := 3 // 默认值
	if r, ok := xyJson.TryGetInt(root, "$.config.maxRetries"); ok {
		maxRetries2 = r
	}

	maxRetries3 := xyJson.GetIntWithDefault(root, "$.config.maxRetries", 3)

	// 验证默认值处理
	if maxRetries1 != 3 || maxRetries2 != 3 || maxRetries3 != 3 {
		t.Errorf("默认值处理不正确: %d, %d, %d", maxRetries1, maxRetries2, maxRetries3)
	}

	// 验证GetWithDefault方法的简洁性
	// 这种方式最简洁，只需要一行代码
	host := xyJson.GetStringWithDefault(root, "$.config.host", "localhost")
	port := xyJson.GetIntWithDefault(root, "$.config.port", 8080)
	debug := xyJson.GetBoolWithDefault(root, "$.config.debug", false)

	if host != "localhost" {
		t.Errorf("期望host为'localhost'，实际为'%s'", host)
	}
	if port != 8080 {
		t.Errorf("期望port为8080，实际为%d", port)
	}
	if !debug {
		t.Error("期望debug为true")
	}
}
