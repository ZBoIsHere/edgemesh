package tunnel

import (
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/edgemesh/agent/pkg/tunnel/agentaddr"
	"github.com/kubeedge/edgemesh/agent/pkg/tunnel/tunnelagent"
	"github.com/kubeedge/edgemesh/common/certificate"
	"github.com/kubeedge/edgemesh/common/constants"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"os"
)

func (t *TunnelAgent) Run() {
	certificateConfig := t.Config.TunnelCertificate
	if certificateConfig.Token == "" {
		os.LookupEnv("CloudCoreToken")
	}
	t.certManager = certificate.NewCertManager(certificateConfig, t.Config.NodeName)
	t.certManager.Start()

	go tunnelagent.StartTunnelAgent()

	for {
		select {
		case <-beehiveContext.Done():
			klog.Warning("EdgeMesh stop")
			return
		default:
		}
		msg, err := beehiveContext.Receive(constants.AgentTunnelModuleName)
		if err != nil {
			klog.Warningf("Module %s receive msg error %v", constants.AgentTunnelModuleName, err)
			continue
		}
		klog.Infof("Module %s get message: %T", constants.AgentTunnelModuleName, msg)
		process(msg)
	}

	// TODO ifRotationDone() ????, 后面要添加这个东西，如果证书轮换了，要重新进行连接
}

func process(msg model.Message) {
	resource := msg.GetResource()
	switch resource {
	case constants.ResourceTypeSecret:
		handleSecretMessage(msg)
	}
}

func handleSecretMessage(msg model.Message) {
	secret, ok := msg.GetContent().(*v1.Secret)
	if !ok {
		klog.Warningf("object type: %T unsupported", secret)
		return
	}
	if secret.GetNamespace() == agentaddr.NewPeerAgentAddr().SecretNameSpace() && secret.GetName() == agentaddr.NewPeerAgentAddr().SecretName() {
		operation := msg.GetOperation()
		switch operation {
		case model.InsertOperation, model.UpdateOperation:
			agentaddr.NewPeerAgentAddr().Reset(secret.Data)
		}
	}
}