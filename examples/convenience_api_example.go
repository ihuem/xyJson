package main

import (
	"fmt"
	xyJson "github/ihuem/xyJson"
	"log"
)

func main3() {
	// 示例JSON数据
	data := `{
		"user": {
			"name": "Alice",
			"age": 30,
			"height": 165.5,
			"active": true,
			"profile": {
				"email": "alice@example.com",
				"phone": "+1234567890"
			},
			"hobbies": ["reading", "swimming", "coding"]
		},
		"product": {
			"name": "Laptop",
			"price": 999.99,
			"inStock": true,
			"quantity": 50
		}
	}`

	// 解析JSON
	root, err := xyJson.ParseString(data)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Println("=== xyJson便利API使用示例 ===")
	fmt.Println()

	// 1. 旧的方式：需要类型断言
	fmt.Println("1. 旧的方式（需要类型断言）:")
	priceValue, err := xyJson.Get(root, "$.product.price")
	if err != nil {
		log.Printf("Get failed: %v", err)
	} else {
		// 需要类型断言
		if scalarValue, ok := priceValue.(xyJson.IScalarValue); ok {
			if price, err := scalarValue.Float64(); err == nil {
				fmt.Printf("  产品价格: %.2f\n", price)
			} else {
				log.Printf("Float64 conversion failed: %v", err)
			}
		} else {
			log.Println("Failed to cast to IScalarValue")
		}
	}
	fmt.Println()

	// 2. 新的便利方式：直接获取类型化的值
	fmt.Println("2. 新的便利方式（直接获取类型化值）:")
	
	// 获取字符串值
	if name, err := xyJson.GetString(root, "$.user.name"); err == nil {
		fmt.Printf("  用户姓名: %s\n", name)
	} else {
		fmt.Printf("  获取用户姓名失败: %v\n", err)
	}

	// 获取整数值
	if age, err := xyJson.GetInt(root, "$.user.age"); err == nil {
		fmt.Printf("  用户年龄: %d岁\n", age)
	} else {
		fmt.Printf("  获取用户年龄失败: %v\n", err)
	}

	// 获取浮点数值
	if height, err := xyJson.GetFloat64(root, "$.user.height"); err == nil {
		fmt.Printf("  用户身高: %.1fcm\n", height)
	} else {
		fmt.Printf("  获取用户身高失败: %v\n", err)
	}

	// 获取布尔值
	if active, err := xyJson.GetBool(root, "$.user.active"); err == nil {
		fmt.Printf("  用户状态: %s\n", map[bool]string{true: "活跃", false: "非活跃"}[active])
	} else {
		fmt.Printf("  获取用户状态失败: %v\n", err)
	}

	// 获取对象
	if profile, err := xyJson.GetObject(root, "$.user.profile"); err == nil {
		fmt.Printf("  用户档案包含 %d 个字段\n", profile.Size())
		profile.Range(func(key string, value xyJson.IValue) bool {
			fmt.Printf("    %s: %s\n", key, value.String())
			return true
		})
	} else {
		fmt.Printf("  获取用户档案失败: %v\n", err)
	}

	// 获取数组
	if hobbies, err := xyJson.GetArray(root, "$.user.hobbies"); err == nil {
		fmt.Printf("  用户有 %d 个爱好:\n", hobbies.Length())
		hobbies.Range(func(index int, value xyJson.IValue) bool {
			fmt.Printf("    %d. %s\n", index+1, value.String())
			return true
		})
	} else {
		fmt.Printf("  获取用户爱好失败: %v\n", err)
	}
	fmt.Println()

	// 3. 最简洁的方式：Must版本（适用于确信路径存在的场景）
	fmt.Println("3. 最简洁的方式（Must版本，适用于确信路径存在的场景）:")
	
	// 注意：Must版本在路径不存在或转换失败时会panic，只在确信数据正确时使用
	productName := xyJson.MustGetString(root, "$.product.name")
	productPrice := xyJson.MustGetFloat64(root, "$.product.price")
	productQuantity := xyJson.MustGetInt(root, "$.product.quantity")
	inStock := xyJson.MustGetBool(root, "$.product.inStock")
	
	fmt.Printf("  产品: %s\n", productName)
	fmt.Printf("  价格: $%.2f\n", productPrice)
	fmt.Printf("  库存: %d件\n", productQuantity)
	fmt.Printf("  状态: %s\n", map[bool]string{true: "有库存", false: "缺货"}[inStock])
	fmt.Println()

	// 4. 错误处理示例
	fmt.Println("4. 错误处理示例:")
	
	// 尝试获取不存在的路径
	if _, err := xyJson.GetString(root, "$.user.nonexistent"); err != nil {
		fmt.Printf("  预期的错误 - 路径不存在: %v\n", err)
	}
	
	// 尝试进行错误的类型转换
	if _, err := xyJson.GetInt(root, "$.user.name"); err != nil {
		fmt.Printf("  预期的错误 - 类型转换失败: %v\n", err)
	}
	fmt.Println()

	// 5. 性能对比示例
	fmt.Println("5. 代码简洁性对比:")
	fmt.Println("  旧方式需要的代码行数: ~8行（包括错误处理和类型断言）")
	fmt.Println("  新方式需要的代码行数: ~3行（包括错误处理）")
	fmt.Println("  Must方式需要的代码行数: ~1行（无错误处理）")
	fmt.Println()

	fmt.Println("=== 总结 ===")
	fmt.Println("新的便利API提供了以下优势:")
	fmt.Println("1. 无需手动类型断言")
	fmt.Println("2. 代码更简洁易读")
	fmt.Println("3. 类型安全")
	fmt.Println("4. 提供Must版本用于简化确定场景")
	fmt.Println("5. 保持与原有API的完全兼容性")
}