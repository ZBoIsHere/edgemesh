package tunnel

import (
	"fmt"
	"github.com/kubeedge/beehive/pkg/core"
	"github.com/kubeedge/edgemesh/agent/pkg/tunnel/config"
	"github.com/kubeedge/edgemesh/agent/pkg/tunnel/controller"
	"github.com/kubeedge/edgemesh/agent/pkg/tunnel/protocol/tcp"
	"github.com/kubeedge/edgemesh/common/informers"
	"github.com/kubeedge/edgemesh/common/modules"
	"github.com/libp2p/go-libp2p-core/host"
)

var Agent *TunnelAgent

// TunnelAgent is used for solving cross subset communication
type TunnelAgent struct {
	Config      *config.TunnelAgentConfig
	Host        host.Host
	TCPProxySvc *tcp.TCPProxyService
}

func newTunnelAgent(c *config.TunnelAgentConfig, ifm *informers.Manager) (agent *TunnelAgent, err error) {
	agent = &TunnelAgent{Config: c}
	if !c.Enable {
		return agent, nil
	}

	controller.Init(ifm)
	Agent = agent
	return agent, nil
}

// Register register tunnelagent to beehive modules
func Register(c *config.TunnelAgentConfig, ifm *informers.Manager) error {
	agent, err := newTunnelAgent(c, ifm)
	if err != nil {
		return fmt.Errorf("register module tunnelagent error: %v", err)
	}
	core.Register(agent)
	return nil
}

// Name of tunnelagent
func (t *TunnelAgent) Name() string {
	return modules.TunnelAgentModuleName
}

// Group of tunnelagent
func (t *TunnelAgent) Group() string {
	return modules.TunnelAgentModuleName
}

// Enable indicates whether enable this module
func (t *TunnelAgent) Enable() bool {
	return t.Config.Enable
}

// Start tunnelserver
func (t *TunnelAgent) Start() {
	t.Run()
}
