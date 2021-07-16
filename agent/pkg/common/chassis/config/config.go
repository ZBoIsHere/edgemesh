package config

import "sync"

// GoChassisConfig defines some configurations related to go-chassis
type GoChassisConfig struct {
	// Protocol indicates the network protocols config supported in edgemesh
	Protocol *Protocol `json:"protocol,omitempty"`
	// LoadBalancer indicates the load balance strategy
	LoadBalancer *LoadBalancer `json:"loadBalancer,omitempty"`
}

// Protocol indicates the network protocols config supported in edgemesh
type Protocol struct {
	// TCPBufferSize indicates 4-layer tcp buffer size
	// default 8192
	TCPBufferSize int `json:"tcpBufferSize,omitempty"`
	// TCPClientTimeout indicates 4-layer tcp client timeout, the unit is second.
	// default 2
	TCPClientTimeout int `json:"tcpClientTimeout,omitempty"`
	// TCPReconnectTimes indicates 4-layer tcp reconnect times
	// default 3
	TCPReconnectTimes int `json:"tcpReconnectTimes,omitempty"`
}

// LoadBalancer indicates the loadbalance strategy in edgemesh
type LoadBalancer struct {
	// DefaultLBStrategy indicates default load balance strategy name
	// default "RoundRobin"
	DefaultLBStrategy string `json:"defaultLBStrategy,omitempty"`
	// SupportedLBStrategies indicates supported load balance strategies name
	// default []string{"RoundRobin", "Random", "ConsistentHash"}
	SupportedLBStrategies []string `json:"supportLBStrategies,omitempty"`
	// ConsistentHash indicates the extension of the go-chassis loadbalancer
	ConsistentHash *ConsistentHash `json:"consistentHash,omitempty"`
}

// ConsistentHash strategy is an extension of the go-chassis loadbalancer
// For more information about the consistentHash algorithm, please take a look at
// https://research.googleblog.com/2017/04/consistent-hashing-with-bounded-loads.html
type ConsistentHash struct {
	// PartitionCount indicates the hash ring partition count
	// default 100
	PartitionCount int `json:"partitionCount,omitempty"`
	// ReplicationFactor indicates the hash ring replication factor
	// default 10
	ReplicationFactor int `json:"replicationFactor,omitempty"`
	// Load indicates the hash ring bounded loads
	// default 1.25
	Load float64 `json:"load,omitempty"`
}

var (
	once    sync.Once
	Chassis Configure
)

type Configure struct {
	GoChassisConfig
}

// InitConfigure init go-chassis configures
func InitConfigure(c *GoChassisConfig) {
	once.Do(func() {
		Chassis = Configure{
			GoChassisConfig: *c,
		}
	})
}
