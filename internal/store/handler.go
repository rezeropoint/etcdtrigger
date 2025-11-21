package store

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/rezeropoint/etcdtrigger/v2/core"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary

// storeManager 配置存储管理器实现
type storeManager struct {
	client         *clientv3.Client
	logCtx         *core.LogContext
	data           sync.Map // 存储不同类型的配置实例
	typeCaches     sync.Map // 缓存结构体类型
	prefixWatchers sync.Map // 前缀监听器
}

// newManager 创建配置存储管理器实例
func newManager(client *clientv3.Client, logCtx *core.LogContext, config *Config) *storeManager {
	manager := &storeManager{
		client: client,
		logCtx: logCtx,
	}

	// 初始化预配置的监听
	for _, cfg := range config.Configs {
		if cfg.Struct != nil {
			manager.initTypeStore(cfg.Struct)
			if err := manager.initConfig(context.Background(), cfg); err != nil {
				manager.log("init_config").WithFields(logx.Field("path", cfg.Path), logx.Field("error", err.Error())).Error("初始化配置失败")
			}
			go manager.watchConfigChanges(context.Background(), cfg)
		}
	}

	return manager
}

// GetConfig 从缓存获取配置
func (m *storeManager) GetConfig(key string, result any) bool {
	t := reflect.TypeOf(result)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		m.log("get_config").WithFields(logx.Field("key", key), logx.Field("error", "result 必须是指向结构体的指针")).Error("参数类型错误")
		return false
	}

	instanceMap, ok := m.data.Load(t)
	if !ok {
		m.log("get_config").WithFields(logx.Field("key", key), logx.Field("error", fmt.Sprintf("未找到类型 %v", t))).Error("类型未找到")
		return false
	}

	typedMap := instanceMap.(*sync.Map)
	value, ok := typedMap.Load(key)
	if !ok {
		return false
	}

	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(value).Elem())
	return true
}

// GetAllKeys 获取指定前缀的所有键
func (m *storeManager) GetAllKeys(prefix string) []string {
	keys := make([]string, 0)

	m.data.Range(func(_, value any) bool {
		instanceMap := value.(*sync.Map)
		instanceMap.Range(func(key, _ any) bool {
			keyStr, ok := key.(string)
			if ok && len(keyStr) >= len(prefix) && keyStr[:len(prefix)] == prefix {
				keys = append(keys, keyStr)
			}
			return true
		})
		return true
	})

	return keys
}

// PutConfig 写入配置
func (m *storeManager) PutConfig(ctx context.Context, key string, config any) error {
	value, err := jsonIter.Marshal(config)
	if err != nil {
		m.log("put_config").WithFields(logx.Field("key", key), logx.Field("error", err.Error())).Error("序列化失败")
		return fmt.Errorf("%w: %v", core.ErrMarshalFailed, err)
	}

	_, err = m.client.Put(ctx, key, string(value))
	if err != nil {
		m.log("put_config").WithFields(logx.Field("key", key), logx.Field("error", err.Error())).Error("写入失败")
		return fmt.Errorf("%w: %v", core.ErrPutFailed, err)
	}

	m.log("put_config").WithFields(logx.Field("key", key)).Info("写入成功")
	return nil
}

// DeleteConfig 删除配置
func (m *storeManager) DeleteConfig(ctx context.Context, key string) error {
	_, err := m.client.Delete(ctx, key)
	if err != nil {
		m.log("delete_config").WithFields(logx.Field("key", key), logx.Field("error", err.Error())).Error("删除失败")
		return fmt.Errorf("%w: %v", core.ErrDeleteFailed, err)
	}

	m.log("delete_config").WithFields(logx.Field("key", key)).Info("删除成功")
	return nil
}

// AddPrefixWatcher 添加前缀监听器
func (m *storeManager) AddPrefixWatcher(prefix string, callback core.PrefixWatchCallback) {
	m.prefixWatchers.Store(prefix, callback)

	// 触发已存在的配置
	m.data.Range(func(_, value any) bool {
		instanceMap := value.(*sync.Map)
		instanceMap.Range(func(key, _ any) bool {
			keyStr, ok := key.(string)
			if ok && len(keyStr) >= len(prefix) && keyStr[:len(prefix)] == prefix {
				callback(keyStr, core.EventTypePut)
			}
			return true
		})
		return true
	})

	m.log("add_prefix_watcher").WithFields(logx.Field("prefix", prefix)).Info("添加成功")
}

// log 创建结构化日志
func (m *storeManager) log(operation string) logx.Logger {
	return m.logCtx.WithModule("store", operation)
}
