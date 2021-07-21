package simulator

type Controller interface {
	Init(elevators []ElevatorID)
	Called(floor FloorID)
	FloorSelected(elevatorID ElevatorID, floor FloorID)
	CompletedMove(elevatorID ElevatorID)
}

type ControlledElevators interface {
	// MoveTo instructs the given elevator to go to the specified target floor.
	MoveTo(elevatorID ElevatorID, floor FloorID)
}
