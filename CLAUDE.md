# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

这是一个 Go 语言编写的 etcd 配置监听库，提供了对 etcd 键值变更的实时监听功能。项目使用 go-zero 框架进行日志记录。

## 代码架构

### 核心组件

- **EtcdClient 接口** (`etcd.go:8-13`): 定义了与 etcd 交互的核心方法
  - `Subscribe()`: 订阅配置变更，支持前缀匹配
  - `Put()`: 写入键值对
  - `Delete()`: 删除键值对
  - `Close()`: 关闭客户端连接

- **etcdClient 实现** (`handler.go:12-17`): EtcdClient 接口的具体实现
  - 使用 etcd clientv3 进行底层通信
  - 支持用户名/密码认证
  - 自动处理连接状态检查

- **Config 结构体** (`config.go:5-12`): etcd 客户端配置
  - 支持多个 etcd 端点
  - 可配置连接超时时间
  - 支持认证信息

### 错误处理

所有错误定义在 `error.go` 中，包括：
- 连接相关错误 (ErrEtcdConnectionFailed 等)
- 键值操作错误 (ErrEtcdPutOperation 等)
- 监听相关错误 (ErrEtcdWatchFailed 等)
- 配置相关错误 (ErrInvalidConfig 等)

## 开发命令

```bash
# 构建项目
go build

# 运行测试
go test ./...

# 运行单个测试
go test -run TestName

# 格式化代码
go fmt ./...

# 代码检查
go vet ./...

# 下载依赖
go mod download

# 整理依赖
go mod tidy
```

## 关键设计模式

1. **订阅模式**: Subscribe 方法会先获取当前所有匹配的键值对，然后启动 goroutine 监听后续变更
2. **上下文管理**: 使用 context.Context 进行生命周期管理和优雅关闭
3. **错误包装**: 使用 fmt.Errorf 和 errors.Is 进行错误包装和检查
4. **资源管理**: 确保 etcd 客户端连接正确关闭以防止资源泄漏

## 注意事项

- 所有数据库相关结构体应使用 sql.NullString 或 sql.NullTime 防止 NULL 值扫描错误
- Subscribe 方法会在独立的 goroutine 中运行，需要妥善处理上下文取消
- 错误日志使用 go-zero 的 logx 包进行记录