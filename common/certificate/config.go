package certificate

type TunnelCertificate struct {
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
}