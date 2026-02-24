package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

type ClientConfig struct {
	ServerAddr  string
	PacketSize  int
	SendRate    int // 每秒发送包数
	Duration    time.Duration
	PayloadType string // "random" 或 "sequential"
}

type ClientStats struct {
	SentCount  int64
	ErrorCount int64
	StartTime  time.Time
	mutex      sync.RWMutex
}

type Client struct {
	conn   quic.Connection
	config ClientConfig
	stats  ClientStats
}

func (c *Client) connect() error {
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

	c.conn = conn
	c.stats.StartTime = time.Now()
	return nil
}

func (c *Client) generatePayload(seqNum uint64) []byte {
	payload := make([]byte, c.config.PacketSize)

	// 前8字节：序列号
	binary.BigEndian.PutUint64(payload[:8], seqNum)

	// 接下来8字节：时间戳（纳秒）
	timestamp := time.Now().UnixNano()
	binary.BigEndian.PutUint64(payload[8:16], uint64(timestamp))

	// 剩余部分：根据配置生成数据
	switch c.config.PayloadType {
	case "random":
		rand.Read(payload[16:])
	case "sequential":
		for i := 16; i < len(payload); i++ {
			payload[i] = byte(i % 256)
		}
	default:
		// 默认填充固定模式
		for i := 16; i < len(payload); i++ {
			payload[i] = byte(seqNum % 256)
		}
	}

	return payload
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
			payload := c.generatePayload(seqNum)

			err := c.conn.SendDatagram(payload)
			if err != nil {
				c.stats.mutex.Lock()
				c.stats.ErrorCount++
				c.stats.mutex.Unlock()
				fmt.Printf("发送包 #%d 失败: %v\n", seqNum, err)
			} else {
				c.stats.mutex.Lock()
				c.stats.SentCount++
				c.stats.mutex.Unlock()

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

func (c *Client) printFinalStats() {
	c.stats.mutex.RLock()
	defer c.stats.mutex.RUnlock()

	duration := time.Since(c.stats.StartTime)
	actualRate := float64(c.stats.SentCount) / duration.Seconds()

	fmt.Printf("\n=== 客户端发送统计 ===\n")
	fmt.Printf("发送包数: %d\n", c.stats.SentCount)
	fmt.Printf("发送错误: %d\n", c.stats.ErrorCount)
	fmt.Printf("实际发送速率: %.2f pps\n", actualRate)
	fmt.Printf("总发送时间: %v\n", duration)
	fmt.Printf("总数据量: %.2f MB\n", float64(c.stats.SentCount*int64(c.config.PacketSize))/1024/1024)
	fmt.Printf("=====================\n")
}

func main() {
	var config ClientConfig

	flag.StringVar(&config.ServerAddr, "server", "localhost:8080", "服务器地址")
	flag.IntVar(&config.PacketSize, "size", 1024, "数据包大小（字节）")
	flag.IntVar(&config.SendRate, "rate", 100, "发送速率（包/秒）")
	flag.DurationVar(&config.Duration, "duration", 30*time.Second, "发送持续时间")
	flag.StringVar(&config.PayloadType, "payload", "random", "负载类型 (random/sequential)")
	flag.Parse()

	client := &Client{config: config}

	fmt.Printf("连接到服务器: %s\n", config.ServerAddr)
	if err := client.connect(); err != nil {
		log.Fatal("连接失败:", err)
	}
	defer client.conn.CloseWithError(0, "")

	fmt.Printf("连接成功，开始性能测试...\n")

	// 发送数据包
	client.sendPackets()

	// 等待一小段时间确保最后的包被发送
	time.Sleep(100 * time.Millisecond)

	// 打印最终统计
	client.printFinalStats()

	fmt.Printf("测试完成，保持连接5秒以查看服务器统计...\n")
	time.Sleep(5 * time.Second)
}
