package main

import (
	"context"

	"github.com/libp2p/go-libp2p/core/connmgr"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/pnet"
	tpt "github.com/libp2p/go-libp2p/core/transport"
	libp2pquic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/quicreuse"
	ma "github.com/multiformats/go-multiaddr"
)

// DatagramTransport 包装libp2pquic.Transport并返回DatagramConns
type DatagramTransport struct {
	tpt.Transport
	connManager *quicreuse.ConnManager
}

// NewDatagramTransport 创建支持datagram的QUIC transport
func NewDatagramTransport(key crypto.PrivKey, connManager *quicreuse.ConnManager, psk pnet.PSK, gater connmgr.ConnectionGater, rcmgr network.ResourceManager) (tpt.Transport, error) {
	baseTransport, err := libp2pquic.NewTransport(key, connManager, psk, gater, rcmgr)
	if err != nil {
		return nil, err
	}
	return &DatagramTransport{
		Transport:   baseTransport,
		connManager: connManager,
	}, nil
}

func (t *DatagramTransport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (tpt.CapableConn, error) {
	c, err := t.Transport.Dial(ctx, raddr, p)
	if err != nil {
		return nil, err
	}
	qc, err := getQuicConn(c)
	if err != nil {
		c.Close()
		return nil, err
	}
	return &datagramConn{CapableConn: c, quicConn: qc}, nil
}

func (t *DatagramTransport) Listen(addr ma.Multiaddr) (tpt.Listener, error) {
	ln, err := t.Transport.Listen(addr)
	if err != nil {
		return nil, err
	}
	return &datagramListener{Listener: ln}, nil
}

// datagramListener 包装tpt.Listener并升级接受的连接
type datagramListener struct {
	tpt.Listener
}

func (l *datagramListener) Accept() (tpt.CapableConn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	qc, err := getQuicConn(c)
	if err != nil {
		c.Close()
		return nil, err
	}
	return &datagramConn{CapableConn: c, quicConn: qc}, nil
}
