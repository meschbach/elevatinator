package simulator

import "fmt"

const (
	Idle = iota
	MovingUp
	MovingDown
)

type Elevator struct {
	state       int
	moveToFloor int

	capacity     int8
	currentFloor int
	//TODO: Probably better as event stream for controller
	desiredFloors []int
}

func NewElevator(capacity int8) *Elevator {
	return &Elevator{
		state:         Idle,
		capacity:      capacity,
		currentFloor:  0,
		desiredFloors: make([]int, 0),
	}
}

func (e *Elevator) Tick(s *Simulation, id int, tick Tick) {
	switch e.state {
	case MovingUp:
		e.currentFloor++
		e.maybeDoneMoving(s, id)
	case MovingDown:
		e.currentFloor--
		e.maybeDoneMoving(s, id)
	}
}

func (e *Elevator) maybeDoneMoving(s *Simulation, id int) {
	if e.currentFloor == e.moveToFloor {
		fmt.Printf("Elevator{id: %d} -- Finished moving to floor %d\n", id, e.currentFloor)
		e.state = Idle
		s.elevatorDoneMoving(ElevatorID(id))
	}
}

func (e *Elevator) moveTo(s *Simulation, id int, floor int) {
	distance := floor - e.currentFloor
	e.move(s, id, distance)
}

func (e *Elevator) move(s *Simulation, id int, floors int) {
	switch e.state {
	case Idle:
		e.moveToFloor = e.currentFloor + floors
		if floors > 0 {
			e.state = MovingUp
		} else if floors < 0 {
			e.state = MovingDown
		} else {
			e.maybeDoneMoving(s, id)
		}
	}
}

func (e *Elevator) isAtFloor(s *Simulation, floor FloorID) bool {
	switch e.state {
	case Idle:
		return e.currentFloor == int(floor)
	}
	return false
}
