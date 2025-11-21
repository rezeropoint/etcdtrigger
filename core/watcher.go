package core

// WatchCallback 监听回调函数类型
// 用于处理单个键值变更事件
type WatchCallback func(event *WatchEvent) error

// PrefixWatchCallback 前缀监听回调函数类型
// 用于处理某个前缀下的键值变更事件
type PrefixWatchCallback func(key string, eventType EventType)
