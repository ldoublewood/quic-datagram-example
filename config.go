package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	configDirName      = ".quic-datagram-test"
	privateKeyFileName = "private_key"
	peerIDFileName     = "peer_id"
	bootstrapFileName  = "bootstrap_nodes"
)

// Config 存储节点配置
type Config struct {
	PrivateKey     crypto.PrivKey
	PeerID         peer.ID
	BootstrapNodes []string
	ConfigDir      string
}

// LoadOrCreateConfig 加载或创建配置
func LoadOrCreateConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("获取用户目录失败: %w", err)
	}

	configDir := filepath.Join(homeDir, configDirName)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %w", err)
	}

	config := &Config{ConfigDir: configDir}

	// 加载或生成私钥
	if err := config.loadOrCreatePrivateKey(); err != nil {
		return nil, err
	}

	// 加载bootstrap节点
	if err := config.loadBootstrapNodes(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) loadOrCreatePrivateKey() error {
	keyPath := filepath.Join(c.ConfigDir, privateKeyFileName)
	peerIDPath := filepath.Join(c.ConfigDir, peerIDFileName)

	// 尝试加载现有私钥
	if data, err := os.ReadFile(keyPath); err == nil {
		keyBytes, err := base64.StdEncoding.DecodeString(string(data))
		if err != nil {
			return fmt.Errorf("解码私钥失败: %w", err)
		}

		privKey, err := crypto.UnmarshalPrivateKey(keyBytes)
		if err != nil {
			return fmt.Errorf("解析私钥失败: %w", err)
		}

		peerID, err := peer.IDFromPrivateKey(privKey)
		if err != nil {
			return fmt.Errorf("从私钥生成Peer ID失败: %w", err)
		}

		c.PrivateKey = privKey
		c.PeerID = peerID
		fmt.Printf("已加载现有密钥，Peer ID: %s\n", peerID)
		return nil
	}

	// 生成新私钥
	fmt.Println("未找到现有密钥，正在生成新密钥...")
	privKey, err := GenerateRandomKey()
	if err != nil {
		return fmt.Errorf("生成密钥对失败: %w", err)
	}

	peerID, err := peer.IDFromPrivateKey(privKey)
	if err != nil {
		return fmt.Errorf("从私钥生成Peer ID失败: %w", err)
	}

	// 保存私钥
	keyBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return fmt.Errorf("序列化私钥失败: %w", err)
	}

	keyStr := base64.StdEncoding.EncodeToString(keyBytes)
	if err := os.WriteFile(keyPath, []byte(keyStr), 0600); err != nil {
		return fmt.Errorf("保存私钥失败: %w", err)
	}

	// 保存Peer ID
	if err := os.WriteFile(peerIDPath, []byte(peerID.String()), 0644); err != nil {
		return fmt.Errorf("保存Peer ID失败: %w", err)
	}

	c.PrivateKey = privKey
	c.PeerID = peerID
	fmt.Printf("已生成新密钥，Peer ID: %s\n", peerID)
	fmt.Printf("配置已保存到: %s\n", c.ConfigDir)
	return nil
}

func (c *Config) loadBootstrapNodes() error {
	bootstrapPath := filepath.Join(c.ConfigDir, bootstrapFileName)

	data, err := os.ReadFile(bootstrapPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，创建示例文件
			example := `# Bootstrap节点配置文件
# 每行一个multiaddr地址，例如：
# /ip4/127.0.0.1/udp/4363/quic-v1/p2p/12D3KooWExamplePeerID
# /ip4/192.168.1.100/udp/4363/quic-v1/p2p/12D3KooWAnotherPeerID
`
			if err := os.WriteFile(bootstrapPath, []byte(example), 0644); err != nil {
				return fmt.Errorf("创建bootstrap配置文件失败: %w", err)
			}
			fmt.Printf("已创建bootstrap配置文件: %s\n", bootstrapPath)
			fmt.Println("请编辑该文件添加bootstrap节点地址")
			return nil
		}
		return fmt.Errorf("读取bootstrap配置失败: %w", err)
	}

	// 解析bootstrap节点
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		c.BootstrapNodes = append(c.BootstrapNodes, line)
	}

	if len(c.BootstrapNodes) > 0 {
		fmt.Printf("已加载 %d 个bootstrap节点\n", len(c.BootstrapNodes))
	}

	return scanner.Err()
}

// GenerateRandomKey 生成随机密钥（用于测试）
func GenerateRandomKey() (crypto.PrivKey, error) {
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, 0)
	return priv, err
}
