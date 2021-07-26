package tunnel

import (
	"context"
	"fmt"
	"github.com/kubeedge/edgemesh/agent/pkg/tunnel/controller"
	"github.com/kubeedge/edgemesh/agent/pkg/tunnel/protocol/tcp"
	"github.com/kubeedge/edgemesh/common/certificate"
	"github.com/kubeedge/edgemesh/common/constants"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"k8s.io/klog/v2"
	"time"
)

func (t *TunnelAgent) Run() {
	privateKey, err := certificate.GetPrivateKey(t.Config.TunnelCertificate, t.Config.NodeName)
	if err != nil {
		klog.Errorln(err)
		return
	}

	relayAddr, err := controller.APIConn.Get(constants.SERVER_ADDR_NAME)
	if err != nil {
		klog.Errorln(err)
		return
	}

	raddrInfo, err := peer.AddrInfoFromP2pAddr(relayAddr)
	if err != nil {
		klog.Errorln(err)
		return
	}

	host, err := libp2p.New(context.Background(),
		libp2p.EnableRelay(circuit.OptActive),
		libp2p.EnableAutoRelay(),
		libp2p.ForceReachabilityPrivate(),
		libp2p.StaticRelays([]peer.AddrInfo{*raddrInfo}),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", t.Config.ListenPort)),
		libp2p.EnableHolePunching(),
		libp2p.Identity(privateKey),
	)
	if err != nil {
		errMsg := fmt.Errorf("Start tunnel server failed, %v", err)
		klog.Errorln(errMsg)
		return
	}

	t.Host = host
	t.TCPProxySvc = tcp.NewTCPProxyService(host)
	klog.Infoln("Start tunnel agent success")

	isStop := false
	for !isStop {
		klog.Warningf("Tunnel agent connecting to tunnel server %s", t.Config.TunnelServer)
		time.Sleep(2 * time.Second)
		for _, v := range host.Addrs() {
			if _, err := v.ValueForProtocol(ma.P_CIRCUIT); err == nil {
				klog.Infof("Tunnel agent connected to tunnel server %s", t.Config.TunnelServer)
				isStop = true
				break
			}
		}
	}

	nodeName := t.Config.NodeName
	controller.APIConn.SetSelfAddr2Secret(nodeName, host.ID(), host.Addrs())

	// set tcp proxy handler
	host.SetStreamHandler(tcp.TCPProxyProtocol, t.TCPProxySvc.ProxyStreamHandler)

	select {}
	// TODO ifRotationDone() ????, 后面要添加这个东西，如果证书轮换了，要重新进行连接
}