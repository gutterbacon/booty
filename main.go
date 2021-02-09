package main

import (
	"go.amplifyedge.org/booty-v2/dep/orchestrator"
)

var (
	// best variable name
	conductor *orchestrator.Orchestrator
)

func init() {
	conductor = orchestrator.NewOrchestrator("booty")
}

func main() {
	logger := conductor.Logger()
	if err := conductor.Command().Execute(); err != nil {
		logger.Errorf("error: %v", err)
	}
}
