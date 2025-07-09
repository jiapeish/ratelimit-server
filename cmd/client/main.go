package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ratelimit-server/pkg"
)

func main() {
	// 命令行参数
	var (
		serverURL   = flag.String("server", "http://localhost:8080", "服务器URL")
		concurrency = flag.Int("concurrency", 1, "并发数")
		duration    = flag.Duration("duration", 30*time.Second, "测试持续时间")
		requestRate = flag.Int("rate", 20, "请求速率（请求/秒）")
		timeout     = flag.Duration("timeout", 10*time.Second, "请求超时时间")
	)
	flag.Parse()

	// 创建客户端
	client := pkg.NewClient(*serverURL, *timeout)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动负载测试（在goroutine中）
	go func() {
		client.LoadTest(ctx, *concurrency, *duration, *requestRate)
	}()

	// 等待信号或测试完成
	select {
	case <-sigChan:
		fmt.Println("\n正在停止测试...")
		cancel()
	case <-time.After(*duration + 5*time.Second):
		fmt.Println("测试完成")
	}

	fmt.Println("客户端已退出")
}
