package watcher

import (
	"github.com/rezeropoint/etcdtrigger/v2/core"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Manager 原始监听管理器接口
type Manager interface {
	Watch(key string, callback core.WatchCallback) error // 订阅配置变更
	WatchPut(key string, value []byte) error             // 写入原始数据
	WatchDelete(key string) error                        // 删除数据
	WatchGet(key string) ([]byte, error)                 // 获取原始数据
}

// NewManager 创建监听管理器
func NewManager(client *clientv3.Client, logCtx *core.LogContext, config *Config) Manager {
	return newManager(client, logCtx, config)
}
