package tunnelserver

import (
	"context"
	"encoding/pem"
	"fmt"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/crypto"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
)

func StartTunnelServer() {

	// todo
	certBytes, err := ioutil.ReadFile("test")
	//certBytes, err := ioutil.ReadFile(config.Config.Tunnel.TLSPrivateKeyFile)
	if err != nil {
		klog.Errorln(err)
		return
	}

	block, _ := pem.Decode(certBytes)

	priv, err := crypto.UnmarshalECDSAPrivateKey(block.Bytes)
	if err != nil {
		klog.Errorln(err)
		return
	}

	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", 10006)),
		libp2p.EnableRelay(circuit.OptHop),
		libp2p.ForceReachabilityPrivate(),
		libp2p.Identity(priv),
	)
	if err != nil {
		klog.Fatalf("Start tunnel server failed, %v", err)
		os.Exit(1)
	}

	klog.Infoln("Start tunnel server success")
	for _, v := range host.Addrs() {
		klog.Infof("%s : %v/p2p/%s\n", "Tunnel server addr", v, host.ID().Pretty())
	}

	select {}
}
