package main

import (
	"context"

	"github.com/quic-go/quic-go"
)

// Connection 是一个通用的连接接口，支持native和libp2p两种模式
type Connection interface {
	SendDatagram(data []byte) error
	ReceiveDatagram(ctx context.Context) ([]byte, error)
	Close() error
	RemoteAddr() string
}

// NativeConnection 包装原生QUIC连接
type NativeConnection struct {
	conn *quic.Conn
}

func (c *NativeConnection) SendDatagram(data []byte) error {
	return c.conn.SendDatagram(data)
}

func (c *NativeConnection) ReceiveDatagram(ctx context.Context) ([]byte, error) {
	return c.conn.ReceiveDatagram(ctx)
}

func (c *NativeConnection) Close() error {
	return c.conn.CloseWithError(0, "")
}

func (c *NativeConnection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}
