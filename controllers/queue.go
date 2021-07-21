package controllers

import (
	"fmt"
	"github.com/meschbach/elevatinator/simulator"
)

const (
	QueueIdle = iota
	QueuePickingUpCall
	QueueDroppingOff
)

type queueRequest struct {
	requestType int
	floor simulator.FloorID
}

type QueueController struct {
	state int
	id simulator.ElevatorID
	elevator simulator.ControlledElevators
	pending []queueRequest
}

func NewQueueController(elevators simulator.ControlledElevators) simulator.Controller {
	return &QueueController{
		state: QueueIdle,
		id: -1,
		elevator: elevators,
		pending: make([]queueRequest, 0),
	}
}

func (m *QueueController) Init(elevators []simulator.ElevatorID) {
	m.id = elevators[0]
}

func (m *QueueController) Called(floor simulator.FloorID) {
	fmt.Printf("Call at %d\n", floor)
	m.enqueueOrPerform(queueRequest{
		requestType: QueuePickingUpCall,
		floor:       floor,
	})
}

func (m *QueueController) FloorSelected(elevatorID simulator.ElevatorID, floor simulator.FloorID) {
	fmt.Printf("Elevator call...state: %d\n", m.state)
	m.enqueueOrPerform(queueRequest{
		requestType: QueueDroppingOff,
		floor:       floor,
	})
}

func (m *QueueController) CompletedMove(elevatorID simulator.ElevatorID) {
	fmt.Printf("Completed move to %d with %d pending\n", m.state, m.pending)
	m.dequeueOrIdle()
}

func (m *QueueController) perform(what queueRequest)  {
	fmt.Printf("Performing %#v", what)
	m.state = what.requestType
	m.elevator.MoveTo(m.id, what.floor)
}

func (m *QueueController) enqueueOrPerform(what queueRequest)  {
	if m.state == QueueIdle {
		m.perform(what)
	} else {
		m.pending = append(m.pending, what)
	}
}

func (m *QueueController) dequeueOrIdle() {
	if len(m.pending) > 0 {
		next := m.pending[0]
		m.pending = m.pending[1:]
		m.perform(next)
	} else {
		m.state = QueueIdle
	}
}
