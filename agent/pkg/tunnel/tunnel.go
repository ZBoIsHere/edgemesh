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
	"github.com/libp2p/go-libp2p-core/host"
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

	relay, err := controller.APIConn.GetPeerAddrInfo(constants.SERVER_ADDR_NAME)
	if err != nil {
		klog.Errorln(err)
		return
	}

	h, err := libp2p.New(context.Background(),
		libp2p.EnableRelay(circuit.OptActive),
		libp2p.EnableAutoRelay(),
		libp2p.ForceReachabilityPrivate(),
		libp2p.StaticRelays([]peer.AddrInfo{*relay}),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", t.Config.ListenPort)),
		libp2p.EnableHolePunching(),
		libp2p.Identity(privateKey),
	)
	if err != nil {
		errMsg := fmt.Errorf("Start tunnel server failed, %v", err)
		klog.Errorln(errMsg)
		return
	}

	t.Host = h
	t.TCPProxySvc = tcp.NewTCPProxyService(h)
	klog.Infoln("Start tunnel agent success")

	isStop := false
	for !isStop {
		klog.Warningf("Tunnel agent connecting to tunnel server %s", t.Config.TunnelServer)
		time.Sleep(2 * time.Second)
		for _, v := range h.Addrs() {
			if _, err := v.ValueForProtocol(ma.P_CIRCUIT); err == nil {
				klog.Infof("Tunnel agent connected to tunnel server %s", t.Config.TunnelServer)
				isStop = true
				break
			}
		}
	}

	nodeName := t.Config.NodeName
	controller.APIConn.SetPeerAddrInfo(nodeName, host.InfoFromHost(h))

	// set tcp proxy handler
	h.SetStreamHandler(tcp.TCPProxyProtocol, t.TCPProxySvc.ProxyStreamHandler)

	select {}
	// TODO ifRotationDone() ????, 后面要添加这个东西，如果证书轮换了，要重新进行连接
}