package controller

import (
	"fmt"

	listers "k8s.io/client-go/listers/core/v1"

	"github.com/kubeedge/edgemesh/agent/pkg/common/informers"
)

type DNSController interface {
	Init()
	GetSvcIP(namespace, name string) (ip string, err error)
}

type controller struct {
	svcLister listers.ServiceLister
}

func New(ifm *informers.Manager) *controller {
	c := &controller{}
	kubeFactory := ifm.GetKubeFactory()
	// get lister
	c.svcLister = kubeFactory.Core().V1().Services().Lister()
	return c
}

func (c *controller) Init() {}

// GetSvcIP get service cluster ip
func (c *controller) GetSvcIP(namespace, name string) (ip string, err error) {
	svc, err := c.svcLister.Services(namespace).Get(name)
	if err != nil {
		return "", err
	}
	ip = svc.Spec.ClusterIP
	if ip == "" || ip == "None" {
		return "", fmt.Errorf("service %s.%s no cluster ip", name, namespace)
	}
	return ip, nil
}
