package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ratelimit-server/pkg"
)

func main() {
	// 统一定义所有参数
	mode := flag.String("mode", "server", "运行模式: server 或 client")
	addr := flag.String("addr", ":8080", "服务器监听地址")
	rateLimit := flag.Float64("rate", 10.0, "速率限制（请求/秒）")
	capacity := flag.Int("capacity", 20, "令牌桶容量")
	serverURL := flag.String("server", "http://localhost:8080", "服务器URL")
	concurrency := flag.Int("concurrency", 5, "并发数")
	duration := flag.Duration("duration", 30*time.Second, "测试持续时间")
	requestRate := flag.Int("rateClient", 5, "请求速率（请求/秒）")
	timeout := flag.Duration("timeout", 10*time.Second, "请求超时时间")
	flag.Parse()

	switch *mode {
	case "server":
		runServer(*addr, *rateLimit, *capacity)
	case "client":
		runClient(*serverURL, *concurrency, *duration, *requestRate, *timeout)
	default:
		fmt.Printf("未知模式: %s\n", *mode)
		fmt.Println("支持的模式: server, client")
		os.Exit(1)
	}
}

func runServer(addr string, rateLimit float64, capacity int) {
	server := pkg.NewServer(addr, rateLimit, capacity)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	<-sigChan
	fmt.Println("\n正在关闭服务器...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Stop(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	fmt.Println("服务器已关闭")
}

func runClient(serverURL string, concurrency int, duration time.Duration, requestRate int, timeout time.Duration) {
	client := pkg.NewClient(serverURL, timeout)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		client.LoadTest(ctx, concurrency, duration, requestRate)
	}()

	select {
	case <-sigChan:
		fmt.Println("\n正在停止测试...")
		cancel()
	case <-time.After(duration + 5*time.Second):
		fmt.Println("测试完成")
	}

	fmt.Println("客户端已退出")
}
