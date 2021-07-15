package config

import "github.com/kubeedge/edgemesh/common/certificate"

// Tunnel indicates networking module config
type TunnelServerConfig struct {
	// Enable indicates whether Tunnel is enabled,
	// if set to false (for debugging etc.), skip checking other Networking configs.
	// default true
	Enable bool `json:"enable,omitempty"`
	// TunnelCertificate indicates the set of tunnel server config about certificate
	certificate.TunnelCertificate
	// NodeName indicates the node name of tunnel server
	NodeName string `json:"nodeName"`
}
