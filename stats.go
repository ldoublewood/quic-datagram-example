package main

import (
	"encoding/binary"
	"fmt"
	"sync"
	"time"
)

// ClientStats 客户端统计信息
type ClientStats struct {
	SentCount  int64
	ErrorCount int64
	StartTime  time.Time
	mutex      sync.RWMutex
}

func (s *ClientStats) IncrementSent() {
	s.mutex.Lock()
	s.SentCount++
	s.mutex.Unlock()
}

func (s *ClientStats) IncrementError() {
	s.mutex.Lock()
	s.ErrorCount++
	s.mutex.Unlock()
}

func (s *ClientStats) PrintFinal(packetSize int) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	duration := time.Since(s.StartTime)
	actualRate := float64(s.SentCount) / duration.Seconds()

	fmt.Printf("\n=== 客户端发送统计 ===\n")
	fmt.Printf("发送包数: %d\n", s.SentCount)
	fmt.Printf("发送错误: %d\n", s.ErrorCount)
	fmt.Printf("实际发送速率: %.2f pps\n", actualRate)
	fmt.Printf("总发送时间: %v\n", duration)
	fmt.Printf("总数据量: %.2f MB\n", float64(s.SentCount*int64(packetSize))/1024/1024)
	fmt.Printf("=====================\n")
}

// ServerStats 服务端统计信息
type ServerStats struct {
	ReceivedCount int64
	LostCount     int64
	TotalLatency  time.Duration
	MinLatency    time.Duration
	MaxLatency    time.Duration
	LastSeqNum    uint64
	mutex         sync.RWMutex
}

func (s *ServerStats) ProcessPacket(data []byte) {
	if len(data) < 16 {
		return
	}

	seqNum := binary.BigEndian.Uint64(data[:8])
	timestamp := int64(binary.BigEndian.Uint64(data[8:16]))
	sendTime := time.Unix(0, timestamp)

	now := time.Now()
	latency := now.Sub(sendTime)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.ReceivedCount++
	s.TotalLatency += latency

	if s.MinLatency == 0 || latency < s.MinLatency {
		s.MinLatency = latency
	}
	if latency > s.MaxLatency {
		s.MaxLatency = latency
	}

	// 检测丢包
	if seqNum > s.LastSeqNum+1 {
		lost := seqNum - s.LastSeqNum - 1
		s.LostCount += int64(lost)
		fmt.Printf("检测到丢包: 序列号 %d-%d (丢失 %d 个包)\n",
			s.LastSeqNum+1, seqNum-1, lost)
	}
	s.LastSeqNum = seqNum

	fmt.Printf("收到包 #%d, 延迟: %v, 大小: %d 字节\n",
		seqNum, latency, len(data))
}

func (s *ServerStats) Print() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.ReceivedCount > 0 {
		avgLatency := s.TotalLatency / time.Duration(s.ReceivedCount)
		lossRate := float64(s.LostCount) / float64(s.ReceivedCount+s.LostCount) * 100

		fmt.Printf("\n=== 统计信息 ===\n")
		fmt.Printf("接收包数: %d\n", s.ReceivedCount)
		fmt.Printf("丢失包数: %d\n", s.LostCount)
		fmt.Printf("丢包率: %.2f%%\n", lossRate)
		fmt.Printf("平均延迟: %v\n", avgLatency)
		fmt.Printf("最小延迟: %v\n", s.MinLatency)
		fmt.Printf("最大延迟: %v\n", s.MaxLatency)
		fmt.Printf("================\n\n")
	}
}

// GeneratePayload 生成测试数据包
func GeneratePayload(seqNum uint64, packetSize int, payloadType string) []byte {
	payload := make([]byte, packetSize)

	// 前8字节：序列号
	binary.BigEndian.PutUint64(payload[:8], seqNum)

	// 接下来8字节：时间戳（纳秒）
	timestamp := time.Now().UnixNano()
	binary.BigEndian.PutUint64(payload[8:16], uint64(timestamp))

	// 剩余部分：根据配置生成数据
	switch payloadType {
	case "random":
		// 使用简单的伪随机填充
		for i := 16; i < len(payload); i++ {
			payload[i] = byte((seqNum*uint64(i) + uint64(i)) % 256)
		}
	case "sequential":
		for i := 16; i < len(payload); i++ {
			payload[i] = byte(i % 256)
		}
	default:
		for i := 16; i < len(payload); i++ {
			payload[i] = byte(seqNum % 256)
		}
	}

	return payload
}
