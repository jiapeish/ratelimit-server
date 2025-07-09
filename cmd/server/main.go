package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ratelimit-server/pkg"
)

func main() {
	// 命令行参数
	var (
		addr      = flag.String("addr", ":8080", "服务器监听地址")
		rateLimit = flag.Float64("rate", 30.0, "速率限制（请求/秒）")
		capacity  = flag.Int("capacity", 50, "令牌桶容量")
	)
	flag.Parse()

	// 创建服务器
	server := pkg.NewServer(*addr, *rateLimit, *capacity)

	// 创建上下文用于优雅关闭
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器（在goroutine中）
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// 等待信号
	<-sigChan
	fmt.Println("\n正在关闭服务器...")

	// 优雅关闭
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5)
	defer shutdownCancel()

	if err := server.Stop(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	fmt.Println("服务器已关闭")
}
