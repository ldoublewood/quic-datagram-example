#!/bin/bash

echo "=== QUIC Datagram 性能测试 ==="
echo

# 构建程序
echo "1. 构建程序..."
go mod tidy
go build -o bin/quic-test *.go
if [ $? -ne 0 ]; then
    echo "构建失败！"
    exit 1
fi

echo "2. 程序构建成功！"
echo

echo "=== Native模式使用说明 ==="
echo "1. 在一个终端运行服务端:"
echo "   ./bin/quic-test -mode native -addr 0.0.0.0:4363"
echo
echo "2. 在另一个终端运行客户端:"
echo "   ./bin/quic-test -mode native -server localhost:4363"
echo

echo "=== LibP2P模式使用说明 ==="
echo "1. 在一个终端运行服务端:"
echo "   ./bin/quic-test -mode libp2p -listen /ip4/0.0.0.0/udp/4363/quic-v1"
echo
echo "2. 记录服务端显示的完整multiaddr地址"
echo
echo "3. 在另一个终端运行客户端（替换为实际地址）:"
echo "   ./bin/quic-test -mode libp2p -peer /ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooW..."
echo

echo "=== 测试示例 ==="
echo "Native模式:"
echo "  - 基本测试: ./bin/quic-test -mode native -server localhost:4363"
echo "  - 高频小包: ./bin/quic-test -mode native -server localhost:4363 -size 64 -rate 1000 -duration 10s"
echo "  - 大包测试: ./bin/quic-test -mode native -server localhost:4363 -size 8192 -rate 50 -duration 20s"
echo
echo "LibP2P模式:"
echo "  - 基本测试: ./bin/quic-test -mode libp2p -peer <multiaddr>"
echo "  - 性能测试: ./bin/quic-test -mode libp2p -peer <multiaddr> -size 1024 -rate 200 -duration 30s"
echo

echo "=== 配置文件 ==="
echo "配置目录: ~/.quic-datagram-test/"
echo "  - private_key: 节点私钥"
echo "  - peer_id: 节点Peer ID"
echo "  - bootstrap_nodes: Bootstrap节点列表"
echo

echo "程序已准备就绪！"
echo "更多信息请查看 README.md 或运行 ./examples.sh"
