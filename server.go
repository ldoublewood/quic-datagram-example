package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"net"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

type PacketStats struct {
	ReceivedCount int64
	LostCount     int64
	TotalLatency  time.Duration
	MinLatency    time.Duration
	MaxLatency    time.Duration
	LastSeqNum    uint64
}

type Server struct {
	stats PacketStats
	mutex sync.RWMutex
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(time.Hour),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{certDER},
			PrivateKey:  key,
		}},
		NextProtos: []string{"quic-datagram-test"},
	}
}

func (s *Server) handleConnection(conn quic.Connection) {
	fmt.Printf("客户端连接: %s\n", conn.RemoteAddr())

	ctx := context.Background()

	for {
		data, err := conn.ReceiveDatagram(ctx)
		if err != nil {
			fmt.Printf("接收数据报错误: %v\n", err)
			return
		}

		s.processPacket(data)
	}
}

func (s *Server) processPacket(data []byte) {
	if len(data) < 16 { // 8字节序列号 + 8字节时间戳
		return
	}

	seqNum := binary.BigEndian.Uint64(data[:8])
	timestamp := int64(binary.BigEndian.Uint64(data[8:16]))
	sendTime := time.Unix(0, timestamp)

	now := time.Now()
	latency := now.Sub(sendTime)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.stats.ReceivedCount++
	s.stats.TotalLatency += latency

	if s.stats.MinLatency == 0 || latency < s.stats.MinLatency {
		s.stats.MinLatency = latency
	}
	if latency > s.stats.MaxLatency {
		s.stats.MaxLatency = latency
	}

	// 检测丢包
	if seqNum > s.stats.LastSeqNum+1 {
		lost := seqNum - s.stats.LastSeqNum - 1
		s.stats.LostCount += int64(lost)
		fmt.Printf("检测到丢包: 序列号 %d-%d (丢失 %d 个包)\n",
			s.stats.LastSeqNum+1, seqNum-1, lost)
	}
	s.stats.LastSeqNum = seqNum

	fmt.Printf("收到包 #%d, 延迟: %v, 大小: %d 字节\n",
		seqNum, latency, len(data))
}

func (s *Server) printStats() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.RLock()
		stats := s.stats
		s.mutex.RUnlock()

		if stats.ReceivedCount > 0 {
			avgLatency := stats.TotalLatency / time.Duration(stats.ReceivedCount)
			lossRate := float64(stats.LostCount) / float64(stats.ReceivedCount+stats.LostCount) * 100

			fmt.Printf("\n=== 统计信息 ===\n")
			fmt.Printf("接收包数: %d\n", stats.ReceivedCount)
			fmt.Printf("丢失包数: %d\n", stats.LostCount)
			fmt.Printf("丢包率: %.2f%%\n", lossRate)
			fmt.Printf("平均延迟: %v\n", avgLatency)
			fmt.Printf("最小延迟: %v\n", stats.MinLatency)
			fmt.Printf("最大延迟: %v\n", stats.MaxLatency)
			fmt.Printf("================\n\n")
		}
	}
}

func main() {
	addr := "0.0.0.0:4363"

	listener, err := quic.ListenAddr(addr, generateTLSConfig(), &quic.Config{
		EnableDatagrams: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	server := &Server{}

	fmt.Printf("QUIC Datagram 服务器启动，监听地址: %s\n", addr)

	// 启动统计信息打印
	go server.printStats()

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("接受连接错误: %v", err)
			continue
		}

		go server.handleConnection(conn)
	}
}
