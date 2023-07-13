package simulator

type Actor struct {
	floorGoal         int
	startingFloor     int
	startingTick      Tick
	completedGoalTick Tick

	state   int
	actorID int
}

const (
	Unstarted = iota
	Finished
	WaitingOnFloor
	EnteringElevator
	WaitingInElevator
)

func (a *Actor) Tick(simulation *Simulation, tick Tick) {
	switch a.state {
	case Finished:
		return
	case Unstarted:
		if tick < a.startingTick {
			return
		}
		a.actorID = simulation.StartAt(a, a.startingFloor)
		simulation.callElevator(a.startingFloor)
		a.state = WaitingOnFloor
	case WaitingOnFloor:
		elevatorIDs := simulation.ElevatorsAt(a.actorID)
		if len(elevatorIDs) == 0 {
			return
		}
		simulation.Enter(a.actorID, elevatorIDs[0])
		a.state = EnteringElevator
	case EnteringElevator:
		simulation.PressButton(a.actorID, a.floorGoal)
		a.state = WaitingInElevator
	default:
	}
}

func (a *Actor) elevatorStopped(simulation *Simulation, tick Tick, floor int) {
	switch a.state {
	case WaitingInElevator:
		if simulation.ElevatorAtFloor(a.actorID, FloorID(a.floorGoal)) {
			a.completedGoalTick = tick
			a.state = Finished
			simulation.exitElevator(a.actorID)
			simulation.dispatchControllerEvent(OnActorFinished(simulation.tick, 1))
		}
	}
}

func (a *Actor) done() bool {
	return a.state == Finished
}

func NewActor(goal int, startingFloor int, startingTick Tick) *Actor {
	return &Actor{
		floorGoal:         goal,
		startingFloor:     startingFloor,
		startingTick:      startingTick,
		completedGoalTick: -1,
		state:             Unstarted,
	}
}
