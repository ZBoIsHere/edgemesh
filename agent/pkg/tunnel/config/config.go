package config

import "github.com/kubeedge/edgemesh/common/certificate"

type TunnelAgentConfig struct {
	// Enable indicates whether EdgeHub is enabled,
	// if set to false (for debugging etc.), skip checking other EdgeHub configs.
	// default true
	Enable bool `json:"enable"`
	// TunnelServer indicates the server address of edgemesh server
	TunnelServer string `json:"tunnelServer"`
	// TunnelCertificate indicates the set of tunnel server config about certificate
	certificate.TunnelCertificate
	// NodeName indicates the node name of tunnel server
	NodeName string `json:"nodeName"`
}
