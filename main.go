package main

import (
	"fmt"
	"os"
)

func main() {
	// 检查是否是服务端还是客户端模式
	// 通过检查命令行参数来判断
	
	// 先解析一次flag来判断模式
	isServer := false
	for _, arg := range os.Args[1:] {
		if arg == "-addr" || arg == "-listen" {
			isServer = true
			break
		}
		if arg == "-server" || arg == "-peer" {
			isServer = false
			break
		}
	}
	
	// 如果没有明确的参数，检查是否有-h或-help
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "-help" || os.Args[1] == "--help") {
		printHelp()
		return
	}
	
	// 如果参数中包含mode，检查后续是否有server相关参数
	for i, arg := range os.Args[1:] {
		if arg == "-mode" && i+1 < len(os.Args)-1 {
			// mode参数存在
			_ = i
		}
	}
	
	// 默认行为：如果没有明确的server/client参数，显示帮助
	if !isServer && !hasServerOrClientFlag() {
		fmt.Println("QUIC Datagram 性能测试工具")
		fmt.Println()
		fmt.Println("请指定运行模式：")
		fmt.Println()
		fmt.Println("服务端模式（使用 -addr 或 -listen）：")
		fmt.Println("  go run *.go -mode native -addr 0.0.0.0:4363")
		fmt.Println("  go run *.go -mode libp2p -listen /ip4/0.0.0.0/udp/4363/quic-v1")
		fmt.Println()
		fmt.Println("客户端模式（使用 -server 或 -peer）：")
		fmt.Println("  go run *.go -mode native -server localhost:4363")
		fmt.Println("  go run *.go -mode libp2p -peer /ip4/127.0.0.1/udp/4363/quic-v1/p2p/...")
		fmt.Println()
		fmt.Println("使用 -h 查看完整帮助")
		return
	}
	
	if isServer {
		runServer()
	} else {
		runClient()
	}
}

func hasServerOrClientFlag() bool {
	for _, arg := range os.Args[1:] {
		if arg == "-addr" || arg == "-listen" || arg == "-server" || arg == "-peer" {
			return true
		}
	}
	return false
}

func printHelp() {
	fmt.Println("QUIC Datagram 性能测试工具 v" + Version)
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  服务端: go run *.go [服务端选项]")
	fmt.Println("  客户端: go run *.go [客户端选项]")
	fmt.Println()
	fmt.Println("通用选项:")
	fmt.Println("  -mode string")
	fmt.Println("        连接模式: native 或 libp2p (默认 \"native\")")
	fmt.Println()
	fmt.Println("服务端选项 (Native模式):")
	fmt.Println("  -addr string")
	fmt.Println("        监听地址 (默认 \"0.0.0.0:4363\")")
	fmt.Println()
	fmt.Println("服务端选项 (LibP2P模式):")
	fmt.Println("  -listen string")
	fmt.Println("        监听的multiaddr (默认 \"/ip4/0.0.0.0/udp/4363/quic-v1\")")
	fmt.Println()
	fmt.Println("客户端选项 (Native模式):")
	fmt.Println("  -server string")
	fmt.Println("        服务器地址 (默认 \"localhost:4363\")")
	fmt.Println()
	fmt.Println("客户端选项 (LibP2P模式):")
	fmt.Println("  -peer string")
	fmt.Println("        目标节点的完整multiaddr (必需)")
	fmt.Println()
	fmt.Println("客户端测试选项:")
	fmt.Println("  -size int")
	fmt.Println("        数据包大小（字节） (默认 1024)")
	fmt.Println("  -rate int")
	fmt.Println("        发送速率（包/秒） (默认 100)")
	fmt.Println("  -duration duration")
	fmt.Println("        测试持续时间 (默认 30s)")
	fmt.Println("  -payload string")
	fmt.Println("        负载类型: random 或 sequential (默认 \"random\")")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println()
	fmt.Println("  Native模式:")
	fmt.Println("    服务端: go run *.go -mode native -addr 0.0.0.0:4363")
	fmt.Println("    客户端: go run *.go -mode native -server localhost:4363 -size 1024 -rate 100")
	fmt.Println()
	fmt.Println("  LibP2P模式:")
	fmt.Println("    服务端: go run *.go -mode libp2p -listen /ip4/0.0.0.0/udp/4363/quic-v1")
	fmt.Println("    客户端: go run *.go -mode libp2p -peer /ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooW...")
	fmt.Println()
	fmt.Println("配置文件位置: ~/.quic-datagram-test/")
	fmt.Println()
	fmt.Println("更多信息请查看 README.md")
}
