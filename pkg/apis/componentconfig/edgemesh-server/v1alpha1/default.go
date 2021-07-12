package v1alpha1

import (
	"os"
	"path"

	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/cloudcore/v1alpha1"
	"github.com/kubeedge/kubeedge/common/constants"
	meshConst "github.com/kubeedge/edgemesh/pkg/common/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func NewDefaultEdgeMeshServerConfig() *Config {
	token := os.Getenv(meshConst.CloudCoreToken)
	if token == "" {
		klog.Fatal("CloudCore Token is empty, Please provide")
	}
	hostnameOverride, err := os.Hostname()
	if err != nil {
		hostnameOverride = constants.DefaultHostnameOverride
	}

	c := &Config{
		TypeMeta:      metav1.TypeMeta{
			Kind: Kind,
			APIVersion: path.Join(GroupName, APIVersion),
		},
		KubeAPIConfig: &v1alpha1.KubeAPIConfig{
			Master:      "",
			ContentType: constants.DefaultKubeContentType,
			QPS:         constants.DefaultKubeQPS,
			Burst:       constants.DefaultKubeBurst,
			KubeConfig:  "",
		},
		Modules:       &Modules{
			Tunnel: &Tunnel{
				Enable:             true,
				Heartbeat:          15,
				TLSCAFile:          constants.DefaultCAFile,
				TLSCertFile:        constants.DefaultCertFile,
				TLSPrivateKeyFile:  constants.DefaultKeyFile,
				RotateCertificates: true,
				HostnameOverride:   hostnameOverride,
				// TODO fetch token from env or file ,which come from the tokensecret
				Token: token,
			}},
	}
	return c
}