package tunnel

import (
	"fmt"
	"github.com/kubeedge/beehive/pkg/core"
	"github.com/kubeedge/edgemesh/common/certificate"
	"github.com/kubeedge/edgemesh/common/constants"
	"github.com/kubeedge/edgemesh/server/pkg/tunnel/config"
)

type TunnelServer struct {
	Config      *config.TunnelServerConfig
	certManager certificate.CertManager
	enable      bool
}

func newTunnelServer(c *config.TunnelServerConfig) (server *TunnelServer, err error) {
	server = &TunnelServer{
		Config: c,
		enable: c.Enable,
	}
	return server, nil
}

func Register(c *config.TunnelServerConfig) error {
	server, err := newTunnelServer(c)
	if err != nil {
		return fmt.Errorf("register module tunnelserver error: %v", err)
	}
	core.Register(server)
	return nil
}

func (t *TunnelServer) Name() string {
	return constants.AgentTunnelModuleName
}

func (t *TunnelServer) Group() string {
	return constants.AgentTunnelGroupName
}

func (t *TunnelServer) Enable() bool {
	return t.enable
}

func (t *TunnelServer) Start() {
	//t.Run()
}
