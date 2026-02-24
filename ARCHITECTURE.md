# 架构说明

## 项目结构

```
.
├── client.go              # 客户端主程序
├── server.go              # 服务端主程序
├── connection.go          # 通用连接接口和Native实现
├── libp2p_connection.go   # LibP2P连接实现
├── libp2p_transport.go    # LibP2P自定义Transport
├── config.go              # 配置管理（密钥、Peer ID、Bootstrap节点）
├── stats.go               # 统计信息处理
├── go.mod                 # Go模块依赖
├── Makefile               # 构建和运行脚本
├── README.md              # 使用文档
├── ARCHITECTURE.md        # 架构说明（本文件）
├── examples.sh            # 使用示例脚本
└── quick-test.sh          # 快速测试脚本
```

## 核心设计

### 1. 连接抽象层

为了支持Native和LibP2P两种模式，我们定义了统一的`Connection`接口：

```go
type Connection interface {
    SendDatagram(data []byte) error
    ReceiveDatagram(ctx context.Context) ([]byte, error)
    Close() error
    RemoteAddr() string
}
```

#### Native实现 (NativeConnection)
- 直接包装`quic.Connection`
- 通过IP地址建立连接
- 使用TLS进行安全传输

#### LibP2P实现 (LibP2PConnection)
- 包装libp2p的`network.Conn`
- 通过Peer ID建立连接
- 使用libp2p的安全传输层

### 2. LibP2P Transport层

为了在libp2p中使用QUIC datagram功能，我们实现了自定义的Transport：

```go
type DatagramTransport struct {
    tpt.Transport
    connManager *quicreuse.ConnManager
}
```

这个Transport包装了libp2p的标准QUIC transport，并在连接建立时注入datagram支持。

### 3. 配置管理

配置文件存储在`~/.quic-datagram-test/`目录下：

- `private_key`: Ed25519私钥（base64编码）
- `peer_id`: 从私钥派生的Peer ID
- `bootstrap_nodes`: Bootstrap节点列表（每行一个multiaddr）

首次运行时自动生成，后续运行重用相同配置。

### 4. 性能测试逻辑

客户端和服务端共享相同的性能测试逻辑：

#### 数据包格式
```
[0-7]   序列号 (uint64, big-endian)
[8-15]  时间戳 (int64 nanoseconds, big-endian)
[16-]   负载数据
```

#### 统计指标
- **客户端**: 发送包数、错误数、发送速率、总数据量
- **服务端**: 接收包数、丢包数、延迟统计（最小/最大/平均）

### 5. 模式切换

通过`-mode`参数在Native和LibP2P模式间切换：

```bash
# Native模式
go run *.go -mode native -server localhost:4363

# LibP2P模式
go run *.go -mode libp2p -peer /ip4/.../p2p/...
```

## 代码重用

为了最大化代码重用，我们将以下功能抽象为独立模块：

1. **stats.go**: 统计信息处理，两种模式共享
2. **connection.go**: 连接接口定义，统一API
3. **config.go**: 配置管理，libp2p模式使用

## 依赖关系

```
client.go / server.go
    ↓
connection.go (接口)
    ↓
├── NativeConnection (quic-go)
└── LibP2PConnection (libp2p)
        ↓
    libp2p_transport.go
        ↓
    libp2p_connection.go
```

## 扩展性

### 添加新的连接模式

1. 实现`Connection`接口
2. 在client.go和server.go中添加新的连接方法
3. 更新命令行参数解析

### 添加新的统计指标

1. 在`stats.go`中添加新字段
2. 在`ProcessPacket`方法中更新统计逻辑
3. 在`Print`方法中显示新指标

## 性能考虑

1. **并发安全**: 所有统计信息使用`sync.RWMutex`保护
2. **内存分配**: 数据包payload预分配，避免频繁GC
3. **时间精度**: 使用纳秒级时间戳进行延迟测量
4. **丢包检测**: 基于序列号的简单丢包检测算法

## 测试建议

1. **本地测试**: 使用localhost进行基本功能验证
2. **网络测试**: 在不同网络环境下测试性能
3. **压力测试**: 逐步增加发送速率，找到系统瓶颈
4. **长时间测试**: 运行数小时以验证稳定性

## 已知限制

1. LibP2P模式需要手动配置bootstrap节点
2. 丢包检测基于序列号，不处理乱序包
3. 没有实现自动重连机制
4. 统计信息只在内存中，不持久化
