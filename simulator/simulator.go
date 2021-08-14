package simulator

import "fmt"

type Tick int64

// Simulation encapsulates game state over time.
//
// An instance of Simulation is not guaranteed to be thread safe.
type Simulation struct {
	tick                Tick
	elevators           []*Elevator
	floors              []*Floor
	actors              []*Actor
	enteredActors       []*actorState
	controller          Controller
	controllerListeners []ControllerListener
}

// Tick advances the simulation by a single tick.  For each tick the following occurs:
//  * elevators are notified
//  * actors are notified
//
// True is returned if all existing actors have completed their objectives
func (s *Simulation) Tick() bool {
	currentTick := s.tick
	s.tick++
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

// TickUpTo advances the Simulation by up to the additional count of ticks or all actors have completed their objectives,
// which ever occurs first.
//
// Current simulation tick is returned as a result.
func (s *Simulation) TickUpTo(additional Tick) Tick {
	endTick := s.tick + additional
	for s.tick < endTick && s.Tick() {
	}
	return s.tick
}

// ActorsCompletedObjectives checks if all registered actors have completed their objectives.  If all actors have
// completed their objectives then true is returned, otherwise false.
func (s *Simulation) ActorsCompletedObjectives() bool {
	//TODO: Ideally there is a better way to structure this
	for _, actor := range s.actors {
		if !actor.done() {
			return false
		}
	}
	return true
}

// CurrentTick provides the tick the simulator is at.
//
// NOTE: An instance of Simulation provides no multithreading consistency guarantees.
func (s *Simulation) CurrentTick() Tick {
	return s.tick
}

func (s *Simulation) AttachActor(actor *Actor) {
	s.actors = append(s.actors, actor)
}

type ControllerFunc func(ControlledElevators) Controller

func (s *Simulation) AttachControllerFunc(factory ControllerFunc) {
	s.controller = factory(s)
	ids := make([]ElevatorID, len(s.elevators))
	for i := range s.elevators {
		ids[i] = ElevatorID(i)
	}
	s.controller.Init(ids)
}

func (s *Simulation) MoveTo(elevatorID ElevatorID, floor FloorID) {
	fmt.Printf("Simulation{tick: %d} -- Moving elevator %d to %d\n",s.tick, elevatorID,floor)
	elevator := s.elevators[elevatorID]
	elevator.moveTo(s, int(elevatorID), int(floor))
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
		s.controller.FloorSelected(ElevatorID(state.placeIndex), FloorID(floor))
		elevator.desiredFloors = append(elevator.desiredFloors, floor)
		s.dispatchControllerEvent(OnElevatorFloorRequest(s.tick,ElevatorID(state.placeIndex),FloorID(floor)))
	}
}

func (s *Simulation) elevatorDoneMoving(elevatorID ElevatorID) {
	for i, a := range s.enteredActors {
		if a.placeType == PlaceElevator && a.placeIndex == int(elevatorID) {
			s.actors[i].elevatorStopped(s, s.tick, s.elevators[elevatorID].currentFloor)
			s.dispatchControllerEvent(OnElevatorArrived(s.tick,elevatorID,FloorID(s.elevators[elevatorID].currentFloor)))
		}
	}
	fmt.Printf("Simulation{tick: %d} -- Elevator{id: %d} is at floor %d\n", s.tick, elevatorID, s.elevators[elevatorID].currentFloor)
	s.controller.CompletedMove(elevatorID)
}

func (s *Simulation) callElevator(floor int) {
	s.dispatchControllerEvent(OnElevatorCalled(s.tick, FloorID(floor)))
	s.controller.Called(FloorID(floor))
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
