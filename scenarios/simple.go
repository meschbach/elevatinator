package scenarios

import "github.com/meschbach/elevatinator/simulator"

func SinglePersonUp(simulation *simulator.Simulation) simulator.Tick  {
	simulation.AttachActor(simulator.NewActor(4,0,0))
	simulation.Initialize(1, 5)
	return 20
}

func SinglePersonDown(simulation *simulator.Simulation) simulator.Tick {
	simulation.AttachActor(simulator.NewActor(2,4,0))
	simulation.Initialize(1, 5)
	return 20
}

func MultipleUpAndBack(simulation *simulator.Simulation) simulator.Tick {
	simulation.AttachActor(simulator.NewActor(3,0,0))
	simulation.AttachActor(simulator.NewActor(2,0,8))
	simulation.AttachActor(simulator.NewActor(0,1,17))
	simulation.Initialize(1, 5)
	return 40
}
