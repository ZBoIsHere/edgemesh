package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/kubeedge/edgemesh/server/cmd/edgemesh-server/app"
	"k8s.io/component-base/logs"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	command := app.NewEdgeMeshServerCommand()

	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
