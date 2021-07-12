package controller

import (
	"fmt"
	"strings"
	"sync"

	v1 "k8s.io/api/core/v1"
	listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	"github.com/kubeedge/edgemesh/agent/pkg/common/informers"
)

type ProxyController interface {
	Init()
	GetSvcPorts(ip string) string
	GetSvcIP(svcName string) string
}

type controller struct {
	svcLister   listers.ServiceLister
	svcInformer cache.SharedIndexInformer

	sync.RWMutex
	SvcPortsByIP map[string]string // key: clusterIP, value: SvcPorts
	IPBySvc      map[string]string // key: svcName.svcNamespace, value: clusterIP
}

func New(ifm *informers.Manager) *controller {
	c := &controller{
		SvcPortsByIP: make(map[string]string),
		IPBySvc:      make(map[string]string),
	}
	kubeFactory := ifm.GetKubeFactory()
	// get lister
	c.svcLister = kubeFactory.Core().V1().Services().Lister()
	// get informer
	c.svcInformer = kubeFactory.Core().V1().Services().Informer()
	// register informers
	ifm.Register(c.svcInformer)
	return c
}

func (c *controller) Init() {
	// set informers event handler
	c.svcInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.svcAdd, UpdateFunc: c.svcUpdate, DeleteFunc: c.svcDelete})
}

func getSvcPorts(svc *v1.Service) string {
	svcPorts := ""
	svcName := svc.Namespace + "." + svc.Name
	for _, p := range svc.Spec.Ports {
		pro := strings.Split(p.Name, "-")
		sub := fmt.Sprintf("%s,%d,%d|", pro[0], p.Port, p.TargetPort.IntVal)
		svcPorts = svcPorts + sub
	}
	svcPorts += svcName
	return svcPorts
}

func (c *controller) svcAdd(obj interface{}) {
	svc, ok := obj.(*v1.Service)
	if !ok {
		klog.Errorf("invalid type %v", obj)
		return
	}
	svcPorts := getSvcPorts(svc)
	svcName := svc.Namespace + "." + svc.Name
	ip := svc.Spec.ClusterIP
	if ip == "" || ip == "None" {
		return
	}
	c.addOrUpdateService(svcName, ip, svcPorts)
}

func (c *controller) svcUpdate(oldObj, newObj interface{}) {
	svc, ok := newObj.(*v1.Service)
	if !ok {
		klog.Errorf("invalid type %v", newObj)
		return
	}
	svcPorts := getSvcPorts(svc)
	svcName := svc.Namespace + "." + svc.Name
	ip := svc.Spec.ClusterIP
	if ip == "" || ip == "None" {
		return
	}
	c.addOrUpdateService(svcName, ip, svcPorts)
}

func (c *controller) svcDelete(obj interface{}) {
	svc, ok := obj.(*v1.Service)
	if !ok {
		klog.Errorf("invalid type %v", obj)
		return
	}
	svcName := svc.Namespace + "." + svc.Name
	ip := svc.Spec.ClusterIP
	if ip == "" || ip == "None" {
		return
	}
	c.deleteService(svcName, ip)
}

// AddOrUpdateService add or updates a service
func (c *controller) addOrUpdateService(svcName, ip, svcPorts string) {
	c.Lock()
	defer c.Unlock()
	c.IPBySvc[svcName] = ip
	c.SvcPortsByIP[ip] = svcPorts
}

// DeleteService deletes a service
func (c *controller) deleteService(svcName, ip string) {
	c.Lock()
	defer c.Unlock()
	delete(c.IPBySvc, svcName)
	delete(c.SvcPortsByIP, ip)
}

// GetSvcIP returns the ip by given service name
func (c *controller) GetSvcIP(svcName string) string {
	c.RLock()
	defer c.RUnlock()
	ip := c.IPBySvc[svcName]
	return ip
}

// GetSvcPorts is a thread-safe operation to get from map
func (c *controller) GetSvcPorts(ip string) string {
	c.RLock()
	defer c.RUnlock()
	svcPorts := c.SvcPortsByIP[ip]
	return svcPorts
}
