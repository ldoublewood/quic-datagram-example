# QUIC Datagram 性能测试工具

这是一个基于 Go 和 quic-go 库的 QUIC datagram 模式性能测试工具，用于评估数据传输性能、延迟和丢包情况。

## 功能特性

- **服务端**: 接收 datagram 数据包并统计性能指标
- **客户端**: 发送可配置的数据包进行性能测试
- **性能指标**:
  - 数据传输延迟（最小/最大/平均）
  - 丢包检测和丢包率统计
  - 吞吐量统计
  - 实时性能监控

## 安装依赖

```bash
go mod tidy
```

## 使用方法

### 1. 启动服务端

```bash
go run server.go
```

服务端将在 `localhost:8080` 启动，并每5秒打印一次统计信息。

### 2. 启动客户端

基本用法：
```bash
go run client.go
```

自定义参数：
```bash
go run client.go -server localhost:8080 -size 1024 -rate 100 -duration 30s -payload random
```

### 客户端参数说明

- `-server`: 服务器地址 (默认: localhost:8080)
- `-size`: 数据包大小，字节 (默认: 1024)
- `-rate`: 发送速率，包/秒 (默认: 100)
- `-duration`: 测试持续时间 (默认: 30s)
- `-payload`: 负载类型 (random/sequential，默认: random)

## 测试场景示例

### 高频小包测试
```bash
go run client.go -size 64 -rate 1000 -duration 60s
```

### 大包吞吐量测试
```bash
go run client.go -size 8192 -rate 50 -duration 30s
```

### 长时间稳定性测试
```bash
go run client.go -size 1024 -rate 200 -duration 300s
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
2. 高发送速率可能受到系统网络缓冲区限制
3. 测试结果会受到网络条件和系统负载影响
4. 建议在相同网络环境下进行多次测试以获得稳定结果