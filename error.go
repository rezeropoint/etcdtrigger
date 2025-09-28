package etcdtrigger

import "errors"

var (
	// 连接相关错误
	ErrEtcdConnectionFailed = errors.New("连接Etcd失败")     // ErrEtcdConnectionFailed 表示连接Etcd失败
	ErrEtcdEndpointsEmpty   = errors.New("etcd端点列表不能为空") // ErrEtcdEndpointsEmpty 表示etcd端点列表为空
	ErrEtcdClientCreation   = errors.New("创建etcd客户端失败")  // ErrEtcdClientCreation 表示创建etcd客户端失败
	ErrEtcdConnectionCheck  = errors.New("etcd连接状态检查失败") // ErrEtcdConnectionCheck 表示etcd连接状态检查失败
	ErrEtcdClientClose      = errors.New("关闭etcd客户端失败")  // ErrEtcdClientClose 表示关闭etcd客户端失败

	// 键值操作相关错误
	ErrInvalidEtcdKey      = errors.New("etcd键不能为空")    // ErrInvalidEtcdKey 表示etcd键不能为空
	ErrEtcdPutOperation    = errors.New("向etcd写入键值对失败") // ErrEtcdPutOperation 表示向etcd写入键值对失败
	ErrEtcdGetOperation    = errors.New("从etcd获取键值对失败") // ErrEtcdGetOperation 表示从etcd获取键值对失败
	ErrEtcdDeleteOperation = errors.New("从etcd删除键值对失败") // ErrEtcdDeleteOperation 表示从etcd删除键值对失败

	// 监听相关错误
	ErrEtcdWatchFailed   = errors.New("etcd监听失败") // ErrEtcdWatchFailed 表示etcd监听失败
	ErrCallbackExecution = errors.New("回调函数执行失败") // ErrCallbackExecution 表示回调函数执行失败

	// 上下文相关错误
	ErrContextCancelled   = errors.New("上下文已取消")  // ErrContextCancelled 表示上下文已取消
	ErrWatchChannelClosed = errors.New("监听通道已关闭") // ErrWatchChannelClosed 表示监听通道已关闭
	ErrOperationTimeout   = errors.New("操作超时")    // ErrOperationTimeout 表示操作超时

	// 配置相关错误
	ErrInvalidConfig      = errors.New("无效的配置")     // ErrInvalidConfig 表示配置无效
	ErrInvalidDialTimeout = errors.New("无效的连接超时时间") // ErrInvalidDialTimeout 表示连接超时时间无效
)
