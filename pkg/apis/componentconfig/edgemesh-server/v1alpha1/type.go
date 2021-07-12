package v1alpha1

import (
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/cloudcore/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	metav1.TypeMeta
	// KubeAPIConfig indicates the kubernetes cluster info which cloudCore will connected
	// default use InClusterConfig
	// +Optional
	KubeAPIConfig *v1alpha1.KubeAPIConfig `json:"kubeAPIConfig,omitempty"`
	// Modules indicates edgeMesh modules config
	// +Required
	Modules *Modules `json:"modules,omitempty"`
}

// Modules indicates the modules of EdgeMesh-Server will be use
type Modules struct {
	// Tunnel indicates CloudHub module config
	Tunnel *Tunnel `json:"tunnel,omitempty"`
}

// Tunnel indicates networking module config
type Tunnel struct {
	// Enable indicates whether Tunnel is enabled,
	// if set to false (for debugging etc.), skip checking other Networking configs.
	// default true
	Enable bool `json:"enable,omitempty"`
	// Heartbeat indicates heart beat (second)
	// default 15
	Heartbeat int32 `json:"heartbeat,omitempty"`
	// TLSCAFile set ca file path
	// default "/etc/kubeedge/ca/rootCA.crt"
	TLSCAFile string `json:"tlsCaFile,omitempty"`
	// TLSCertFile indicates the file containing x509 Certificate for HTTPS
	// default "/etc/kubeedge/certs/server.crt"
	TLSCertFile string `json:"tlsCertFile,omitempty"`
	// TLSPrivateKeyFile indicates the file containing x509 private key matching tlsCertFile
	// default "/etc/kubeedge/certs/server.key"
	TLSPrivateKeyFile string `json:"tlsPrivateKeyFile,omitempty"`
	// Token indicates the priority of joining the cluster for the edge
	Token string `json:"token"`
	// HTTPServer indicates the server for edge to apply for the certificate.
	HTTPServer string `json:"httpServer,omitempty"`
	// RotateCertificates indicates whether edge certificate can be rotated
	// default true
	RotateCertificates bool `json:"rotateCertificates,omitempty"`
	// HostnameOverride indicates hostname
	// default os.Hostname()
	HostnameOverride string `json:"hostnameOverride,omitempty"`
}
