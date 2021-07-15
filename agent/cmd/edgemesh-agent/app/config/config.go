package config

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/kubeedge/edgemesh/common/certificate"

	chassisconfig "github.com/kubeedge/edgemesh/agent/pkg/common/chassis/config"
	dnsconfig "github.com/kubeedge/edgemesh/agent/pkg/dns/config"
	gwconfig "github.com/kubeedge/edgemesh/agent/pkg/gateway/config"
	proxyconfig "github.com/kubeedge/edgemesh/agent/pkg/proxy/config"
	tunnelconfig "github.com/kubeedge/edgemesh/agent/pkg/tunnel/config"
	meshConstants "github.com/kubeedge/edgemesh/common/constants"
	"github.com/kubeedge/kubeedge/common/constants"
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/cloudcore/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
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
	// TunnelAgentConfig indicates tunnelagent module config
	TunnelAgentConfig *tunnelconfig.TunnelAgentConfig `json:"tunnel,omitempty"`
}

// NewEdgeMeshAgentConfig returns a full EdgeMeshAgentConfig object
func NewEdgeMeshAgentConfig() *EdgeMeshAgentConfig {
	nodeName, isExist := os.LookupEnv(meshConstants.MY_NODE_NAME)
	if !isExist {
		klog.Fatalf("env %s not exist", meshConstants.MY_NODE_NAME)
		os.Exit(1)
	}
	cloudcoreToken, isExist := os.LookupEnv(meshConstants.CLOUDCORE_TOKEN)
	if !isExist {
		klog.Fatalf("env %s not exist", meshConstants.CLOUDCORE_TOKEN)
		os.Exit(1)
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
			TunnelAgentConfig: &tunnelconfig.TunnelAgentConfig{
				Enable: true,
				TunnelCertificate: certificate.TunnelCertificate{
					TLSCAFile:          constants.DefaultCAFile,
					TLSCertFile:        constants.DefaultCertFile,
					TLSPrivateKeyFile:  constants.DefaultKeyFile,
					Token:              cloudcoreToken,
					HTTPServer:         "https://127.0.0.1:10002",
					RotateCertificates: true,
				},
				NodeName: nodeName,
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
