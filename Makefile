# QUIC Datagram 性能测试工具 Makefile

.PHONY: all build clean deps server-native client-native server-libp2p client-libp2p test-native-small test-native-large help

all: build

# 安装依赖
deps:
	go mod tidy

# 构建程序
build: deps
	mkdir -p bin
	go build -o bin/quic-test *.go

# 清理构建文件
clean:
	rm -rf bin/
	go clean

# Native模式 - 服务端
server-native:
	go run *.go -mode native -addr 0.0.0.0:4363

# Native模式 - 客户端
client-native:
	go run *.go -mode native -server localhost:4363 -size 1024 -rate 100 -duration 30s

# LibP2P模式 - 服务端
server-libp2p:
	go run *.go -mode libp2p -listen /ip4/0.0.0.0/udp/4363/quic-v1

# LibP2P模式 - 客户端
client-libp2p:
	@echo "请使用: make client-libp2p PEER=<multiaddr>"
	@echo "例如: make client-libp2p PEER=/ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooW..."
ifdef PEER
	go run *.go -mode libp2p -peer $(PEER) -size 1024 -rate 100 -duration 30s
endif

# 测试场景 - Native模式小包高频
test-native-small:
	go run *.go -mode native -server localhost:4363 -size 64 -rate 1000 -duration 10s

# 测试场景 - Native模式大包吞吐
test-native-large:
	go run *.go -mode native -server localhost:4363 -size 8192 -rate 50 -duration 10s

# 帮助信息
help:
	@echo "可用命令:"
	@echo "  make deps              - 安装依赖"
	@echo "  make build             - 编译程序"
	@echo "  make clean             - 清理编译文件"
	@echo ""
	@echo "Native模式:"
	@echo "  make server-native     - 启动Native模式服务端"
	@echo "  make client-native     - 启动Native模式客户端"
	@echo "  make test-native-small - Native模式小包测试"
	@echo "  make test-native-large - Native模式大包测试"
	@echo ""
	@echo "LibP2P模式:"
	@echo "  make server-libp2p     - 启动LibP2P模式服务端"
	@echo "  make client-libp2p PEER=<addr> - 启动LibP2P模式客户端"
