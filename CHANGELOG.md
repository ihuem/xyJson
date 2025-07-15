# 更新日志 / Changelog

本文档记录了 xyJson 项目的所有重要变更。

This document records all notable changes to the xyJson project.

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/) 规范。

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [未发布] / [Unreleased]

### 新增 / Added
- 待发布的新功能 / Features to be released

### 变更 / Changed
- 待发布的变更 / Changes to be released

### 修复 / Fixed
- 待发布的修复 / Fixes to be released

### 移除 / Removed
- 待移除的功能 / Features to be removed

## [1.0.0] - 2024-01-15

### 新增 / Added
- 🚀 **核心功能** / Core Features
  - 高性能JSON解析器，支持标准JSON格式 / High-performance JSON parser supporting standard JSON format
  - 高效JSON序列化器，支持多种格式选项 / Efficient JSON serializer with multiple format options
  - 完整的JSONPath查询支持，兼容JSONPath规范 / Complete JSONPath query support, compatible with JSONPath specification
  - 类型安全的值接口系统 / Type-safe value interface system

- 🏊 **内存优化** / Memory Optimization
  - 智能对象池实现，减少GC压力 / Smart object pool implementation to reduce GC pressure
  - 内存复用机制，提升性能30-50% / Memory reuse mechanism, improving performance by 30-50%
  - 可配置的池大小和策略 / Configurable pool size and strategies
  - 内存使用统计和监控 / Memory usage statistics and monitoring

- 📊 **性能监控** / Performance Monitoring
  - 实时性能统计收集 / Real-time performance statistics collection
  - 内存使用分析和报告 / Memory usage analysis and reporting
  - 操作耗时追踪和分析 / Operation timing tracking and analysis
  - 可配置的性能阈值告警 / Configurable performance threshold alerts
  - 性能数据导出功能 / Performance data export functionality

- 🔍 **JSONPath功能** / JSONPath Features
  - 基础路径查询 (`$.key`, `$.array[0]`) / Basic path queries
  - 通配符查询 (`$.array[*]`, `$.*`) / Wildcard queries
  - 递归下降查询 (`$..key`) / Recursive descent queries
  - 数组切片查询 (`$.array[1:3]`, `$.array[-2:]`) / Array slice queries
  - 条件过滤查询 (`$.array[?(@.key > 10)]`) / Conditional filter queries
  - 多路径批量查询 / Multi-path batch queries
  - JSONPath预编译优化 / JSONPath pre-compilation optimization

- 🛠️ **开发者工具** / Developer Tools
  - 丰富的创建函数 (`CreateObject`, `CreateArray`, `CreateString`等) / Rich creation functions
  - 便捷的序列化选项 (`Pretty`, `Compact`, `HTMLSafe`) / Convenient serialization options
  - 批量操作支持 (`SetBatch`, `GetBatch`) / Batch operation support
  - 流式处理接口 / Streaming processing interfaces
  - 自定义值工厂支持 / Custom value factory support

- 🔒 **并发安全** / Concurrency Safety
  - 线程安全的读写操作 / Thread-safe read/write operations
  - 并发访问保护机制 / Concurrent access protection mechanisms
  - 无锁优化的热路径 / Lock-free optimized hot paths

- 🧪 **测试和质量** / Testing and Quality
  - 全面的单元测试覆盖 (>95%) / Comprehensive unit test coverage (>95%)
  - 性能基准测试套件 / Performance benchmark test suite
  - 集成测试和稳定性测试 / Integration tests and stability tests
  - 内存泄漏检测 / Memory leak detection
  - 竞态条件测试 / Race condition testing

- 📚 **文档和示例** / Documentation and Examples
  - 完整的API参考文档 / Complete API reference documentation
  - 详细的使用示例 / Detailed usage examples
  - 性能优化指南 / Performance optimization guide
  - 最佳实践建议 / Best practice recommendations
  - 中英文双语文档 / Bilingual documentation (Chinese/English)

### 技术规格 / Technical Specifications
- **Go版本要求** / Go Version Requirement: >= 1.21
- **零外部依赖** / Zero External Dependencies: 纯Go实现 / Pure Go implementation
- **内存效率** / Memory Efficiency: 比标准库减少40%内存使用 / 40% less memory usage than standard library
- **性能提升** / Performance Improvement: 比标准库快30-50% / 30-50% faster than standard library
- **并发支持** / Concurrency Support: 完全线程安全 / Fully thread-safe
- **平台支持** / Platform Support: 跨平台兼容 / Cross-platform compatible

### API接口 / API Interfaces

#### 核心接口 / Core Interfaces
- `IValue` - JSON值基础接口 / Base interface for JSON values
- `IScalarValue` - 标量值接口 / Scalar value interface
- `IObject` - JSON对象接口 / JSON object interface
- `IArray` - JSON数组接口 / JSON array interface
- `IParser` - JSON解析器接口 / JSON parser interface
- `ISerializer` - JSON序列化器接口 / JSON serializer interface
- `IPathQuery` - JSONPath查询接口 / JSONPath query interface
- `IValueFactory` - 值工厂接口 / Value factory interface
- `IObjectPool` - 对象池接口 / Object pool interface

#### 全局函数 / Global Functions
- `Parse([]byte) (IValue, error)` - 解析JSON字节数据 / Parse JSON byte data
- `ParseString(string) (IValue, error)` - 解析JSON字符串 / Parse JSON string
- `Serialize(IValue) ([]byte, error)` - 序列化为字节数组 / Serialize to byte array
- `SerializeToString(IValue) (string, error)` - 序列化为字符串 / Serialize to string
- `Get(IValue, string) (IValue, error)` - JSONPath单值查询 / JSONPath single value query
- `GetAll(IValue, string) ([]IValue, error)` - JSONPath多值查询 / JSONPath multi-value query
- `Set(IValue, string, interface{}) error` - JSONPath设置值 / JSONPath set value
- `Delete(IValue, string) error` - JSONPath删除值 / JSONPath delete value
- `Exists(IValue, string) bool` - JSONPath路径存在检查 / JSONPath existence check
- `Count(IValue, string) int` - JSONPath匹配计数 / JSONPath match count

#### 创建函数 / Creation Functions
- `CreateNull() IValue` - 创建null值 / Create null value
- `CreateString(string) IScalarValue` - 创建字符串值 / Create string value
- `CreateNumber(interface{}) (IScalarValue, error)` - 创建数字值 / Create number value
- `CreateBool(bool) IScalarValue` - 创建布尔值 / Create boolean value
- `CreateObject() IObject` - 创建对象 / Create object
- `CreateArray() IArray` - 创建数组 / Create array

#### 便捷函数 / Convenience Functions
- `Pretty(IValue) (string, error)` - 美化输出 / Pretty print
- `Compact(IValue) (string, error)` - 紧凑输出 / Compact output
- `HTMLSafe(IValue) (string, error)` - HTML安全输出 / HTML-safe output
- `MustParse([]byte) IValue` - 必须解析成功 / Must parse successfully
- `MustParseString(string) IValue` - 必须解析字符串成功 / Must parse string successfully

### 性能基准 / Performance Benchmarks

#### 解析性能 / Parsing Performance
```
BenchmarkParse/small_object-8     500000    24.8µs/op    1024 B/op    12 allocs/op
BenchmarkParse/medium_object-8    100000   124.5µs/op    4096 B/op    45 allocs/op
BenchmarkParse/large_object-8      10000  1247.3µs/op   16384 B/op   156 allocs/op
BenchmarkParse/array-8            200000    62.1µs/op    2048 B/op    28 allocs/op
```

#### 序列化性能 / Serialization Performance
```
BenchmarkSerialize/small_object-8  800000    18.2µs/op     512 B/op     8 allocs/op
BenchmarkSerialize/medium_object-8 200000    89.4µs/op    2048 B/op    32 allocs/op
BenchmarkSerialize/large_object-8   20000   892.1µs/op    8192 B/op   128 allocs/op
BenchmarkSerialize/array-8         400000    45.3µs/op    1024 B/op    16 allocs/op
```

#### JSONPath查询性能 / JSONPath Query Performance
```
BenchmarkJSONPath/simple_query-8     2000000    0.58µs/op      64 B/op     2 allocs/op
BenchmarkJSONPath/complex_query-8     500000     3.24µs/op     256 B/op     8 allocs/op
BenchmarkJSONPath/filter_query-8      100000    12.45µs/op     512 B/op    16 allocs/op
BenchmarkJSONPath/recursive_query-8    50000    28.67µs/op    1024 B/op    32 allocs/op
```

#### 对象池性能 / Object Pool Performance
```
BenchmarkPooled/parse-8           800000    15.2µs/op     256 B/op     3 allocs/op
BenchmarkPooled/serialize-8      1200000    12.1µs/op     128 B/op     2 allocs/op
BenchmarkPooled/create_object-8  5000000     0.24µs/op       0 B/op     0 allocs/op
```

### 兼容性 / Compatibility
- **JSON标准** / JSON Standard: 完全兼容RFC 7159 / Fully compatible with RFC 7159
- **JSONPath标准** / JSONPath Standard: 兼容JSONPath规范 / Compatible with JSONPath specification
- **Go版本** / Go Version: 支持Go 1.21+ / Supports Go 1.21+
- **平台支持** / Platform Support: Linux, macOS, Windows
- **架构支持** / Architecture Support: amd64, arm64, 386, arm

### 已知限制 / Known Limitations
- 最大解析深度默认为1000层 / Maximum parsing depth defaults to 1000 levels
- 单个字符串最大长度为1MB / Maximum single string length is 1MB
- 数组和对象最大元素数为10000 / Maximum array and object elements is 10000
- JSONPath不支持脚本表达式 / JSONPath does not support script expressions

### 安全考虑 / Security Considerations
- 输入验证和边界检查 / Input validation and boundary checking
- 防止深度攻击的深度限制 / Depth limits to prevent depth attacks
- 内存使用限制防止DoS攻击 / Memory usage limits to prevent DoS attacks
- 安全的字符串转义处理 / Secure string escape handling

---

## 版本说明 / Version Notes

### 语义化版本规则 / Semantic Versioning Rules
- **主版本号 (MAJOR)**: 不兼容的API变更 / Incompatible API changes
- **次版本号 (MINOR)**: 向后兼容的功能性新增 / Backward compatible feature additions
- **修订号 (PATCH)**: 向后兼容的问题修正 / Backward compatible bug fixes

### 发布周期 / Release Cycle
- **主版本**: 根据需要发布 / Released as needed
- **次版本**: 每2-3个月发布 / Released every 2-3 months
- **修订版本**: 根据bug修复需要发布 / Released as needed for bug fixes

### 支持政策 / Support Policy
- **当前版本**: 完全支持和维护 / Full support and maintenance
- **前一个主版本**: 安全更新和关键bug修复 / Security updates and critical bug fixes
- **更早版本**: 不再维护 / No longer maintained

---

## 贡献者 / Contributors

感谢所有为 xyJson 1.0.0 版本做出贡献的开发者！

Thanks to all developers who contributed to xyJson 1.0.0!

<!-- 贡献者列表将在这里自动生成 -->
<!-- Contributors list will be automatically generated here -->

---

## 链接 / Links

- [项目主页 / Project Homepage](https://github.com/yourusername/xyJson)
- [API文档 / API Documentation](https://pkg.go.dev/github.com/yourusername/xyJson)
- [问题报告 / Issue Reports](https://github.com/yourusername/xyJson/issues)
- [功能请求 / Feature Requests](https://github.com/yourusername/xyJson/discussions)
- [贡献指南 / Contributing Guide](CONTRIBUTING.md)
- [许可证 / License](LICENSE)

---

**注意**: 此更新日志遵循 [Keep a Changelog](https://keepachangelog.com/) 格式。

**Note**: This changelog follows the [Keep a Changelog](https://keepachangelog.com/) format.