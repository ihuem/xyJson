package main

import (
	"fmt"

	xyJson "github.com/ihuem/xyJson"
)

func main5() {
	fmt.Println("=== GetWithDefault 方法示例 ===")
	fmt.Println()

	// 模拟配置文件JSON
	configJSON := `{
		"server": {
			"host": "api.example.com",
			"port": 8080,
			"ssl": true
		},
		"database": {
			"host": "db.example.com",
			"timeout": 30
		},
		"features": {
			"cache": true,
			"logging": true
		}
	}`

	root, err := xyJson.ParseString(configJSON)
	if err != nil {
		fmt.Printf("解析配置失败: %v\n", err)
		return
	}

	fmt.Println("1. 配置读取场景 - GetWithDefault 的优势")
	fmt.Println("-------------------------------------------")

	// 使用GetWithDefault方法读取配置，代码非常简洁
	serverHost := xyJson.GetStringWithDefault(root, "$.server.host", "localhost")
	serverPort := xyJson.GetIntWithDefault(root, "$.server.port", 3000)
	serverSSL := xyJson.GetBoolWithDefault(root, "$.server.ssl", false)

	// 读取可能不存在的配置项，自动使用默认值
	maxConnections := xyJson.GetIntWithDefault(root, "$.server.maxConnections", 100)
	readTimeout := xyJson.GetIntWithDefault(root, "$.server.readTimeout", 5000)
	debugMode := xyJson.GetBoolWithDefault(root, "$.server.debug", false)

	fmt.Printf("服务器配置:\n")
	fmt.Printf("  主机: %s\n", serverHost)
	fmt.Printf("  端口: %d\n", serverPort)
	fmt.Printf("  SSL: %t\n", serverSSL)
	fmt.Printf("  最大连接数: %d (默认值)\n", maxConnections)
	fmt.Printf("  读取超时: %d毫秒 (默认值)\n", readTimeout)
	fmt.Printf("  调试模式: %t (默认值)\n", debugMode)
	fmt.Println()

	fmt.Println("2. 对比不同方法的代码复杂度")
	fmt.Println("----------------------------")

	// 方法1：使用Get方法（需要错误处理）
	fmt.Println("方法1 - Get方法（需要错误处理）:")
	dbTimeout1, err := xyJson.GetInt(root, "$.database.timeout")
	if err != nil {
		dbTimeout1 = 60 // 默认值
		fmt.Printf("  数据库超时: %d (使用默认值，因为获取失败)\n", dbTimeout1)
	} else {
		fmt.Printf("  数据库超时: %d\n", dbTimeout1)
	}

	// 方法2：使用TryGet方法（需要判断布尔值）
	fmt.Println("方法2 - TryGet方法（需要判断布尔值）:")
	dbTimeout2 := 60 // 默认值
	if timeout, ok := xyJson.TryGetInt(root, "$.database.timeout"); ok {
		dbTimeout2 = timeout
		fmt.Printf("  数据库超时: %d\n", dbTimeout2)
	} else {
		fmt.Printf("  数据库超时: %d (使用默认值)\n", dbTimeout2)
	}

	// 方法3：使用GetWithDefault方法（最简洁）
	fmt.Println("方法3 - GetWithDefault方法（最简洁）:")
	dbTimeout3 := xyJson.GetIntWithDefault(root, "$.database.timeout", 60)
	fmt.Printf("  数据库超时: %d\n", dbTimeout3)
	fmt.Println()

	fmt.Println("3. 实际应用场景示例")
	fmt.Println("-------------------")

	// 场景1：Web服务器配置
	fmt.Println("场景1 - Web服务器配置:")
	config := map[string]interface{}{
		"host":           xyJson.GetStringWithDefault(root, "$.server.host", "0.0.0.0"),
		"port":           xyJson.GetIntWithDefault(root, "$.server.port", 8080),
		"ssl":            xyJson.GetBoolWithDefault(root, "$.server.ssl", false),
		"maxConnections": xyJson.GetIntWithDefault(root, "$.server.maxConnections", 1000),
		"requestTimeout": xyJson.GetIntWithDefault(root, "$.server.requestTimeout", 30000),
	}

	for key, value := range config {
		fmt.Printf("  %s: %v\n", key, value)
	}
	fmt.Println()

	// 场景2：功能开关配置
	fmt.Println("场景2 - 功能开关配置:")
	features := map[string]bool{
		"cache":      xyJson.GetBoolWithDefault(root, "$.features.cache", true),
		"logging":    xyJson.GetBoolWithDefault(root, "$.features.logging", true),
		"monitoring": xyJson.GetBoolWithDefault(root, "$.features.monitoring", false),
		"analytics":  xyJson.GetBoolWithDefault(root, "$.features.analytics", false),
		"backup":     xyJson.GetBoolWithDefault(root, "$.features.backup", true),
	}

	for feature, enabled := range features {
		status := "禁用"
		if enabled {
			status = "启用"
		}
		fmt.Printf("  %s: %s\n", feature, status)
	}
	fmt.Println()

	fmt.Println("4. 复杂类型的默认值处理")
	fmt.Println("------------------------")

	// 创建默认的数据库配置对象
	defaultDBConfig := xyJson.CreateObject()
	defaultDBConfig.Set("host", xyJson.CreateString("localhost"))
	defaultDBConfig.Set("port", xyJson.CreateNumber(5432))
	defaultDBConfig.Set("database", xyJson.CreateString("myapp"))

	// 尝试获取Redis配置，如果不存在则使用默认配置
	redisConfig := xyJson.GetObjectWithDefault(root, "$.redis", defaultDBConfig)
	if redisConfig != nil {
		fmt.Println("Redis配置（使用默认值）:")
		host, _ := xyJson.GetString(redisConfig, "$.host")
		port, _ := xyJson.GetInt(redisConfig, "$.port")
		db, _ := xyJson.GetString(redisConfig, "$.database")
		fmt.Printf("  主机: %s\n", host)
		fmt.Printf("  端口: %d\n", port)
		fmt.Printf("  数据库: %s\n", db)
	}
	fmt.Println()

	// 演示当defaultValue为nil时返回空对象/数组
	fmt.Println("演示nil默认值处理:")
	emptyObj := xyJson.GetObjectWithDefault(root, "$.nonexistent.object", nil)
	fmt.Printf("空对象大小: %d\n", emptyObj.Size())

	emptyArr := xyJson.GetArrayWithDefault(root, "$.nonexistent.array", nil)
	fmt.Printf("空数组长度: %d\n", emptyArr.Length())
	fmt.Println()

	fmt.Println("5. GetWithDefault 方法的优势总结")
	fmt.Println("--------------------------------")
	fmt.Println("✓ 只返回一个值，无需错误处理")
	fmt.Println("✓ 无需判断布尔值，代码更简洁")
	fmt.Println("✓ 用户可以指定合理的默认值")
	fmt.Println("✓ 特别适合配置文件读取场景")
	fmt.Println("✓ 减少样板代码，提高开发效率")
	fmt.Println("✓ 代码可读性更强，意图更明确")
	fmt.Println()

	fmt.Println("6. 使用建议")
	fmt.Println("-----------")
	fmt.Println("• 配置读取：优先使用 GetWithDefault")
	fmt.Println("• 需要区分错误类型：使用 Get")
	fmt.Println("• 需要知道是否成功：使用 TryGet")
	fmt.Println("• 确信数据正确：使用 Must（谨慎使用）")
}
