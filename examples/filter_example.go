package main

import (
	"fmt"
	"log"
	"strconv"

	xyJson "github.com/ihuem/xyJson"
)

func main7() {
	// 创建示例数据
	jsonData := `{
		"employees": [
			{"name": "张三", "age": 28, "department": "Engineering", "salary": 120000, "active": true},
			{"name": "李四", "age": 32, "department": "Marketing", "salary": 85000, "active": true},
			{"name": "王五", "age": 45, "department": "Engineering", "salary": 150000, "active": false},
			{"name": "赵六", "age": 29, "department": "Sales", "salary": 95000, "active": true},
			{"name": "钱七", "age": 38, "department": "Engineering", "salary": 135000, "active": true}
		],
		"products": [
			{"name": "笔记本电脑", "price": 8999, "category": "Electronics", "inStock": true, "rating": 4.5},
			{"name": "无线鼠标", "price": 199, "category": "Electronics", "inStock": false, "rating": 4.2},
			{"name": "办公椅", "price": 1299, "category": "Furniture", "inStock": true, "rating": 4.8},
			{"name": "台灯", "price": 299, "category": "Furniture", "inStock": true, "rating": 4.0}
		]
	}`

	// 解析JSON
	root, err := xyJson.ParseString(jsonData)
	if err != nil {
		log.Fatalf("解析JSON失败: %v", err)
	}

	fmt.Println("=== Filter函数使用示例 ===")

	// 示例1: 过滤高薪员工 (薪资 > 100000)
	fmt.Println("\n1. 过滤高薪员工 (薪资 > 100000):")
	highEarners, err := xyJson.Filter(root, "$.employees[*]", func(emp xyJson.IValue) bool {
		salary, err := xyJson.Get(emp, "$.salary")
		if err != nil {
			return false
		}
		if scalarValue, ok := salary.(xyJson.IScalarValue); ok {
			num, _ := scalarValue.Float64()
			return num > 100000
		}
		return false
	})
	if err != nil {
		log.Printf("过滤失败: %v", err)
	} else {
		for i, emp := range highEarners {
			name, _ := xyJson.GetString(emp, "$.name")
			salary := xyJson.MustGetFloat64(emp, "$.salary")
			fmt.Printf("  %d. %s - 薪资: %.0f\n", i+1, name, salary)
		}
	}

	// 示例2: 过滤工程部门的员工
	fmt.Println("\n2. 过滤工程部门的员工:")
	engineers, err := xyJson.Filter(root, "$.employees[*]", func(emp xyJson.IValue) bool {
		dept, err := xyJson.GetString(emp, "$.department")
		return err == nil && dept == "Engineering"
	})
	if err != nil {
		log.Printf("过滤失败: %v", err)
	} else {
		for i, emp := range engineers {
			name, _ := xyJson.GetString(emp, "$.name")
			age := xyJson.MustGetFloat64(emp, "$.age")
			fmt.Printf("  %d. %s - 年龄: %.0f\n", i+1, name, age)
		}
	}

	// 示例3: 过滤在职且年龄小于35的员工
	fmt.Println("\n3. 过滤在职且年龄小于35的员工:")
	youngActiveEmployees := xyJson.MustFilter(root, "$.employees[*]", func(emp xyJson.IValue) bool {
		active, _ := xyJson.GetBool(emp, "$.active")
		age := xyJson.MustGetInt64(emp, "$.age")
		return active && age < 35
	})
	for i, emp := range youngActiveEmployees {
		name, _ := xyJson.GetString(emp, "$.name")
		age := xyJson.MustGetFloat64(emp, "$.age")
		dept, _ := xyJson.GetString(emp, "$.department")
		fmt.Printf("  %d. %s - 年龄: %.0f, 部门: %s\n", i+1, name, age, dept)
	}

	// 示例4: 过滤库存商品
	fmt.Println("\n4. 过滤库存商品:")
	inStockProducts := xyJson.MustFilter(root, "$.products[*]", func(product xyJson.IValue) bool {
		inStock, _ := xyJson.GetBool(product, "$.inStock")
		return inStock
	})
	for i, product := range inStockProducts {
		name, _ := xyJson.GetString(product, "$.name")
		price := xyJson.MustGetFloat64(product, "$.price")
		rating := xyJson.MustGetFloat64(product, "$.rating")
		fmt.Printf("  %d. %s - 价格: %.0f, 评分: %.1f\n", i+1, name, price, rating)
	}

	// 示例5: 过滤高评分且价格适中的产品 (评分 >= 4.5, 价格 <= 5000)
	fmt.Println("\n5. 过滤高评分且价格适中的产品:")
	goodValueProducts, err := xyJson.Filter(root, "$.products[*]", func(product xyJson.IValue) bool {
		rating, err1 := xyJson.Get(product, "$.rating")
		price, err2 := xyJson.Get(product, "$.price")
		if err1 != nil || err2 != nil {
			return false
		}
		ratingScalar, ok1 := rating.(xyJson.IScalarValue)
		priceScalar, ok2 := price.(xyJson.IScalarValue)
		if !ok1 || !ok2 {
			return false
		}
		ratingVal, _ := ratingScalar.Float64()
		priceVal, _ := priceScalar.Float64()
		return ratingVal >= 4.5 && priceVal <= 5000
	})
	if err != nil {
		log.Printf("过滤失败: %v", err)
	} else {
		for i, product := range goodValueProducts {
			name, _ := xyJson.GetString(product, "$.name")
			price := xyJson.MustGetFloat64(product, "$.price")
			rating := xyJson.MustGetFloat64(product, "$.rating")
			category, _ := xyJson.GetString(product, "$.category")
			fmt.Printf("  %d. %s - 价格: %.0f, 评分: %.1f, 类别: %s\n", i+1, name, price, rating, category)
		}
	}

	// 示例6: 复杂过滤 - 根据多个条件过滤员工
	fmt.Println("\n6. 复杂过滤 - 工程部门的高薪在职员工:")
	seniorEngineers := xyJson.MustFilter(root, "$.employees[*]", func(emp xyJson.IValue) bool {
		dept, _ := xyJson.GetString(emp, "$.department")
		salary := xyJson.MustGetFloat64(emp, "$.salary")
		active, _ := xyJson.GetBool(emp, "$.active")
		return dept == "Engineering" && salary > 120000 && active
	})
	for i, emp := range seniorEngineers {
		name, _ := xyJson.GetString(emp, "$.name")
		salary := xyJson.MustGetFloat64(emp, "$.salary")
		age := xyJson.MustGetFloat64(emp, "$.age")
		fmt.Printf("  %d. %s - 薪资: %.0f, 年龄: %.0f\n", i+1, name, salary, age)
	}

	// 示例7: 错误处理示例
	fmt.Println("\n7. 错误处理示例:")
	_, err = xyJson.Filter(root, "$.invalid[*]", func(item xyJson.IValue) bool {
		return true
	})
	if err != nil {
		fmt.Printf("  预期的错误: %v\n", err)
	}

	// 使用MustFilter处理错误路径（返回空数组）
	emptyResult := xyJson.MustFilter(root, "$.invalid[*]", func(item xyJson.IValue) bool {
		return true
	})
	fmt.Printf("  MustFilter结果长度: %d (空数组)\n", len(emptyResult))

	fmt.Println("\n=== Filter函数示例完成 ===")
}

// 辅助函数：格式化数字
func formatNumber(num float64) string {
	return strconv.FormatFloat(num, 'f', 0, 64)
}
