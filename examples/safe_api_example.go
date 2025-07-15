package main

import (
	"fmt"

	xyJson "github.com/ihuem/xyJson"
)

func main4() {
	fmt.Println("=== xyJson 安全API使用示例 ===")
	fmt.Println()

	// 准备测试数据
	data := `{
		"user": {
			"name": "Alice",
			"age": 30,
			"height": 165.5,
			"active": true,
			"profile": {
				"email": "alice@example.com"
			}
		},
		"settings": {
			"theme": "dark",
			"notifications": false
		}
	}`

	root, err := xyJson.ParseString(data)
	if err != nil {
		fmt.Printf("解析JSON失败: %v\n", err)
		return
	}

	fmt.Println("1. 传统Get方法 - 需要错误处理")
	fmt.Println("================================")

	// 使用Get方法获取存在的值
	name, err := xyJson.GetString(root, "$.user.name")
	if err != nil {
		fmt.Printf("获取姓名失败: %v\n", err)
	} else {
		fmt.Printf("姓名: %s\n", name)
	}

	// 使用Get方法获取不存在的值
	city, err := xyJson.GetString(root, "$.user.city")
	if err != nil {
		fmt.Printf("获取城市失败: %v\n", err)
	} else {
		fmt.Printf("城市: %s\n", city)
	}

	fmt.Println()
	fmt.Println("2. 新的TryGet方法 - 更安全简洁")
	fmt.Println("==================================")

	// 使用TryGet方法获取存在的值
	if name, ok := xyJson.TryGetString(root, "$.user.name"); ok {
		fmt.Printf("姓名: %s\n", name)
	} else {
		fmt.Println("获取姓名失败")
	}

	// 使用TryGet方法获取不存在的值
	if city, ok := xyJson.TryGetString(root, "$.user.city"); ok {
		fmt.Printf("城市: %s\n", city)
	} else {
		fmt.Println("城市信息不存在")
	}

	// 使用TryGet方法获取不同类型的值
	if age, ok := xyJson.TryGetInt(root, "$.user.age"); ok {
		fmt.Printf("年龄: %d岁\n", age)
	}

	if height, ok := xyJson.TryGetFloat64(root, "$.user.height"); ok {
		fmt.Printf("身高: %.1fcm\n", height)
	}

	if active, ok := xyJson.TryGetBool(root, "$.user.active"); ok {
		fmt.Printf("活跃状态: %v\n", active)
	}

	fmt.Println()
	fmt.Println("3. Must方法 - 有风险但简洁（已添加安全警告）")
	fmt.Println("=============================================")

	// 安全使用Must方法（确信数据存在）
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Must方法panic: %v\n", r)
			}
		}()

		// 这个会成功
		name := xyJson.MustGetString(root, "$.user.name")
		fmt.Printf("姓名（Must方法）: %s\n", name)

		// 这个会panic
		city := xyJson.MustGetString(root, "$.user.city")
		fmt.Printf("城市（Must方法）: %s\n", city) // 不会执行到这里
	}()

	fmt.Println()
	fmt.Println("4. 实际应用场景对比")
	fmt.Println("==================")

	// 场景1：配置读取（推荐使用TryGet + 默认值）
	fmt.Println("\n场景1: 配置读取")
	theme := "light" // 默认主题
	if t, ok := xyJson.TryGetString(root, "$.settings.theme"); ok {
		theme = t
	}
	fmt.Printf("当前主题: %s\n", theme)

	notifications := true // 默认开启通知
	if n, ok := xyJson.TryGetBool(root, "$.settings.notifications"); ok {
		notifications = n
	}
	fmt.Printf("通知设置: %v\n", notifications)

	// 场景2：数据验证（推荐使用Get方法获取详细错误）
	fmt.Println("\n场景2: 数据验证")
	email, err := xyJson.GetString(root, "$.user.profile.email")
	if err != nil {
		fmt.Printf("邮箱验证失败: %v\n", err)
	} else {
		fmt.Printf("邮箱: %s\n", email)
	}

	// 场景3：快速原型开发（可以谨慎使用Must方法）
	fmt.Println("\n场景3: 快速原型开发（谨慎使用Must）")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("原型代码出错: %v\n", r)
			}
		}()

		// 在确信数据结构的情况下使用Must方法
		userName := xyJson.MustGetString(root, "$.user.name")
		userAge := xyJson.MustGetInt(root, "$.user.age")
		fmt.Printf("用户信息: %s, %d岁\n", userName, userAge)
	}()

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("1. Get方法: 适合需要详细错误信息的场景")
	fmt.Println("2. TryGet方法: 最安全，适合大多数场景")
	fmt.Println("3. Must方法: 有风险，仅在确信数据正确时使用")
	fmt.Println("4. 推荐优先使用TryGet方法，它既安全又简洁")
}
