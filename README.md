# EtcdTrigger

一个基于 Go 语言的 etcd 配置监听与管理库，提供实时配置变更监听（Watcher）和强类型配置缓存（Store）功能。

## 功能特性

- **双模式支持**:
  - 🚀 **Watcher 模式**: 原始回调监听，处理字节数组数据，适合底层事件处理
  - 💾 **Store 模式**: 强类型配置缓存，自动 JSON 序列化/反序列化，支持从内存直接读取配置
- 📋 **前缀匹配**: 支持按目录前缀监听配置变更
- 🔄 **自动同步**: 初始化时自动加载现有配置，后续变更实时同步
- 🔌 **依赖注入**: 灵活集成，支持传入外部管理的 etcd 客户端
- ⚡ **高性能**: 基于 go-zero 框架和 etcd client v3

## 安装

```bash
go get github.com/rezeropoint/etcdtrigger/v2
```

## 快速开始

### 基本用法

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/rezeropoint/etcdtrigger/v2/core"
    "github.com/rezeropoint/etcdtrigger/v2/engine"
    clientv3 "go.etcd.io/etcd/client/v3"
)

// 定义配置结构体
type DatabaseConfig struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

func main() {
    // 1. 创建 etcd 客户端（由调用方管理生命周期）
    etcdClient, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{"localhost:2379"},
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer etcdClient.Close()

    // 2. 创建引擎
    eng := engine.NewEngine(etcdClient, &engine.Config{
        PodName:     "my-pod",
        ServiceName: "my-service",
        // 预加载配置：自动监听并缓存到内存（Store 功能）
        Configs: []core.WatchConfig{
            {Path: "/app/config/db", Struct: &DatabaseConfig{}},
        },
    })

    // 3. 使用 Watcher 功能（原始操作）
    eng.Watch("/app/events/", func(event *core.WatchEvent) error {
        log.Printf("收到事件: Type=%s, Key=%s, Value=%s", event.EventType, event.Key, string(event.Value))
        return nil
    })

    // 4. 使用 Store 功能（强类型读写）
    // 写入配置
    dbConfig := &DatabaseConfig{Host: "localhost", Port: 3306}
    eng.PutConfig(context.Background(), "/app/config/db", dbConfig)

    // 等待同步...
    time.Sleep(100 * time.Millisecond)

    // 从本地缓存读取配置（不访问 etcd）
    var cachedDB DatabaseConfig
    if eng.GetConfig("/app/config/db", &cachedDB) {
        log.Printf("读取配置: %+v", cachedDB)
    }
}
```

## 核心概念

EtcdTrigger 提供了两种核心交互模式：

### 1. Watcher (原始监听)
适用于需要直接处理 etcd 原始数据的场景。
- `Watch`: 监听变更
- `WatchPut`: 写入原始字节
- `WatchGet`: 获取原始字节
- `WatchDelete`: 删除键

### 2. Store (强类型缓存)
适用于应用程序配置管理。配置数据被自动缓存到内存中，读取操作极其高效（无网络开销）。
- `PutConfig`: 序列化并写入配置
- `GetConfig`: 从内存缓存读取反序列化后的对象
- `AddPrefixWatcher`: 监听前缀变更
- `Configs` (初始化参数): 启动时自动加载并缓存的配置项

## API 文档

### Engine 接口

```go
type Engine interface {
    // Watcher 功能
    Watch(key string, callback core.WatchCallback) error
    WatchPut(key string, value []byte) error
    WatchDelete(key string) error
    WatchGet(key string) ([]byte, error)

    // Store 功能
    GetConfig(key string, result any) bool
    PutConfig(ctx context.Context, key string, config any) error
    DeleteConfig(ctx context.Context, key string) error
    GetAllKeys(prefix string) []string
    AddPrefixWatcher(prefix string, callback core.PrefixWatchCallback)

    // 获取底层客户端
    Client() *clientv3.Client
}
```

### Config 结构体

```go
type Config struct {
    PodName     string             // Pod 标识
    ServiceName string             // 服务名称
    Configs     []core.WatchConfig // 预加载配置列表
}
```

## 依赖

- [go-zero](https://github.com/zeromicro/go-zero)
- [etcd client v3](https://go.etcd.io/etcd/client/v3)

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件