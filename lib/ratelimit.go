package lib

import (
	"context"
	"sync"
	"time"
)

// RateLimiter 速率限制器接口
type RateLimiter interface {
	Allow(ctx context.Context) bool
	Wait(ctx context.Context) error
}

// TokenBucketRateLimiter 基于令牌桶的速率限制器
type TokenBucketRateLimiter struct {
	rate       float64 // 令牌产生速率（个/秒）
	capacity   int     // 桶容量
	tokens     float64 // 当前令牌数
	lastRefill time.Time
	mu         sync.Mutex
}

// GetRate 获取速率
func (tb *TokenBucketRateLimiter) GetRate() float64 {
	return tb.rate
}

// GetCapacity 获取容量
func (tb *TokenBucketRateLimiter) GetCapacity() int {
	return tb.capacity
}

// NewTokenBucketRateLimiter 创建新的令牌桶速率限制器
func NewTokenBucketRateLimiter(rate float64, capacity int) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     float64(capacity),
		lastRefill: time.Now(),
	}
}

// refill 补充令牌
func (tb *TokenBucketRateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()

	// 计算需要补充的令牌数
	newTokens := elapsed * tb.rate
	tb.tokens = min(float64(tb.capacity), tb.tokens+newTokens)
	tb.lastRefill = now
}

// Allow 检查是否允许请求（非阻塞）
func (tb *TokenBucketRateLimiter) Allow(ctx context.Context) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

// Wait 等待直到允许请求（阻塞）
func (tb *TokenBucketRateLimiter) Wait(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if tb.Allow(ctx) {
				return nil
			}
			time.Sleep(10 * time.Millisecond) // 短暂等待后重试
		}
	}
}

// min 返回两个float64中的较小值
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
