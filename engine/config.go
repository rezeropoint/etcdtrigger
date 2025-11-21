package engine

import "github.com/rezeropoint/etcdtrigger/v2/core"

// Config 引擎配置
type Config struct {
	PodName     string             `json:",optional"` // Pod 标识（日志用）
	ServiceName string             `json:",optional"` // 服务名称（日志用）
	Configs     []core.WatchConfig `json:",optional"` // 预加载配置（强类型缓存用）
}
