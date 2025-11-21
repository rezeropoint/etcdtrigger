package store

import "github.com/rezeropoint/etcdtrigger/v2/core"

// Config 配置存储管理器配置
type Config struct {
	Configs []core.WatchConfig // 预加载配置列表
}
