package constants

// Resources
const (
	ResourceTypeSecret = "secret"

	MY_NODE_NAME = "MY_NODE_NAME"
	CLOUDCORE_TOKEN = "CLOUDCORE_TOKEN"

	AgentTunnelModuleName = "agent-tunnel"
	AgentTunnelGroupName  = "agent-tunnel"

	NamespaceSystem = "kubeedge"
	// Tunnel modules
	DefaultCAURL            = "/ca.crt"
	DefaultAgentCertURL     = "/agent.crt"
	DefaultHostnameOverride = "default-agent-node"
	ServerDefaultCAFile     = "/etc/kubeedge/edgemesh/server/ca/rootCA.crt"
	ServerDefaultCertFile   = "/etc/kubeedge/edgemesh/server/certs/server.crt"
	ServerDefaultKeyFile    = "/etc/kubeedge/edgemesh/server/certs/server.key"
	AgentDefaultCAFile      = "/etc/kubeedge/edgemesh/agent/ca/rootCA.crt"
	AgentDefaultCertFile    = "/etc/kubeedge/edgemesh/agent/certs/server.crt"
	AgentDefaultKeyFile     = "/etc/kubeedge/edgemesh/agent/certs/server.key"
)
