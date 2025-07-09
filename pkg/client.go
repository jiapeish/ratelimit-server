package pkg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Client HTTP客户端结构
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient 创建新的HTTP客户端
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
	}
}

// SendRequest 发送单个请求
func (c *Client) SendRequest(ctx context.Context) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}

	return resp, nil
}

// LoadTest 执行负载测试
func (c *Client) LoadTest(ctx context.Context, concurrency int, duration time.Duration, requestRate int) {
	fmt.Printf("开始负载测试:\n")
	fmt.Printf("- 目标服务器: %s\n", c.baseURL)
	fmt.Printf("- 并发数: %d\n", concurrency)
	fmt.Printf("- 持续时间: %v\n", duration)
	fmt.Printf("- 请求速率: %d 请求/秒\n", requestRate)

	// 计算请求间隔
	interval := time.Second / time.Duration(requestRate)

	// 统计信息
	var (
		successCount   int64
		errorCount     int64
		rateLimitCount int64
		mu             sync.Mutex
		wg             sync.WaitGroup
	)

	// 创建定时器
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// 设置测试结束时间
	endTime := time.Now().Add(duration)

	// 启动工作协程
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if time.Now().After(endTime) {
						return
					}

					resp, err := c.SendRequest(ctx)
					if err != nil {
						mu.Lock()
						errorCount++
						mu.Unlock()
						fmt.Printf("Worker %d: 请求失败 - %v\n", workerID, err)
						continue
					}

					// 读取响应体
					body, _ := io.ReadAll(resp.Body)
					resp.Body.Close()

					mu.Lock()
					if resp.StatusCode == http.StatusTooManyRequests {
						rateLimitCount++
						fmt.Printf("Worker %d: 速率限制触发 (429) - %s\n", workerID, string(body))
					} else if resp.StatusCode == http.StatusOK {
						successCount++
						fmt.Printf("Worker %d: 请求成功 (200) - %s\n", workerID, string(body))
					} else {
						errorCount++
						fmt.Printf("Worker %d: 请求失败 (%d) - %s\n", workerID, resp.StatusCode, string(body))
					}
					mu.Unlock()
				}
			}
		}(i)
	}

	// 等待测试完成
	wg.Wait()

	// 打印统计结果
	fmt.Printf("\n=== 测试结果 ===\n")
	fmt.Printf("成功请求: %d\n", successCount)
	fmt.Printf("失败请求: %d\n", errorCount)
	fmt.Printf("速率限制触发: %d\n", rateLimitCount)
	fmt.Printf("总请求数: %d\n", successCount+errorCount+rateLimitCount)
}
