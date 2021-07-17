package main

import (
	"github.com/meschbach/elevatinator/simulator"
	"fmt"
)

func main() {
	stream := simulator.NewEventLog()
	maxTicks := simulator.Tick(100)

	simulation := simulator.NewSimulation()
	simulation.AttachControllerListener(stream)
	simulation.AttachActor(simulator.NewActor(3,0,0))
	//simulation.AttachActor(simulator.NewActor(2,0,0))
	//simulation.AttachActor(simulator.NewActor(0,1,1))
	simulation.Initialize(1, 4)
	simulation.AttachControllerFunc(simulator.NewMoveController)

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
