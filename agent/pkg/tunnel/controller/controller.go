package controller

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kubeedge/edgemesh/common/constants"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/client-go/kubernetes/typed/core/v1"
	k8slisters "k8s.io/client-go/listers/core/v1"
	"sync"

	"github.com/kubeedge/edgemesh/common/informers"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

var (
	APIConn *TunnelAgentController
	once    sync.Once
)

type TunnelAgentController struct {
	// used to register secret call function
	secretInformer cache.SharedIndexInformer
	// used to get or list secret
	secretLister k8slisters.SecretLister
	// used to add or update or delete secret
	secretOperator v12.SecretInterface
}

func Init(ifm *informers.Manager) *TunnelAgentController {
	once.Do(func() {
		kubeFactor := ifm.GetKubeFactory()
		APIConn = &TunnelAgentController{
			secretInformer:    kubeFactor.Core().V1().Secrets().Informer(),
			secretLister: kubeFactor.Core().V1().Secrets().Lister(),
			secretOperator: ifm.GetKubeClient().CoreV1().Secrets(constants.SECRET_NAMESPACE),
		}
		ifm.RegisterInformer(APIConn.secretInformer)
	})
	return APIConn
}

func (c *TunnelAgentController) Get(nodeName string) (peerAddr ma.Multiaddr, err error) {
	secret, err := c.secretLister.Secrets(constants.SECRET_NAMESPACE).Get(constants.SECRET_NAME)
	if err != nil {
		return nil, fmt.Errorf("Get %s addr from api server err: %v", nodeName, err)
	}
	addr := secret.Data[nodeName]
	if len(addr) == 0 {
		return nil, fmt.Errorf("Get %s addr from api server err: %v", nodeName, err)
	}

	peerAddr, err = ma.NewMultiaddrBytes(addr)
	if err != nil {
		return nil, fmt.Errorf("%s transfer to multiAddr err: %v", string(addr), err)
	}
	return peerAddr, nil
}

func (c *TunnelAgentController) SetSelfAddr2Secret(nodeName string, id peer.ID, addrs []ma.Multiaddr) error {
	for k, v := range addrs {
		newAddr := fmt.Sprintf("%v/p2p/%v", v, id)
		newMultiAddr, err := ma.NewMultiaddr(newAddr)
		if err != nil {
			klog.Errorf("%s transfer to multiaddr err: %v", newAddr, err)
			return err
		}
		addrs[k] = newMultiAddr
	}

	secret, err := c.secretLister.Secrets(constants.SECRET_NAMESPACE).Get(constants.SECRET_NAME)
	if err != nil {
		return fmt.Errorf("Get secret %s in %s failed: %v", constants.SECRET_NAME, constants.SECRET_NAMESPACE, err)
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	} else if bytes.Equal(secret.Data[nodeName], ma.Join(addrs...).Bytes()) {
		return nil
	}

	secret.Data[nodeName] = ma.Join(addrs...).Bytes()
	secret, err = c.secretOperator.Update(context.Background(), secret, v13.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("Update secret %v err: ", secret, err)
	}
	return nil
}
