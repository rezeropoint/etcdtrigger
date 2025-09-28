package etcdtrigger

import (
	"context"
)

// EtcdClient etcd客户端接口
type EtcdClient interface {
	Subscribe(key string, callback func(string, []byte) error) error // Subscribe 订阅配置变更
	Put(key string, value []byte) error                              // Put 向etcd写入键值对
	Delete(key string) error                                         // Delete 删除etcd中的键值对
	Close() error                                                    // Close 关闭etcd客户端
}

// NewEtcdClientWithConfig 使用配置创建etcd客户端
func NewEtcdClient(ctx context.Context, cancel context.CancelFunc, config *Config) (EtcdClient, error) {
	if config == nil {
		return nil, ErrInvalidConfig
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newEtcdClient(ctx, cancel, config)
}
