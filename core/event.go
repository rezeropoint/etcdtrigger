package core

// EventType 事件类型
type EventType string

const (
	EventTypePut    EventType = "PUT"
	EventTypeDelete EventType = "DELETE"
)

// WatchEvent 监听事件
type WatchEvent struct {
	Key       string    // 键
	Value     []byte    // 值（原始字节）
	EventType EventType // 事件类型
}

// String 返回事件类型的字符串表示
func (e EventType) String() string {
	return string(e)
}

// IsPut 是否为 PUT 事件
func (e EventType) IsPut() bool {
	return e == EventTypePut
}

// IsDelete 是否为 DELETE 事件
func (e EventType) IsDelete() bool {
	return e == EventTypeDelete
}
