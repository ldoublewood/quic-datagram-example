# 使用指南

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. Native模式测试

这是最简单的测试方式，直接通过IP地址连接。

**终端1 - 启动服务端:**
```bash
go run *.go -mode native -addr 0.0.0.0:4363
```

**终端2 - 启动客户端:**
```bash
go run *.go -mode native -server localhost:4363 -size 1024 -rate 100 -duration 30s
```

### 3. LibP2P模式测试

通过Peer ID建立连接，支持NAT穿透和更复杂的网络拓扑。

**终端1 - 启动服务端:**
```bash
go run *.go -mode libp2p -listen /ip4/0.0.0.0/udp/4363/quic-v1
```

服务端会显示类似以下信息：
```
已加载现有密钥，Peer ID: 12D3KooWExamplePeerID123456789
LibP2P QUIC Datagram 服务器启动
Peer ID: 12D3KooWExamplePeerID123456789
监听地址:
  /ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooWExamplePeerID123456789
  /ip4/192.168.1.100/udp/4363/quic-v1/p2p/12D3KooWExamplePeerID123456789
```

**终端2 - 启动客户端（使用上面显示的完整地址）:**
```bash
go run *.go -mode libp2p \
  -peer /ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooWExamplePeerID123456789 \
  -size 1024 -rate 100 -duration 30s
```

## 配置文件

### 配置目录位置

配置文件存储在: `~/.quic-datagram-test/`

### 配置文件说明

1. **private_key** - 节点私钥（base64编码）
   - 首次运行自动生成
   - 用于libp2p身份认证

2. **peer_id** - 节点Peer ID
   - 从私钥派生
   - 用于libp2p网络中的节点标识

3. **bootstrap_nodes** - Bootstrap节点列表
   - 每行一个multiaddr地址
   - 用于libp2p节点发现（当前版本未使用，预留）

### 查看配置

```bash
# 查看Peer ID
cat ~/.quic-datagram-test/peer_id

# 查看所有配置文件
ls -la ~/.quic-datagram-test/
```

### 重置配置

如果需要重新生成密钥和Peer ID：

```bash
rm -rf ~/.quic-datagram-test/
# 下次运行时会自动生成新配置
```

## 命令行参数

### 通用参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-mode` | 连接模式: native 或 libp2p | native |
| `-size` | 数据包大小（字节） | 1024 |
| `-rate` | 发送速率（包/秒） | 100 |
| `-duration` | 测试持续时间 | 30s |
| `-payload` | 负载类型: random 或 sequential | random |

### Native模式参数

**服务端:**
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-addr` | 监听地址 | 0.0.0.0:4363 |

**客户端:**
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-server` | 服务器地址 | localhost:4363 |

### LibP2P模式参数

**服务端:**
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-listen` | 监听的multiaddr | /ip4/0.0.0.0/udp/4363/quic-v1 |

**客户端:**
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-peer` | 目标节点的完整multiaddr | (必需) |

## 测试场景

### 1. 延迟测试（低频率）

测试网络延迟，适合评估基础网络质量。

```bash
# Native模式
go run *.go -mode native -server localhost:4363 -size 256 -rate 10 -duration 60s

# LibP2P模式
go run *.go -mode libp2p -peer <multiaddr> -size 256 -rate 10 -duration 60s
```

### 2. 吞吐量测试（大包）

测试最大数据传输能力。

```bash
# Native模式
go run *.go -mode native -server localhost:4363 -size 8192 -rate 100 -duration 60s

# LibP2P模式
go run *.go -mode libp2p -peer <multiaddr> -size 8192 -rate 100 -duration 60s
```

### 3. 高频小包测试

测试高并发场景下的性能。

```bash
# Native模式
go run *.go -mode native -server localhost:4363 -size 64 -rate 1000 -duration 60s

# LibP2P模式
go run *.go -mode libp2p -peer <multiaddr> -size 64 -rate 1000 -duration 60s
```

### 4. 压力测试

找到系统瓶颈。

```bash
# Native模式
go run *.go -mode native -server localhost:4363 -size 1024 -rate 2000 -duration 120s

# LibP2P模式
go run *.go -mode libp2p -peer <multiaddr> -size 1024 -rate 2000 -duration 120s
```

### 5. 长时间稳定性测试

验证长时间运行的稳定性。

```bash
# Native模式 - 运行1小时
go run *.go -mode native -server localhost:4363 -size 1024 -rate 100 -duration 3600s

# LibP2P模式 - 运行1小时
go run *.go -mode libp2p -peer <multiaddr> -size 1024 -rate 100 -duration 3600s
```

## 使用Makefile

项目提供了Makefile简化常用操作。

### 查看所有命令

```bash
make help
```

### Native模式

```bash
# 启动服务端
make server-native

# 启动客户端
make client-native

# 小包测试
make test-native-small

# 大包测试
make test-native-large
```

### LibP2P模式

```bash
# 启动服务端
make server-libp2p

# 启动客户端（需要指定PEER参数）
make client-libp2p PEER=/ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooW...
```

## 输出解读

### 服务端输出

```
=== 统计信息 ===
接收包数: 3000
丢失包数: 5
丢包率: 0.17%
平均延迟: 1.234ms
最小延迟: 0.856ms
最大延迟: 5.678ms
================
```

- **接收包数**: 成功接收的数据包总数
- **丢失包数**: 检测到的丢包数量（基于序列号）
- **丢包率**: 丢包数 / (接收包数 + 丢包数)
- **平均延迟**: 所有包的平均往返延迟
- **最小/最大延迟**: 延迟的范围

### 客户端输出

```
=== 客户端发送统计 ===
发送包数: 3000
发送错误: 0
实际发送速率: 99.87 pps
总发送时间: 30.039s
总数据量: 2.93 MB
=====================
```

- **发送包数**: 成功发送的数据包总数
- **发送错误**: 发送失败的次数
- **实际发送速率**: 实际达到的发送速率（可能略低于目标速率）
- **总发送时间**: 实际测试时长
- **总数据量**: 发送的总字节数

## 故障排查

### 问题1: 编译失败

```bash
# 清理并重新安装依赖
go clean -modcache
go mod tidy
go build *.go
```

### 问题2: Native模式连接失败

检查：
1. 服务端是否正在运行
2. 防火墙是否允许UDP流量
3. 端口是否被占用

```bash
# 检查端口占用
netstat -an | grep 4363
# 或
lsof -i :4363
```

### 问题3: LibP2P模式连接失败

检查：
1. multiaddr地址是否正确（包含完整的/p2p/部分）
2. 网络是否可达
3. 配置文件是否正确生成

```bash
# 验证配置
ls -la ~/.quic-datagram-test/
cat ~/.quic-datagram-test/peer_id
```

### 问题4: 高丢包率

可能原因：
1. 发送速率过高，超过网络容量
2. 系统缓冲区不足
3. CPU负载过高

解决方法：
1. 降低发送速率 (`-rate`)
2. 增加系统UDP缓冲区
3. 减少其他系统负载

### 问题5: 延迟异常高

可能原因：
1. 网络拥塞
2. 系统时钟不准确
3. CPU调度延迟

解决方法：
1. 在更好的网络环境测试
2. 使用NTP同步时钟
3. 降低系统负载

## 高级用法

### 跨主机测试

**主机A（服务端）:**
```bash
# Native模式
go run *.go -mode native -addr 0.0.0.0:4363

# LibP2P模式
go run *.go -mode libp2p -listen /ip4/0.0.0.0/udp/4363/quic-v1
```

**主机B（客户端）:**
```bash
# Native模式（替换为主机A的IP）
go run *.go -mode native -server 192.168.1.100:4363

# LibP2P模式（使用主机A显示的完整multiaddr）
go run *.go -mode libp2p -peer /ip4/192.168.1.100/udp/4363/quic-v1/p2p/12D3KooW...
```

### 编译为独立二进制

```bash
# 编译
go build -o quic-test *.go

# 运行
./quic-test -mode native -server localhost:4363
```

### 性能调优

1. **增加系统UDP缓冲区:**
```bash
# Linux
sudo sysctl -w net.core.rmem_max=26214400
sudo sysctl -w net.core.wmem_max=26214400
```

2. **使用更高的发送速率:**
```bash
go run *.go -mode native -server localhost:4363 -rate 5000
```

3. **调整包大小以匹配MTU:**
```bash
# 以太网MTU通常是1500，减去IP和UDP头部
go run *.go -mode native -server localhost:4363 -size 1400
```

## 脚本工具

### examples.sh
显示各种使用示例：
```bash
./examples.sh
```

### quick-test.sh
快速验证安装和基本功能：
```bash
./quick-test.sh
```

### test.sh
构建程序并显示使用说明：
```bash
./test.sh
```

## 更多信息

- 架构说明: [ARCHITECTURE.md](ARCHITECTURE.md)
- 项目README: [README.md](README.md)
