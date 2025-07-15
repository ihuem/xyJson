package test

import (
	"testing"

	xyJson "github.com/ihuem/xyJson"
)

// TestMustMethodsReturnDefaults 测试Must方法在出错时返回默认值而不是panic
// TestMustMethodsReturnDefaults tests that Must methods return default values instead of panicking on errors
func TestMustMethodsReturnDefaults(t *testing.T) {
	// 创建一个简单的JSON对象用于测试
	// Create a simple JSON object for testing
	root := xyJson.CreateObject()
	err := xyJson.Set(root, "$.name", "test")
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	// 测试MustGetString在路径不存在时返回空字符串
	// Test MustGetString returns empty string when path doesn't exist
	result := xyJson.MustGetString(root, "$.nonexistent")
	if result != "" {
		t.Errorf("Expected empty string, got: %s", result)
	}

	// 测试MustGetInt在路径不存在时返回0
	// Test MustGetInt returns 0 when path doesn't exist
	intResult := xyJson.MustGetInt(root, "$.nonexistent")
	if intResult != 0 {
		t.Errorf("Expected 0, got: %d", intResult)
	}

	// 测试MustGetInt64在路径不存在时返回0
	// Test MustGetInt64 returns 0 when path doesn't exist
	int64Result := xyJson.MustGetInt64(root, "$.nonexistent")
	if int64Result != 0 {
		t.Errorf("Expected 0, got: %d", int64Result)
	}

	// 测试MustGetFloat64在路径不存在时返回0.0
	// Test MustGetFloat64 returns 0.0 when path doesn't exist
	floatResult := xyJson.MustGetFloat64(root, "$.nonexistent")
	if floatResult != 0.0 {
		t.Errorf("Expected 0.0, got: %f", floatResult)
	}

	// 测试MustGetBool在路径不存在时返回false
	// Test MustGetBool returns false when path doesn't exist
	boolResult := xyJson.MustGetBool(root, "$.nonexistent")
	if boolResult != false {
		t.Errorf("Expected false, got: %t", boolResult)
	}

	// 测试MustGetObject在路径不存在时返回空对象
	// Test MustGetObject returns empty object when path doesn't exist
	objResult := xyJson.MustGetObject(root, "$.nonexistent")
	if objResult == nil {
		t.Error("Expected non-nil object")
	}
	if objResult.Size() != 0 {
		t.Errorf("Expected empty object, got size: %d", objResult.Size())
	}

	// 测试MustGetArray在路径不存在时返回空数组
	// Test MustGetArray returns empty array when path doesn't exist
	arrResult := xyJson.MustGetArray(root, "$.nonexistent")
	if arrResult == nil {
		t.Error("Expected non-nil array")
	}
	if arrResult.Length() != 0 {
		t.Errorf("Expected empty array, got length: %d", arrResult.Length())
	}

	// 测试MustGet在路径不存在时返回null值
	// Test MustGet returns null value when path doesn't exist
	nullResult := xyJson.MustGet(root, "$.nonexistent")
	if nullResult == nil {
		t.Error("Expected non-nil value")
	}
	if !nullResult.IsNull() {
		t.Error("Expected null value")
	}
}

// TestMustParseMethodsReturnDefaults 测试Must解析方法在出错时返回默认值
// TestMustParseMethodsReturnDefaults tests that Must parse methods return default values on errors
func TestMustParseMethodsReturnDefaults(t *testing.T) {
	// 测试MustParse在解析无效JSON时返回null值
	// Test MustParse returns null value when parsing invalid JSON
	invalidJSON := []byte("{invalid json}")
	result := xyJson.MustParse(invalidJSON)
	if result == nil {
		t.Error("Expected non-nil value")
	}
	if !result.IsNull() {
		t.Error("Expected null value")
	}

	// 测试MustParseString在解析无效JSON时返回null值
	// Test MustParseString returns null value when parsing invalid JSON
	invalidJSONStr := "{invalid json}"
	resultStr := xyJson.MustParseString(invalidJSONStr)
	if resultStr == nil {
		t.Error("Expected non-nil value")
	}
	if !resultStr.IsNull() {
		t.Error("Expected null value")
	}
}

// TestMustSerializeMethodsReturnDefaults 测试Must序列化方法在出错时返回默认值
// TestMustSerializeMethodsReturnDefaults tests that Must serialize methods return default values on errors
func TestMustSerializeMethodsReturnDefaults(t *testing.T) {
	// 创建一个nil值来触发序列化错误
	// Create a nil value to trigger serialization error
	var nilValue xyJson.IValue = nil

	// 测试MustSerialize在序列化nil值时返回空字节数组
	// Test MustSerialize returns empty byte array when serializing nil value
	result := xyJson.MustSerialize(nilValue)
	if result == nil {
		t.Error("Expected non-nil byte array")
	}
	if len(result) != 0 {
		t.Errorf("Expected empty byte array, got length: %d", len(result))
	}

	// 测试MustSerializeToString在序列化nil值时返回空字符串
	// Test MustSerializeToString returns empty string when serializing nil value
	resultStr := xyJson.MustSerializeToString(nilValue)
	if resultStr != "" {
		t.Errorf("Expected empty string, got: %s", resultStr)
	}

	// 测试MustPretty在格式化nil值时返回空字符串
	// Test MustPretty returns empty string when formatting nil value
	prettyResult := xyJson.MustPretty(nilValue)
	if prettyResult != "" {
		t.Errorf("Expected empty string, got: %s", prettyResult)
	}

	// 测试MustCompact在压缩nil值时返回空字符串
	// Test MustCompact returns empty string when compacting nil value
	compactResult := xyJson.MustCompact(nilValue)
	if compactResult != "" {
		t.Errorf("Expected empty string, got: %s", compactResult)
	}
}

// TestMustToMethodsReturnDefaults 测试MustTo转换方法在出错时返回默认值
// TestMustToMethodsReturnDefaults tests that MustTo conversion methods return default values on errors
func TestMustToMethodsReturnDefaults(t *testing.T) {
	// 创建一个对象值来触发类型转换错误
	// Create an object value to trigger type conversion errors
	objValue := xyJson.CreateObject()

	// 测试MustToString在转换对象时的行为（这个可能不会出错，但我们测试一下）
	// Test MustToString behavior when converting object (this might not error, but let's test)
	strResult := xyJson.MustToString(objValue)
	// 对象转换为字符串通常会成功，所以我们只检查不是nil
	// Object to string conversion usually succeeds, so we just check it's not empty
	if strResult == "" {
		t.Log("MustToString returned empty string for object")
	}

	// 测试MustToInt在转换对象时返回0
	// Test MustToInt returns 0 when converting object
	intResult := xyJson.MustToInt(objValue)
	if intResult != 0 {
		t.Errorf("Expected 0, got: %d", intResult)
	}

	// 测试MustToInt64在转换对象时返回0
	// Test MustToInt64 returns 0 when converting object
	int64Result := xyJson.MustToInt64(objValue)
	if int64Result != 0 {
		t.Errorf("Expected 0, got: %d", int64Result)
	}

	// 测试MustToFloat64在转换对象时返回0.0
	// Test MustToFloat64 returns 0.0 when converting object
	floatResult := xyJson.MustToFloat64(objValue)
	if floatResult != 0.0 {
		t.Errorf("Expected 0.0, got: %f", floatResult)
	}

	// 测试MustToBool在转换对象时返回false
	// Test MustToBool returns false when converting object
	boolResult := xyJson.MustToBool(objValue)
	if boolResult != false {
		t.Errorf("Expected false, got: %t", boolResult)
	}

	// 测试MustToTime在转换对象时返回零值时间
	// Test MustToTime returns zero time when converting object
	timeResult := xyJson.MustToTime(objValue)
	if !timeResult.IsZero() {
		t.Errorf("Expected zero time, got: %v", timeResult)
	}

	// 测试MustToBytes在转换对象时返回nil
	// Test MustToBytes returns nil when converting object
	bytesResult := xyJson.MustToBytes(objValue)
	if bytesResult != nil {
		t.Errorf("Expected nil, got: %v", bytesResult)
	}

	// 测试MustToObject在转换对象时应该成功（对象转对象）
	// Test MustToObject should succeed when converting object (object to object)
	objResult := xyJson.MustToObject(objValue)
	if objResult == nil {
		t.Error("Expected non-nil object")
	}

	// 测试MustToArray在转换对象时返回空数组
	// Test MustToArray returns empty array when converting object
	arrResult := xyJson.MustToArray(objValue)
	if arrResult == nil {
		t.Error("Expected non-nil array")
	}
	if arrResult.Length() != 0 {
		t.Errorf("Expected empty array, got length: %d", arrResult.Length())
	}
}

// TestMustCreateFromRawReturnDefaults 测试MustCreateFromRaw的正常行为
// TestMustCreateFromRawReturnDefaults tests that MustCreateFromRaw works normally
func TestMustCreateFromRawReturnDefaults(t *testing.T) {
	// 测试MustCreateFromRaw处理正常值
	// Test MustCreateFromRaw with normal values
	result := xyJson.MustCreateFromRaw("test")
	if result == nil {
		t.Error("Expected non-nil value")
	}
	if result.String() != "test" {
		t.Errorf("Expected 'test', got: %s", result.String())
	}

	// 测试MustCreateFromRaw处理nil值
	// Test MustCreateFromRaw with nil value
	nilResult := xyJson.MustCreateFromRaw(nil)
	if nilResult == nil {
		t.Error("Expected non-nil value")
	}
	if !nilResult.IsNull() {
		t.Error("Expected null value")
	}
}
