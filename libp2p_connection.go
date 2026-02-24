package main

import (
	"context"
	"errors"

	"github.com/libp2p/go-libp2p/core/network"
	tpt "github.com/libp2p/go-libp2p/core/transport"
	"github.com/quic-go/quic-go"
)

// DatagramConn 是支持datagram的连接接口
type DatagramConn interface {
	tpt.CapableConn
	SendDatagram([]byte) error
	ReceiveDatagram(context.Context) ([]byte, error)
}

// datagramConn 包装CapableConn并添加datagram方法
type datagramConn struct {
	tpt.CapableConn
	quicConn *quic.Conn
}

func (c *datagramConn) SendDatagram(data []byte) error {
	return c.quicConn.SendDatagram(data)
}

func (c *datagramConn) ReceiveDatagram(ctx context.Context) ([]byte, error) {
	return c.quicConn.ReceiveDatagram(ctx)
}

func (c *datagramConn) As(target any) bool {
	if t, ok := target.(**quic.Conn); ok {
		*t = c.quicConn
		return true
	}
	if t, ok := target.(*DatagramConn); ok {
		*t = c
		return true
	}
	return c.CapableConn.As(target)
}

// getQuicConn 从CapableConn中提取底层的quic.Conn
func getQuicConn(c tpt.CapableConn) (*quic.Conn, error) {
	var qc *quic.Conn
	if ok := c.As(&qc); ok && qc != nil {
		return qc, nil
	}
	return nil, errors.New("underlying connection does not support quic.Conn")
}

// LibP2PConnection 包装libp2p连接
type LibP2PConnection struct {
	conn     network.Conn
	dgConn   DatagramConn
	peerAddr string
}

func NewLibP2PConnection(conn network.Conn) (*LibP2PConnection, error) {
	var dgConn DatagramConn
	if ok := conn.As(&dgConn); !ok {
		return nil, errors.New("connection does not support DatagramConn")
	}
	
	return &LibP2PConnection{
		conn:     conn,
		dgConn:   dgConn,
		peerAddr: conn.RemotePeer().String(),
	}, nil
}

func (c *LibP2PConnection) SendDatagram(data []byte) error {
	return c.dgConn.SendDatagram(data)
}

func (c *LibP2PConnection) ReceiveDatagram(ctx context.Context) ([]byte, error) {
	return c.dgConn.ReceiveDatagram(ctx)
}

func (c *LibP2PConnection) Close() error {
	return c.conn.Close()
}

func (c *LibP2PConnection) RemoteAddr() string {
	return c.peerAddr
}
