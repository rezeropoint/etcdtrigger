package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rezeropoint/etcdtrigger"
)

func main() {
	// 创建配置
	config := &etcdtrigger.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
		// 如果需要认证，可以添加以下配置
		// Username: "your_username",
		// Password: "your_password",
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建客户端
	client, err := etcdtrigger.NewEtcdClient(ctx, cancel, config)
	if err != nil {
		log.Fatal("创建etcd客户端失败:", err)
	}
	defer client.Close()

	// 订阅配置变更
	err = client.Subscribe("/app/config/", func(key string, value []byte) error {
		if value == nil {
			log.Printf("⚠️  配置被删除: %s", key)
		} else {
			log.Printf("🔄 配置变更: %s = %s", key, string(value))
		}
		return nil
	})

	if err != nil {
		log.Fatal("订阅失败:", err)
	}

	log.Println("🚀 开始监听配置变更...")
	log.Println("💡 提示: 你可以使用以下命令测试配置变更:")
	log.Println("   etcdctl put /app/config/database/host localhost")
	log.Println("   etcdctl put /app/config/database/port 3306")
	log.Println("   etcdctl del /app/config/database/host")

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 示例：写入一些配置
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("📝 写入示例配置...")

		if err := client.Put("/app/config/database/host", []byte("localhost")); err != nil {
			log.Printf("写入配置失败: %v", err)
		}

		if err := client.Put("/app/config/database/port", []byte("3306")); err != nil {
			log.Printf("写入配置失败: %v", err)
		}

		if err := client.Put("/app/config/redis/host", []byte("127.0.0.1")); err != nil {
			log.Printf("写入配置失败: %v", err)
		}
	}()

	// 等待退出信号
	<-sigChan
	log.Println("🛑 收到退出信号，正在关闭...")
}
