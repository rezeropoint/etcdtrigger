# EtcdTrigger

一个基于 Go 语言的 etcd 配置监听库，提供实时配置变更监听功能。

## 功能特性

- 🚀 实时监听 etcd 配置变更
- 📋 支持前缀匹配监听
- 🔄 自动处理初始化配置加载
- 🔐 支持用户名密码认证
- 🛡️ 完善的错误处理机制
- ⚡ 基于 go-zero 框架的高性能日志

## 安装

```bash
go get github.com/rezeropoint/etcdtrigger
```

## 快速开始

### 基本用法

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/rezeropoint/etcdtrigger"
)

func main() {
    // 创建配置
    config := &etcdtrigger.Config{
        Endpoints:   []string{"localhost:2379"},
        DialTimeout: 5 * time.Second,
    }

    // 创建上下文
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 创建客户端
    client, err := etcdtrigger.NewEtcdClient(ctx, cancel, config)
    if err != nil {
        log.Fatal("创建etcd客户端失败:", err)
    }
    defer client.Close()

    // 订阅配置变更
    err = client.Subscribe("/config/", func(key string, value []byte) error {
        if value == nil {
            log.Printf("配置被删除: %s", key)
        } else {
            log.Printf("配置变更: %s = %s", key, string(value))
        }
        return nil
    })

    if err != nil {
        log.Fatal("订阅失败:", err)
    }

    // 阻塞等待
    select {}
}
```

### 带认证的用法

```go
config := &etcdtrigger.Config{
    Endpoints:   []string{"localhost:2379"},
    DialTimeout: 5 * time.Second,
    Username:    "your_username",
    Password:    "your_password",
}
```

## API 文档

### Config 结构体

```go
type Config struct {
    Key         string        // 监听的配置键前缀
    Endpoints   []string      // Etcd服务器端点列表
    DialTimeout time.Duration // 连接超时时间
    Username    string        // 用户名（可选）
    Password    string        // 密码（可选）
}
```

### EtcdClient 接口

```go
type EtcdClient interface {
    Subscribe(key string, callback func(string, []byte) error) error
    Put(key string, value []byte) error
    Delete(key string) error
    Close() error
}
```

#### Subscribe

订阅指定前缀的配置变更。会先加载所有现有配置，然后监听后续变更。

**参数:**
- `key`: 监听的键前缀
- `callback`: 配置变更回调函数，参数为键名和值（删除时值为 nil）

#### Put

向 etcd 写入键值对。

**参数:**
- `key`: 键名
- `value`: 值的字节数组

#### Delete

从 etcd 删除指定键。

**参数:**
- `key`: 要删除的键名

#### Close

关闭 etcd 客户端连接。

## 错误处理

库定义了详细的错误类型，便于错误处理和调试：

- `ErrEtcdConnectionFailed`: 连接 etcd 失败
- `ErrEtcdEndpointsEmpty`: etcd 端点列表为空
- `ErrInvalidEtcdKey`: etcd 键不能为空
- `ErrEtcdPutOperation`: 写入操作失败
- 更多错误类型请查看 `error.go`

## 依赖

- [go-zero](https://github.com/zeromicro/go-zero) - 高性能微服务框架
- [etcd client v3](https://go.etcd.io/etcd/client/v3) - etcd 官方客户端

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 贡献

欢迎提交 Issue 和 Pull Request！

## 支持

如果有任何问题，请提交 Issue 或联系维护者。