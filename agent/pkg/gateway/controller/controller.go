package controller

import (
	istioapi "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	"github.com/kubeedge/edgemesh/agent/pkg/common/informers"
	"github.com/kubeedge/edgemesh/agent/pkg/gateway/config"
	"github.com/kubeedge/edgemesh/agent/pkg/gateway/manager"
)

type GatewayController interface {
	Init()
}

type controller struct {
	gwInformer cache.SharedIndexInformer
	gwManager  *manager.Manager
}

func New(ifm *informers.Manager, cfg *config.EdgeGatewayConfig) *controller {
	c := &controller{
		gwManager: manager.NewGatewayManager(cfg),
	}
	istioFactory := ifm.GetIstioFactory()
	// get informer
	c.gwInformer = istioFactory.Networking().V1alpha3().Gateways().Informer()
	// register informers
	ifm.Register(c.gwInformer)
	return c
}

func (c *controller) Init() {
	// set informers event handler
	c.gwInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.gwAdd, UpdateFunc: c.gwUpdate, DeleteFunc: c.gwDelete})
}

func (c *controller) gwAdd(obj interface{}) {
	gw, ok := obj.(*istioapi.Gateway)
	if !ok {
		klog.Errorf("invalid type %v", obj)
		return
	}
	c.gwManager.AddGateway(gw)
}

func (c *controller) gwUpdate(oldObj, newObj interface{}) {
	gw, ok := newObj.(*istioapi.Gateway)
	if !ok {
		klog.Errorf("invalid type %v", newObj)
		return
	}
	c.gwManager.UpdateGateway(gw)
}

func (c *controller) gwDelete(obj interface{}) {
	gw, ok := obj.(*istioapi.Gateway)
	if !ok {
		klog.Errorf("invalid type %v", obj)
		return
	}
	c.gwManager.DeleteGateway(gw)
}
