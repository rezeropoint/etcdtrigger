# 贡献指南

感谢您对 EtcdTrigger 项目的关注！我们欢迎各种形式的贡献。

## 如何贡献

### 报告问题

如果您发现了 bug 或有功能建议，请：

1. 首先检查 [Issues](https://github.com/rezeropoint/etcdtrigger/issues) 确认问题是否已经被报告
2. 如果没有，请创建新的 Issue，包含：
   - 清晰的问题描述
   - 复现步骤（如果是 bug）
   - 期望的行为
   - 您的环境信息（Go 版本、操作系统等）

### 提交代码

1. Fork 此仓库
2. 创建您的功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 代码规范

- 遵循 Go 语言官方代码规范
- 运行 `go fmt` 格式化代码
- 运行 `go vet` 检查代码
- 为新功能添加测试
- 确保所有测试通过

### 测试

在提交之前，请确保：

```bash
# 运行所有测试
go test ./...

# 运行带覆盖率的测试
go test -cover ./...

# 代码检查
go vet ./...

# 格式化代码
go fmt ./...
```

## 开发环境设置

1. 确保安装了 Go 1.25.1 或更高版本
2. 克隆仓库：`git clone https://github.com/rezeropoint/etcdtrigger.git`
3. 安装依赖：`go mod download`
4. 运行测试：`go test ./...`

## 问题和讨论

如有任何问题，欢迎：
- 创建 Issue
- 发起 Discussion
- 联系维护者

感谢您的贡献！