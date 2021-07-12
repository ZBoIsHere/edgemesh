package lister

import (
	istiolisters "istio.io/client-go/pkg/listers/networking/v1alpha3"
	k8slisters "k8s.io/client-go/listers/core/v1"

	"github.com/kubeedge/edgemesh/agent/pkg/common/informers"
)

var mgr *Manager

type Manager struct {
	podLister    k8slisters.PodLister
	secretLister k8slisters.SecretLister
	svcLister    k8slisters.ServiceLister
	epLister     k8slisters.EndpointsLister
	drLister     istiolisters.DestinationRuleLister
	gwLister     istiolisters.GatewayLister
	vsLister     istiolisters.VirtualServiceLister
}

// Init lister manager
func Init(ifm *informers.Manager) {
	mgr = new(Manager)
	kubeFactor := ifm.GetKubeFactory()
	istioFactor := ifm.GetIstioFactory()

	mgr.podLister = kubeFactor.Core().V1().Pods().Lister()
	mgr.secretLister = kubeFactor.Core().V1().Secrets().Lister()
	mgr.svcLister = kubeFactor.Core().V1().Services().Lister()
	mgr.epLister = kubeFactor.Core().V1().Endpoints().Lister()
	mgr.drLister = istioFactor.Networking().V1alpha3().DestinationRules().Lister()
	mgr.gwLister = istioFactor.Networking().V1alpha3().Gateways().Lister()
	mgr.vsLister = istioFactor.Networking().V1alpha3().VirtualServices().Lister()
}

func GetPodLister() k8slisters.PodLister {
	return mgr.podLister
}

func GetSecretLister() k8slisters.SecretLister {
	return mgr.secretLister
}

func GetSvcLister() k8slisters.ServiceLister {
	return mgr.svcLister
}

func GetEpLister() k8slisters.EndpointsLister {
	return mgr.epLister
}

func GetDrLister() istiolisters.DestinationRuleLister {
	return mgr.drLister
}

func GetGwLister() istiolisters.GatewayLister {
	return mgr.gwLister
}

func GetVsLister() istiolisters.VirtualServiceLister {
	return mgr.vsLister
}
