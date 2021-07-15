package tunnel

import (
	"github.com/kubeedge/edgemesh/common/certificate"
	"github.com/kubeedge/edgemesh/server/pkg/tunnel/tunnelserver"
)

func (t *TunnelServer) Run() {
	certificateConfig := certificate.TunnelCertificate{
		TLSCAFile:          t.Config.TLSCAFile,
		TLSCertFile:        t.Config.TLSCertFile,
		TLSPrivateKeyFile:  t.Config.TLSPrivateKeyFile,
		Token:              t.Config.Token,
		HTTPServer:         t.Config.HTTPServer,
		RotateCertificates: t.Config.RotateCertificates,
	}
	t.certManager = certificate.NewCertManager(certificateConfig, t.Config.NodeName)
	t.certManager.Start()
	// TunnelServer mainly used to help hole punch or relay between edgemesh-agent
	go tunnelserver.StartTunnelServer()

	// TODO ifRotationDone() ????, 后面要添加这个东西，如果证书轮换了，要重新进行连接
	select {}
}
