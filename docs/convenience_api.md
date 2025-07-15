# xyJson 便利API使用指南

## 概述

xyJson 提供了一套便利的API，让您可以直接获取特定类型的值，而无需进行手动类型断言。这些API大大简化了JSON数据的访问和处理。

## 问题背景

在使用原始的 `xyJson.Get` 方法时，您需要进行额外的类型断言才能使用返回的值。这个过程虽然灵活，但在日常使用中显得繁琐：

```go
// 旧的方式：需要类型断言
priceValue, err := xyJson.Get(root, "$.product.price")
if err != nil {
    return err
}

// 需要类型断言
scalarValue, ok := priceValue.(xyJson.IScalarValue)
if !ok {
    return errors.New("failed to cast to IScalarValue")
}

price, err := scalarValue.Float64()
if err != nil {
    return err
}

fmt.Printf("Price: %.2f\n", price)
```

## 解决方案

我们新增了四套便利API，满足不同的安全需求和使用场景：

### 1. Get系列方法
返回 `(值, error)` 格式，提供详细的错误信息：

```go
name, err := xyJson.GetString(root, "$.user.name")
if err != nil {
    fmt.Printf("获取失败: %v\n", err)
    return
}
```

### 2. TryGet系列方法 ⭐ 推荐
返回 `(值, bool)` 格式，最安全的选择，不会panic：

```go
if name, ok := xyJson.TryGetString(root, "$.user.name"); ok {
    fmt.Println("姓名:", name)
} else {
    fmt.Println("姓名不存在")
}
```

### 3. MustGet系列方法 ⚠️ 谨慎使用
直接返回值，失败时panic，仅在确信数据正确时使用：

```go
// 警告：失败时会panic
name := xyJson.MustGetString(root, "$.user.name")
```

### 4. GetWithDefault系列方法 ✨ 便利选择
失败时返回默认值，最适合处理可选字段：

```go
// 失败时返回默认值
name := xyJson.GetStringWithDefault(root, "$.user.name", "Unknown")
port := xyJson.GetIntWithDefault(root, "$.server.port", 8080)
```

## 可用的便利方法

### 完整方法列表

| 基础类型 | Get系列 | TryGet系列 | Must系列 | GetWithDefault系列 |
|---------|---------|------------|----------|--------------------|
| String | `GetString(root, path)` → `(string, error)` | `TryGetString(root, path)` → `(string, bool)` | `MustGetString(root, path)` → `string` | `GetStringWithDefault(root, path, defaultValue)` → `string` |
| Int | `GetInt(root, path)` → `(int, error)` | `TryGetInt(root, path)` → `(int, bool)` | `MustGetInt(root, path)` → `int` | `GetIntWithDefault(root, path, defaultValue)` → `int` |
| Int64 | `GetInt64(root, path)` → `(int64, error)` | `TryGetInt64(root, path)` → `(int64, bool)` | `MustGetInt64(root, path)` → `int64` | `GetInt64WithDefault(root, path, defaultValue)` → `int64` |
| Float64 | `GetFloat64(root, path)` → `(float64, error)` | `TryGetFloat64(root, path)` → `(float64, bool)` | `MustGetFloat64(root, path)` → `float64` | `GetFloat64WithDefault(root, path, defaultValue)` → `float64` |
| Bool | `GetBool(root, path)` → `(bool, error)` | `TryGetBool(root, path)` → `(bool, bool)` | `MustGetBool(root, path)` → `bool` | `GetBoolWithDefault(root, path, defaultValue)` → `bool` |
| Object | `GetObject(root, path)` → `(IObject, error)` | `TryGetObject(root, path)` → `(IObject, bool)` | `MustGetObject(root, path)` → `IObject` | `GetObjectWithDefault(root, path, defaultValue)` → `IObject` |
| Array | `GetArray(root, path)` → `(IArray, error)` | `TryGetArray(root, path)` → `(IArray, bool)` | `MustGetArray(root, path)` → `IArray` | `GetArrayWithDefault(root, path, defaultValue)` → `IArray` |

### 方法特点对比

| 特性 | Get系列 | TryGet系列 | Must系列 | GetWithDefault系列 |
|------|---------|------------|----------|--------------------|
| **安全性** | ✅ 安全 | ✅ 最安全 | ❌ 会panic | ✅ 安全 |
| **错误信息** | ✅ 详细 | ❌ 无详细信息 | ❌ 直接panic | ❌ 无详细信息 |
| **代码简洁性** | 🔶 中等 | ✅ 简洁 | ✅ 最简洁 | ✅ 最简洁 |
| **推荐场景** | 调试、详细错误处理 | 日常使用、生产环境 | 原型开发、确信数据正确 | 可选字段、配置默认值 |
| **失败处理** | 返回error | 返回false | panic | 返回默认值 |
| **零值返回** | 需检查error | 自动返回零值 | 不适用 | 返回指定默认值 |

## 使用示例

### 基本用法对比

```go
package main

import (
    "fmt"
    "log"
    xyJson "github/ihuem/xyJson"
)

func main() {
    data := `{
        "user": {
            "name": "Alice",
            "age": 30,
            "height": 165.5,
            "active": true,
            "salary": 75000.50,
            "profile": {
                "email": "alice@example.com"
            },
            "skills": ["Go", "Python", "JavaScript"]
        }
    }`

    root, err := xyJson.ParseString(data)
    if err != nil {
        log.Fatal(err)
    }

    // 1. Get系列方法 - 详细错误处理
    fmt.Println("=== Get系列方法 ===")
    name, err := xyJson.GetString(root, "$.user.name")
    if err != nil {
        fmt.Printf("获取姓名失败: %v\n", err)
        return
    }
    fmt.Printf("姓名: %s\n", name) // 姓名: Alice

    age, err := xyJson.GetInt(root, "$.user.age")
    if err != nil {
        fmt.Printf("获取年龄失败: %v\n", err)
        return
    }
    fmt.Printf("年龄: %d\n", age) // 年龄: 30

    // 2. TryGet系列方法 - 推荐使用
    fmt.Println("\n=== TryGet系列方法（推荐） ===")
    if height, ok := xyJson.TryGetFloat64(root, "$.user.height"); ok {
        fmt.Printf("身高: %.1f\n", height) // 身高: 165.5
    } else {
        fmt.Println("身高信息不存在")
    }

    if active, ok := xyJson.TryGetBool(root, "$.user.active"); ok {
        fmt.Printf("活跃状态: %t\n", active) // 活跃状态: true
    } else {
        fmt.Println("活跃状态不存在")
    }

    if salary, ok := xyJson.TryGetFloat64(root, "$.user.salary"); ok {
        fmt.Printf("薪资: %.2f\n", salary) // 薪资: 75000.50
    }

    // 处理不存在的字段
    if city, ok := xyJson.TryGetString(root, "$.user.city"); ok {
        fmt.Printf("城市: %s\n", city)
    } else {
        fmt.Println("城市信息不存在") // 这行会被执行
    }

    // 3. Must系列方法 - 谨慎使用
    fmt.Println("\n=== Must系列方法（谨慎使用） ===")
    // 仅在确信数据存在时使用
    userName := xyJson.MustGetString(root, "$.user.name")
    userAge := xyJson.MustGetInt(root, "$.user.age")
    fmt.Printf("用户: %s, %d岁\n", userName, userAge)

    // 获取复杂类型
    fmt.Println("\n=== 复杂类型处理 ===")
    if profile, ok := xyJson.TryGetObject(root, "$.user.profile"); ok {
        if email, ok := xyJson.TryGetString(profile, "$.email"); ok {
            fmt.Printf("邮箱: %s\n", email) // 邮箱: alice@example.com
        }
    }

    if skills, ok := xyJson.TryGetArray(root, "$.user.skills"); ok {
        fmt.Printf("技能数量: %d\n", skills.Length()) // 技能数量: 3
    }
}
```

### TryGet系列方法详细示例

```go
func demonstrateTryGetMethods(root xyJson.IValue) {
    fmt.Println("=== TryGet方法详细演示 ===")
    
    // 基本类型的安全获取
    if name, ok := xyJson.TryGetString(root, "$.user.name"); ok {
        fmt.Printf("✓ 姓名: %s\n", name)
    } else {
        fmt.Println("✗ 姓名不存在或类型错误")
    }
    
    if age, ok := xyJson.TryGetInt(root, "$.user.age"); ok {
        fmt.Printf("✓ 年龄: %d\n", age)
    } else {
        fmt.Println("✗ 年龄不存在或类型错误")
    }
    
    // 处理可能不存在的字段
    if phone, ok := xyJson.TryGetString(root, "$.user.phone"); ok {
        fmt.Printf("✓ 电话: %s\n", phone)
    } else {
        fmt.Println("✗ 电话号码未提供")
    }
    
    // 类型转换失败的情况
    if invalidAge, ok := xyJson.TryGetInt(root, "$.user.name"); ok {
        fmt.Printf("年龄: %d\n", invalidAge)
    } else {
        fmt.Println("✗ 无法将姓名转换为整数（预期行为）")
    }
    
    // 嵌套对象的安全访问
    if profile, ok := xyJson.TryGetObject(root, "$.user.profile"); ok {
        fmt.Println("✓ 找到用户档案")
        if email, ok := xyJson.TryGetString(profile, "$.email"); ok {
            fmt.Printf("  邮箱: %s\n", email)
        }
        if bio, ok := xyJson.TryGetString(profile, "$.bio"); ok {
            fmt.Printf("  简介: %s\n", bio)
        } else {
            fmt.Println("  简介未提供")
        }
    } else {
        fmt.Println("✗ 用户档案不存在")
    }
}

// 使用Must版本（适用于确信数据正确的场景）
func processUserDataWithMust(root xyJson.IValue) {
    // 当您确信这些路径存在且类型正确时，可以使用Must版本
    name := xyJson.MustGetString(root, "$.user.name")
    age := xyJson.MustGetInt(root, "$.user.age")
    height := xyJson.MustGetFloat64(root, "$.user.height")
    active := xyJson.MustGetBool(root, "$.user.active")

    fmt.Printf("User: %s, Age: %d, Height: %.1f, Active: %t\n", 
               name, age, height, active)
}
```

### 处理复杂数据结构

```go
func processComplexData(root xyJson.IValue) {
    fmt.Println("=== 复杂数据结构处理 ===")
    
    // 使用TryGet安全获取对象
    if profile, ok := xyJson.TryGetObject(root, "$.user.profile"); ok {
        fmt.Printf("✓ 用户档案包含 %d 个字段:\n", profile.Size())
        profile.Range(func(key string, value xyJson.IValue) bool {
            fmt.Printf("  %s: %s\n", key, value.String())
            return true
        })
    } else {
        fmt.Println("✗ 用户档案不存在")
    }

    // 使用TryGet安全获取数组
    if skills, ok := xyJson.TryGetArray(root, "$.user.skills"); ok {
        fmt.Printf("✓ 用户掌握 %d 项技能:\n", skills.Length())
        skills.Range(func(index int, value xyJson.IValue) bool {
            if skill, ok := xyJson.TryGetString(value, "$"); ok {
                fmt.Printf("  %d. %s\n", index+1, skill)
            }
            return true
        })
    } else {
        fmt.Println("✗ 技能列表不存在")
    }
    
    // 对比：使用Get方法处理相同数据
    fmt.Println("\n=== 使用Get方法对比 ===")
    profile, err := xyJson.GetObject(root, "$.user.profile")
    if err != nil {
        fmt.Printf("获取档案失败: %v\n", err)
    } else {
        fmt.Printf("档案字段数: %d\n", profile.Size())
    }

    hobbies, err := xyJson.GetArray(root, "$.user.hobbies")
    if err != nil {
        fmt.Printf("获取爱好失败: %v\n", err)
    } else {
        fmt.Printf("爱好数量: %d\n", hobbies.Length())
    }
}
```

## 错误处理

便利API会返回以下类型的错误：

1. **路径不存在错误**：当指定的JSONPath不存在时
2. **类型转换错误**：当值无法转换为请求的类型时

### 三种方法的错误处理对比

```go
func demonstrateErrorHandling(root xyJson.IValue) {
    fmt.Println("=== 错误处理演示 ===")
    
    // 1. Get方法 - 详细错误信息
    if _, err := xyJson.GetString(root, "$.nonexistent.path"); err != nil {
        fmt.Printf("Get方法错误: %v\n", err)
    }
    
    if _, err := xyJson.GetInt(root, "$.user.name"); err != nil {
        fmt.Printf("类型转换错误: %v\n", err)
    }
    
    // 2. TryGet方法 - 简洁的布尔返回
    if value, ok := xyJson.TryGetString(root, "$.nonexistent.path"); ok {
        fmt.Printf("值: %s\n", value)
    } else {
        fmt.Println("TryGet: 路径不存在或类型错误")
    }
    
    if value, ok := xyJson.TryGetInt(root, "$.user.name"); ok {
        fmt.Printf("值: %d\n", value)
    } else {
        fmt.Println("TryGet: 无法转换为整数")
    }
    
    // 3. Must方法 - 会panic（仅演示，实际使用需谨慎）
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Must方法panic: %v\n", r)
        }
    }()
    
    // 这行会导致panic
    // value := xyJson.MustGetString(root, "$.nonexistent.path")
    fmt.Println("Must方法：仅在确信数据正确时使用")
}

// 实际项目中的错误处理模式
func practicalErrorHandling(root xyJson.IValue) error {
    // 推荐：使用TryGet进行安全访问
    name, ok := xyJson.TryGetString(root, "$.user.name")
    if !ok {
        return fmt.Errorf("用户姓名缺失或格式错误")
    }
    
    age, ok := xyJson.TryGetInt(root, "$.user.age")
    if !ok {
        return fmt.Errorf("用户年龄缺失或格式错误")
    }
    
    // 可选字段使用默认值
    email := "未提供"
    if e, ok := xyJson.TryGetString(root, "$.user.email"); ok {
        email = e
    }
    
    fmt.Printf("处理用户: %s, %d岁, 邮箱: %s\n", name, age, email)
    return nil
}
```

## 性能考虑

- 便利API在内部调用原始的 `Get` 方法，然后进行类型转换
- 性能开销主要来自类型转换，通常是可以忽略的
- Must版本的方法在性能上与普通版本相同，只是错误处理方式不同

## 最佳实践

### 1. 选择合适的方法版本

- **TryGet系列**：⭐ 推荐用于大多数场景，安全且简洁
- **Get系列**：用于需要详细错误信息的场景
- **Must系列**：⚠️ 仅用于确信数据正确的场景（如配置文件解析）

### 2. 推荐的使用模式

```go
// 最佳实践：优先使用TryGet
func processUserInfo(root xyJson.IValue) error {
    // 必需字段验证
    name, ok := xyJson.TryGetString(root, "$.user.name")
    if !ok {
        return fmt.Errorf("用户姓名是必需的")
    }
    
    age, ok := xyJson.TryGetInt(root, "$.user.age")
    if !ok {
        return fmt.Errorf("用户年龄是必需的")
    }
    
    // 可选字段处理
    email := "未提供"
    if e, ok := xyJson.TryGetString(root, "$.user.email"); ok {
        email = e
    }
    
    // 带默认值的字段
    active := true // 默认值
    if a, ok := xyJson.TryGetBool(root, "$.user.active"); ok {
        active = a
    }
    
    fmt.Printf("用户: %s, %d岁, 邮箱: %s, 活跃: %v\n", name, age, email, active)
    return nil
}

// 需要详细错误信息时使用Get
func validateUserData(root xyJson.IValue) error {
    if _, err := xyJson.GetString(root, "$.user.name"); err != nil {
        return fmt.Errorf("姓名验证失败: %w", err)
    }
    
    if _, err := xyJson.GetInt(root, "$.user.age"); err != nil {
        return fmt.Errorf("年龄验证失败: %w", err)
    }
    
    return nil
}

// 配置文件等确信数据正确的场景
func loadConfig(root xyJson.IValue) {
    // 配置文件通常结构固定，可以使用Must
    appName := xyJson.MustGetString(root, "$.app.name")
    port := xyJson.MustGetInt(root, "$.server.port")
    debug := xyJson.MustGetBool(root, "$.debug")
    
    fmt.Printf("应用: %s, 端口: %d, 调试: %v\n", appName, port, debug)
}
```

### 3. 类型安全

```go
// ✅ 推荐：使用TryGet类型特定方法
if age, ok := xyJson.TryGetInt(root, "$.user.age"); ok {
    fmt.Printf("年龄: %d\n", age)
}

// ✅ 可以：使用Get方法处理错误
age, err := xyJson.GetInt(root, "$.user.age")
if err != nil {
    return fmt.Errorf("获取年龄失败: %w", err)
}

// ❌ 避免：手动类型断言
ageValue, _ := xyJson.Get(root, "$.user.age")
scalar, _ := ageValue.(xyJson.IScalarValue)
age, _ := scalar.Int()
```

### 4. 性能优化建议

```go
// 批量获取时，先获取父对象
if user, ok := xyJson.TryGetObject(root, "$.user"); ok {
    // 在子对象上操作，避免重复路径解析
    name, _ := xyJson.TryGetString(user, "$.name")
    age, _ := xyJson.TryGetInt(user, "$.age")
    email, _ := xyJson.TryGetString(user, "$.email")
    
    // 处理数据...
}
```

## 兼容性

- 新的便利API与现有的API完全兼容
- 您可以在同一个项目中混合使用新旧API
- 现有代码无需修改即可继续工作

## 总结

新的便利API提供了以下优势：

1. **简化代码**：无需手动类型断言
2. **提高可读性**：代码意图更加明确
3. **类型安全**：编译时类型检查
4. **多种选择**：Get、TryGet、Must三种风格满足不同需求
5. **兼容性**：与现有API完全兼容
6. **安全性**：TryGet方法提供最安全的访问方式

### 选择指南

| 场景 | 推荐方法 | 原因 |
|------|----------|------|
| 日常开发 | **TryGet系列** | 安全、简洁、不会panic |
| 调试分析 | Get系列 | 提供详细错误信息 |
| 配置解析 | Must系列 | 结构固定，失败应该终止程序 |
| 可选字段 | TryGet系列 | 可以优雅地处理缺失字段 |
| 数据验证 | Get系列 | 需要具体的错误信息 |

### 迁移建议

```go
// 旧代码
value, err := xyJson.Get(root, "$.user.name")
if err != nil {
    return err
}
scalar, ok := value.(xyJson.IScalarValue)
if !ok {
    return errors.New("not a scalar")
}
name, err := scalar.String()
if err != nil {
    return err
}

// 新代码（推荐）
if name, ok := xyJson.TryGetString(root, "$.user.name"); ok {
    // 使用name
} else {
    // 处理缺失或错误
}

// 或者
name, err := xyJson.GetString(root, "$.user.name")
if err != nil {
    return fmt.Errorf("获取用户姓名失败: %w", err)
}
```

通过使用这些便利方法，您可以编写更简洁、更安全、更易维护的JSON处理代码。特别推荐在新项目中优先使用 **TryGet系列方法**。

### GetWithDefault系列方法详细示例 ✨ 新增

`GetWithDefault`系列方法是最新添加的便利API，专门用于处理可选字段和提供默认值的场景。

```go
func demonstrateGetWithDefault(root xyJson.IValue) {
    fmt.Println("=== GetWithDefault方法演示 ===\n")
    
    // 1. 基本类型的默认值处理
    fmt.Println("1. 基本类型默认值:")
    name := xyJson.GetStringWithDefault(root, "$.user.name", "Unknown")
    age := xyJson.GetIntWithDefault(root, "$.user.age", 0)
    height := xyJson.GetFloat64WithDefault(root, "$.user.height", 170.0)
    active := xyJson.GetBoolWithDefault(root, "$.user.active", true)
    
    fmt.Printf("姓名: %s\n", name)
    fmt.Printf("年龄: %d\n", age)
    fmt.Printf("身高: %.1f\n", height)
    fmt.Printf("活跃: %t\n", active)
    
    // 2. 配置读取场景（最佳用例）
    fmt.Println("\n2. 配置读取场景:")
    serverHost := xyJson.GetStringWithDefault(root, "$.server.host", "localhost")
    serverPort := xyJson.GetIntWithDefault(root, "$.server.port", 8080)
    maxConnections := xyJson.GetIntWithDefault(root, "$.server.maxConnections", 100)
    sslEnabled := xyJson.GetBoolWithDefault(root, "$.server.ssl", false)
    timeout := xyJson.GetFloat64WithDefault(root, "$.server.timeout", 30.0)
    
    fmt.Printf("服务器配置:\n")
    fmt.Printf("  主机: %s\n", serverHost)
    fmt.Printf("  端口: %d\n", serverPort)
    fmt.Printf("  最大连接: %d\n", maxConnections)
    fmt.Printf("  SSL: %t\n", sslEnabled)
    fmt.Printf("  超时: %.1f秒\n", timeout)
    
    // 3. 复杂类型的默认值
    fmt.Println("\n3. 复杂类型默认值:")
    
    // 创建默认对象
    defaultConfig := xyJson.CreateObject()
    defaultConfig.Set("host", xyJson.CreateString("localhost"))
    defaultConfig.Set("port", xyJson.CreateNumber(5432))
    
    dbConfig := xyJson.GetObjectWithDefault(root, "$.database", defaultConfig)
    fmt.Printf("数据库配置大小: %d\n", dbConfig.Size())
    
    // 当defaultValue为nil时，返回空对象/数组
    emptyObj := xyJson.GetObjectWithDefault(root, "$.missing.object", nil)
    emptyArr := xyJson.GetArrayWithDefault(root, "$.missing.array", nil)
    fmt.Printf("空对象大小: %d\n", emptyObj.Size())
    fmt.Printf("空数组长度: %d\n", emptyArr.Length())
}

// 对比三种方法的代码简洁性
func compareMethodsSimplicity(root xyJson.IValue) {
    fmt.Println("=== 代码简洁性对比 ===\n")
    
    // 场景：获取服务器端口，默认值8080
    
    // 方法1：Get方法（最繁琐）
    fmt.Println("方法1 - Get方法:")
    port1 := 8080 // 默认值
    if p, err := xyJson.GetInt(root, "$.server.port"); err == nil {
        port1 = p
    }
    fmt.Printf("端口: %d\n", port1)
    
    // 方法2：TryGet方法（中等复杂度）
    fmt.Println("\n方法2 - TryGet方法:")
    port2 := 8080 // 默认值
    if p, ok := xyJson.TryGetInt(root, "$.server.port"); ok {
        port2 = p
    }
    fmt.Printf("端口: %d\n", port2)
    
    // 方法3：GetWithDefault方法（最简洁）
    fmt.Println("\n方法3 - GetWithDefault方法:")
    port3 := xyJson.GetIntWithDefault(root, "$.server.port", 8080)
    fmt.Printf("端口: %d\n", port3)
    
    fmt.Println("\n✨ GetWithDefault方法最简洁，只需一行代码！")
}

// 实际应用场景：Web服务器配置
func loadWebServerConfig(root xyJson.IValue) {
    fmt.Println("=== Web服务器配置加载 ===\n")
    
    // 使用GetWithDefault加载配置，代码简洁且安全
    config := struct {
        Host           string
        Port           int
        SSL            bool
        MaxConnections int
        Timeout        float64
        Debug          bool
        LogLevel       string
    }{
        Host:           xyJson.GetStringWithDefault(root, "$.server.host", "0.0.0.0"),
        Port:           xyJson.GetIntWithDefault(root, "$.server.port", 8080),
        SSL:            xyJson.GetBoolWithDefault(root, "$.server.ssl", false),
        MaxConnections: xyJson.GetIntWithDefault(root, "$.server.maxConnections", 1000),
        Timeout:        xyJson.GetFloat64WithDefault(root, "$.server.timeout", 30.0),
        Debug:          xyJson.GetBoolWithDefault(root, "$.debug", false),
        LogLevel:       xyJson.GetStringWithDefault(root, "$.logging.level", "info"),
    }
    
    fmt.Printf("服务器配置:\n")
    fmt.Printf("  监听地址: %s:%d\n", config.Host, config.Port)
    fmt.Printf("  SSL启用: %t\n", config.SSL)
    fmt.Printf("  最大连接: %d\n", config.MaxConnections)
    fmt.Printf("  超时时间: %.1f秒\n", config.Timeout)
    fmt.Printf("  调试模式: %t\n", config.Debug)
    fmt.Printf("  日志级别: %s\n", config.LogLevel)
}
```

### GetWithDefault方法的优势

1. **代码最简洁**：只需一行代码，无需if判断
2. **类型安全**：直接返回正确类型，无需类型断言
3. **默认值灵活**：可以指定任意合理的默认值
4. **特别适合配置**：配置文件读取的最佳选择
5. **无panic风险**：安全可靠，不会导致程序崩溃
6. **语义清晰**：代码意图一目了然

### 四种方法选择指南

| 场景 | 推荐方法 | 原因 |
|------|----------|------|
| 配置文件读取 | `GetWithDefault` | 代码最简洁，支持默认值 |
| 可选字段处理 | `GetWithDefault` | 无需判断，直接使用默认值 |
| 日常开发 | `TryGet` | 安全可靠，代码简洁 |
| 错误调试 | `Get` | 提供详细错误信息 |
| 确信数据正确 | `Must` | 代码最简洁，但有panic风险 |

## JSONPath预编译功能 🚀 新增

### 概述

JSONPath预编译功能是xyJson库的重要性能优化特性，通过预编译JSONPath表达式，可以显著提升重复查询的性能。当您需要多次使用相同的JSONPath表达式时，预编译功能可以带来约58%的性能提升。

### 问题背景

在传统的JSONPath查询中，每次调用都需要重新解析路径表达式：

```go
// 传统方式：每次都要解析路径
for i := 0; i < 1000; i++ {
    name, _ := xyJson.GetString(root, "$.user.name")  // 每次都解析"$.user.name"
    age, _ := xyJson.GetInt(root, "$.user.age")      // 每次都解析"$.user.age"
    // 处理数据...
}
```

这种方式在大量重复查询时会产生不必要的性能开销。

### 解决方案

预编译功能允许您一次编译路径，多次使用：

```go
// 预编译方式：一次编译，多次使用
namePath := xyJson.CompilePath("$.user.name")
agePath := xyJson.CompilePath("$.user.age")

for i := 0; i < 1000; i++ {
    name, _ := namePath.Query(root)  // 直接使用预编译的路径
    age, _ := agePath.Query(root)    // 直接使用预编译的路径
    // 处理数据...
}
```

### 核心API

#### 1. 路径编译

```go
// 编译JSONPath表达式
func CompilePath(path string) (*CompiledPath, error)

// 使用指定工厂编译路径
func CompilePathWithFactory(path string, factory IValueFactory) (*CompiledPath, error)
```

#### 2. CompiledPath方法

```go
type CompiledPath struct {
    // 私有字段...
}

// 查询方法
func (cp *CompiledPath) Query(root IValue) (IValue, error)           // 查询单个值
func (cp *CompiledPath) QueryAll(root IValue) ([]IValue, error)      // 查询所有匹配值

// 修改方法
func (cp *CompiledPath) Set(root IValue, value IValue) error         // 设置值
func (cp *CompiledPath) Delete(root IValue) error                    // 删除值

// 检查方法
func (cp *CompiledPath) Exists(root IValue) bool                     // 检查路径是否存在
func (cp *CompiledPath) Count(root IValue) int                       // 计算匹配数量

// 工具方法
func (cp *CompiledPath) Path() string                                // 获取原始路径字符串
```

#### 3. 缓存管理

```go
// 清空路径缓存
func ClearPathCache()

// 获取缓存统计信息
func GetPathCacheStats() (hits, misses, size int)

// 设置缓存最大大小
func SetPathCacheMaxSize(size int)
```

### 使用示例

#### 基本用法

```go
package main

import (
    "fmt"
    "log"
    xyJson "github/ihuem/xyJson"
)

func main() {
    data := `{
        "users": [
            {"name": "Alice", "age": 30, "email": "alice@example.com"},
            {"name": "Bob", "age": 25, "email": "bob@example.com"},
            {"name": "Charlie", "age": 35, "email": "charlie@example.com"}
        ]
    }`

    root, err := xyJson.ParseString(data)
    if err != nil {
        log.Fatal(err)
    }

    // 1. 编译常用路径
    fmt.Println("=== 编译JSONPath ===\n")
    
    userNamesPath, err := xyJson.CompilePath("$.users[*].name")
    if err != nil {
        log.Fatal(err)
    }
    
    userAgesPath, err := xyJson.CompilePath("$.users[*].age")
    if err != nil {
        log.Fatal(err)
    }
    
    firstUserPath, err := xyJson.CompilePath("$.users[0]")
    if err != nil {
        log.Fatal(err)
    }

    // 2. 使用预编译路径查询
    fmt.Println("=== 使用预编译路径查询 ===\n")
    
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
    
    // 查询所有用户年龄
    ages, err := userAgesPath.QueryAll(root)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("\n所有用户年龄:")
    for i, ageValue := range ages {
        if age, ok := xyJson.TryGetInt(ageValue, "$"); ok {
            fmt.Printf("  %d. %d岁\n", i+1, age)
        }
    }
    
    // 查询第一个用户
    firstUser, err := firstUserPath.Query(root)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("\n第一个用户信息:")
    if name, ok := xyJson.TryGetString(firstUser, "$.name"); ok {
        fmt.Printf("  姓名: %s\n", name)
    }
    if age, ok := xyJson.TryGetInt(firstUser, "$.age"); ok {
        fmt.Printf("  年龄: %d\n", age)
    }
    if email, ok := xyJson.TryGetString(firstUser, "$.email"); ok {
        fmt.Printf("  邮箱: %s\n", email)
    }
}
```

#### 高级用法：批量数据处理

```go
func processBatchData(root xyJson.IValue) {
    fmt.Println("=== 批量数据处理示例 ===\n")
    
    // 预编译常用路径
    paths := map[string]*xyJson.CompiledPath{
        "userNames":  xyJson.MustCompilePath("$.users[*].name"),
        "userAges":   xyJson.MustCompilePath("$.users[*].age"),
        "userEmails": xyJson.MustCompilePath("$.users[*].email"),
        "activeUsers": xyJson.MustCompilePath("$.users[?(@.active == true)]"),
        "adminUsers":  xyJson.MustCompilePath("$.users[?(@.role == 'admin')]"),
    }
    
    // 批量查询
    results := make(map[string][]xyJson.IValue)
    for name, path := range paths {
        values, err := path.QueryAll(root)
        if err != nil {
            fmt.Printf("查询 %s 失败: %v\n", name, err)
            continue
        }
        results[name] = values
        fmt.Printf("%s: 找到 %d 个结果\n", name, len(values))
    }
    
    // 处理结果
    if names, ok := results["userNames"]; ok {
        fmt.Println("\n用户列表:")
        for i, nameValue := range names {
            if name, ok := xyJson.TryGetString(nameValue, "$"); ok {
                fmt.Printf("  %d. %s\n", i+1, name)
            }
        }
    }
}

// 便利方法：MustCompilePath
func xyJson.MustCompilePath(path string) *CompiledPath {
    compiled, err := xyJson.CompilePath(path)
    if err != nil {
        panic(fmt.Sprintf("编译路径失败: %v", err))
    }
    return compiled
}
```

#### 修改操作示例

```go
func demonstrateModificationOperations(root xyJson.IValue) {
    fmt.Println("=== 修改操作示例 ===\n")
    
    // 编译修改路径
    firstUserAgePath, _ := xyJson.CompilePath("$.users[0].age")
    firstUserEmailPath, _ := xyJson.CompilePath("$.users[0].email")
    newUserPath, _ := xyJson.CompilePath("$.users[3]")
    
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
    newUser.Set("name", xyJson.CreateString("David"))
    newUser.Set("age", xyJson.CreateNumber(28))
    newUser.Set("email", xyJson.CreateString("david@example.com"))
    
    if err := newUserPath.Set(root, newUser); err != nil {
        fmt.Printf("添加失败: %v\n", err)
    } else {
        fmt.Println("✓ 新用户添加成功")
    }
    
    // 4. 检查操作结果
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
}
```

#### 缓存管理示例

```go
func demonstrateCacheManagement() {
    fmt.Println("=== 缓存管理示例 ===\n")
    
    // 1. 查看初始缓存状态
    hits, misses, size := xyJson.GetPathCacheStats()
    fmt.Printf("初始缓存状态: 命中=%d, 未命中=%d, 大小=%d\n", hits, misses, size)
    
    // 2. 编译一些路径（会被缓存）
    paths := []string{
        "$.user.name",
        "$.user.age",
        "$.user.email",
        "$.users[*].name",
        "$.users[0].profile",
    }
    
    fmt.Println("\n编译路径（首次编译，会缓存）:")
    compiledPaths := make([]*xyJson.CompiledPath, len(paths))
    for i, path := range paths {
        compiled, err := xyJson.CompilePath(path)
        if err != nil {
            fmt.Printf("编译 %s 失败: %v\n", path, err)
            continue
        }
        compiledPaths[i] = compiled
        fmt.Printf("✓ 编译: %s\n", path)
    }
    
    // 3. 查看缓存状态
    hits, misses, size = xyJson.GetPathCacheStats()
    fmt.Printf("\n编译后缓存状态: 命中=%d, 未命中=%d, 大小=%d\n", hits, misses, size)
    
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
    hits, misses, size = xyJson.GetPathCacheStats()
    fmt.Printf("\n最终缓存状态: 命中=%d, 未命中=%d, 大小=%d\n", hits, misses, size)
    
    // 6. 设置缓存大小限制
    fmt.Println("\n设置缓存最大大小为3:")
    xyJson.SetPathCacheMaxSize(3)
    
    // 7. 编译更多路径，触发缓存清理
    morePaths := []string{
        "$.config.database.host",
        "$.config.database.port",
        "$.config.redis.host",
        "$.config.redis.port",
    }
    
    for _, path := range morePaths {
        xyJson.CompilePath(path)
        hits, misses, size = xyJson.GetPathCacheStats()
        fmt.Printf("编译 %s 后缓存大小: %d\n", path, size)
    }
    
    // 8. 清空缓存
    fmt.Println("\n清空缓存:")
    xyJson.ClearPathCache()
    hits, misses, size = xyJson.GetPathCacheStats()
    fmt.Printf("清空后缓存状态: 命中=%d, 未命中=%d, 大小=%d\n", hits, misses, size)
}
```

### 性能对比

基准测试结果显示，预编译路径在重复查询场景下具有显著的性能优势：

```
BenchmarkCompiledPathVsRegular/Compiled_Path-8         	 2264063	   529.3 ns/op
BenchmarkCompiledPathVsRegular/Regular_Path-8          	  945694	  1267 ns/op
BenchmarkCompiledPathVsRegular/Compiled_Path_with_Compilation-8 	  106134	 11282 ns/op

BenchmarkPathCachePerformance/Cache_Miss-8             	  105129	 11406 ns/op
BenchmarkPathCachePerformance/Cache_Hit-8              	  125127	  9584 ns/op
```

**性能分析：**
- **预编译路径 vs 常规路径**：约58%的性能提升（529.3ns vs 1267ns）
- **缓存命中 vs 缓存未命中**：约16%的性能提升（9584ns vs 11406ns）
- **编译开销**：首次编译需要额外时间（11282ns），但在重复使用时迅速摊销

### 最佳实践

#### 1. 何时使用预编译

✅ **推荐使用场景：**
- 重复查询相同路径（循环处理、批量操作）
- 性能敏感的应用
- 固定的JSONPath表达式
- 长时间运行的服务

❌ **不推荐使用场景：**
- 一次性查询
- 动态生成的路径
- 内存受限的环境
- 路径表达式经常变化

#### 2. 缓存管理策略

```go
// 应用启动时设置合理的缓存大小
func init() {
    // 根据应用规模设置缓存大小
    xyJson.SetPathCacheMaxSize(100) // 缓存100个常用路径
}

// 定期清理缓存（可选）
func periodicCacheCleanup() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        hits, misses, size := xyJson.GetPathCacheStats()
        
        // 如果命中率过低，清理缓存
        if hits > 0 && float64(hits)/(float64(hits+misses)) < 0.5 {
            xyJson.ClearPathCache()
            log.Println("缓存命中率过低，已清理缓存")
        }
        
        log.Printf("缓存统计: 命中=%d, 未命中=%d, 大小=%d\n", hits, misses, size)
    }
}
```

#### 3. 错误处理

```go
// 安全的路径编译
func safeCompilePath(path string) (*xyJson.CompiledPath, error) {
    compiled, err := xyJson.CompilePath(path)
    if err != nil {
        return nil, fmt.Errorf("编译JSONPath失败 '%s': %w", path, err)
    }
    return compiled, nil
}

// 批量编译路径
func compileMultiplePaths(paths []string) (map[string]*xyJson.CompiledPath, error) {
    compiled := make(map[string]*xyJson.CompiledPath)
    
    for _, path := range paths {
        cp, err := safeCompilePath(path)
        if err != nil {
            return nil, err
        }
        compiled[path] = cp
    }
    
    return compiled, nil
}
```

#### 4. 与便利API结合使用

```go
// 结合预编译路径和便利API
func efficientDataProcessing(root xyJson.IValue) {
    // 预编译常用路径
    userPath := xyJson.MustCompilePath("$.users[0]")
    
    // 查询用户对象
    user, err := userPath.Query(root)
    if err != nil {
        log.Printf("查询用户失败: %v", err)
        return
    }
    
    // 在用户对象上使用便利API
    name := xyJson.GetStringWithDefault(user, "$.name", "Unknown")
    age := xyJson.GetIntWithDefault(user, "$.age", 0)
    email := xyJson.GetStringWithDefault(user, "$.email", "")
    
    fmt.Printf("用户信息: %s, %d岁, %s\n", name, age, email)
}
```

### 技术实现细节

#### 1. 线程安全
- `CompiledPath` 结构体是线程安全的，可以在多个goroutine中并发使用
- 内置缓存使用读写锁保护，支持并发访问
- 所有公共方法都是线程安全的

#### 2. 内存管理
- 缓存使用LRU策略，自动清理最少使用的条目
- 支持设置最大缓存大小，防止内存泄漏
- `CompiledPath` 对象可以安全地被垃圾回收

#### 3. 向后兼容
- 预编译功能完全向后兼容
- 现有代码无需修改即可继续工作
- 可以渐进式地迁移到预编译API

### 总结

JSONPath预编译功能为xyJson库带来了显著的性能提升，特别适合需要重复查询的场景。通过合理使用预编译功能和缓存管理，您可以：

1. **提升性能**：重复查询性能提升约58%
2. **简化代码**：一次编译，多次使用
3. **节省资源**：避免重复解析开销
4. **保持安全**：线程安全的设计
5. **易于维护**：清晰的API设计

**推荐使用模式：**
- 对于重复查询，优先使用预编译路径
- 合理设置缓存大小，平衡内存和性能
- 结合便利API使用，获得最佳开发体验
- 在性能敏感的场景中，预编译是必备选择

## 相关文档

- [JSONPath预编译功能详细指南](compiled_path.md) - 深入了解预编译功能的技术原理和最佳实践
- [性能优化指南](performance_guide.md) - 全面的性能优化建议
- [API参考文档](api_reference.md) - 完整的API文档