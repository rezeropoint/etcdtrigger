// Package engine 提供 etcd 配置管理的统一接口层。
//
// Engine 是对外暴露的核心接口，组合了两种功能模式：
//   - Watcher: 原始回调监听，处理字节数组数据
//   - Store: 强类型配置缓存，自动 JSON 序列化/反序列化
//
// 使用示例：
//
//	// 调用方创建和管理 etcd 客户端
//	etcdClient, _ := clientv3.New(clientv3.Config{
//	    Endpoints:   []string{"localhost:2379"},
//	    DialTimeout: 5 * time.Second,
//	})
//	defer etcdClient.Close()
//
//	// 创建 Engine
//	eng := engine.NewEngine(etcdClient, &engine.Config{
//	    PodName:     "my-pod",
//	    ServiceName: "my-service",
//	})
//
//	// 使用 Store 功能
//	var cfg MyConfig
//	if eng.GetConfig("/app/config", &cfg) {
//	    fmt.Println(cfg)
//	}
package engine

import (
	"context"

	"github.com/rezeropoint/etcdtrigger/v2/core"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Engine 是 etcd 配置管理引擎的核心接口
// 说明：
//   - 提供 Watcher 功能：原始数据操作，适用于需要处理字节数组的场景
//   - 提供 Store 功能：强类型配置缓存，适用于结构化配置管理
//   - etcd 客户端由调用方创建和管理，Engine 仅持有引用
type Engine interface {
	// Watch 订阅指定前缀的配置变更
	// 参数：
	//   - key: 监听的键或前缀
	//   - callback: 配置变更时的回调函数
	// 返回：
	//   - error: 订阅失败时返回错误
	// 说明：
	//   - 支持前缀匹配，会先触发当前已存在的值
	//   - 后续变更会异步触发回调
	Watch(key string, callback core.WatchCallback) error

	// WatchPut 写入原始字节数据到 etcd
	// 参数：
	//   - key: 键名
	//   - value: 原始字节数据
	// 返回：
	//   - error: 写入失败时返回错误
	WatchPut(key string, value []byte) error

	// WatchDelete 从 etcd 删除指定 key
	// 参数：
	//   - key: 要删除的键名
	// 返回：
	//   - error: 删除失败时返回错误
	WatchDelete(key string) error

	// WatchGet 从 etcd 获取原始字节数据
	// 参数：
	//   - key: 键名
	// 返回：
	//   - []byte: 原始字节数据
	//   - error: 获取失败或 key 不存在时返回错误
	WatchGet(key string) ([]byte, error)

	// GetConfig 从内存缓存获取强类型配置
	// 参数：
	//   - key: 配置键名
	//   - result: 指向结构体的指针，用于接收配置
	// 返回：
	//   - bool: true 表示获取成功，false 表示配置不存在
	// 说明：
	//   - result 必须是指向结构体的指针
	//   - 从内存缓存读取，不会访问 etcd
	GetConfig(key string, result any) bool

	// GetAllKeys 返回指定前缀下的所有缓存键
	// 参数：
	//   - prefix: 键前缀
	// 返回：
	//   - []string: 匹配前缀的所有键列表
	GetAllKeys(prefix string) []string

	// PutConfig 写入配置到 etcd
	// 参数：
	//   - ctx: 上下文
	//   - key: 配置键名
	//   - config: 配置对象，会自动 JSON 序列化
	// 返回：
	//   - error: 序列化或写入失败时返回错误
	PutConfig(ctx context.Context, key string, config any) error

	// DeleteConfig 从 etcd 删除配置
	// 参数：
	//   - ctx: 上下文
	//   - key: 要删除的配置键名
	// 返回：
	//   - error: 删除失败时返回错误
	DeleteConfig(ctx context.Context, key string) error

	// AddPrefixWatcher 添加前缀监听器
	// 参数：
	//   - prefix: 要监听的键前缀
	//   - callback: 配置变更时的回调函数
	// 说明：
	//   - 添加时会立即触发已存在配置的回调
	//   - 后续匹配前缀的配置变更都会触发回调
	AddPrefixWatcher(prefix string, callback core.PrefixWatchCallback)

	// Client 返回底层的 etcd 客户端
	// 返回：
	//   - *clientv3.Client: etcd 客户端实例
	// 说明：
	//   - 用于需要直接操作 etcd 的高级场景
	//   - 客户端生命周期由调用方管理
	Client() *clientv3.Client
}

// NewEngine 创建新的 Engine
// 参数：
//   - client: etcd 客户端（由调用方管理生命周期）
//   - config: 引擎配置
// 返回：
//   - Engine: 配置管理引擎实例
func NewEngine(client *clientv3.Client, config *Config) Engine {
	return newEngine(client, config)
}
