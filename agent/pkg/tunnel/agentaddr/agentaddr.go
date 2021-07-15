package agentaddr

import (
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"k8s.io/klog/v2"
	"sync"
)

const (
	DEFAULT_SECRET_NAMESPACE = "kubeedge"
	DEFAULT_SECRET_NAME      = "edgemeshagentsecret"
)

var (
	once sync.Once
	aa   *PeerAgentAddr
)

type PeerAgentAddr struct {
	sync.RWMutex
	nodeName2PeerAddr map[string][]byte
	// todo add secret getter
}

func NewPeerAgentAddr() *PeerAgentAddr {
	once.Do(func() {
		aa = &PeerAgentAddr{
			nodeName2PeerAddr: make(map[string][]byte),
		}
	})
	return aa
}

func (na *PeerAgentAddr) SetSelfAddr2Secret(nodeName string, id peer.ID, addrs []ma.Multiaddr) error {
	na.Lock()
	defer na.Unlock()

	for k, v := range addrs {
		newAddr := fmt.Sprintf("%v/p2p/%v", v, id)
		newMultiAddr, err := ma.NewMultiaddr(newAddr)
		if err != nil {
			klog.Errorf("%s transfer to multiaddr err: %v", newAddr, err)
			return err
		}
		addrs[k] = newMultiAddr
	}
	na.nodeName2PeerAddr[nodeName] = ma.Join(addrs...).Bytes()

	// todo
	/*
		secret, _ := client.GetKubeClient().CoreV1().Secrets(na.SecretNameSpace()).Get(context.Background(), na.SecretName(), v1.GetOptions{})
		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}
		secret.Data[nodeName] = ma.Join(addrs...).Bytes()
		secret, err := client.GetKubeClient().CoreV1().Secrets("kubeedge").Update(context.Background(), secret, v1.UpdateOptions{})
		if err != nil {
			klog.Errorf("Update secret %v err: ", secret, err)
			return err
		}
	*/
	return nil
}

func (na *PeerAgentAddr) Get(nodeName string) (peerAddr ma.Multiaddr, err error) {
	na.RLock()
	defer na.RUnlock()
	addr := na.nodeName2PeerAddr[nodeName]
	if addr == nil {
		// todo
		/*
			klog.Warningf("Get %s addr from cache failed, get from api server", nodeName)
			secret, err := client.GetKubeClient().CoreV1().Secrets(na.SecretNameSpace()).Get(context.Background(), na.SecretName(), v1.GetOptions{})
			if err != nil {
				klog.Errorf("Get %s addr from api server err: %v", nodeName, err)
				return nil, err
			}
			addr = secret.Data[nodeName]
		*/
	}

	peerAddr, err = ma.NewMultiaddrBytes(addr)
	if err != nil {
		klog.Errorf("%s transfer to multiAddr err: %v", string(addr), err)
		return nil, err
	}
	return peerAddr, nil
}

func (na *PeerAgentAddr) Reset(agentAddrs map[string][]byte) {
	na.Lock()
	defer na.Unlock()
	for k, v := range agentAddrs {
		na.nodeName2PeerAddr[k] = v
	}
}

func (na *PeerAgentAddr) SecretNameSpace() string {
	// todo read from config file
	return DEFAULT_SECRET_NAMESPACE
}

func (na *PeerAgentAddr) SecretName() string {
	return DEFAULT_SECRET_NAME
}
