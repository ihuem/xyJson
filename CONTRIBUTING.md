# 贡献指南 / Contributing Guide

感谢您对 xyJson 项目的关注！我们欢迎所有形式的贡献，包括但不限于代码、文档、测试、问题报告和功能建议。

Thank you for your interest in the xyJson project! We welcome all forms of contributions, including but not limited to code, documentation, tests, issue reports, and feature suggestions.

## 📋 目录 / Table of Contents

- [行为准则](#行为准则--code-of-conduct)
- [如何贡献](#如何贡献--how-to-contribute)
- [开发环境设置](#开发环境设置--development-setup)
- [代码规范](#代码规范--coding-standards)
- [提交规范](#提交规范--commit-conventions)
- [测试要求](#测试要求--testing-requirements)
- [文档要求](#文档要求--documentation-requirements)
- [发布流程](#发布流程--release-process)

## 🤝 行为准则 / Code of Conduct

### 我们的承诺 / Our Pledge

为了营造一个开放和友好的环境，我们作为贡献者和维护者承诺，无论年龄、体型、残疾、种族、性别认同和表达、经验水平、国籍、个人形象、种族、宗教或性取向如何，参与我们项目和社区的每个人都能获得无骚扰的体验。

In the interest of fostering an open and welcoming environment, we as contributors and maintainers pledge to make participation in our project and our community a harassment-free experience for everyone.

### 我们的标准 / Our Standards

积极行为的例子包括：
Examples of positive behavior include:

- 使用友好和包容的语言 / Using welcoming and inclusive language
- 尊重不同的观点和经验 / Being respectful of differing viewpoints and experiences
- 优雅地接受建设性批评 / Gracefully accepting constructive criticism
- 专注于对社区最有利的事情 / Focusing on what is best for the community
- 对其他社区成员表示同情 / Showing empathy towards other community members

## 🚀 如何贡献 / How to Contribute

### 1. 报告问题 / Reporting Issues

在报告问题之前，请：
Before reporting an issue, please:

- 检查现有的 [Issues](https://github.com/yourusername/xyJson/issues) 确保问题未被报告
- 使用最新版本测试问题是否仍然存在
- 收集相关信息（Go版本、操作系统、错误信息等）

#### 问题报告模板 / Issue Report Template

```markdown
**问题描述 / Bug Description**
简洁清晰地描述问题

**重现步骤 / Steps to Reproduce**
1. 执行 '...'
2. 点击 '....'
3. 滚动到 '....'
4. 看到错误

**期望行为 / Expected Behavior**
描述您期望发生的情况

**实际行为 / Actual Behavior**
描述实际发生的情况

**环境信息 / Environment**
- Go版本: [例如 1.21.0]
- 操作系统: [例如 Ubuntu 20.04]
- xyJson版本: [例如 v1.0.0]

**附加信息 / Additional Context**
添加任何其他相关信息、截图等
```

### 2. 功能请求 / Feature Requests

我们欢迎新功能的建议！请：
We welcome suggestions for new features! Please:

- 检查是否已有类似的功能请求
- 详细描述功能的用途和价值
- 提供使用场景和示例
- 考虑向后兼容性

#### 功能请求模板 / Feature Request Template

```markdown
**功能描述 / Feature Description**
简洁清晰地描述您想要的功能

**问题背景 / Problem Statement**
描述这个功能要解决的问题

**建议解决方案 / Proposed Solution**
描述您希望如何实现这个功能

**替代方案 / Alternative Solutions**
描述您考虑过的其他解决方案

**使用场景 / Use Cases**
提供具体的使用场景和示例代码

**优先级 / Priority**
- [ ] 低 / Low
- [ ] 中 / Medium  
- [ ] 高 / High
- [ ] 紧急 / Critical
```

### 3. 代码贡献 / Code Contributions

#### 贡献流程 / Contribution Workflow

1. **Fork 仓库** / Fork the repository
   ```bash
   # 在 GitHub 上点击 Fork 按钮
   # Click the Fork button on GitHub
   ```

2. **克隆您的 Fork** / Clone your fork
   ```bash
   git clone https://github.com/yourusername/xyJson.git
   cd xyJson
   ```

3. **添加上游仓库** / Add upstream repository
   ```bash
   git remote add upstream https://github.com/originalowner/xyJson.git
   ```

4. **创建功能分支** / Create a feature branch
   ```bash
   git checkout -b feature/your-feature-name
   ```

5. **进行更改** / Make your changes
   - 遵循代码规范
   - 添加测试
   - 更新文档

6. **运行测试** / Run tests
   ```bash
   go test ./...
   go test -race ./...
   go test -bench=. ./benchmark/
   ```

7. **提交更改** / Commit changes
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

8. **推送到您的 Fork** / Push to your fork
   ```bash
   git push origin feature/your-feature-name
   ```

9. **创建 Pull Request** / Create a Pull Request
   - 在 GitHub 上创建 PR
   - 填写 PR 模板
   - 等待代码审查

## 🛠️ 开发环境设置 / Development Setup

### 系统要求 / System Requirements

- Go 1.21 或更高版本 / Go 1.21 or higher
- Git
- Make (可选 / optional)

### 环境设置步骤 / Setup Steps

1. **克隆仓库** / Clone the repository
   ```bash
   git clone https://github.com/yourusername/xyJson.git
   cd xyJson
   ```

2. **安装依赖** / Install dependencies
   ```bash
   go mod download
   go mod tidy
   ```

3. **验证安装** / Verify installation
   ```bash
   go test ./...
   ```

4. **运行示例** / Run examples
   ```bash
   go run examples/basic_usage.go
   go run examples/advanced_features.go
   ```

### 开发工具推荐 / Recommended Development Tools

- **IDE**: VS Code, GoLand, Vim/Neovim
- **Go工具** / Go tools:
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  go install github.com/securecodewarrior/sast-scan@latest
  ```

## 📝 代码规范 / Coding Standards

### Go 代码规范 / Go Code Standards

1. **遵循官方规范** / Follow official standards
   - 使用 `gofmt` 格式化代码
   - 使用 `goimports` 管理导入
   - 遵循 [Effective Go](https://golang.org/doc/effective_go.html)

2. **命名约定** / Naming Conventions
   ```go
   // 包名：小写，简洁
   package xyJson
   
   // 公开接口：大写开头，清晰描述
   type IValue interface {}
   
   // 私有类型：小写开头
   type parser struct {}
   
   // 常量：大写，使用下划线分隔
   const MAX_DEPTH = 1000
   
   // 变量：驼峰命名
   var defaultParser *parser
   ```

3. **注释规范** / Comment Standards
   ```go
   // Package xyJson 提供高性能的JSON处理功能。
   // 包括解析、序列化、JSONPath查询等特性。
   package xyJson
   
   // IValue 表示一个JSON值的接口。
   // 所有JSON值类型都实现此接口。
   type IValue interface {
       // Type 返回值的类型
       Type() ValueType
       
       // String 返回值的字符串表示
       String() string
   }
   
   // Parse 解析JSON字节数据并返回对应的值。
   // 参数 data 是要解析的JSON字节数据。
   // 返回解析后的值和可能的错误。
   func Parse(data []byte) (IValue, error) {
       // 实现细节...
   }
   ```

4. **错误处理** / Error Handling
   ```go
   // 定义自定义错误类型
   type ParseError struct {
       Message  string
       Line     int
       Column   int
       Position int64
   }
   
   func (e *ParseError) Error() string {
       return fmt.Sprintf("parse error at line %d, column %d: %s", 
           e.Line, e.Column, e.Message)
   }
   
   // 使用有意义的错误信息
   if len(data) == 0 {
       return nil, &ParseError{
           Message: "empty input data",
           Line:    1,
           Column:  1,
       }
   }
   ```

5. **性能考虑** / Performance Considerations
   ```go
   // 避免不必要的内存分配
   func (p *parser) parseString() string {
       // 重用缓冲区
       p.buffer.Reset()
       // ...
   }
   
   // 使用对象池
   var stringPool = sync.Pool{
       New: func() interface{} {
           return make([]byte, 0, 64)
       },
   }
   ```

### 代码质量检查 / Code Quality Checks

运行以下命令确保代码质量：
Run the following commands to ensure code quality:

```bash
# 格式化代码
go fmt ./...

# 整理导入
goimports -w .

# 静态分析
go vet ./...

# Lint 检查
golangci-lint run

# 安全扫描
gosec ./...
```

## 📋 提交规范 / Commit Conventions

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：
We use [Conventional Commits](https://www.conventionalcommits.org/) specification:

### 提交消息格式 / Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### 类型 / Types

- `feat`: 新功能 / New feature
- `fix`: 错误修复 / Bug fix
- `docs`: 文档更新 / Documentation update
- `style`: 代码格式化 / Code formatting
- `refactor`: 代码重构 / Code refactoring
- `perf`: 性能优化 / Performance improvement
- `test`: 测试相关 / Test related
- `chore`: 构建过程或辅助工具的变动 / Build process or auxiliary tool changes
- `ci`: CI配置文件和脚本的变动 / CI configuration files and scripts changes

### 示例 / Examples

```bash
# 新功能
git commit -m "feat(parser): add support for custom number formats"

# 错误修复
git commit -m "fix(serializer): handle null values correctly"

# 文档更新
git commit -m "docs: update JSONPath query examples"

# 性能优化
git commit -m "perf(pool): optimize object reuse strategy"

# 破坏性变更
git commit -m "feat!: change IValue interface signature

BREAKING CHANGE: IValue.String() now returns (string, error)"
```

## 🧪 测试要求 / Testing Requirements

### 测试覆盖率 / Test Coverage

- 新代码必须有相应的测试 / New code must have corresponding tests
- 总体测试覆盖率应保持在 90% 以上 / Overall test coverage should remain above 90%
- 关键路径必须有 100% 覆盖率 / Critical paths must have 100% coverage

### 测试类型 / Test Types

1. **单元测试** / Unit Tests
   ```bash
   go test ./test/...
   ```

2. **集成测试** / Integration Tests
   ```bash
   go test ./test/integration_test.go
   ```

3. **基准测试** / Benchmark Tests
   ```bash
   go test -bench=. ./benchmark/
   ```

4. **竞态条件测试** / Race Condition Tests
   ```bash
   go test -race ./...
   ```

### 测试编写指南 / Test Writing Guidelines

1. **测试命名** / Test Naming
   ```go
   func TestParser_ParseString_ValidJSON(t *testing.T) {}
   func TestParser_ParseString_InvalidJSON(t *testing.T) {}
   func BenchmarkParser_ParseString(b *testing.B) {}
   ```

2. **表格驱动测试** / Table-Driven Tests
   ```go
   func TestParser_ParseString(t *testing.T) {
       tests := []struct {
           name     string
           input    string
           expected interface{}
           wantErr  bool
       }{
           {
               name:     "valid object",
               input:    `{"key":"value"}`,
               expected: map[string]interface{}{"key": "value"},
               wantErr:  false,
           },
           // 更多测试用例...
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               result, err := ParseString(tt.input)
               if (err != nil) != tt.wantErr {
                   t.Errorf("ParseString() error = %v, wantErr %v", err, tt.wantErr)
                   return
               }
               // 验证结果...
           })
       }
   }
   ```

3. **测试辅助函数** / Test Helper Functions
   ```go
   func assertJSONEqual(t *testing.T, expected, actual IValue) {
       t.Helper()
       if !expected.Equals(actual) {
           t.Errorf("JSON values not equal:\nexpected: %s\nactual: %s", 
               expected.String(), actual.String())
       }
   }
   ```

## 📚 文档要求 / Documentation Requirements

### 代码文档 / Code Documentation

1. **包级别文档** / Package-level Documentation
   ```go
   // Package xyJson 提供高性能的JSON处理功能。
   //
   // 主要特性：
   //   - 高性能解析和序列化
   //   - JSONPath查询支持
   //   - 内存池优化
   //   - 性能监控
   //
   // 基本用法：
   //   value, err := xyJson.ParseString(`{"key":"value"}`)
   //   if err != nil {
   //       log.Fatal(err)
   //   }
   //   result, _ := xyJson.SerializeToString(value)
   package xyJson
   ```

2. **公开API文档** / Public API Documentation
   - 所有公开的类型、函数、方法都必须有文档注释
   - 文档应包括用途、参数、返回值、使用示例
   - 特殊情况和错误条件应该说明

3. **示例代码** / Example Code
   ```go
   // Example_basicUsage 演示基本的JSON操作
   func Example_basicUsage() {
       // 解析JSON
       value, _ := ParseString(`{"name":"Alice","age":30}`)
       
       // 访问数据
       obj := value.(IObject)
       name := obj.Get("name").String()
       
       fmt.Println(name)
       // Output: "Alice"
   }
   ```

### 文档更新 / Documentation Updates

当您的更改影响以下内容时，请更新相应文档：
When your changes affect the following, please update the corresponding documentation:

- API 接口变更 → 更新 API 文档
- 新功能添加 → 更新 README 和示例
- 性能改进 → 更新性能指南
- 配置选项变更 → 更新配置文档

## 🔄 发布流程 / Release Process

### 版本号规范 / Version Numbering

我们使用 [语义化版本](https://semver.org/) (SemVer)：
We use [Semantic Versioning](https://semver.org/) (SemVer):

- `MAJOR.MINOR.PATCH` (例如 `1.2.3`)
- `MAJOR`: 不兼容的API变更 / Incompatible API changes
- `MINOR`: 向后兼容的功能性新增 / Backward compatible functionality additions
- `PATCH`: 向后兼容的问题修正 / Backward compatible bug fixes

### 发布检查清单 / Release Checklist

- [ ] 所有测试通过 / All tests pass
- [ ] 测试覆盖率 ≥ 90% / Test coverage ≥ 90%
- [ ] 文档已更新 / Documentation updated
- [ ] CHANGELOG.md 已更新 / CHANGELOG.md updated
- [ ] 版本号已更新 / Version number updated
- [ ] 性能基准测试通过 / Performance benchmarks pass
- [ ] 安全扫描通过 / Security scan passes

## 🆘 获取帮助 / Getting Help

如果您在贡献过程中遇到问题，可以通过以下方式获取帮助：
If you encounter issues during the contribution process, you can get help through:

- 📧 **邮件** / Email: support@xyJson.dev
- 💬 **讨论** / Discussions: [GitHub Discussions](https://github.com/yourusername/xyJson/discussions)
- 🐛 **问题** / Issues: [GitHub Issues](https://github.com/yourusername/xyJson/issues)
- 📱 **社区** / Community: [Discord/Slack链接]

## 🙏 致谢 / Acknowledgments

感谢所有为 xyJson 项目做出贡献的开发者！您的贡献让这个项目变得更好。

Thanks to all developers who have contributed to the xyJson project! Your contributions make this project better.

### 贡献者列表 / Contributors List

<!-- 这里会自动生成贡献者列表 -->
<!-- Contributors list will be automatically generated here -->

---

**再次感谢您的贡献！🎉**

**Thank you again for your contribution! 🎉**