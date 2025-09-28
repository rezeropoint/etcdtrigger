package etcdtrigger

import "time"

// Config etcd客户端配置
type Config struct {
	Key         string        // 监听的配置键前缀
	Endpoints   []string      `json:"endpoints" yaml:"endpoints"`       // Etcd服务器端点列表
	DialTimeout time.Duration `json:"dial_timeout" yaml:"dial_timeout"` // 连接超时时间
	Username    string        `json:"username" yaml:"username"`         // 用户名（可选）
	Password    string        `json:"password" yaml:"password"`         // 密码（可选）
}

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	if len(c.Endpoints) == 0 {
		return ErrEtcdEndpointsEmpty
	}

	if c.DialTimeout <= 0 {
		c.DialTimeout = 5 * time.Second
	}

	return nil
}
