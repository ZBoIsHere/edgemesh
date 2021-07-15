package tunnel

import (
	"fmt"
	"github.com/kubeedge/beehive/pkg/core"
	"github.com/kubeedge/edgemesh/agent/pkg/common/informers"
	"github.com/kubeedge/edgemesh/agent/pkg/tunnel/config"
	"github.com/kubeedge/edgemesh/common/certificate"
	"github.com/kubeedge/edgemesh/common/constants"
)

// TunnelAgent is used for solving cross subset communication
type TunnelAgent struct {
	Config      *config.TunnelAgentConfig
	certManager certificate.CertManager
	enable      bool
}

func newTunnelAgent(c *config.TunnelAgentConfig, ifm *informers.Manager) (agent *TunnelAgent, err error) {
	agent = &TunnelAgent{
		Config: c,
		enable: c.Enable,
	}
	return agent, nil
}

func Register(c *config.TunnelAgentConfig, ifm *informers.Manager) error {
	agent, err := newTunnelAgent(c, ifm)
	if err != nil {
		return fmt.Errorf("register module tunnelagent error: %v", err)
	}
	core.Register(agent)
	return nil
}

func (t *TunnelAgent) Name() string {
	return constants.AgentTunnelModuleName
}

func (t *TunnelAgent) Group() string {
	return constants.AgentTunnelGroupName
}

func (t *TunnelAgent) Enable() bool {
	return t.enable
}

func (t *TunnelAgent) Start() {
	//t.Run()
}
