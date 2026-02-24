#!/bin/bash

# QUIC Datagram 性能测试工具 - 使用示例

echo "=== QUIC Datagram 性能测试工具使用示例 ==="
echo ""

show_example() {
    echo "----------------------------------------"
    echo "$1"
    echo "----------------------------------------"
    echo "$2"
    echo ""
}

show_example "1. Native模式 - 基本测试" \
"# 终端1: 启动服务端
go run *.go -mode native -addr 0.0.0.0:4363

# 终端2: 启动客户端
go run *.go -mode native -server localhost:4363"

show_example "2. Native模式 - 高频小包测试" \
"# 服务端
go run *.go -mode native

# 客户端
go run *.go -mode native -server localhost:4363 -size 64 -rate 1000 -duration 60s"

show_example "3. Native模式 - 大包吞吐量测试" \
"# 服务端
go run *.go -mode native

# 客户端
go run *.go -mode native -server localhost:4363 -size 8192 -rate 50 -duration 30s"

show_example "4. LibP2P模式 - 基本测试" \
"# 终端1: 启动服务端
go run *.go -mode libp2p -listen /ip4/0.0.0.0/udp/4363/quic-v1

# 记录服务端显示的完整multiaddr地址，例如：
# /ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooWExamplePeerID

# 终端2: 启动客户端（使用上面的multiaddr）
go run *.go -mode libp2p -peer /ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooW..."

show_example "5. LibP2P模式 - 性能测试" \
"# 服务端
go run *.go -mode libp2p

# 客户端（替换为实际的peer地址）
go run *.go -mode libp2p \\
  -peer /ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooW... \\
  -size 1024 -rate 200 -duration 60s"

show_example "6. 查看配置文件" \
"# 配置文件位置
ls -la ~/.quic-datagram-test/

# 查看Peer ID
cat ~/.quic-datagram-test/peer_id

# 编辑bootstrap节点
nano ~/.quic-datagram-test/bootstrap_nodes"

show_example "7. 使用Makefile快捷命令" \
"# 查看所有可用命令
make help

# Native模式服务端
make server-native

# Native模式客户端
make client-native

# LibP2P模式服务端
make server-libp2p

# LibP2P模式客户端
make client-libp2p PEER=/ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooW..."

echo "更多信息请查看 README.md"
