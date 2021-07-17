package simulator

type Controller interface {
	Init(elevators []int)
	Called(floor int)
	FloorSelected(elevatorID int, floor int)
	CompletedMove(elevatorID int)
}

type ControlledElevators interface {
	// MoveTo instructs the given elevator to go to the specified target floor.
	MoveTo(elevatorID int, floor int)
}
