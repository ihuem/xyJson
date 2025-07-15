# xyJson 示例 / Examples

本目录包含了 xyJson 库的使用示例，帮助您快速了解和使用各种功能。

This directory contains usage examples for the xyJson library to help you quickly understand and use various features.

## 文件列表 / File List

### fast_path_example.go

演示 xyJson 快速路径功能的完整示例，包括：

A complete example demonstrating xyJson's fast path functionality, including:

- **快速路径解析** / **Fast Path Parsing**: 使用 `UnmarshalToStructFast` 进行高性能JSON解析
- **标准路径解析** / **Standard Path Parsing**: 使用 `UnmarshalToStruct` 进行功能丰富的解析
- **性能对比** / **Performance Comparison**: 展示两种方法的性能差异
- **使用建议** / **Usage Recommendations**: 何时选择哪种方法

#### 运行示例 / Running the Example

```bash
cd examples
go run fast_path_example.go
```

#### 示例输出 / Example Output

```
=== xyJson 快速路径示例 ===
=== xyJson Fast Path Example ===

1. 快速路径解析 (Fast Path Parsing):
   用户: {ID:123 Name:Alice Johnson Email:alice@example.com Active:true Balance:1250.75 Created:2023-01-15 10:30:00 +0000 UTC Address:{Street:123 Main St City:New York Zip:10001} Tags:[premium verified] Scores:[95 88 92]}
   地址: {Street:123 Main St City:New York Zip:10001}
   标签: [premium verified]

2. 标准路径解析 (Standard Path Parsing):
   用户: {ID:123 Name:Alice Johnson Email:alice@example.com Active:true Balance:1250.75 Created:2023-01-15 10:30:00 +0000 UTC Address:{Street:123 Main St City:New York Zip:10001} Tags:[premium verified] Scores:[95 88 92]}

3. 结果对比 (Result Comparison):
   ✓ 快速路径和标准路径结果一致
   ✓ Fast path and standard path produce identical results
```

## 性能特点 / Performance Characteristics

### 快速路径 / Fast Path

- **性能优势** / **Performance Advantage**: 接近官方 `encoding/json` 包的性能
- **内存优化** / **Memory Optimization**: 减少内存分配和垃圾回收压力
- **适用场景** / **Use Cases**: 简单的JSON到struct转换，高性能需求

### 标准路径 / Standard Path

- **功能丰富** / **Feature Rich**: 支持JSONPath查询、复杂操作
- **灵活性** / **Flexibility**: 提供更多的JSON处理选项
- **适用场景** / **Use Cases**: 需要高级功能的复杂JSON处理

## 基准测试结果 / Benchmark Results

基于最新的性能测试：

Based on latest performance tests:

```
BenchmarkCompareFastPath/xyJson_Fast-8         500000    2144 ns/op     296 B/op       7 allocs/op
BenchmarkCompareFastPath/xyJson_Standard-8       5000  225708 ns/op    1890 B/op      30 allocs/op
BenchmarkCompareFastPath/Official_JSON-8       500000    2361 ns/op     296 B/op       7 allocs/op
```

**性能提升** / **Performance Improvement**:
- 快速路径比标准路径快约 **105倍** / Fast path is about **105x faster** than standard path
- 快速路径性能接近官方JSON包 / Fast path performance is close to official JSON package

## 使用建议 / Usage Recommendations

### 选择快速路径的场景 / Choose Fast Path When:

✅ 高性能需求 / High performance requirements  
✅ 简单JSON到struct转换 / Simple JSON to struct conversion  
✅ 大量数据处理 / Large data processing  
✅ 内存使用敏感 / Memory usage sensitive  

### 选择标准路径的场景 / Choose Standard Path When:

✅ 需要JSONPath查询 / Need JSONPath queries  
✅ 复杂JSON操作 / Complex JSON operations  
✅ 需要高级功能 / Need advanced features  
✅ 灵活性优先于性能 / Flexibility over performance  

## 更多示例 / More Examples

更多示例正在开发中，敬请期待！

More examples are under development, stay tuned!

---

如有问题或建议，请提交 Issue 或 Pull Request。

For questions or suggestions, please submit an Issue or Pull Request.