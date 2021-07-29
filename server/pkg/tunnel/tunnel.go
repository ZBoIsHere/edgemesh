package tunnel

import (
	"github.com/kubeedge/edgemesh/common/constants"
	"github.com/kubeedge/edgemesh/server/pkg/tunnel/controller"
	"github.com/libp2p/go-libp2p-core/host"
	"k8s.io/klog/v2"
)

func (t *TunnelServer) Run() {
	klog.Infoln("Start tunnel server success")
	for _, v := range t.Host.Addrs() {
		klog.Infof("%s : %v/p2p/%s\n", "Tunnel server addr", v, t.Host.ID().Pretty())
	}

	err := controller.APIConn.SetPeerAddrInfo(constants.SERVER_ADDR_NAME, host.InfoFromHost(t.Host))
	if err != nil {
		klog.Errorf("failed update [%s] addr %v to secret: %v", constants.SERVER_ADDR_NAME, t.Host.Addrs(), err)
	}
	klog.Infof("success update [%s] addr %v to secret", constants.SERVER_ADDR_NAME, t.Host.Addrs())

	// TODO ifRotationDone() ????, 后面要添加这个东西，如果证书轮换了，要重新进行连接
	select {}
}
