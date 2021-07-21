package scenarios

import (
	"fmt"
	"github.com/meschbach/elevatinator/simulator"
)

type Scenario func(simulation *simulator.Simulation)simulator.Tick

func RunScenario(factory simulator.ControllerFunc, scenario Scenario)  {
	stream := simulator.NewEventLog()

	simulation := simulator.NewSimulation()
	simulation.AttachControllerListener(stream)
	maxTicks := scenario(simulation)
	simulation.AttachControllerFunc(factory)

	tick := simulation.TickUpTo(maxTicks)
	if simulation.ActorsCompletedObjectives() {
		fmt.Printf("WIN!!! All actors completed objectives at tick %d\n", tick)
	} else {
		fmt.Printf(":-( Some actors did not make it to their objectives @ tick %d\n", tick)
		fmt.Println("Event stream:")
		for _, e := range stream.Events {
			fmt.Printf("\t- %s\n", e.ToString())
		}
	}
}
