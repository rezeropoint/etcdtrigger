package core

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

// LogContext 日志上下文
type LogContext struct {
	PodName     string
	ServiceName string
}

// WithModule 创建带模块和操作的日志记录器
func (c *LogContext) WithModule(module, operation string) logx.Logger {
	return logx.WithContext(context.Background()).WithFields(
		logx.Field("service", c.ServiceName),
		logx.Field("pod", c.PodName),
		logx.Field("module", module),
		logx.Field("operation", operation),
	)
}
