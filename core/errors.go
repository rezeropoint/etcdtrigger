package core

import "errors"

// 预定义错误 - 连接相关
var (
	ErrConnectionFailed    = errors.New("etcd connection failed")
	ErrConnectionClosed    = errors.New("etcd connection closed")
	ErrConnectionTimeout   = errors.New("etcd connection timeout")
	ErrAuthenticationFailed = errors.New("etcd authentication failed")
)

// 预定义错误 - 配置相关
var (
	ErrInvalidConfig      = errors.New("invalid config")
	ErrConfigNotFound     = errors.New("config not found in cache")
	ErrConfigEmpty        = errors.New("config is empty")
	ErrEndpointsEmpty     = errors.New("etcd endpoints cannot be empty")
)

// 预定义错误 - 操作相关
var (
	ErrPutFailed      = errors.New("etcd put operation failed")
	ErrDeleteFailed   = errors.New("etcd delete operation failed")
	ErrGetFailed      = errors.New("etcd get operation failed")
	ErrWatchFailed    = errors.New("etcd watch failed")
	ErrWatchCanceled  = errors.New("etcd watch canceled")
)

// 预定义错误 - 序列化相关
var (
	ErrMarshalFailed   = errors.New("json marshal failed")
	ErrUnmarshalFailed = errors.New("json unmarshal failed")
)
