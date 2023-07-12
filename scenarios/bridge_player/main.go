package main

import (
	"fmt"
	"github.com/meschbach/elevatinator/ipc/grpc/telepathy"
	"github.com/meschbach/elevatinator/scenarios"
)

func main() {
	bridge, err := telepathy.DialLanding("localhost:9998")
	if err != nil {
		fmt.Printf("Failed to dail bridge because %s\n", err)
		return
	}

	//scenarios := []scenarios.Scenario{scenarios.SinglePersonUp}
	scenario := scenarios.SinglePersonDown
	scenarios.RunScenario(bridge.ControllerAdapter(), scenario)
}
