package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rezeropoint/etcdtrigger/v2/core"
	"github.com/rezeropoint/etcdtrigger/v2/engine"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// RedisConfig Redis 配置结构体
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func main() {
	// 1. 创建 etcd 客户端（由调用方管理生命周期）
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
		// Username: "your_username",
		// Password: "your_password",
	})
	if err != nil {
		log.Fatal("创建 etcd 客户端失败:", err)
	}
	defer etcdClient.Close()

	// 2. 创建引擎（传入 etcd 客户端）
	eng := engine.NewEngine(etcdClient, &engine.Config{
		PodName:     "example-pod",
		ServiceName: "example-service",
		// 预加载配置：自动监听并缓存到内存（Store 功能）
		Configs: []core.WatchConfig{
			{Path: "/app/config/database/", Struct: &DatabaseConfig{}},
			{Path: "/app/config/redis/", Struct: &RedisConfig{}},
		},
	})

	// ---- Watcher 功能演示（原始操作）----

	// 使用 Watch 原始回调模式
	err = eng.Watch("/app/events/", func(event *core.WatchEvent) error {
		if event.EventType.IsDelete() {
			log.Printf("[Watcher] 键被删除: %s", event.Key)
		} else {
			log.Printf("[Watcher] 键变更: %s = %s", event.Key, string(event.Value))
		}
		return nil
	})
	if err != nil {
		log.Fatal("订阅失败:", err)
	}

	// ---- Store 功能演示（强类型缓存）----

	// 使用 AddPrefixWatcher 前缀监听器
	eng.AddPrefixWatcher("/app/config/", func(key string, eventType core.EventType) {
		log.Printf("[Store] %s: %s", eventType, key)
	})

	log.Println("开始监听配置变更...")

	// 演示写入和读取
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("写入示例配置...")

		// Store: 写入数据库配置（自动 JSON 序列化）
		dbConfig := &DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "secret",
		}
		if err := eng.PutConfig(context.Background(), "/app/config/database/main", dbConfig); err != nil {
			log.Printf("写入数据库配置失败: %v", err)
		}

		// Store: 写入 Redis 配置
		redisConfig := &RedisConfig{
			Host:     "127.0.0.1",
			Port:     6379,
			Password: "",
			DB:       0,
		}
		if err := eng.PutConfig(context.Background(), "/app/config/redis/main", redisConfig); err != nil {
			log.Printf("写入 Redis 配置失败: %v", err)
		}

		// Watcher: 写入事件（原始字节）
		if err := eng.WatchPut("/app/events/test", []byte("hello world")); err != nil {
			log.Printf("写入事件失败: %v", err)
		}

		// 等待配置同步
		time.Sleep(1 * time.Second)

		// Store: 从缓存获取强类型配置
		var cachedDB DatabaseConfig
		if eng.GetConfig("/app/config/database/main", &cachedDB) {
			log.Printf("[缓存读取] 数据库配置: Host=%s, Port=%d", cachedDB.Host, cachedDB.Port)
		}

		// Store: 获取所有键
		keys := eng.GetAllKeys("/app/config/")
		log.Printf("[所有键] %v", keys)

		// Watcher: 获取原始数据
		if data, err := eng.WatchGet("/app/events/test"); err == nil {
			log.Printf("[原始读取] /app/events/test = %s", string(data))
		}
	}()

	// 等待退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("收到退出信号，正在关闭...")
}
