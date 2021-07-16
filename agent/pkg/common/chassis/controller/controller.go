package controller

import (
	"fmt"
	"strings"

	"github.com/buraksezer/consistent"
	apiv1alpha3 "istio.io/api/networking/v1alpha3"
	istioapi "istio.io/client-go/pkg/apis/networking/v1alpha3"
	istiolisters "istio.io/client-go/pkg/listers/networking/v1alpha3"
	v1 "k8s.io/api/core/v1"
	k8slisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	"github.com/kubeedge/edgemesh/agent/pkg/common/chassis/loadbalancer/consistenthash/hashring"
	"github.com/kubeedge/edgemesh/agent/pkg/common/informers"
	"github.com/kubeedge/edgemesh/agent/pkg/common/util"
)

type ChassisController interface {
	Init()
}

type controller struct {
	podLister k8slisters.PodLister
	svcLister k8slisters.ServiceLister
	drLister  istiolisters.DestinationRuleLister

	epInformer cache.SharedIndexInformer
	drInformer cache.SharedIndexInformer
}

func New(ifm *informers.Manager) *controller {
	c := &controller{}
	kubeFactory := ifm.GetKubeFactory()
	istioFactory := ifm.GetIstioFactory()
	// get lister
	c.podLister = kubeFactory.Core().V1().Pods().Lister()
	c.svcLister = kubeFactory.Core().V1().Services().Lister()
	c.drLister = istioFactory.Networking().V1alpha3().DestinationRules().Lister()
	// get informer
	c.epInformer = kubeFactory.Core().V1().Endpoints().Informer()
	c.drInformer = istioFactory.Networking().V1alpha3().DestinationRules().Informer()
	// register informers
	ifm.Register(c.epInformer)
	ifm.Register(c.drInformer)
	return c
}

func (c *controller) Init() {
	// set informers event handler
	c.epInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{UpdateFunc: c.epUpdate})
	c.drInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.drAdd, UpdateFunc: c.drUpdate, DeleteFunc: c.drDelete})
}

func (c *controller) epUpdate(oldObj, newObj interface{}) {
	ep, ok := newObj.(*v1.Endpoints)
	if !ok {
		klog.Errorf("invalid type %v", newObj)
		return
	}
	// service, endpoints and destinationRule should have the same name
	dr, err := c.drLister.DestinationRules(ep.Namespace).Get(ep.Name)
	if err == nil && dr != nil && isConsistentHashLB(dr) {
		key := fmt.Sprintf("%s.%s", ep.Namespace, ep.Name)
		svc, err := c.svcLister.Services(ep.Namespace).Get(ep.Name)
		if err != nil || svc == nil {
			klog.Errorf("no service %s exists", key)
			return
		}
		// may need to update hash ring
		c.updateHashRingByService(key, svc)
	}
}

func (c *controller) drAdd(obj interface{}) {
	dr, ok := obj.(*istioapi.DestinationRule)
	if !ok {
		klog.Errorf("invalid type %v", obj)
		return
	}
	// If the loadbalance strategy is `consistentHash`, we need to create a hash ring.
	if isConsistentHashLB(dr) {
		c.createHashRing(dr.Namespace, dr.Name)
	}
}

func (c *controller) drUpdate(oldObj, newObj interface{}) {
	oldDr, ok := oldObj.(*istioapi.DestinationRule)
	if !ok {
		klog.Errorf("invalid type %v", oldObj)
		return
	}
	newDr, ok := newObj.(*istioapi.DestinationRule)
	if !ok {
		klog.Errorf("invalid type %v", newObj)
		return
	}
	key := fmt.Sprintf("%s.%s", newDr.Namespace, newDr.Name)
	if isConsistentHashLB(oldDr) && !isConsistentHashLB(newDr) {
		// If the loadbalance strategy is updated, if it is no longer a `consistentHash` strategy,
		// we need to delete the exists hash ring.
		hashring.DeleteHashRing(key)
	} else if !isConsistentHashLB(oldDr) && isConsistentHashLB(newDr) {
		// If the loadbalance strategy is updated, and it is a `consistentHash` strategy,
		// we need to create a hash ring.
		c.createHashRing(newDr.Namespace, newDr.Name)
	}
}

func (c *controller) drDelete(obj interface{}) {
	dr, ok := obj.(*istioapi.DestinationRule)
	if !ok {
		klog.Errorf("invalid type %v", obj)
		return
	}
	if isConsistentHashLB(dr) {
		key := fmt.Sprintf("%s.%s", dr.Namespace, dr.Name)
		hashring.DeleteHashRing(key)
	}
}

func isConsistentHashLB(dr *istioapi.DestinationRule) (ok bool) {
	switch dr.Spec.TrafficPolicy.LoadBalancer.LbPolicy.(type) {
	case *apiv1alpha3.LoadBalancerSettings_ConsistentHash:
		ok = true
	default:
		ok = false
	}
	return ok
}

// createHashRing create hash ring if needed
func (c *controller) createHashRing(namespace, name string) {
	key := fmt.Sprintf("%s.%s", namespace, name)
	if _, ok := hashring.GetHashRing(key); ok {
		klog.Warningf("hash ring %s already exists in cache", key)
		return
	}
	svc, err := c.svcLister.Services(namespace).Get(name)
	if err != nil || svc == nil {
		klog.Errorf("unable to get the service bound to the destinationRule, reason: %v", err)
		return
	}
	// create new hash ring
	c.createHashRingByService(svc)
}

// createHashRingByService create and store hash ring by Service
func (c *controller) createHashRingByService(svc *v1.Service) {
	// get pods
	pods, err := c.podLister.Pods(svc.Namespace).List(util.GetPodsSelector(svc))
	if err != nil || pods == nil {
		klog.Errorf("failed to get pod list, reason: %v", err)
		return
	}
	// create service instances
	var instances []hashring.ServiceInstance
	for _, p := range pods {
		if p.Status.Phase == v1.PodRunning {
			instances = append(instances, hashring.ServiceInstance{
				Namespace:  svc.Namespace,
				Name:       svc.Name,
				InstanceIP: p.Status.HostIP,
			})
		}
	}
	// create hash ring
	hr := hashring.NewServiceInstanceHashRing(instances)
	// store hash ring
	key := fmt.Sprintf("%s.%s", svc.Namespace, svc.Name)
	hashring.AddOrUpdateHashRing(key, hr)
}

// updateHashRingByService update hash ring by Service
func (c *controller) updateHashRingByService(ring string, svc *v1.Service) {
	// get hash ring
	hr, ok := hashring.GetHashRing(ring)
	if !ok {
		klog.Errorf("cannot find hash ring %s", ring)
		return
	}
	// get pods
	pods, err := c.podLister.Pods(svc.Namespace).List(util.GetPodsSelector(svc))
	if err != nil || pods == nil {
		klog.Errorf("failed to get pod list, reason: %v", err)
		return
	}
	added, deleted := lookForDifference(hr, pods, ring)
	for _, key := range added {
		klog.Infof("add ServiceInstance %s to hash ring %s", key, ring)
		namespace, name, instanceIP := splitServiceInstanceKey(key)
		hr.Add(hashring.ServiceInstance{
			Namespace:  namespace,
			Name:       name,
			InstanceIP: instanceIP,
		})
	}
	for _, key := range deleted {
		klog.Infof("delete ServiceInstance %s from hash ring %s", key, ring)
		hr.Remove(key)
	}
	// refresh cache
	if len(added) != 0 || len(deleted) != 0 {
		hashring.AddOrUpdateHashRing(ring, hr)
	}
}

func splitServiceInstanceKey(key string) (namespace, name, instanceIP string) {
	parts := strings.Split(key, "#")
	if len(parts) != 3 {
		klog.Errorf("invalid key format")
		return
	}
	return parts[0], parts[1], parts[2]
}

// lookForDifference look for the difference between v1.Pods and HashRing.Members
func lookForDifference(hr *consistent.Consistent, pods []*v1.Pod, key string) ([]string, []string) {
	var src, dest []string
	// get source array from hr.Members
	for _, member := range hr.GetMembers() {
		si, ok := member.(hashring.ServiceInstance)
		if !ok {
			klog.Errorf("can't convert to ServiceInstance")
			continue
		}
		src = append(src, si.String())
	}
	klog.Infof("src: %+v", src)
	// build destination array from v1.Pods
	namespace, name := util.SplitServiceKey(key)
	for _, p := range pods {
		if p.DeletionTimestamp != nil {
			continue
		}
		if p.Status.Phase == v1.PodRunning {
			key := fmt.Sprintf("%s#%s#%s", namespace, name, p.Status.HostIP)
			dest = append(dest, key)
		}
	}
	klog.Infof("dest: %+v", dest)
	return arrayCompare(src, dest)
}

// arrayCompare finds the difference between two arrays.
func arrayCompare(src []string, dest []string) ([]string, []string) {
	msrc := make(map[string]byte) // source array set
	mall := make(map[string]byte) // union set
	var set []string              // intersection set

	// 1.Create a map for the source array.
	for _, v := range src {
		msrc[v] = 0
		mall[v] = 0
	}
	// 2.Elements that cannot be stored in the destination array are duplicate elements.
	for _, v := range dest {
		mall[v] = 1
		set = append(set, v)
	}
	// 3.union - intersection = all variable elements
	for _, v := range set {
		delete(mall, v)
	}
	// 4.Now, mall is a complement set, then we use mall to traverse the source array.
	// The element that can be found is the deleted element, and the element that cannot be found is the added element.
	var added, deleted []string
	for v := range mall {
		_, exist := msrc[v]
		if exist {
			deleted = append(deleted, v)
		} else {
			added = append(added, v)
		}
	}
	return added, deleted
}
