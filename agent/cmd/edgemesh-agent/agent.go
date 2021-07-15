package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/kubeedge/edgemesh/agent/cmd/edgemesh-agent/app"
	"k8s.io/component-base/logs"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	command := app.NewEdgeMeshAgentCommand()

	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
