# 贡献指南

感谢你对本项目的关注！

## 项目结构

```
.
├── client.go              # 客户端主程序
├── server.go              # 服务端主程序
├── connection.go          # 连接接口定义
├── libp2p_connection.go   # LibP2P连接实现
├── libp2p_transport.go    # LibP2P Transport
├── config.go              # 配置管理
├── stats.go               # 统计处理
├── version.go             # 版本信息
└── go.mod                 # 依赖管理
```

## 开发环境

### 要求

- Go 1.21 或更高版本
- 基本的网络知识
- 了解QUIC和libp2p（可选）

### 设置

```bash
# 克隆项目
git clone <repository-url>
cd quic-datagram-test

# 安装依赖
go mod tidy

# 运行测试
./quick-test.sh
```

## 代码规范

### Go代码风格

遵循标准的Go代码规范：

```bash
# 格式化代码
go fmt ./...

# 检查代码
go vet ./...
```

### 命名约定

- 接口名：使用名词，如 `Connection`
- 实现类：描述性名称，如 `NativeConnection`
- 方法名：使用动词，如 `SendDatagram`
- 变量名：简洁明了，如 `conn`, `stats`

### 注释

- 所有导出的类型和函数都应有注释
- 复杂逻辑应添加解释性注释
- 使用中文或英文，保持一致

## 添加新功能

### 1. 添加新的连接模式

如果要添加新的连接方式（如WebTransport）：

1. 实现`Connection`接口：
```go
type MyConnection struct {
    // 你的字段
}

func (c *MyConnection) SendDatagram(data []byte) error {
    // 实现
}

func (c *MyConnection) ReceiveDatagram(ctx context.Context) ([]byte, error) {
    // 实现
}

func (c *MyConnection) Close() error {
    // 实现
}

func (c *MyConnection) RemoteAddr() string {
    // 实现
}
```

2. 在`client.go`和`server.go`中添加连接方法

3. 更新命令行参数解析

4. 更新文档

### 2. 添加新的统计指标

在`stats.go`中：

1. 添加新字段到`ClientStats`或`ServerStats`
2. 在`ProcessPacket`中更新统计逻辑
3. 在`Print`方法中显示新指标

### 3. 添加新的配置选项

在`config.go`中：

1. 添加新字段到`Config`结构
2. 实现加载和保存逻辑
3. 更新文档

## 测试

### 单元测试

```bash
go test ./...
```

### 集成测试

```bash
./quick-test.sh
```

### 手动测试

1. Native模式测试
2. LibP2P模式测试
3. 不同参数组合测试
4. 跨主机测试

## 提交代码

### 提交信息格式

```
<类型>: <简短描述>

<详细描述>

<相关Issue>
```

类型：
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建/工具相关

示例：
```
feat: 添加WebTransport支持

实现了基于WebTransport的连接模式，支持浏览器客户端。

Closes #123
```

### Pull Request流程

1. Fork项目
2. 创建特性分支：`git checkout -b feature/my-feature`
3. 提交更改：`git commit -am 'feat: 添加新功能'`
4. 推送分支：`git push origin feature/my-feature`
5. 创建Pull Request

### PR检查清单

- [ ] 代码已格式化（`go fmt`）
- [ ] 代码已检查（`go vet`）
- [ ] 添加了必要的测试
- [ ] 更新了相关文档
- [ ] 提交信息清晰明了
- [ ] 没有破坏现有功能

## 报告问题

### Bug报告

请包含以下信息：

1. **环境信息**
   - 操作系统和版本
   - Go版本
   - 项目版本

2. **重现步骤**
   - 详细的操作步骤
   - 使用的命令和参数

3. **预期行为**
   - 你期望发生什么

4. **实际行为**
   - 实际发生了什么
   - 错误信息或日志

5. **其他信息**
   - 截图（如果适用）
   - 相关配置文件

### 功能请求

请描述：

1. **需求背景**
   - 为什么需要这个功能
   - 使用场景

2. **建议方案**
   - 你认为应该如何实现

3. **替代方案**
   - 其他可能的实现方式

## 文档贡献

文档同样重要！你可以：

- 修正错别字
- 改进说明
- 添加示例
- 翻译文档

## 代码审查

我们欢迎代码审查！在审查时请关注：

1. **功能正确性**
   - 代码是否实现了预期功能
   - 是否有边界情况未处理

2. **代码质量**
   - 是否遵循Go最佳实践
   - 是否有重复代码
   - 错误处理是否完善

3. **性能**
   - 是否有性能问题
   - 是否有不必要的内存分配

4. **可维护性**
   - 代码是否易于理解
   - 是否有足够的注释
   - 接口设计是否合理

## 社区准则

- 尊重他人
- 建设性反馈
- 保持专业
- 欢迎新手

## 获取帮助

如有问题，可以：

1. 查看文档：
   - [README.md](README.md)
   - [USAGE_GUIDE.md](USAGE_GUIDE.md)
   - [ARCHITECTURE.md](ARCHITECTURE.md)

2. 查看示例：
   - 运行 `./examples.sh`
   - 查看测试脚本

3. 提交Issue

## 许可证

贡献的代码将采用与项目相同的许可证。

---

再次感谢你的贡献！
