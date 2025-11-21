package core

// WatchConfig 监听配置项
// 用于定义需要监听的键及其绑定的结构体类型
type WatchConfig struct {
	Path   string // 监听路径（支持前缀）
	Struct any    // 绑定的结构体实例（用于 JSON 反序列化，可为 nil）
}
