package config

import (
	"os"

	"github.com/kubeedge/edgemesh/common/certificate"
	meshConstants "github.com/kubeedge/edgemesh/common/constants"
	"github.com/kubeedge/kubeedge/common/constants"
	"k8s.io/klog/v2"
)

type TunnelAgentConfig struct {
	// Enable indicates whether TunnelAgent is enabled,
	// if set to false (for debugging etc.), skip checking other TunnelAgent configs.
	// default true
	Enable bool `json:"enable"`
	// TunnelServer indicates the server address of edgemesh server
	TunnelServer string `json:"tunnelServer"`
	// TunnelCertificate indicates the set of tunnel agent config about certificate
	certificate.TunnelCertificate
	// NodeName indicates the node name of tunnel agent
	NodeName string `json:"nodeName"`
	// ListenPort indicates the listen port of tunnel agent
	// default 10006
	ListenPort int `json:"listenPort"`
}

func NewTunnelAgentConfig() *TunnelAgentConfig {
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

	return &TunnelAgentConfig{
		Enable: true,
		TunnelCertificate: certificate.TunnelCertificate{
			TLSCAFile:          constants.DefaultCAFile,
			TLSCertFile:        constants.DefaultCertFile,
			TLSPrivateKeyFile:  constants.DefaultKeyFile,
			Token:              cloudcoreToken,
			HTTPServer:         "https://127.0.0.1:10002",
			RotateCertificates: true,
		},
		NodeName:   nodeName,
		ListenPort: 10006,
	}
}
