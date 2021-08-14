package main

import (
	"github.com/meschbach/elevatinator/ipc/grpc/telepathy"
	"github.com/meschbach/elevatinator/scenarios"
)

func main()  {
	bridge, err := telepathy.DialLanding("localhost:9998")
	if err != nil {
		panic(err)
	}

	//scenarios := []scenarios.Scenario{scenarios.SinglePersonUp}
	scenario := scenarios.SinglePersonDown
	scenarios.RunScenario(bridge.ControllerAdapter(), scenario)
}
