package store

import (
	"context"

	"github.com/rezeropoint/etcdtrigger/v2/core"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Manager 配置存储管理器接口
type Manager interface {
	GetConfig(key string, result any) bool                           // 从缓存获取配置（强类型）
	GetAllKeys(prefix string) []string                               // 获取指定前缀的所有键
	PutConfig(ctx context.Context, key string, config any) error     // 写入配置（自动序列化）
	DeleteConfig(ctx context.Context, key string) error              // 删除配置
	AddPrefixWatcher(prefix string, callback core.PrefixWatchCallback) // 添加前缀监听器
}

// NewManager 创建配置存储管理器
func NewManager(client *clientv3.Client, logCtx *core.LogContext, config *Config) Manager {
	return newManager(client, logCtx, config)
}
