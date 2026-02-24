#!/bin/bash

# 快速测试脚本 - 验证Native和LibP2P模式

set -e

echo "=== QUIC Datagram 快速测试 ==="
echo ""

# 测试编译
echo "1. 测试编译..."
go build -o /tmp/quic-test *.go
if [ $? -eq 0 ]; then
    echo "✓ 编译成功"
else
    echo "✗ 编译失败"
    exit 1
fi
echo ""

# 测试配置文件生成
echo "2. 测试配置文件生成..."
/tmp/quic-test -mode libp2p -listen /ip4/127.0.0.1/udp/0/quic-v1 &
SERVER_PID=$!
sleep 2
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

if [ -f ~/.quic-datagram-test/peer_id ]; then
    echo "✓ 配置文件已生成"
    echo "  Peer ID: $(cat ~/.quic-datagram-test/peer_id)"
else
    echo "✗ 配置文件生成失败"
    exit 1
fi
echo ""

# 测试Native模式
echo "3. 测试Native模式..."
echo "  启动服务端..."
/tmp/quic-test -mode native -addr 127.0.0.1:14363 > /tmp/native-server.log 2>&1 &
NATIVE_SERVER_PID=$!
sleep 2

echo "  启动客户端..."
timeout 5 /tmp/quic-test -mode native -server 127.0.0.1:14363 -size 512 -rate 10 -duration 3s > /tmp/native-client.log 2>&1 || true

kill $NATIVE_SERVER_PID 2>/dev/null || true
wait $NATIVE_SERVER_PID 2>/dev/null || true

if grep -q "连接成功" /tmp/native-client.log; then
    echo "✓ Native模式测试成功"
else
    echo "✗ Native模式测试失败"
    echo "服务端日志:"
    cat /tmp/native-server.log
    echo "客户端日志:"
    cat /tmp/native-client.log
fi
echo ""

echo "=== 测试完成 ==="
echo ""
echo "查看详细日志:"
echo "  服务端: cat /tmp/native-server.log"
echo "  客户端: cat /tmp/native-client.log"
echo ""
echo "配置目录: ~/.quic-datagram-test/"

# 清理
rm -f /tmp/quic-test
