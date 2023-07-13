package simulator

// MoveController is designed to operate a single elevator with a single actor at a time.  Designed specifically for
// testing and not actual use.  Demonstrates basic interactions with of an Controller with the Simulation.
type MoveController struct {
	simulation ControlledElevators
	elevatorID ElevatorID
}

func NewMoveController(elevators ControlledElevators) Controller {
	return &MoveController{
		simulation: elevators,
		elevatorID: -1,
	}
}

func (m *MoveController) Init(elevators []ElevatorID) {
	m.elevatorID = elevators[0]
}

func (m *MoveController) Called(floor FloorID) {
	m.simulation.MoveTo(m.elevatorID, floor)
}
func (m *MoveController) FloorSelected(elevatorID ElevatorID, floor FloorID) {
	m.simulation.MoveTo(m.elevatorID, floor)
}
func (m *MoveController) CompletedMove(elevatorID ElevatorID) {}
