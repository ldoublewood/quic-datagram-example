package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	tpt "github.com/libp2p/go-libp2p/core/transport"
	"github.com/libp2p/go-libp2p/p2p/transport/quicreuse"
	"github.com/quic-go/quic-go"
)

type Server struct {
	stats ServerStats
	mode  string
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

func (s *Server) handleConnection(conn Connection) {
	fmt.Printf("客户端连接: %s\n", conn.RemoteAddr())

	ctx := context.Background()

	for {
		data, err := conn.ReceiveDatagram(ctx)
		if err != nil {
			fmt.Printf("接收数据报错误: %v\n", err)
			return
		}

		s.stats.ProcessPacket(data)
	}
}

func (s *Server) printStats() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.stats.Print()
	}
}

func runNativeServer(addr string) error {
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), &quic.Config{
		EnableDatagrams: true,
	})
	if err != nil {
		return err
	}
	defer listener.Close()

	server := &Server{mode: "native"}
	fmt.Printf("Native QUIC Datagram 服务器启动，监听地址: %s\n", addr)

	go server.printStats()

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("接受连接错误: %v", err)
			continue
		}

		nativeConn := &NativeConnection{conn: conn}
		go server.handleConnection(nativeConn)
	}
}

func makeDatagramTransport(config *Config) (tpt.Transport, error) {
	var resetKey quic.StatelessResetKey
	var tokenKey quic.TokenGeneratorKey
	connManager, err := quicreuse.NewConnManager(resetKey, tokenKey)
	if err != nil {
		return nil, err
	}
	return NewDatagramTransport(config.PrivateKey, connManager, nil, nil, nil)
}

func runLibP2PServer(listenAddr string, config *Config) error {
	transport, err := makeDatagramTransport(config)
	if err != nil {
		return fmt.Errorf("创建transport失败: %w", err)
	}

	h, err := libp2p.New(
		libp2p.Identity(config.PrivateKey),
		libp2p.Transport(func() (tpt.Transport, error) { return transport, nil }),
		libp2p.ListenAddrStrings(listenAddr),
		libp2p.DisableRelay(),
	)
	if err != nil {
		return fmt.Errorf("创建libp2p host失败: %w", err)
	}
	defer h.Close()

	server := &Server{mode: "libp2p"}
	
	fmt.Printf("LibP2P QUIC Datagram 服务器启动\n")
	fmt.Printf("Peer ID: %s\n", h.ID())
	fmt.Printf("监听地址:\n")
	for _, addr := range h.Addrs() {
		fmt.Printf("  %s/p2p/%s\n", addr, h.ID())
	}

	// 设置连接通知器
	notifee := &network.NotifyBundle{
		ConnectedF: func(n network.Network, conn network.Conn) {
			fmt.Printf("新连接来自: %s\n", conn.RemotePeer())
			libp2pConn, err := NewLibP2PConnection(conn)
			if err != nil {
				fmt.Printf("创建LibP2P连接失败: %v\n", err)
				return
			}
			go server.handleConnection(libp2pConn)
		},
	}
	h.Network().Notify(notifee)

	go server.printStats()

	// 保持运行
	select {}
}

func runServer() {
	mode := flag.String("mode", "native", "连接模式: native 或 libp2p")
	addr := flag.String("addr", "0.0.0.0:4363", "监听地址 (native模式)")
	listenAddr := flag.String("listen", "/ip4/0.0.0.0/udp/4363/quic-v1", "监听地址 (libp2p模式)")
	flag.Parse()

	var err error
	if *mode == "libp2p" {
		config, err := LoadOrCreateConfig()
		if err != nil {
			log.Fatalf("加载配置失败: %v", err)
		}
		err = runLibP2PServer(*listenAddr, config)
	} else {
		err = runNativeServer(*addr)
	}

	if err != nil {
		log.Fatal(err)
	}
}
