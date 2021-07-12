package config

import (
	"io/ioutil"
	"os"
	"path"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"

	"github.com/kubeedge/kubeedge/common/constants"
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/cloudcore/v1alpha1"

	chassisconfig "github.com/kubeedge/edgemesh/agent/pkg/common/chassis/config"
	dnsconfig "github.com/kubeedge/edgemesh/agent/pkg/dns/config"
	gwconfig "github.com/kubeedge/edgemesh/agent/pkg/gateway/config"
	proxyconfig "github.com/kubeedge/edgemesh/agent/pkg/proxy/config"
	tunnelconfig "github.com/kubeedge/edgemesh/agent/pkg/tunnel/config"
)

const (
	GroupName  = "agent.edgemesh.config.kubeedge.io"
	APIVersion = "v1alpha1"
	Kind       = "EdgeMeshAgent"
)

// EdgeMeshAgentConfig indicates the config of edgeMeshAgent which get from edgeMeshAgent config file
type EdgeMeshAgentConfig struct {
	metav1.TypeMeta
	// KubeAPIConfig indicates the kubernetes cluster info which edgeMeshAgent will connected
	// +Required
	KubeAPIConfig *v1alpha1.KubeAPIConfig `json:"kubeAPIConfig,omitempty"`
	// GoChassisConfig defines some configurations related to go-chassis
	// +Required
	GoChassisConfig *chassisconfig.GoChassisConfig `json:"goChassisConfig,omitempty"`
	// Modules indicates edgeMeshAgent modules config
	// +Required
	Modules *Modules `json:"modules,omitempty"`
}

// Modules indicates the modules of edgeMeshAgent will be use
type Modules struct {
	// EdgeDNSConfig indicates edgedns module config
	EdgeDNSConfig *dnsconfig.EdgeDNSConfig `json:"edgeDNS,omitempty"`
	// EdgeProxyConfig indicates edgeproxy module config
	EdgeProxyConfig *proxyconfig.EdgeProxyConfig `json:"edgeProxy,omitempty"`
	// EdgeGatewayConfig indicates edgegateway module config
	EdgeGatewayConfig *gwconfig.EdgeGatewayConfig `json:"edgeGateway,omitempty"`
	Tunnel *tunnelconfig.Tunnel `json:"tunnel,omitempty"`
}

// NewEdgeMeshAgentConfig returns a full EdgeMeshAgentConfig object
func NewEdgeMeshAgentConfig() *EdgeMeshAgentConfig {
	token := os.Getenv("test")
	if token == "" {
		klog.Fatal("CloudCore Token is empty, Please provide")
	}
	hostnameOverride, err := os.Hostname()
	if err != nil {
		hostnameOverride = constants.DefaultHostnameOverride
	}

	c := &EdgeMeshAgentConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       Kind,
			APIVersion: path.Join(GroupName, APIVersion),
		},
		KubeAPIConfig: &v1alpha1.KubeAPIConfig{
			Master:      "",
			ContentType: constants.DefaultKubeContentType,
			QPS:         constants.DefaultKubeQPS,
			Burst:       constants.DefaultKubeBurst,
			KubeConfig:  constants.DefaultKubeConfig,
		},
		GoChassisConfig: &chassisconfig.GoChassisConfig{
			Protocol: &chassisconfig.Protocol{
				TCPBufferSize:     8192,
				TCPClientTimeout:  2,
				TCPReconnectTimes: 3,
			},
			LoadBalancer: &chassisconfig.LoadBalancer{
				DefaultLBStrategy:     "RoundRobin",
				SupportedLBStrategies: []string{"RoundRobin", "Random", "ConsistentHash"},
				ConsistentHash: &chassisconfig.ConsistentHash{
					PartitionCount:    100,
					ReplicationFactor: 10,
					Load:              1.25,
				},
			},
		},
		Modules: &Modules{
			EdgeDNSConfig: &dnsconfig.EdgeDNSConfig{
				Enable:          true,
				ListenInterface: "docker0",
				ListenPort:      53,
			},
			EdgeProxyConfig: &proxyconfig.EdgeProxyConfig{
				Enable:          true,
				SubNet:          "10.0.0.0/24",
				ListenInterface: "docker0",
				ListenPort:      40001,
			},
			EdgeGatewayConfig: &gwconfig.EdgeGatewayConfig{
				Enable:    true,
				NIC:       "*",
				IncludeIP: "*",
				ExcludeIP: "*",
			},
			Tunnel: &tunnelconfig.Tunnel{
				Enable:             true,
				Heartbeat:          15,
				TLSCAFile:          constants.DefaultCAFile,
				TLSCertFile:        constants.DefaultCertFile,
				TLSPrivateKeyFile:  constants.DefaultKeyFile,
				RotateCertificates: true,
				HostnameOverride:   hostnameOverride,
				// TODO fetch token from env or file ,which come from the tokensecret
				Token: token,
			},
		},
	}

	return c
}

// Parse unmarshal config file into *EdgeMeshAgentConfig
func (c *EdgeMeshAgentConfig) Parse(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		klog.Errorf("Failed to read config file %s: %v", filename, err)
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		klog.Errorf("Failed to unmarshal config file %s: %v", filename, err)
		return err
	}
	return nil
}
