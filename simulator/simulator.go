package simulator

type Tick int64

type Simulation struct {
	tick                Tick
	elevators           []*Elevator
	floors              []*Floor
	actors              []*Actor
	enteredActors       []*actorState
	controller          Controller
	controllerListeners []ControllerListener
}

func (s *Simulation) Tick() bool {
	s.tick++
	currentTick := s.tick
	s.dispatchControllerEvent(OnTickStart(currentTick))
	for i, elevator := range s.elevators {
		elevator.Tick(s, i, currentTick)
	}

	for _, actor := range s.actors {
		actor.Tick(s, currentTick)
	}
	s.dispatchControllerEvent(OnTickDone(currentTick))
	return !s.ActorsCompletedObjectives()
}

func (s *Simulation) TickUpTo(additional Tick) Tick {
	endTick := s.tick + additional
	for s.tick < endTick && s.Tick() {
	}
	return s.tick
}

func (s *Simulation) ActorsCompletedObjectives() bool {
	//TODO: Ideally there is a better way to structure this
	for _, actor := range s.actors {
		if !actor.done() {
			return false
		}
	}
	return true
}

func (s *Simulation) CurrentTick() Tick {
	return s.tick
}

func (s *Simulation) AttachActor(actor *Actor) {
	s.actors = append(s.actors, actor)
}

type ControllerFunc func(ControlledElevators) Controller

func (s *Simulation) AttachControllerFunc(factory ControllerFunc) {
	s.controller = factory(s)
	ids := make([]int, len(s.elevators))
	for i := range s.elevators {
		ids[i] = i
	}
	s.controller.Init(ids)
}

func (s *Simulation) MoveTo(elevatorID int, floor int) {
	elevator := s.elevators[elevatorID]
	elevator.moveTo(s, elevatorID, floor)
}

const (
	PlaceFloor = iota
	PlaceElevator
)

type actorState struct {
	placeType  int
	placeIndex int
}

func (s *Simulation) StartAt(a *Actor, floor int) int {
	state := &actorState{
		placeType:  PlaceFloor,
		placeIndex: floor,
	}
	id := len(s.enteredActors)
	s.enteredActors = append(s.enteredActors, state)
	return id
}

func (s *Simulation) ElevatorsAt(id int) []int {
	state := s.enteredActors[id]
	found := make([]int, 0)
	switch state.placeType {
	case PlaceFloor:
		for id, e := range s.elevators {
			if e.currentFloor == state.placeIndex {
				found = append(found, id)
			}
		}
	}
	return found
}

func (s *Simulation) ElevatorAtFloor(actorID int, floor FloorID) bool {
	state := s.enteredActors[actorID]
	switch state.placeType {
	case PlaceElevator:
		elevator := s.elevators[state.placeIndex]
		return elevator.isAtFloor(s, floor)
	}
	return false
}

func (s *Simulation) Enter(actorID int, elevatorID int) {
	state := s.enteredActors[actorID]
	switch state.placeType {
	case PlaceFloor:
		if s.elevators[elevatorID].currentFloor == state.placeIndex {
			state.placeType = PlaceElevator
			state.placeIndex = elevatorID
		}
	}
}

func (s *Simulation) exitElevator(actorID int) {
	state := s.enteredActors[actorID]
	switch state.placeType {
	case PlaceElevator:
		state.placeType = PlaceFloor
		state.placeIndex = s.elevators[state.placeIndex].currentFloor
	default:
		panic("not in elevator")
	}
}

func (s *Simulation) PressButton(actorID int, floor int) {
	state := s.enteredActors[actorID]
	switch state.placeType {
	case PlaceElevator:
		elevator := s.elevators[state.placeIndex]
		s.controller.FloorSelected(state.placeIndex, floor)
		elevator.desiredFloors = append(elevator.desiredFloors, floor)
		s.dispatchControllerEvent(OnElevatorFloorRequest(s.tick,ElevatorID(state.placeIndex),FloorID(floor)))
	}
}

func (s *Simulation) elevatorDoneMoving(elevatorID int) {
	for i, a := range s.enteredActors {
		if a.placeType == PlaceElevator && a.placeIndex == elevatorID {
			s.actors[i].elevatorStopped(s, s.tick, s.elevators[elevatorID].currentFloor)
			s.dispatchControllerEvent(OnElevatorArrived(s.tick,ElevatorID(elevatorID),FloorID(s.elevators[elevatorID].currentFloor)))
		}
	}
	s.controller.CompletedMove(elevatorID)
}

func (s *Simulation) callElevator(floor int) {
	s.dispatchControllerEvent(OnElevatorCalled(s.tick, FloorID(floor)))
	s.controller.Called(floor)
}

func (s *Simulation) Initialize(elevators int, floors int) {
	s.dispatchControllerEvent(OnInitStart())
	s.elevators = make([]*Elevator, elevators)
	for i := range s.elevators {
		s.elevators[i] = NewElevator(5)
		s.dispatchControllerEvent(OnInformElevator(ElevatorID(i)))
	}
	s.floors = make([]*Floor, floors)
	for i := range s.floors {
		s.floors[i] = NewFloor()
		s.dispatchControllerEvent(OnInformFloor(FloorID(i)))
	}
	s.dispatchControllerEvent(OnInitDone())
}

//TODO: Multithreading
func (s *Simulation) dispatchControllerEvent(event Event) {
	for _, l := range s.controllerListeners {
		l.OnControllerEvent(event)
	}
}

//TODO: Multithreading
func (s *Simulation) AttachControllerListener(listener ControllerListener) {
	s.controllerListeners = append(s.controllerListeners, listener)
}

func NewSimulation() *Simulation {
	s := &Simulation{
		tick:                0,
		elevators:           make([]*Elevator, 0),
		floors:              make([]*Floor, 0),
		actors:              make([]*Actor, 0),
		controllerListeners: make([]ControllerListener, 0),
	}
	return s
}
