#!/bin/bash

echo "=== QUIC Datagram 性能测试 ==="
echo

# 构建程序
echo "1. 构建程序..."
make build
if [ $? -ne 0 ]; then
    echo "构建失败！"
    exit 1
fi

echo "2. 程序构建成功！"
echo

echo "使用说明："
echo "1. 在一个终端运行: ./bin/server"
echo "2. 在另一个终端运行: ./bin/client"
echo
echo "测试示例："
echo "- 基本测试: ./bin/client"
echo "- 高频小包: ./bin/client -size 64 -rate 500 -duration 10s"
echo "- 大包测试: ./bin/client -size 8192 -rate 50 -duration 20s"
echo "- 延迟测试: ./bin/client -size 256 -rate 10 -duration 30s"
echo

echo "程序已准备就绪！"