package queue

import (
	"fmt"
	simulator2 "github.com/meschbach/elevatinator/pkg/simulator"
)

const (
	ControllerIdle = iota
	ControllerPickingUpCall
	ControllerDroppingOff
)

type request struct {
	requestType int
	floor       simulator2.FloorID
}

type Controller struct {
	state    int
	id       simulator2.ElevatorID
	elevator simulator2.ControlledElevators
	pending  []request
}

func NewController(elevators simulator2.ControlledElevators) simulator2.Controller {
	return &Controller{
		state:    ControllerIdle,
		id:       -1,
		elevator: elevators,
		pending:  make([]request, 0),
	}
}

func (m *Controller) Init(elevators []simulator2.ElevatorID) {
	m.id = elevators[0]
}

func (m *Controller) Called(floor simulator2.FloorID) {
	fmt.Printf("Call at %d\n", floor)
	m.enqueueOrPerform(request{
		requestType: ControllerPickingUpCall,
		floor:       floor,
	})
}

func (m *Controller) FloorSelected(elevatorID simulator2.ElevatorID, floor simulator2.FloorID) {
	fmt.Printf("Elevator call...state: %d\n", m.state)
	m.enqueueOrPerform(request{
		requestType: ControllerDroppingOff,
		floor:       floor,
	})
}

func (m *Controller) CompletedMove(elevatorID simulator2.ElevatorID) {
	fmt.Printf("Completed move to %d with %d pending\n", m.state, m.pending)
	m.dequeueOrIdle()
}

func (m *Controller) perform(what request) {
	fmt.Printf("Performing %#v", what)
	m.state = what.requestType
	m.elevator.MoveTo(m.id, what.floor)
}

func (m *Controller) enqueueOrPerform(what request) {
	if m.state == ControllerIdle {
		m.perform(what)
	} else {
		m.pending = append(m.pending, what)
	}
}

func (m *Controller) dequeueOrIdle() {
	if len(m.pending) > 0 {
		next := m.pending[0]
		m.pending = m.pending[1:]
		m.perform(next)
	} else {
		m.state = ControllerIdle
	}
}
