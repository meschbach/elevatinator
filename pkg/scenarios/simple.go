package scenarios

import (
	simulator2 "github.com/meschbach/elevatinator/pkg/simulator"
)

// SinglePersonUp is a scenario where an actor starts on a lower floor and moves up to a higher floor.  Despite the
// simplicity this is an isolated test case to ensure a controller properly moves a single occupant in the intended
// direction.
func SinglePersonUp(simulation *simulator2.Simulation) simulator2.Tick {
	simulation.AttachActor(simulator2.NewActor(4, 0, 0))
	simulation.Initialize(1, 5)
	return 20
}

// SinglePersonDown is a scenario where an actor starts on a higher floor and moves up to a lower floor.  Despite the
// simplicity this is an isolated test case to ensure a controller properly moves a single occupant in the intended
// direction.
func SinglePersonDown(simulation *simulator2.Simulation) simulator2.Tick {
	simulation.AttachActor(simulator2.NewActor(2, 4, 0))
	simulation.Initialize(1, 5)
	return 20
}

// MultipleUpAndBack is a scenario where elevators would have to move in multiple directions in order to service the
// actors.  This can be solved in such a way only a single elevator is operating.
func MultipleUpAndBack(simulation *simulator2.Simulation) simulator2.Tick {
	simulation.AttachActor(simulator2.NewActor(3, 0, 0))
	simulation.AttachActor(simulator2.NewActor(2, 0, 8))
	simulation.AttachActor(simulator2.NewActor(0, 1, 17))
	simulation.Initialize(1, 5)
	return 40
}
