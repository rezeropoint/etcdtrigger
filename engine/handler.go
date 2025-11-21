package engine

import (
	"context"

	"github.com/rezeropoint/etcdtrigger/v2/core"
	"github.com/rezeropoint/etcdtrigger/v2/internal/store"
	"github.com/rezeropoint/etcdtrigger/v2/internal/watcher"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// engine Engine 实现
type engine struct {
	client     *clientv3.Client
	watcherMgr watcher.Manager
	storeMgr   store.Manager
}

// newEngine 创建 Engine 实例
func newEngine(client *clientv3.Client, config *Config) *engine {
	logCtx := &core.LogContext{
		PodName:     config.PodName,
		ServiceName: config.ServiceName,
	}

	return &engine{
		client:     client,
		watcherMgr: watcher.NewManager(client, logCtx, &watcher.Config{}),
		storeMgr:   store.NewManager(client, logCtx, &store.Config{Configs: config.Configs}),
	}
}

// Watch 订阅配置变更（原始回调模式）
func (e *engine) Watch(key string, callback core.WatchCallback) error {
	return e.watcherMgr.Watch(key, callback)
}

// WatchPut 写入原始数据
func (e *engine) WatchPut(key string, value []byte) error {
	return e.watcherMgr.WatchPut(key, value)
}

// WatchDelete 删除数据
func (e *engine) WatchDelete(key string) error {
	return e.watcherMgr.WatchDelete(key)
}

// WatchGet 获取原始数据
func (e *engine) WatchGet(key string) ([]byte, error) {
	return e.watcherMgr.WatchGet(key)
}

// GetConfig 从缓存获取配置（强类型）
func (e *engine) GetConfig(key string, result any) bool {
	return e.storeMgr.GetConfig(key, result)
}

// GetAllKeys 获取指定前缀的所有键
func (e *engine) GetAllKeys(prefix string) []string {
	return e.storeMgr.GetAllKeys(prefix)
}

// PutConfig 写入配置（自动 JSON 序列化）
func (e *engine) PutConfig(ctx context.Context, key string, config any) error {
	return e.storeMgr.PutConfig(ctx, key, config)
}

// DeleteConfig 删除配置
func (e *engine) DeleteConfig(ctx context.Context, key string) error {
	return e.storeMgr.DeleteConfig(ctx, key)
}

// AddPrefixWatcher 添加前缀监听器
func (e *engine) AddPrefixWatcher(prefix string, callback core.PrefixWatchCallback) {
	e.storeMgr.AddPrefixWatcher(prefix, callback)
}

// Client 返回底层 etcd 客户端
func (e *engine) Client() *clientv3.Client {
	return e.client
}
