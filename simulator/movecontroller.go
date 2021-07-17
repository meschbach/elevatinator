package simulator

type MoveController struct {
	simulation ControlledElevators
	elevatorID int
}

func NewMoveController(elevators ControlledElevators) Controller {
	return &MoveController{
		simulation: elevators,
		elevatorID: -1,
	}
}

func (m *MoveController) Init(elevators []int) {
	m.elevatorID = elevators[0]
}

func (m *MoveController) Called(floor int) {
	m.simulation.MoveTo(m.elevatorID, floor)
}
func (m *MoveController) FloorSelected(elevatorID int, floor int) {
	m.simulation.MoveTo(m.elevatorID, floor)
}
func (m *MoveController) CompletedMove(elevatorID int) {}
