# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 快速参考

### 最常用命令
```bash
go build ./...                          # 构建所有模块
go test ./...                           # 运行所有测试
go test ./internal/store -v             # 测试特定包(详细输出)
go test -run TestStoreManager ./...     # 运行特定测试
go fmt ./...                            # 格式化代码
go vet ./...                            # 静态分析
go mod tidy                             # 整理依赖
```

### 核心架构原则(必须遵守)
1. **依赖方向**: Engine → Internal → Core (禁止反向依赖)
2. **Core 层纯净**: 禁止框架类型依赖（仅纯 Go 类型）
3. **客户端外部管理**: `*clientv3.Client` 由调用方创建和关闭
4. **功能分离**: Watcher 处理原始数据，Store 处理强类型配置

### 常见陷阱
- ❌ 在 Core 层引入 etcd 或其他框架依赖
- ❌ 在 Engine 内部创建 etcd 客户端（应由调用方管理）
- ❌ Watcher 和 Store 功能混用（应根据场景选择）
- ❌ 不使用预定义错误（`core/errors.go`）
- ❌ Config 结构体放在 handler.go 而非 config.go
- ❌ 在 Manager 中直接存储 podName/serviceName（应使用 LogContext）

## 项目简介

etcdtrigger v2 是一个 Go 语言编写的 etcd 配置管理库，采用三层架构（Engine → Internal → Core），提供两种功能模式：

- **Watcher**：原始回调监听（Subscribe + 字节数组）
- **Store**：强类型配置缓存（GetConfig + JSON 反序列化）

## 核心架构

### 三层架构与依赖流向

```
Engine Layer (engine/)           ← 对外接口层
    ↓ 依赖
Internal Layer (internal/*)      ← Manager 模式
    ├── watcher/                 ← 原始监听功能
    └── store/                   ← 强类型缓存功能
    ↓ 依赖
Core Layer (core/)               ← 领域模型，纯 Go 类型
```

**关键原则**：
- Core 层不依赖任何框架（无 etcd、无外部库类型）
- Internal 层实现具体业务逻辑
- Engine 层组合 Internal 层的 Manager，对外暴露统一接口

## 核心模块职责

### Engine 层（engine/）

对外接口层，组合 Watcher 和 Store 功能：

| 方法 | 功能 | 所属模块 |
|------|------|----------|
| `Watch` | 订阅配置变更（原始回调） | Watcher |
| `WatchPut/Get/Delete` | 原始数据读写 | Watcher |
| `GetConfig` | 从缓存获取强类型配置 | Store |
| `PutConfig/DeleteConfig` | 写入/删除配置 | Store |
| `AddPrefixWatcher` | 添加前缀监听器 | Store |
| `GetAllKeys` | 获取所有键 | Store |
| `Client` | 返回底层 etcd 客户端 | - |

### Internal 层

#### internal/watcher（原始监听）
- 管理 etcd watch 订阅
- 处理原始字节数据
- 提供回调机制

#### internal/store（强类型缓存）
- 内存缓存配置数据
- JSON 序列化/反序列化
- 前缀监听器支持

### Core 层（core/）

**领域模型**：
- `EventType`：事件类型（Put/Delete）
- `WatchEvent`：监听事件
- `WatchConfig`：配置监听配置
- `WatchCallback`/`PrefixWatchCallback`：回调函数类型

**错误定义**（`core/errors.go`）：
- 连接相关：`ErrConnectionFailed`、`ErrConnectionTimeout` 等
- 配置相关：`ErrInvalidConfig`、`ErrConfigNotFound` 等
- 操作相关：`ErrPutFailed`、`ErrWatchFailed` 等
- 序列化相关：`ErrMarshalFailed`、`ErrUnmarshalFailed`

## Manager 开发规范

### 文件组织（所有 internal/* 和 engine 必须遵守）

| 文件 | 职责 | 必需性 |
|------|-----|--------|
| `<manager>.go` | 接口定义 + NewManager() 构造函数（调用 newManager） | ✅ 必需 |
| `handler.go` | 结构体定义 + newManager() + 接口实现 | ✅ 必需 |
| `config.go` | Config 结构体（即便为空也必须有此文件） | ✅ 必需 |
| `internal.go` | 未导出的私有方法 | 🟡 可选 |

### 构造函数模式（internal 和 engine 层必须遵守）

```go
// config.go - 配置结构体（即便为空也要有）
type Config struct{}

// <name>.go - 接口定义 + 导出构造函数（返回接口）
type Manager interface { ... }

func NewManager(...) Manager {
    return newManager(...)  // 调用 handler.go 的未导出函数
}

// handler.go - 结构体定义 + 未导出构造函数（返回结构体指针）
type xxxManager struct { ... }

func newManager(...) *xxxManager {
    return &xxxManager{ ... }
}
```

**关键点**：
- `NewXxx` 在接口定义文件中，返回**接口**
- `newXxx` 在 handler.go 中，返回**结构体指针**
- engine 层同样遵循：`NewEngine` 返回 `Engine`，`newEngine` 返回 `*engine`

### 日志上下文（LogContext）

日志字段由 `core.LogContext` 统一管理，Engine 创建后传递给各 Manager：

```go
// core/log.go
type LogContext struct {
    PodName     string
    ServiceName string
}

func (c *LogContext) WithModule(module, operation string) logx.Logger

// 使用方式
m.logCtx.WithModule("store", "put_config").WithFields(...).Info("写入成功")
// 或定义简化方法
func (m *storeManager) log(operation string) logx.Logger {
    return m.logCtx.WithModule("store", operation)
}
```

### 关键规范

**注释风格**：
- 对外接口（engine 层）：结构化注释，包含参数、返回、说明
- 内部接口（internal 层）：短代码和短注释可在同一行
- 禁止使用分隔线注释（如 `// ---- xxx ----`）

```go
// 对外接口（engine 层）- 结构化注释
// Watch 订阅指定前缀的配置变更
// 参数：
//   - key: 监听的键或前缀
//   - callback: 配置变更时的回调函数
// 返回：
//   - error: 订阅失败时返回错误
// 说明：
//   - 支持前缀匹配，会先触发当前已存在的值
Watch(key string, callback core.WatchCallback) error

// 内部接口（internal 层）- 简洁风格
type Manager interface {
    Watch(key string, callback core.WatchCallback) error // 订阅配置变更
    Put(key string, value []byte) error                  // 写入原始数据
}
```

**日志使用**：
- 使用 `logx.WithContext().WithFields()` 结构化日志
- logx 使用 `Error`，禁止使用 `Warn`

**错误处理**：
- 使用预定义错误 (`core/errors.go`)
- 示例：`return core.ErrConfigNotFound`

## 两种功能模式选择

| 场景 | 推荐模式 | 原因 |
|------|----------|------|
| 需要原始字节处理 | Watcher | 不涉及序列化 |
| 需要强类型配置 | Store | 自动 JSON 反序列化 |
| 需要配置变更回调 | Store | AddPrefixWatcher |
| 需要直接操作 etcd | Watcher | 原始读写 |

## 常见问题

### Q: 如何添加新的 Manager?
1. 在 `internal/<manager>/` 创建目录
2. 创建 `<manager>.go`、`handler.go`、`config.go`
3. 在 `engine/handler.go` 初始化并组合 Manager
4. 参考 `internal/store/` 实现

### Q: etcd 客户端由谁管理?
由调用方创建和关闭，Engine 只接收客户端引用。
