package main

import (
	"fmt"
	"log"
	"time"

	xyJson "github.com/ihuem/xyJson"
)

// demonstrateBasicCompiledPath 演示基本的预编译路径使用
func demonstrateBasicCompiledPath() {
	fmt.Println("=== 基本预编译路径使用 ===\n")

	// 示例JSON数据
	data := `{
		"users": [
			{"name": "Alice", "age": 30, "email": "alice@example.com", "active": true},
			{"name": "Bob", "age": 25, "email": "bob@example.com", "active": false},
			{"name": "Charlie", "age": 35, "email": "charlie@example.com", "active": true}
		],
		"config": {
			"server": {
				"host": "localhost",
				"port": 8080,
				"ssl": false
			}
		}
	}`

	root, err := xyJson.ParseString(data)
	if err != nil {
		log.Fatal(err)
	}

	// 1. 编译常用路径
	fmt.Println("1. 编译JSONPath表达式:")
	userNamesPath, err := xyJson.CompilePath("$.users[*].name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ 编译用户姓名路径")

	activeUsersPath, err := xyJson.CompilePath("$.users[?(@.active == true)]")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ 编译活跃用户路径")

	serverConfigPath, err := xyJson.CompilePath("$.config.server")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ 编译服务器配置路径")

	// 2. 使用预编译路径查询
	fmt.Println("\n2. 使用预编译路径查询:")

	// 查询所有用户姓名
	names, err := userNamesPath.QueryAll(root)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("所有用户姓名:")
	for i, nameValue := range names {
		if name, ok := xyJson.TryGetString(nameValue, "$"); ok {
			fmt.Printf("  %d. %s\n", i+1, name)
		}
	}

	// 查询活跃用户
	activeUsers, err := activeUsersPath.QueryAll(root)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n活跃用户数量: %d\n", len(activeUsers))
	for i, user := range activeUsers {
		if name, ok := xyJson.TryGetString(user, "$.name"); ok {
			if age, ok := xyJson.TryGetInt(user, "$.age"); ok {
				fmt.Printf("  %d. %s (%d岁)\n", i+1, name, age)
			}
		}
	}

	// 查询服务器配置
	serverConfig, err := serverConfigPath.Query(root)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n服务器配置:")
	host := xyJson.GetStringWithDefault(serverConfig, "$.host", "unknown")
	port := xyJson.GetIntWithDefault(serverConfig, "$.port", 0)
	ssl := xyJson.GetBoolWithDefault(serverConfig, "$.ssl", false)
	fmt.Printf("  主机: %s\n", host)
	fmt.Printf("  端口: %d\n", port)
	fmt.Printf("  SSL: %t\n", ssl)
}

// demonstratePerformanceComparison 演示性能对比
func demonstratePerformanceComparison() {
	fmt.Println("\n=== 性能对比演示 ===\n")

	data := `{
		"users": [
			{"name": "Alice", "age": 30},
			{"name": "Bob", "age": 25},
			{"name": "Charlie", "age": 35}
		]
	}`

	root, err := xyJson.ParseString(data)
	if err != nil {
		log.Fatal(err)
	}

	const iterations = 10000
	path := "$.users[0].name"

	// 1. 传统方式性能测试
	fmt.Printf("1. 传统方式 (%d次查询):\n", iterations)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, err := xyJson.GetString(root, path)
		if err != nil {
			log.Fatal(err)
		}
	}
	traditionalDuration := time.Since(start)
	fmt.Printf("  耗时: %v\n", traditionalDuration)
	fmt.Printf("  平均每次: %v\n", traditionalDuration/iterations)

	// 2. 预编译方式性能测试
	fmt.Printf("\n2. 预编译方式 (%d次查询):\n", iterations)
	compiledPath, err := xyJson.CompilePath(path)
	if err != nil {
		log.Fatal(err)
	}

	start = time.Now()
	for i := 0; i < iterations; i++ {
		_, err := compiledPath.Query(root)
		if err != nil {
			log.Fatal(err)
		}
	}
	compiledDuration := time.Since(start)
	fmt.Printf("  耗时: %v\n", compiledDuration)
	fmt.Printf("  平均每次: %v\n", compiledDuration/iterations)

	// 3. 性能提升计算
	improvementPercent := float64(traditionalDuration-compiledDuration) / float64(traditionalDuration) * 100
	fmt.Printf("\n3. 性能提升:\n")
	fmt.Printf("  提升幅度: %.1f%%\n", improvementPercent)
	fmt.Printf("  速度倍数: %.1fx\n", float64(traditionalDuration)/float64(compiledDuration))
}

// demonstrateCacheManagement 演示缓存管理
func demonstrateCacheManagement() {
	fmt.Println("\n=== 缓存管理演示 ===\n")

	// 1. 查看初始缓存状态
	size, maxSize := xyJson.GetPathCacheStats()
	fmt.Printf("初始缓存状态: 大小=%d, 最大大小=%d\n", size, maxSize)

	// 2. 编译一些路径（会被缓存）
	paths := []string{
		"$.user.name",
		"$.user.age",
		"$.user.email",
		"$.users[*].name",
		"$.users[0].profile",
	}

	fmt.Println("\n编译路径（首次编译，会缓存）:")
	for _, path := range paths {
		_, err := xyJson.CompilePath(path)
		if err != nil {
			fmt.Printf("编译 %s 失败: %v\n", path, err)
			continue
		}
		fmt.Printf("✓ 编译: %s\n", path)
	}

	// 3. 查看缓存状态
	size, maxSize = xyJson.GetPathCacheStats()
	fmt.Printf("\n编译后缓存状态: 大小=%d, 最大大小=%d\n", size, maxSize)

	// 4. 再次编译相同路径（应该命中缓存）
	fmt.Println("\n再次编译相同路径（应该命中缓存）:")
	for _, path := range paths {
		_, err := xyJson.CompilePath(path)
		if err != nil {
			fmt.Printf("编译 %s 失败: %v\n", path, err)
			continue
		}
		fmt.Printf("✓ 缓存命中: %s\n", path)
	}

	// 5. 查看最终缓存状态
	size, maxSize = xyJson.GetPathCacheStats()
	fmt.Printf("\n最终缓存状态: 大小=%d, 最大大小=%d\n", size, maxSize)

	// 6. 缓存信息
	fmt.Printf("当前缓存使用率: %.1f%%\n", float64(size)/float64(maxSize)*100)

	// 7. 演示缓存大小限制
	fmt.Println("\n设置缓存最大大小为3:")
	xyJson.SetPathCacheMaxSize(3)

	// 编译更多路径，触发缓存清理
	morePaths := []string{
		"$.config.database.host",
		"$.config.database.port",
		"$.config.redis.host",
		"$.config.redis.port",
	}

	for _, path := range morePaths {
		xyJson.CompilePath(path)
		size, _ = xyJson.GetPathCacheStats()
		fmt.Printf("编译 %s 后缓存大小: %d\n", path, size)
	}

	// 8. 清空缓存
	fmt.Println("\n清空缓存:")
	xyJson.ClearPathCache()
	size, maxSize = xyJson.GetPathCacheStats()
	fmt.Printf("清空后缓存状态: 大小=%d, 最大大小=%d\n", size, maxSize)
}

// demonstrateModificationOperations 演示修改操作
func demonstrateModificationOperations() {
	fmt.Println("\n=== 修改操作演示 ===\n")

	data := `{
		"users": [
			{"name": "Alice", "age": 30, "email": "alice@example.com"},
			{"name": "Bob", "age": 25, "email": "bob@example.com"}
		]
	}`

	root, err := xyJson.ParseString(data)
	if err != nil {
		log.Fatal(err)
	}

	// 编译修改路径
	firstUserAgePath, _ := xyJson.CompilePath("$.users[0].age")
	firstUserEmailPath, _ := xyJson.CompilePath("$.users[0].email")
	newUserPath, _ := xyJson.CompilePath("$.users[2]")

	// 1. 设置值
	fmt.Println("1. 设置第一个用户的年龄为31:")
	newAge := xyJson.CreateNumber(31)
	if err := firstUserAgePath.Set(root, newAge); err != nil {
		fmt.Printf("设置失败: %v\n", err)
	} else {
		fmt.Println("✓ 年龄设置成功")
	}

	// 2. 更新邮箱
	fmt.Println("\n2. 更新第一个用户的邮箱:")
	newEmail := xyJson.CreateString("alice.updated@example.com")
	if err := firstUserEmailPath.Set(root, newEmail); err != nil {
		fmt.Printf("更新失败: %v\n", err)
	} else {
		fmt.Println("✓ 邮箱更新成功")
	}

	// 3. 添加新用户
	fmt.Println("\n3. 添加新用户:")
	newUser := xyJson.CreateObject()
	newUser.Set("name", xyJson.CreateString("Charlie"))
	newUser.Set("age", xyJson.CreateNumber(28))
	newUser.Set("email", xyJson.CreateString("charlie@example.com"))

	if err := newUserPath.Set(root, newUser); err != nil {
		fmt.Printf("添加失败: %v\n", err)
	} else {
		fmt.Println("✓ 新用户添加成功")
	}

	// 4. 验证修改结果
	fmt.Println("\n4. 验证修改结果:")

	// 检查年龄是否更新
	if age, err := firstUserAgePath.Query(root); err == nil {
		if ageVal, ok := xyJson.TryGetInt(age, "$"); ok {
			fmt.Printf("第一个用户年龄: %d\n", ageVal)
		}
	}

	// 检查邮箱是否更新
	if email, err := firstUserEmailPath.Query(root); err == nil {
		if emailVal, ok := xyJson.TryGetString(email, "$"); ok {
			fmt.Printf("第一个用户邮箱: %s\n", emailVal)
		}
	}

	// 检查新用户是否添加
	if newUser, err := newUserPath.Query(root); err == nil {
		if name, ok := xyJson.TryGetString(newUser, "$.name"); ok {
			fmt.Printf("新用户姓名: %s\n", name)
		}
	}

	// 5. 输出最终JSON
	fmt.Println("\n5. 最终JSON结构:")
	fmt.Println(root.String())
}

// demonstrateAdvancedUsage 演示高级用法
func demonstrateAdvancedUsage() {
	fmt.Println("\n=== 高级用法演示 ===\n")

	data := `{
		"products": [
			{"id": 1, "name": "Laptop", "price": 999.99, "category": "Electronics", "inStock": true},
			{"id": 2, "name": "Book", "price": 29.99, "category": "Education", "inStock": false},
			{"id": 3, "name": "Phone", "price": 699.99, "category": "Electronics", "inStock": true},
			{"id": 4, "name": "Desk", "price": 199.99, "category": "Furniture", "inStock": true}
		]
	}`

	root, err := xyJson.ParseString(data)
	if err != nil {
		log.Fatal(err)
	}

	// 预编译复杂查询路径
	paths := map[string]*xyJson.CompiledPath{
		"allProducts":         mustCompilePath("$.products[*]"),
		"electronicsProducts": mustCompilePath("$.products[?(@.category == 'Electronics')]"),
		"inStockProducts":     mustCompilePath("$.products[?(@.inStock == true)]"),
		"expensiveProducts":   mustCompilePath("$.products[?(@.price > 500)]"),
		"productNames":        mustCompilePath("$.products[*].name"),
		"productPrices":       mustCompilePath("$.products[*].price"),
	}

	// 批量查询和统计
	fmt.Println("1. 产品统计:")
	allProducts, _ := paths["allProducts"].QueryAll(root)
	fmt.Printf("总产品数: %d\n", len(allProducts))

	electronicsProducts, _ := paths["electronicsProducts"].QueryAll(root)
	fmt.Printf("电子产品数: %d\n", len(electronicsProducts))

	inStockProducts, _ := paths["inStockProducts"].QueryAll(root)
	fmt.Printf("有库存产品数: %d\n", len(inStockProducts))

	expensiveProducts, _ := paths["expensiveProducts"].QueryAll(root)
	fmt.Printf("高价产品数 (>$500): %d\n", len(expensiveProducts))

	// 2. 产品列表
	fmt.Println("\n2. 电子产品列表:")
	for i, product := range electronicsProducts {
		name := xyJson.GetStringWithDefault(product, "$.name", "Unknown")
		price := xyJson.GetFloat64WithDefault(product, "$.price", 0.0)
		inStock := xyJson.GetBoolWithDefault(product, "$.inStock", false)
		stockStatus := "缺货"
		if inStock {
			stockStatus = "有货"
		}
		fmt.Printf("  %d. %s - $%.2f (%s)\n", i+1, name, price, stockStatus)
	}

	// 3. 价格分析
	fmt.Println("\n3. 价格分析:")
	prices, _ := paths["productPrices"].QueryAll(root)
	var totalPrice, maxPrice, minPrice float64
	maxPrice = 0
	minPrice = 999999

	for _, priceValue := range prices {
		if price, ok := xyJson.TryGetFloat64(priceValue, "$"); ok {
			totalPrice += price
			if price > maxPrice {
				maxPrice = price
			}
			if price < minPrice {
				minPrice = price
			}
		}
	}

	avgPrice := totalPrice / float64(len(prices))
	fmt.Printf("平均价格: $%.2f\n", avgPrice)
	fmt.Printf("最高价格: $%.2f\n", maxPrice)
	fmt.Printf("最低价格: $%.2f\n", minPrice)
	fmt.Printf("总价值: $%.2f\n", totalPrice)

	// 4. 使用Exists和Count方法
	fmt.Println("\n4. 路径检查:")
	fmt.Printf("是否存在产品: %t\n", paths["allProducts"].Exists(root))
	fmt.Printf("电子产品数量: %d\n", paths["electronicsProducts"].Count(root))
	fmt.Printf("有库存产品数量: %d\n", paths["inStockProducts"].Count(root))
}

// mustCompilePath 便利函数：编译路径，失败时panic
func mustCompilePath(path string) *xyJson.CompiledPath {
	compiled, err := xyJson.CompilePath(path)
	if err != nil {
		panic(fmt.Sprintf("编译路径失败: %v", err))
	}
	return compiled
}

func main() {
	fmt.Println("xyJson JSONPath预编译功能演示")
	fmt.Println("================================")

	// 演示基本用法
	demonstrateBasicCompiledPath()

	// 演示性能对比
	demonstratePerformanceComparison()

	// 演示缓存管理
	demonstrateCacheManagement()

	// 演示修改操作
	demonstrateModificationOperations()

	// 演示高级用法
	demonstrateAdvancedUsage()

	fmt.Println("\n=== 演示完成 ===")
	fmt.Println("\n总结:")
	fmt.Println("- JSONPath预编译功能可以显著提升重复查询的性能")
	fmt.Println("- 内置缓存机制进一步优化了编译开销")
	fmt.Println("- 支持完整的JSONPath操作：查询、修改、删除、检查")
	fmt.Println("- 线程安全设计，适合并发环境使用")
	fmt.Println("- 与现有API完全兼容，可以渐进式迁移")
}
