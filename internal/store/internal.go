package store

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/rezeropoint/etcdtrigger/v2/core"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// initTypeStore 初始化类型存储
func (m *storeManager) initTypeStore(configStruct any) {
	t := reflect.TypeOf(configStruct)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		panic("传入的必须是指向结构体的指针")
	}

	m.data.Store(t, &sync.Map{})
	m.typeCaches.Store(t, reflect.New(t.Elem()).Interface())
}

// initConfig 初始化配置
func (m *storeManager) initConfig(ctx context.Context, cfg core.WatchConfig) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := m.client.Get(ctx, cfg.Path, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range resp.Kvs {
		m.storeConfig(string(kv.Key), kv.Value, cfg.Struct)
	}

	return nil
}

// watchConfigChanges 监听配置变化
func (m *storeManager) watchConfigChanges(ctx context.Context, cfg core.WatchConfig) {
	watchChan := m.client.Watch(ctx, cfg.Path, clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case clientv3.EventTypePut:
				m.storeConfig(string(event.Kv.Key), event.Kv.Value, cfg.Struct)
				m.notifyPrefixWatchers(string(event.Kv.Key), core.EventTypePut)
			case clientv3.EventTypeDelete:
				m.removeConfig(string(event.Kv.Key), cfg.Struct)
				m.notifyPrefixWatchers(string(event.Kv.Key), core.EventTypeDelete)
			}
		}
	}
}

// storeConfig 存储配置
func (m *storeManager) storeConfig(key string, value []byte, configStruct any) {
	t := reflect.TypeOf(configStruct)
	cachedInstance, _ := m.typeCaches.Load(t)
	instance := reflect.New(reflect.TypeOf(cachedInstance).Elem()).Interface()

	if err := jsonIter.Unmarshal(value, instance); err != nil {
		m.log("store_config").WithFields(logx.Field("key", key), logx.Field("error", err.Error())).Error("反序列化失败")
		return
	}

	instanceMap, _ := m.data.Load(t)
	typedMap := instanceMap.(*sync.Map)
	typedMap.Store(key, instance)

	m.log("store_config").WithFields(logx.Field("key", key)).Info("更新成功")
}

// removeConfig 删除配置
func (m *storeManager) removeConfig(key string, configStruct any) {
	t := reflect.TypeOf(configStruct)

	instanceMap, _ := m.data.Load(t)
	typedMap := instanceMap.(*sync.Map)
	typedMap.Delete(key)

	m.log("remove_config").WithFields(logx.Field("key", key)).Info("删除成功")
}

// notifyPrefixWatchers 通知前缀监听器
func (m *storeManager) notifyPrefixWatchers(key string, eventType core.EventType) {
	m.prefixWatchers.Range(func(prefix, value any) bool {
		if k, ok := prefix.(string); ok && len(key) >= len(k) && key[:len(k)] == k {
			if callback, ok := value.(core.PrefixWatchCallback); ok {
				callback(key, eventType)
			}
		}
		return true
	})
}
