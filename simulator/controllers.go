package simulator

type Controller interface {
	Init(elevators []int)
	Called(floor int)
	FloorSelected(elevatorID int, floor int)
	CompletedMove(elevatorID int)
}

type ControlledElevators interface {
	// issues a move command to the elevator, returning ticks until completion
	MoveTo(elevatorID int, floorCount int)
}
