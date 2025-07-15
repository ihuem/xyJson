// Package examples demonstrates basic usage of xyJson library
// 本包演示了 xyJson 库的基本用法
package main

import (
	"fmt"
	"log"

	"github.com/ihuem/xyJson"
)

func basicUsageExample() {
	// 基本解析示例 / Basic parsing example
	basicParsingExample()

	// 创建和操作JSON对象 / Creating and manipulating JSON objects
	objectOperationsExample()

	// 数组操作示例 / Array operations example
	arrayOperationsExample()

	// JSONPath查询示例 / JSONPath query example
	jsonPathExample()

	// 序列化示例 / Serialization example
	serializationExample()
}

// basicParsingExample 演示基本的JSON解析功能
// basicParsingExample demonstrates basic JSON parsing functionality
func basicParsingExample() {
	fmt.Println("=== Basic Parsing Example ===")

	// 解析简单的JSON字符串
	// Parse a simple JSON string
	jsonStr := `{"name":"Alice","age":30,"active":true,"score":95.5}`
	value, err := xyJson.ParseString(jsonStr)
	if err != nil {
		log.Fatal("Parse error:", err)
	}

	// 检查值类型
	// Check value type
	fmt.Printf("Value type: %v\n", value.Type())
	fmt.Printf("Is null: %v\n", value.IsNull())

	// 转换为对象并访问字段
	// Convert to object and access fields
	if obj, ok := value.(xyJson.IObject); ok {
		fmt.Printf("Name: %s\n", obj.Get("name").String())
		fmt.Printf("Age: %s\n", obj.Get("age").String())
		fmt.Printf("Active: %s\n", obj.Get("active").String())
		fmt.Printf("Score: %s\n", obj.Get("score").String())
	}

	fmt.Println()
}

// objectOperationsExample 演示JSON对象的创建和操作
// objectOperationsExample demonstrates JSON object creation and manipulation
func objectOperationsExample() {
	fmt.Println("=== Object Operations Example ===")

	// 创建新的JSON对象
	// Create a new JSON object
	obj := xyJson.CreateObject()

	// 设置各种类型的值
	// Set various types of values
	obj.Set("name", "Bob")
	obj.Set("age", 25)
	obj.Set("height", 175.5)
	obj.Set("married", false)
	obj.Set("address", xyJson.CreateNull())

	// 创建嵌套对象
	// Create nested object
	contact := xyJson.CreateObject()
	contact.Set("email", "bob@example.com")
	contact.Set("phone", "123-456-7890")
	obj.Set("contact", contact)

	// 显示对象信息
	// Display object information
	fmt.Printf("Object size: %d\n", obj.Size())
	fmt.Printf("Has 'name': %v\n", obj.Has("name"))
	fmt.Printf("Has 'salary': %v\n", obj.Has("salary"))

	// 遍历所有键值对
	// Iterate over all key-value pairs
	fmt.Println("All key-value pairs:")
	obj.Range(func(key string, value xyJson.IValue) bool {
		fmt.Printf("  %s: %s\n", key, value.String())
		return true
	})

	// 删除键
	// Delete a key
	deleted := obj.Delete("address")
	fmt.Printf("Deleted 'address': %v\n", deleted)
	fmt.Printf("New size: %d\n", obj.Size())

	fmt.Println()
}

// arrayOperationsExample 演示JSON数组的操作
// arrayOperationsExample demonstrates JSON array operations
func arrayOperationsExample() {
	fmt.Println("=== Array Operations Example ===")

	// 创建新的JSON数组
	// Create a new JSON array
	arr := xyJson.CreateArray()

	// 添加各种类型的元素
	// Add various types of elements
	arr.Append("first")
	arr.Append(42)
	arr.Append(true)
	arr.Append(3.14)

	// 创建嵌套数组
	// Create nested array
	nested := xyJson.CreateArray()
	nested.Append(1)
	nested.Append(2)
	nested.Append(3)
	arr.Append(nested)

	// 显示数组信息
	// Display array information
	fmt.Printf("Array length: %d\n", arr.Length())

	// 访问数组元素
	// Access array elements
	fmt.Println("Array elements:")
	for i := 0; i < arr.Length(); i++ {
		element := arr.Get(i)
		fmt.Printf("  [%d]: %s (type: %v)\n", i, element.String(), element.Type())
	}

	// 使用Range遍历
	// Use Range to iterate
	fmt.Println("Using Range:")
	arr.Range(func(index int, value xyJson.IValue) bool {
		fmt.Printf("  Index %d: %s\n", index, value.String())
		return true
	})

	// 在指定位置插入元素
	// Insert element at specific position
	arr.Insert(1, "inserted")
	fmt.Printf("After insertion, length: %d\n", arr.Length())
	fmt.Printf("Element at index 1: %s\n", arr.Get(1).String())

	fmt.Println()
}

// jsonPathExample 演示JSONPath查询功能
// jsonPathExample demonstrates JSONPath query functionality
func jsonPathExample() {
	fmt.Println("=== JSONPath Example ===")

	// 创建复杂的JSON结构
	// Create complex JSON structure
	jsonStr := `{
		"store": {
			"book": [
				{"title": "Go Programming", "author": "John", "price": 29.99},
				{"title": "JSON Processing", "author": "Jane", "price": 24.99},
				{"title": "Web Development", "author": "Bob", "price": 34.99}
			],
			"bicycle": {
				"color": "red",
				"price": 199.99
			}
		}
	}`

	root, err := xyJson.ParseString(jsonStr)
	if err != nil {
		log.Fatal("Parse error:", err)
	}

	// 基本路径查询
	// Basic path queries
	fmt.Println("Basic path queries:")

	// 获取第一本书的标题
	// Get the title of the first book
	title, err := xyJson.Get(root, "$.store.book[0].title")
	if err == nil {
		fmt.Printf("First book title: %s\n", title.String())
	}

	// 获取自行车颜色
	// Get bicycle color
	color, err := xyJson.Get(root, "$.store.bicycle.color")
	if err == nil {
		fmt.Printf("Bicycle color: %s\n", color.String())
	}

	// 获取所有书的价格
	// Get all book prices
	prices, err := xyJson.GetAll(root, "$.store.book[*].price")
	if err == nil {
		fmt.Println("All book prices:")
		for i, price := range prices {
			fmt.Printf("  Book %d: %s\n", i+1, price.String())
		}
	}

	// 检查路径是否存在
	// Check if path exists
	exists := xyJson.Exists(root, "$.store.magazine")
	fmt.Printf("Magazine section exists: %v\n", exists)

	// 统计书籍数量
	// Count books
	bookCount := xyJson.Count(root, "$.store.book[*]")
	fmt.Printf("Number of books: %d\n", bookCount)

	fmt.Println()
}

// serializationExample 演示序列化功能
// serializationExample demonstrates serialization functionality
func serializationExample() {
	fmt.Println("=== Serialization Example ===")

	// 创建复杂对象
	// Create complex object
	obj := xyJson.CreateObject()
	obj.Set("name", "Charlie")
	obj.Set("age", 28)

	// 添加数组
	// Add array
	hobbies := xyJson.CreateArray()
	hobbies.Append("reading")
	hobbies.Append("swimming")
	hobbies.Append("coding")
	obj.Set("hobbies", hobbies)

	// 添加嵌套对象
	// Add nested object
	address := xyJson.CreateObject()
	address.Set("street", "123 Main St")
	address.Set("city", "New York")
	address.Set("zipcode", "10001")
	obj.Set("address", address)

	// 紧凑序列化
	// Compact serialization
	compact, err := xyJson.Compact(obj)
	if err == nil {
		fmt.Println("Compact JSON:")
		fmt.Println(compact)
	}

	// 美化序列化
	// Pretty serialization
	pretty, err := xyJson.Pretty(obj)
	if err == nil {
		fmt.Println("\nPretty JSON:")
		fmt.Println(pretty)
	}

	// 序列化为字节数组
	// Serialize to byte array
	data, err := xyJson.Serialize(obj)
	if err == nil {
		fmt.Printf("\nSerialized byte length: %d\n", len(data))
	}

	fmt.Println()
}
