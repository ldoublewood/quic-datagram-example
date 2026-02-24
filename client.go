package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	tpt "github.com/libp2p/go-libp2p/core/transport"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/quic-go/quic-go"
)

type ClientConfig struct {
	Mode        string
	ServerAddr  string
	PeerAddr    string // libp2p模式下的multiaddr
	PacketSize  int
	SendRate    int
	Duration    time.Duration
	PayloadType string
}

type Client struct {
	conn   Connection
	config ClientConfig
	stats  ClientStats
}

func (c *Client) connectNative() error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-datagram-test"},
	}

	conn, err := quic.DialAddr(context.Background(), c.config.ServerAddr, tlsConfig, &quic.Config{
		EnableDatagrams: true,
	})
	if err != nil {
		return err
	}

	c.conn = &NativeConnection{conn: conn}
	c.stats.StartTime = time.Now()
	return nil
}

func (c *Client) connectLibP2P(config *Config) error {
	transport, err := makeDatagramTransport(config)
	if err != nil {
		return fmt.Errorf("创建transport失败: %w", err)
	}

	h, err := libp2p.New(
		libp2p.Identity(config.PrivateKey),
		libp2p.Transport(func() (tpt.Transport, error) { return transport, nil }),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/udp/0/quic-v1"),
		libp2p.DisableRelay(),
	)
	if err != nil {
		return fmt.Errorf("创建libp2p host失败: %w", err)
	}

	fmt.Printf("本地 Peer ID: %s\n", h.ID())

	// 解析目标地址
	targetAddr, err := ma.NewMultiaddr(c.config.PeerAddr)
	if err != nil {
		return fmt.Errorf("解析目标地址失败: %w", err)
	}

	// 提取peer ID
	addrInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		return fmt.Errorf("提取peer信息失败: %w", err)
	}

	// 添加到peerstore
	h.Peerstore().AddAddrs(addrInfo.ID, addrInfo.Addrs, peerstore.PermanentAddrTTL)

	// 连接到对等节点
	fmt.Printf("正在连接到: %s\n", addrInfo.ID)
	if err := h.Connect(context.Background(), *addrInfo); err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}

	// 等待连接建立
	time.Sleep(500 * time.Millisecond)

	// 获取连接
	conns := h.Network().ConnsToPeer(addrInfo.ID)
	if len(conns) == 0 {
		return fmt.Errorf("未找到到目标节点的连接")
	}

	libp2pConn, err := NewLibP2PConnection(conns[0])
	if err != nil {
		return fmt.Errorf("创建LibP2P连接失败: %w", err)
	}

	c.conn = libp2pConn
	c.stats.StartTime = time.Now()
	return nil
}

func (c *Client) sendPackets() {
	ctx := context.Background()
	interval := time.Second / time.Duration(c.config.SendRate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var seqNum uint64 = 1
	endTime := time.Now().Add(c.config.Duration)

	fmt.Printf("开始发送数据包，发送速率: %d pps，包大小: %d 字节，持续时间: %v\n",
		c.config.SendRate, c.config.PacketSize, c.config.Duration)

	for time.Now().Before(endTime) {
		select {
		case <-ticker.C:
			payload := GeneratePayload(seqNum, c.config.PacketSize, c.config.PayloadType)

			err := c.conn.SendDatagram(payload)
			if err != nil {
				c.stats.IncrementError()
				fmt.Printf("发送包 #%d 失败: %v\n", seqNum, err)
			} else {
				c.stats.IncrementSent()

				if seqNum%100 == 0 {
					fmt.Printf("已发送 %d 个包\n", seqNum)
				}
			}

			seqNum++
		case <-ctx.Done():
			return
		}
	}
}

func runClient() {
	var config ClientConfig

	flag.StringVar(&config.Mode, "mode", "native", "连接模式: native 或 libp2p")
	flag.StringVar(&config.ServerAddr, "server", "localhost:4363", "服务器地址 (native模式)")
	flag.StringVar(&config.PeerAddr, "peer", "", "对等节点multiaddr (libp2p模式)")
	flag.IntVar(&config.PacketSize, "size", 1024, "数据包大小（字节）")
	flag.IntVar(&config.SendRate, "rate", 100, "发送速率（包/秒）")
	flag.DurationVar(&config.Duration, "duration", 30*time.Second, "发送持续时间")
	flag.StringVar(&config.PayloadType, "payload", "random", "负载类型 (random/sequential)")
	flag.Parse()

	client := &Client{config: config}

	var err error
	if config.Mode == "libp2p" {
		if config.PeerAddr == "" {
			log.Fatal("libp2p模式需要指定 -peer 参数")
		}
		
		cfg, err := LoadOrCreateConfig()
		if err != nil {
			log.Fatalf("加载配置失败: %v", err)
		}

		fmt.Printf("使用LibP2P模式连接到: %s\n", config.PeerAddr)
		if err := client.connectLibP2P(cfg); err != nil {
			log.Fatal("连接失败:", err)
		}
	} else {
		fmt.Printf("使用Native模式连接到服务器: %s\n", config.ServerAddr)
		if err = client.connectNative(); err != nil {
			log.Fatal("连接失败:", err)
		}
	}
	defer client.conn.Close()

	fmt.Printf("连接成功，开始性能测试...\n")

	// 发送数据包
	client.sendPackets()

	// 等待一小段时间确保最后的包被发送
	time.Sleep(100 * time.Millisecond)

	// 打印最终统计
	client.stats.PrintFinal(config.PacketSize)

	fmt.Printf("测试完成，保持连接5秒以查看服务器统计...\n")
	time.Sleep(5 * time.Second)
}
