package main

import (
	"fmt"
	"log"

	"github.com/ihuem/xyJson"
)

func main() {
	// 测试基本的map解析功能
	data := map[string]interface{}{
		"name":   "John Doe",
		"age":    30,
		"active": true,
		"score":  95.5,
		"address": map[string]interface{}{
			"city":    "New York",
			"zip":     "10001",
			"country": "USA",
		},
		"hobbies": []interface{}{"reading", "swimming", "coding"},
		"tags":    []interface{}{"developer", "golang", "json"},
	}

	fmt.Println("=== 测试 ParseFromMap ===")
	value, err := xyJson.ParseFromMap(data)
	if err != nil {
		log.Fatal("ParseFromMap failed:", err)
	}

	fmt.Println("解析成功!")
	fmt.Println("JSON输出:", value.String())

	// 测试访问具体字段
	obj := value.(xyJson.IObject)
	fmt.Println("\n=== 字段访问测试 ===")
	fmt.Println("姓名:", obj.Get("name").String())
	fmt.Println("年龄:", obj.Get("age").String())
	fmt.Println("活跃状态:", obj.Get("active").String())
	fmt.Println("分数:", obj.Get("score").String())

	// 测试嵌套对象
	address := obj.Get("address").(xyJson.IObject)
	fmt.Println("\n=== 嵌套对象测试 ===")
	fmt.Println("城市:", address.Get("city").String())
	fmt.Println("邮编:", address.Get("zip").String())
	fmt.Println("国家:", address.Get("country").String())

	// 测试数组
	hobbies := obj.Get("hobbies").(xyJson.IArray)
	fmt.Println("\n=== 数组测试 ===")
	fmt.Printf("爱好数量: %d\n", hobbies.Length())
	for i := 0; i < hobbies.Length(); i++ {
		fmt.Printf("爱好 %d: %s\n", i+1, hobbies.Get(i).String())
	}

	fmt.Println("\n=== 测试 MustParseFromMap ===")
	simpleData := map[string]interface{}{
		"message": "Hello World",
		"count":   42,
	}

	simpleValue := xyJson.MustParseFromMap(simpleData)
	fmt.Println("简单数据解析结果:", simpleValue.String())

	fmt.Println("\n所有测试完成!")
}
