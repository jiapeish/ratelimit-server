# HTTP 服务器/客户端 速率限制测试系统

这是一个基于 Go 的 HTTP 服务器和客户端系统，用于演示和测试速率限制功能。

## 功能特性

### 服务器端
- 基于令牌桶算法的速率限制
- 支持 GET 请求处理
- 可配置的速率限制参数
- 优雅关闭支持

### 客户端端
- 高并发负载测试
- 可配置的请求速率
- 实时统计和报告
- 支持测试持续时间设置

## 项目结构

```
server/
├── cmd/
│   ├── server/
│   │   └── main.go          # 独立服务器启动文件
│   └── client/
│       └── main.go          # 独立客户端启动文件
├── lib/
│   └── ratelimit.go         # 速率限制库
├── pkg/
│   ├── server.go            # 服务器核心代码
│   └── client.go            # 客户端核心代码
├── main.go                  # 主启动文件（支持服务器/客户端模式）
├── go.mod                   # Go 模块文件
└── README.md               # 项目说明
```

## 使用方法

### 方法一：使用主启动文件（推荐）

#### 启动服务器
```bash
# 使用默认参数启动服务器
go run main.go -mode=server

# 自定义参数启动服务器
go run main.go -mode=server -addr=:9090 -rate=5.0 -capacity=10
```

#### 启动客户端进行测试
```bash
# 使用默认参数测试
go run main.go -mode=client

# 自定义参数测试
go run main.go -mode=client -server=http://localhost:9090 -concurrency=10 -rate=50 -duration=60s
```

### 方法二：使用独立的启动文件

#### 启动服务器
```bash
# 进入服务器目录
cd cmd/server
go run main.go

# 或者从根目录运行
go run cmd/server/main.go -addr=:9090 -rate=5.0 -capacity=10
```

#### 启动客户端
```bash
# 进入客户端目录
cd cmd/client
go run main.go

# 或者从根目录运行
go run cmd/client/main.go -server=http://localhost:9090 -concurrency=10 -rate=50
```

## 命令行参数

### 服务器参数
- `-mode=server`: 运行模式（服务器）
- `-addr=:8080`: 服务器监听地址（默认: :8080）
- `-rate=10.0`: 速率限制，每秒允许的请求数（默认: 10.0）
- `-capacity=20`: 令牌桶容量（默认: 20）

### 客户端参数
- `-mode=client`: 运行模式（客户端）
- `-server=http://localhost:8080`: 目标服务器URL（默认: http://localhost:8080）
- `-concurrency=5`: 并发数（默认: 5）
- `-duration=30s`: 测试持续时间（默认: 30s）
- `-rate=20`: 请求速率，每秒发送的请求数（默认: 20）
- `-timeout=10s`: 请求超时时间（默认: 10s）

## 示例测试场景

### 场景1：正常速率测试
```bash
# 启动服务器，限制每秒5个请求
go run main.go -mode=server -rate=5.0 -capacity=10

# 在另一个终端启动客户端，每秒发送3个请求
go run main.go -mode=client -rate=3 -concurrency=2
```

### 场景2：超速率测试
```bash
# 启动服务器，限制每秒5个请求
go run main.go -mode=server -rate=5.0 -capacity=10

# 在另一个终端启动客户端，每秒发送20个请求
go run main.go -mode=client -rate=20 -concurrency=5
```

### 场景3：高并发测试
```bash
# 启动服务器，限制每秒10个请求
go run main.go -mode=server -rate=10.0 -capacity=20

# 在另一个终端启动客户端，高并发测试
go run main.go -mode=client -rate=50 -concurrency=20 -duration=60s
```

## 速率限制算法

系统使用令牌桶算法实现速率限制：

1. **令牌桶容量**: 最大可存储的令牌数
2. **令牌产生速率**: 每秒产生的令牌数
3. **请求处理**: 每个请求消耗一个令牌
4. **限流机制**: 当令牌不足时，返回 429 状态码

## 输出说明

### 服务器输出
- 启动信息：显示监听地址和速率限制参数
- 请求日志：显示处理的请求信息

### 客户端输出
- 测试配置：显示测试参数
- 实时日志：显示每个请求的结果
- 统计结果：显示成功、失败和限流的请求数量

## 注意事项

1. 确保在测试前先启动服务器
2. 客户端会自动连接到指定的服务器地址
3. 使用 Ctrl+C 可以优雅关闭服务器或客户端
4. 速率限制参数需要根据实际需求调整
5. 高并发测试时注意系统资源使用情况 