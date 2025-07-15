package test

import (
	"testing"
	xyJson "github/ihuem/xyJson"
)

// TestTryGetMethods 测试TryGetXXX系列方法
func TestTryGetMethods(t *testing.T) {
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

	// 测试TryGetString - 成功情况
	if name, ok := xyJson.TryGetString(root, "$.user.name"); !ok {
		t.Error("TryGetString应该成功获取name")
	} else if name != "Alice" {
		t.Errorf("期望name为'Alice'，实际为'%s'", name)
	}

	// 测试TryGetString - 失败情况
	if _, ok := xyJson.TryGetString(root, "$.user.nonexistent"); ok {
		t.Error("TryGetString对不存在的路径应该返回false")
	}

	// 测试TryGetInt - 成功情况
	if age, ok := xyJson.TryGetInt(root, "$.user.age"); !ok {
		t.Error("TryGetInt应该成功获取age")
	} else if age != 30 {
		t.Errorf("期望age为30，实际为%d", age)
	}

	// 测试TryGetInt - 失败情况
	if _, ok := xyJson.TryGetInt(root, "$.user.name"); ok {
		t.Error("TryGetInt对字符串类型应该返回false")
	}

	// 测试TryGetFloat64 - 成功情况
	if height, ok := xyJson.TryGetFloat64(root, "$.user.height"); !ok {
		t.Error("TryGetFloat64应该成功获取height")
	} else if height != 165.5 {
		t.Errorf("期望height为165.5，实际为%f", height)
	}

	// 测试TryGetFloat64 - 失败情况
	if _, ok := xyJson.TryGetFloat64(root, "$.user.name"); ok {
		t.Error("TryGetFloat64对字符串类型应该返回false")
	}

	// 测试TryGetBool - 成功情况
	if active, ok := xyJson.TryGetBool(root, "$.user.active"); !ok {
		t.Error("TryGetBool应该成功获取active")
	} else if !active {
		t.Error("期望active为true")
	}

	// 测试TryGetBool - 失败情况（使用不存在的路径）
	if _, ok := xyJson.TryGetBool(root, "$.user.nonexistent"); ok {
		t.Error("TryGetBool对不存在的路径应该返回false")
	}

	// 测试TryGetObject - 成功情况
	if profile, ok := xyJson.TryGetObject(root, "$.user.profile"); !ok {
		t.Error("TryGetObject应该成功获取profile")
	} else if profile == nil {
		t.Error("profile对象不应该为nil")
	} else if profile.Size() != 1 {
		t.Errorf("期望profile对象有1个属性，实际有%d个", profile.Size())
	}

	// 测试TryGetObject - 失败情况
	if _, ok := xyJson.TryGetObject(root, "$.user.name"); ok {
		t.Error("TryGetObject对字符串类型应该返回false")
	}

	// 测试TryGetArray - 成功情况
	if hobbies, ok := xyJson.TryGetArray(root, "$.user.hobbies"); !ok {
		t.Error("TryGetArray应该成功获取hobbies")
	} else if hobbies == nil {
		t.Error("hobbies数组不应该为nil")
	} else if hobbies.Length() != 2 {
		t.Errorf("期望hobbies数组有2个元素，实际有%d个", hobbies.Length())
	}

	// 测试TryGetArray - 失败情况
	if _, ok := xyJson.TryGetArray(root, "$.user.name"); ok {
		t.Error("TryGetArray对字符串类型应该返回false")
	}
}

// TestTryGetMethodsZeroValues 测试TryGetXXX方法在失败时返回零值
func TestTryGetMethodsZeroValues(t *testing.T) {
	data := `{"test": "value"}`
	root, err := xyJson.ParseString(data)
	if err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	// 测试失败时的零值返回
	if str, ok := xyJson.TryGetString(root, "$.nonexistent"); ok || str != "" {
		t.Errorf("TryGetString失败时应该返回空字符串和false，实际返回'%s'和%v", str, ok)
	}

	if num, ok := xyJson.TryGetInt(root, "$.nonexistent"); ok || num != 0 {
		t.Errorf("TryGetInt失败时应该返回0和false，实际返回%d和%v", num, ok)
	}

	if num64, ok := xyJson.TryGetInt64(root, "$.nonexistent"); ok || num64 != 0 {
		t.Errorf("TryGetInt64失败时应该返回0和false，实际返回%d和%v", num64, ok)
	}

	if float, ok := xyJson.TryGetFloat64(root, "$.nonexistent"); ok || float != 0.0 {
		t.Errorf("TryGetFloat64失败时应该返回0.0和false，实际返回%f和%v", float, ok)
	}

	if boolean, ok := xyJson.TryGetBool(root, "$.nonexistent"); ok || boolean != false {
		t.Errorf("TryGetBool失败时应该返回false和false，实际返回%v和%v", boolean, ok)
	}

	if obj, ok := xyJson.TryGetObject(root, "$.nonexistent"); ok || obj != nil {
		t.Errorf("TryGetObject失败时应该返回nil和false，实际返回%v和%v", obj, ok)
	}

	if arr, ok := xyJson.TryGetArray(root, "$.nonexistent"); ok || arr != nil {
		t.Errorf("TryGetArray失败时应该返回nil和false，实际返回%v和%v", arr, ok)
	}
}

// TestTryGetVsGetComparison 对比TryGet和Get方法的使用
func TestTryGetVsGetComparison(t *testing.T) {
	data := `{"user": {"name": "Bob", "age": 25}}`
	root, err := xyJson.ParseString(data)
	if err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	// 使用Get方法（需要错误处理）
	name1, err1 := xyJson.GetString(root, "$.user.name")
	if err1 != nil {
		t.Errorf("GetString失败: %v", err1)
	}

	// 使用TryGet方法（更简洁）
	name2, ok := xyJson.TryGetString(root, "$.user.name")
	if !ok {
		t.Error("TryGetString失败")
	}

	// 验证结果一致
	if name1 != name2 {
		t.Errorf("Get和TryGet结果不一致: '%s' vs '%s'", name1, name2)
	}

	// 测试不存在的路径
	_, err2 := xyJson.GetString(root, "$.user.nonexistent")
	if err2 == nil {
		t.Error("GetString对不存在路径应该返回错误")
	}

	_, ok2 := xyJson.TryGetString(root, "$.user.nonexistent")
	if ok2 {
		t.Error("TryGetString对不存在路径应该返回false")
	}
}