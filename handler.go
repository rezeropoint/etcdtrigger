package etcdtrigger

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// etcdClient etcd客户端实现
type etcdClient struct {
	client *clientv3.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func newEtcdClient(ctx context.Context, cancel context.CancelFunc, config *Config) (*etcdClient, error) {
	// 创建etcd客户端配置
	clientConfig := clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: config.DialTimeout,
	}

	// 如果提供了认证信息，则添加到配置中
	if config.Username != "" {
		clientConfig.Username = config.Username
		clientConfig.Password = config.Password
	}

	// 创建etcd客户端
	cli, err := clientv3.New(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrEtcdClientCreation, err)
	}

	// 检查连接是否正常
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 3*time.Second)
	defer timeoutCancel()

	// 尝试获取集群状态，确认连接正常
	_, err = cli.Status(timeoutCtx, config.Endpoints[0])
	if err != nil {
		cli.Close() // 关闭连接以防泄漏
		return nil, fmt.Errorf("%w: %v", ErrEtcdConnectionCheck, err)
	}

	return &etcdClient{
		client: cli,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Subscribe 订阅配置变更
// 该函数会同时监听配置的新增、更新和删除事件
// 初始化时会先获取所有现有配置并触发回调
func (c *etcdClient) Subscribe(key string, callback func(string, []byte) error) error {
	if c.client == nil {
		return ErrEtcdConnectionFailed
	}

	if key == "" {
		return ErrInvalidEtcdKey
	}

	// 获取当前值，确保不会漏掉任何更改
	resp, err := c.client.Get(c.ctx, key, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("%w: %v", ErrEtcdGetOperation, err)
	}

	// 处理当前值
	for _, kv := range resp.Kvs {
		if err := callback(string(kv.Key), kv.Value); err != nil {
			logx.Errorf("处理etcd键值失败: %v, 键: %s", err, kv.Key)
		}
	}

	// 监听后续更改
	watchChan := c.client.Watch(c.ctx, key, clientv3.WithPrefix())
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				logx.Infof("etcd watch上下文已取消，停止监听: %s", key)
				return
			case watchResp, ok := <-watchChan:
				if !ok {
					logx.Errorf("etcd watch通道已关闭: %s", key)
					return
				}

				if watchResp.Err() != nil {
					logx.Errorf("etcd watch错误: %v, 键前缀: %s", watchResp.Err(), key)
					continue
				}

				for _, event := range watchResp.Events {
					switch event.Type {
					case clientv3.EventTypePut:
						// 处理PUT事件（创建和更新）
						if err := callback(string(event.Kv.Key), event.Kv.Value); err != nil {
							logx.Errorf("处理etcd PUT事件失败: %v, 键: %s", err, event.Kv.Key)
						}
					case clientv3.EventTypeDelete:
						// 处理DELETE事件
						// 删除事件时，值为空
						if err := callback(string(event.Kv.Key), nil); err != nil {
							logx.Errorf("处理etcd DELETE事件失败: %v, 键: %s", err, event.Kv.Key)
						}
					}
				}
			}
		}
	}()

	return nil
}

// Close 关闭etcd客户端
// 该函数会取消监听上下文并关闭客户端连接
func (c *etcdClient) Close() error {
	if c.cancel != nil {
		c.cancel()
	}

	if c.client != nil {
		if err := c.client.Close(); err != nil {
			return fmt.Errorf("%w: %v", ErrEtcdClientClose, err)
		}
		c.client = nil
	}

	return nil
}

// Put 向etcd写入键值对
// 该函数用于向etcd写入配置信息，触发监听该键的客户端收到通知
func (c *etcdClient) Put(key string, value []byte) error {
	if c.client == nil {
		return ErrEtcdConnectionFailed
	}

	if key == "" {
		return ErrInvalidEtcdKey
	}

	// 向etcd写入键值对
	_, err := c.client.Put(c.ctx, key, string(value))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrEtcdPutOperation, err)
	}

	return nil
}

// Delete 删除etcd中的键值对
// 该函数用于删除etcd中的配置信息，触发监听该键的客户端收到删除通知
func (c *etcdClient) Delete(key string) error {
	if c.client == nil {
		return ErrEtcdConnectionFailed
	}

	if key == "" {
		return ErrInvalidEtcdKey
	}

	// 从etcd删除键值对
	_, err := c.client.Delete(c.ctx, key)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrEtcdDeleteOperation, err)
	}

	return nil
}
