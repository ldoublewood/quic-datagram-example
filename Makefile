# QUIC Datagram 性能测试工具 Makefile

.PHONY: build clean test server client deps

# 构建所有程序
build: deps
	go build -o bin/server server.go
	go build -o bin/client client.go

# 安装依赖
deps:
	go mod tidy

# 清理构建文件
clean:
	rm -rf bin/

# 创建bin目录
bin:
	mkdir -p bin

# 运行服务端
server: build
	./bin/server

# 运行客户端（默认参数）
client: build
	./bin/client

# 快速测试（小包高频）
test-small: build
	./bin/client -size 64 -rate 500 -duration 10s

# 吞吐量测试（大包）
test-throughput: build
	./bin/client -size 8192 -rate 100 -duration 20s

# 延迟测试（低频）
test-latency: build
	./bin/client -size 256 -rate 10 -duration 30s

# 压力测试
test-stress: build
	./bin/client -size 1024 -rate 1000 -duration 60s