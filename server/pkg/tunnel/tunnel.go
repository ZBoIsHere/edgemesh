package tunnel

import (
	"fmt"
	"github.com/kubeedge/edgemesh/common/constants"
	"github.com/kubeedge/edgemesh/server/pkg/tunnel/controller"
	ma "github.com/multiformats/go-multiaddr"
	"k8s.io/klog/v2"
)

func (t *TunnelServer) Run() {
	klog.Infoln("Start tunnel server success")
	for _, v := range t.Host.Addrs() {
		klog.Infof("%s : %v/p2p/%s\n", "Tunnel server addr", v, t.Host.ID().Pretty())
	}

	var addrs []ma.Multiaddr
	publicIPAddrStr := fmt.Sprintf("/ip4/%s/tcp/%d", t.Config.PublicIP, t.Config.ListenPort)
	if t.Config.PublicIP != "" {
		klog.Infof("%s : %s/p2p/%s\n", "Tunnel server addr", publicIPAddrStr, t.Host.ID().Pretty())
	}
	publicIPAddr, _ := ma.NewMultiaddr(publicIPAddrStr)
	addrs = append(addrs, publicIPAddr)

	err := controller.APIConn.SetSelfAddr2Secret(constants.SERVER_ADDR_NAME, t.Host.ID(), addrs)
	if err != nil {
		klog.Errorf("failed update [%s] addr %v to secret: %v", constants.SERVER_ADDR_NAME, addrs, err)
	}
	klog.Infof("success update [%s] addr %v to secret", constants.SERVER_ADDR_NAME, addrs)

	// TODO ifRotationDone() ????, 后面要添加这个东西，如果证书轮换了，要重新进行连接
	select {}
}
