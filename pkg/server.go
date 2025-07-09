package pkg

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"ratelimit-server/lib"
)

// Server HTTP服务器结构
type Server struct {
	addr      string
	rateLimit *lib.TokenBucketRateLimiter
	server    *http.Server
}

// NewServer 创建新的HTTP服务器
func NewServer(addr string, rateLimit float64, capacity int) *Server {
	return &Server{
		addr:      addr,
		rateLimit: lib.NewTokenBucketRateLimiter(rateLimit, capacity),
	}
}

// rateLimitMiddleware 速率限制中间件
func (s *Server) rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// 检查是否允许请求
		if !s.rateLimit.Allow(ctx) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	}
}

// handleRoot 处理根路径请求
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "This is a server with rate limit - Request processed at %s", time.Now().Format(time.RFC3339))
}

// Start 启动服务器
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// 注册路由，应用速率限制中间件
	mux.HandleFunc("/", s.rateLimitMiddleware(s.handleRoot))

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	fmt.Printf("Server starting on %s with rate limit: %.2f requests/second, capacity: %d\n",
		s.addr, s.rateLimit.GetRate(), s.rateLimit.GetCapacity())

	return s.server.ListenAndServe()
}

// Stop 停止服务器
func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
