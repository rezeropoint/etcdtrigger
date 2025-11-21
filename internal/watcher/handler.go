package watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/rezeropoint/etcdtrigger/v2/core"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// watcherManager 监听管理器实现
type watcherManager struct {
	client *clientv3.Client
	logCtx *core.LogContext
}

// newManager 创建监听管理器实例
func newManager(client *clientv3.Client, logCtx *core.LogContext, _ *Config) *watcherManager {
	return &watcherManager{
		client: client,
		logCtx: logCtx,
	}
}

// Watch 订阅配置变更
func (m *watcherManager) Watch(key string, callback core.WatchCallback) error {
	if m.client == nil {
		return core.ErrConnectionClosed
	}

	if key == "" {
		return core.ErrConfigEmpty
	}

	ctx := context.Background()

	// 获取当前值
	resp, err := m.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("%w: %v", core.ErrGetFailed, err)
	}

	// 处理当前值
	for _, kv := range resp.Kvs {
		event := &core.WatchEvent{
			Key:       string(kv.Key),
			Value:     kv.Value,
			EventType: core.EventTypePut,
		}
		if err := callback(event); err != nil {
			m.log("subscribe").WithFields(logx.Field("key", string(kv.Key)), logx.Field("error", err.Error())).Error("处理键值失败")
		}
	}

	// 监听后续变更
	watchChan := m.client.Watch(ctx, key, clientv3.WithPrefix())
	go func() {
		for watchResp := range watchChan {
			if watchResp.Err() != nil {
				m.log("subscribe").WithFields(logx.Field("key", key), logx.Field("error", watchResp.Err().Error())).Error("监听错误")
				continue
			}

			for _, ev := range watchResp.Events {
				event := &core.WatchEvent{
					Key: string(ev.Kv.Key),
				}

				switch ev.Type {
				case clientv3.EventTypePut:
					event.Value = ev.Kv.Value
					event.EventType = core.EventTypePut
				case clientv3.EventTypeDelete:
					event.Value = nil
					event.EventType = core.EventTypeDelete
				}

				if err := callback(event); err != nil {
					m.log("subscribe").WithFields(logx.Field("key", event.Key), logx.Field("event_type", event.EventType), logx.Field("error", err.Error())).Error("处理事件失败")
				}
			}
		}
	}()

	m.log("subscribe").WithFields(logx.Field("key", key)).Info("订阅成功")

	return nil
}

// WatchPut 写入原始数据
func (m *watcherManager) WatchPut(key string, value []byte) error {
	if m.client == nil {
		return core.ErrConnectionClosed
	}

	if key == "" {
		return core.ErrConfigEmpty
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.client.Put(ctx, key, string(value))
	if err != nil {
		return fmt.Errorf("%w: %v", core.ErrPutFailed, err)
	}

	return nil
}

// WatchDelete 删除数据
func (m *watcherManager) WatchDelete(key string) error {
	if m.client == nil {
		return core.ErrConnectionClosed
	}

	if key == "" {
		return core.ErrConfigEmpty
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("%w: %v", core.ErrDeleteFailed, err)
	}

	return nil
}

// WatchGet 获取原始数据
func (m *watcherManager) WatchGet(key string) ([]byte, error) {
	if m.client == nil {
		return nil, core.ErrConnectionClosed
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := m.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", core.ErrGetFailed, err)
	}

	if len(resp.Kvs) == 0 {
		return nil, core.ErrConfigNotFound
	}

	return resp.Kvs[0].Value, nil
}

// log 创建结构化日志
func (m *watcherManager) log(operation string) logx.Logger {
	return m.logCtx.WithModule("watcher", operation)
}
