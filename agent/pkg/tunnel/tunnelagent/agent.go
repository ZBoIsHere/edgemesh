package tunnelagent

import (
	"context"
	"encoding/pem"
	"github.com/kubeedge/edgemesh/agent/pkg/tunnel/agentaddr"
	"github.com/kubeedge/edgemesh/agent/pkg/tunnel/protocol/tcp"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"io/ioutil"
	"k8s.io/klog/v2"
	"log"
	"os"
	"sync"
	"time"
)

var (
	once sync.Once
	ta   *TunnelAgent
)

const (
	MY_NODE_NAME = "MY_NODE_NAME"
)

type TunnelAgent struct {
	NodeName    string
	Host        host.Host
	TCPProxySvc *tcp.TCPProxyService
	AgentAddr   *agentaddr.PeerAgentAddr
}

func NewTunnelAgent() *TunnelAgent {
	once.Do(func() {
		ctx := context.Background()

		// todo
		certBytes, err := ioutil.ReadFile("test")
		//certBytes, err := ioutil.ReadFile(config.Config.Tunnel.TLSPrivateKeyFile)
		if err != nil {
			klog.Errorln(err)
			return
		}

		block, _ := pem.Decode(certBytes)

		priv, err := crypto.UnmarshalECDSAPrivateKey(block.Bytes)
		if err != nil {
			klog.Errorln(err)
			return
		}

		// todo read tunnel server addr from where
		relayAddr := "test"
		//relayAddr := config.Config.TunnelServer
		addr, err := ma.NewMultiaddr(relayAddr)
		if err != nil {
			klog.Errorln(err)
			panic(err)
		}

		raddrInfo, err := peer.AddrInfoFromP2pAddr(addr)

		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		h, err := libp2p.New(ctx,
			libp2p.EnableRelay(circuit.OptActive),
			libp2p.EnableAutoRelay(),
			libp2p.ForceReachabilityPrivate(),
			libp2p.StaticRelays([]peer.AddrInfo{*raddrInfo}),
			libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/10006"),
			libp2p.EnableHolePunching(),
			libp2p.Identity(priv),
		)

		nodeName, isExist := os.LookupEnv(MY_NODE_NAME)
		if !isExist {
			klog.Errorln("Can't find nodeName")
			os.Exit(1)
		}

		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		ta = &TunnelAgent{
			NodeName:    nodeName,
			Host:        h,
			TCPProxySvc: tcp.NewTCPProxyService(h),
			AgentAddr:   agentaddr.NewPeerAgentAddr(),
		}
	})
	return ta
}

func (ta *TunnelAgent) GetHost() host.Host {
	return ta.Host
}

func (ta *TunnelAgent) GetTCPProxyService() *tcp.TCPProxyService {
	return ta.TCPProxySvc
}

func (ta *TunnelAgent) GetAgentAddr() *agentaddr.PeerAgentAddr {
	return ta.AgentAddr
}

func (ta *TunnelAgent) GetSelfNodeName() string {
	return ta.NodeName
}

func StartTunnelAgent() {
	tunnelAgent := NewTunnelAgent()
	host := tunnelAgent.GetHost()

	isStop := false
	for !isStop {
		// todo
		klog.Warningf("Tunnel agent connecting to tunnel server %s", "test")
		//klog.Warningf("Tunnel agent connecting to tunnel server %s", config.Config.TunnelServer)
		time.Sleep(2 * time.Second)
		for _, v := range host.Addrs() {
			if _, err := v.ValueForProtocol(ma.P_CIRCUIT); err == nil {
				// todo
				klog.Infof("Tunnel agent connected to tunnel server %s", "test")
				//klog.Infof("Tunnel agent connected to tunnel server %s", config.Config.TunnelServer)
				isStop = true
				break
			}
		}
	}

	nodeName := tunnelAgent.GetSelfNodeName()
	tunnelAgent.AgentAddr.SetSelfAddr2Secret(nodeName, host.ID(), host.Addrs())

	// set tcp proxy handler
	host.SetStreamHandler(tcp.TCPProxyProtocol, NewTunnelAgent().TCPProxySvc.ProxyStreamHandler)
	// todo set http proxy handler

	select {}
}
