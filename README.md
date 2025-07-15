# xyJson

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)
[![Coverage](https://img.shields.io/badge/coverage-95%25-brightgreen.svg)](#)

高性能Go JSON库，专为现代应用程序设计，提供极致的性能和丰富的功能。

## 🚀 核心特性

### 🔥 极致性能
- **零拷贝解析**：最小化内存分配和数据复制
- **对象池优化**：自动回收和重用对象，减少GC压力
- **并发安全**：内置读写锁，支持高并发访问
- **内存优化**：智能内存管理，减少内存碎片

### 🎯 JSONPath查询
- **完整支持**：符合JSONPath规范的查询语法
- **高级过滤**：支持复杂的条件过滤和表达式
- **批量操作**：一次查询获取多个结果
- **路径缓存**：自动缓存常用路径，提升查询性能

### 🛡️ 类型安全
- **强类型转换**：安全的类型转换，避免运行时错误
- **溢出检测**：自动检测数值溢出和精度丢失
- **空值处理**：优雅处理null值和未定义字段
- **错误恢复**：详细的错误信息和恢复机制

### 📊 性能监控
- **实时统计**：解析、序列化操作的实时性能数据
- **内存分析**：详细的内存使用情况和分配统计
- **性能基准**：内置基准测试和性能对比
- **可视化报告**：生成详细的性能分析报告

### ⚡ 并发安全
- **读写锁**：细粒度的并发控制
- **原子操作**：关键路径使用原子操作
- **无锁设计**：部分操作采用无锁算法
- **竞态检测**：内置竞态条件检测

### ⚙️ 可配置性
- **灵活配置**：支持多种使用场景的配置
- **环境适配**：开发、测试、生产环境的不同配置
- **性能调优**：可调节的性能参数
- **插件扩展**：支持自定义扩展和插件

## 📦 安装

```bash
go get github.com/ihuem/xyJson
```

## 🚀 快速开始

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/ihuem/xyJson"
)

func main() {
    // 解析JSON字符串
    jsonStr := `{"name":"张三","age":30,"city":"北京"}`
    value, err := xyJson.ParseString(jsonStr)
    if err != nil {
        panic(err)
    }

    // 获取值
    name := xyJson.MustGet(value, "$.name").String()
    fmt.Printf("姓名: %s\n", name)

    // 修改值
    xyJson.Set(value, "$.age", xyJson.CreateNumber(31))

    // 序列化回JSON
    result, _ := xyJson.SerializeToString(value)
    fmt.Printf("修改后: %s\n", result)
}
```

## 📚 更多功能

详细的文档和示例正在完善中...

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

---

⭐ 如果这个项目对你有帮助，请给我们一个星标！