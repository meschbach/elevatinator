package scenarios

import (
	"fmt"
	simulator2 "github.com/meschbach/elevatinator/pkg/simulator"
	"testing"
)

// Scenario configures a simulation for a specific puzzle.  A Scenario provides the maximum number of possible ticks to
// allow a controller to run for all simulation.Actors to win.
type Scenario func(simulation *simulator2.Simulation) simulator2.Tick

// RunScenario runs the given scenario against the controller produced via the factory.  If the controller completes the
// scenario in less than the maximum allowed ticks then the count is produced.  Otherwise the event log from the run is
// written out.
func RunScenario(factory simulator2.ControllerFunc, scenario Scenario) {
	stream := simulator2.NewEventLog()

	simulation := simulator2.NewSimulation()
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

// TestScenario integrates with Go's built-in testing framework to assert a given controller is able to complete the
// given scenario.  This is useful for functional level integration testing with Controllers.
func TestScenario(t *testing.T, factory simulator2.ControllerFunc, scenario Scenario) {
	stream := simulator2.NewEventLog()

	simulation := simulator2.NewSimulation()
	simulation.AttachControllerListener(stream)
	maxTicks := scenario(simulation)
	simulation.AttachControllerFunc(factory)

	tick := simulation.TickUpTo(maxTicks)
	if simulation.ActorsCompletedObjectives() {
		t.Logf("WIN!!! All actors completed objectives at tick %d", tick)
	} else {
		t.Errorf(":-( Some actors did not make it to their objectives @ tick %d", tick)
		t.Logf("Event stream:")
		for _, e := range stream.Events {
			switch e.EventType {
			case simulator2.TickStart:
				//do nothing
			case simulator2.TickDone:
				t.Logf("-- Tick %d --", e.Timestamp)
			default:
				t.Logf("\t- %s\n", e.ToString())
			}
		}
	}
}
