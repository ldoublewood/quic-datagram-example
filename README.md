# QUIC Datagram 性能测试工具

这是一个基于 Go 和 quic-go 库的 QUIC datagram 模式性能测试工具，支持原生QUIC和libp2p两种连接方式，用于评估数据传输性能、延迟和丢包情况。

## 功能特性

- **双模式支持**: 
  - Native模式：直接通过IP地址建立QUIC连接
  - LibP2P模式：通过Peer ID建立libp2p QUIC连接
- **服务端**: 接收 datagram 数据包并统计性能指标
- **客户端**: 发送可配置的数据包进行性能测试
- **自动配置管理**: 
  - 自动生成和管理公私钥对
  - 自动生成Peer ID
  - 支持bootstrap节点配置
- **性能指标**:
  - 数据传输延迟（最小/最大/平均）
  - 丢包检测和丢包率统计
  - 吞吐量统计
  - 实时性能监控

## 安装依赖

```bash
go mod tidy
```

## 配置文件

首次运行时，程序会在用户目录下创建 `~/.quic-datagram-test/` 配置目录，包含：

- `private_key`: 节点私钥（base64编码）
- `peer_id`: 节点Peer ID
- `bootstrap_nodes`: Bootstrap节点列表（libp2p模式使用）

### Bootstrap节点配置

编辑 `~/.quic-datagram-test/bootstrap_nodes` 文件，每行一个multiaddr地址：

```
/ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooWExamplePeerID
/ip4/192.168.1.100/udp/4363/quic-v1/p2p/12D3KooWAnotherPeerID
```

## 使用方法

### Native模式

#### 1. 启动服务端

```bash
go run *.go -mode native
# 或指定监听地址
go run *.go -mode native -addr 0.0.0.0:4363
```

#### 2. 启动客户端

```bash
go run *.go -mode native -server localhost:4363
# 自定义参数
go run *.go -mode native -server localhost:4363 -size 1024 -rate 100 -duration 30s
```

### LibP2P模式

#### 1. 启动服务端

```bash
go run *.go -mode libp2p
# 或指定监听地址
go run *.go -mode libp2p -listen /ip4/0.0.0.0/udp/4363/quic-v1
```

服务端启动后会显示Peer ID和完整的multiaddr地址。

#### 2. 启动客户端

```bash
# 使用服务端显示的完整multiaddr地址
go run *.go -mode libp2p -peer /ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooW...
# 自定义参数
go run *.go -mode libp2p -peer /ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooW... -size 1024 -rate 100
```

## 参数说明

### 通用参数

- `-mode`: 连接模式，`native` 或 `libp2p` (默认: native)
- `-size`: 数据包大小，字节 (默认: 1024)
- `-rate`: 发送速率，包/秒 (默认: 100)
- `-duration`: 测试持续时间 (默认: 30s)
- `-payload`: 负载类型 (random/sequential，默认: random)

### Native模式参数

- `-addr`: 服务端监听地址 (默认: 0.0.0.0:4363)
- `-server`: 客户端连接的服务器地址 (默认: localhost:4363)

### LibP2P模式参数

- `-listen`: 服务端监听的multiaddr (默认: /ip4/0.0.0.0/udp/4363/quic-v1)
- `-peer`: 客户端连接的目标节点multiaddr (必需)

## 测试场景示例

### Native模式测试

高频小包测试：
```bash
# 服务端
go run *.go -mode native

# 客户端
go run *.go -mode native -server localhost:4363 -size 64 -rate 1000 -duration 60s
```

大包吞吐量测试：
```bash
go run *.go -mode native -server localhost:4363 -size 8192 -rate 50 -duration 30s
```

### LibP2P模式测试

```bash
# 服务端
go run *.go -mode libp2p
# 记录显示的multiaddr，例如：/ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooW...

# 客户端
go run *.go -mode libp2p -peer /ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooW... -size 1024 -rate 200
```

## 输出指标

### 服务端统计
- 接收包数
- 丢失包数和丢包率
- 延迟统计（最小/最大/平均）

### 客户端统计
- 发送包数和错误数
- 实际发送速率
- 总数据量

## 注意事项

1. 确保防火墙允许 UDP 流量通过
2. LibP2P模式需要正确配置multiaddr地址
3. 高发送速率可能受到系统网络缓冲区限制
4. 测试结果会受到网络条件和系统负载影响
5. 建议在相同网络环境下进行多次测试以获得稳定结果
6. 首次运行会自动生成配置文件，后续运行会重用相同的密钥和Peer ID